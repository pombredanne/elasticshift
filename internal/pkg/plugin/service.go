/*
Copyright 2018 The Elasticshift Authors.
*/
package plugin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/store"
	"gitlab.com/conspico/elasticshift/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

const (
	PLUGIN_DIR = "plugins"
	BUNDLE_EXT = ".bundle"
)

var (
	errTeamNotFound = "Team doesn't exist"
	errPluginExist  = "Plugin with the name and version already exist"
)

type service struct {
	pluginStore  store.Plugin
	teamStore    store.Team
	sysconfStore store.Sysconf
	logger       logrus.Logger
}

// SysconfService ..
type Service interface {
	PushPlugin(w http.ResponseWriter, r *http.Request)
}

// NewVCSService ..
func NewService(logger logrus.Logger, d store.Database, s store.Shift) Service {

	return &service{
		pluginStore:  s.Plugin,
		teamStore:    s.Team,
		sysconfStore: s.Sysconf,
		logger:       logger,
	}
}

func (s service) PushPlugin(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	// data validation
	errors := make([]string, 0)
	name := r.FormValue("name")
	if name == "" {
		errors = append(errors, "Plugin name cannot be empty")
	}

	description := r.FormValue("description")
	if description == "" {
		errors = append(errors, "Short description is required")
	}

	language := r.FormValue("language")
	if language == "" {
		errors = append(errors, "Language should be provided")
	}

	version := r.FormValue("version")
	if version == "" {
		errors = append(errors, "Plugin name cannot be empty")
	}

	author := r.FormValue("author")
	if author == "" {
		errors = append(errors, "Author name cannot be empty")
	}

	email := r.FormValue("email")
	if email == "" {
		errors = append(errors, "Email cannot be empty")
	}

	teamName := r.FormValue("team")
	if teamName == "" {
		errors = append(errors, "Team name must be provided")
	}

	if len(errors) > 0 {
		http.Error(w, strings.Join(errors, "\n"), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("plugin")
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	size, _ := strconv.Atoi(r.Header.Get("Content-Length"))
	if size <= 0 {
		http.Error(w, "multipart plugin file not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// validate the team
	team, err := s.teamStore.GetTeam("", teamName)
	fmt.Printf("Error : %v", err)
	if err != nil && err.Error() != "not found" {
		http.Error(w, "Failed to verify the team :"+err.Error(), http.StatusInternalServerError)
		return
	}

	if team.Name == "" {
		http.Error(w, errTeamNotFound, http.StatusBadRequest)
		return
	}

	// check the plugin existence, to see if it's a newer version
	var p types.Plugin
	err = s.pluginStore.FindOne(bson.M{"name": name, "version": version, "team": teamName}, &p)
	if err != nil && err.Error() == "not found" {
		// plugin doesnot exist, so allow new one
		err = nil
	}

	if err != nil {
		http.Error(w, "Failed to find the plugin: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if p.ID.Hex() != "" {
		http.Error(w, "Plugin with given name and version already exist on your team.", http.StatusBadRequest)
		return
	}

	// find the system storage
	result, err := s.sysconfStore.GetDefaultStorage()
	if err != nil {
		http.Error(w, "Failed to fetch the default storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if result.Kind == "" {
		http.Error(w, "Default storage has not been selected, Use web to configure it.", http.StatusInternalServerError)
		return
	}

	// upload the file to system storage and extract them.
	pluginDir := filepath.Join(result.Path, PLUGIN_DIR)
	pluginTeamDir := filepath.Join(pluginDir, teamName)

	err = utils.Mkdir(pluginTeamDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	plugfile, err := os.Create(filepath.Join(pluginTeamDir, name+"-"+version+BUNDLE_EXT))
	if err != nil {
		s.logger.Errorf("Failed to write plugin bundle to storage :%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(plugfile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer plugfile.Close()

	// store the data to plugin store
	plug := &types.Plugin{}
	plug.ID = bson.NewObjectId()
	plug.Name = name
	plug.Team = teamName
	plug.Description = description
	plug.Language = language
	plug.Version = version
	plug.Author = author
	plug.Email = email

	if sourceURL := r.FormValue("source_url"); sourceURL != "" {
		plug.SourceURL = sourceURL
	}

	err = s.pluginStore.Save(plug)
	if err != nil {
		s.logger.Errorf("Failed to save plugin for team %s: %v", teamName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.StatusText(http.StatusOK)
}
