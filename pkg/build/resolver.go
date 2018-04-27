/*
Copyright 2017 The Elasticshift Authors.
*/
package build

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/cloudprovider/docker"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/keys"
	"gitlab.com/conspico/elasticshift/pkg/shiftfile/parser"
	"gitlab.com/conspico/elasticshift/pkg/sysconf"
	"gitlab.com/conspico/elasticshift/pkg/utils"
	"gitlab.com/conspico/elasticshift/pkg/vcs/repository"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty           = errors.New("Build ID cannot be empty")
	errRepositoryIDCantBeEmpty = errors.New("Repository ID is must in order to trigger the build")
	errInvalidRepositoryID     = errors.New("Please provide the valid repository ID")

	logfile = "logfile"

	// place holders vcs_account, repository, branch
	RAW_GUTHUB_URL = "https://raw.githubusercontent.com/%s/%s/%s/Shiftfile"

	DIR_CODE    = "code"
	DIR_PLUGINS = "plugins"
	DIR_WORKER  = "worker"
	DIR_LOGS    = "logs"

	// TODO check for windows container
	VOL_SHIFT   = "/shift"
	VOL_CODE    = filepath.Join(VOL_SHIFT, DIR_CODE)
	VOL_PLUGINS = filepath.Join(VOL_SHIFT, DIR_PLUGINS)
	VOL_LOGS    = filepath.Join(VOL_SHIFT, DIR_LOGS)
)

const (
	LogType_Embedded = "embedded"
	LogType_File     = "file"
	LogType_NFS      = "nfs"
)

type resolver struct {
	store           Store
	repositoryStore repository.Store
	sysconfStore    sysconf.Store
	logger          logrus.Logger
	Ctx             context.Context
	BuildQueue      chan types.Build
}

func (r *resolver) ContainerLauncher() {

	for b := range r.BuildQueue {

		go func(b types.Build) {

			// start the container
			// TODO select the default orchestration, by config
			opts := &docker.ClientOptions{}
			opts.Host = docker.DefaultHost
			opts.Ctx = r.Ctx

			cli, err := docker.NewClient(opts)
			if err != nil {
				r.SLog(b.ID, fmt.Sprintf("Failed to connect to docker daemon: %v", err))
			}

			imgName, err := r.findImageName(b)
			if err != nil {
				r.SLog(b.ID, fmt.Sprintf("Unable to find the build image from Shiftfile", b.CloneURL))
			}
			fmt.Println("Image name: " + imgName)

			// find the system storage
			storage, err := r.sysconfStore.GetDefaultStorage()
			if err != nil {
				r.SLog(b.ID, "Failed to fetch the default storage: "+err.Error())
				return
			}

			err = utils.Mkdir(filepath.Join(storage.Path, "code", b.Team))
			if err != nil {
				r.SLog(b.ID, "Unable to create directory for cloning the project:"+err.Error())
			}

			hostIp := utils.GetIP()
			if hostIp == "" {
				hostIp = "127.0.0.1"
			}

			env := []string{
				"SHIFT_HOST=shiftserver",
				"SHIFT_PORT=5051",
				"SHIFT_LOGGER=" + LogType_File,
				"SHIFT_BUILDID=" + b.ID.Hex(),
				"SHIFT_TIMEOUT=120m",
				"WORKER_PORT=" + "6060",
			}

			filepath.Join(storage.Path, b.Team, DIR_CODE)

			hc := &container.HostConfig{}
			hc.Binds = []string{
				filepath.Join(storage.Path, b.Team, DIR_CODE) + ":" + VOL_CODE,
				filepath.Join(storage.Path, b.Team, DIR_LOGS) + ":" + VOL_LOGS,
				filepath.Join(storage.Path, DIR_PLUGINS) + ":" + VOL_PLUGINS,
				filepath.Join(storage.Path, DIR_WORKER) + ":" + VOL_SHIFT,
			}

			workerPort, _ := nat.NewPort("tcp", "6060")
			serverPort, _ := nat.NewPort("tcp", "5051")

			exposedPorts := map[nat.Port]struct{}{
				serverPort: struct{}{},
				workerPort: struct{}{},
			}

			c := &container.Config{
				Image:        imgName,
				Entrypoint:   strslice.StrSlice{"./shift/worker"},
				Env:          env,
				AttachStdout: true,
				ExposedPorts: exposedPorts,
			}

			containerID, err := cli.CreateContainer(c, hc, b.ID.Hex())
			if err != nil {
				str := fmt.Sprintf("Unable to create the container %v", err)
				r.SLog(b.ID, str)
			}

			fmt.Println("Container ID =", containerID)
			err = r.store.UpdateContainerID(b.ID, containerID)
			if err != nil {
				r.logger.Errorln("Failed to update the container id: ", containerID)
			}

			err = cli.StartContainer(containerID)
			if err != nil {
				r.logger.Errorln("Failed to start the container: %v", err)
			}
		}(b)
	}
}

func (r *resolver) TriggerBuild(params graphql.ResolveParams) (interface{}, error) {

	repository_id, _ := params.Args["repository_id"].(string)
	if repository_id == "" {
		return nil, errRepositoryIDCantBeEmpty
	}

	repo, err := r.repositoryStore.GetRepositoryByID(repository_id)
	if err != nil {
		return nil, errInvalidRepositoryID
	}

	branch, _ := params.Args["branch"].(string)
	if branch == "" {
		branch = repo.DefaultBranch
	}

	status := types.BS_RUNNING
	rb, err := r.store.FetchBuild(repo.Team, repository_id, branch, types.BS_RUNNING)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate if there are any build running", err)
	}

	if len(rb) > 0 {
		status = types.BS_WAITING
	}

	b := types.Build{}
	b.ID = bson.NewObjectId()
	b.RepositoryID = repository_id
	b.VcsID = repo.VcsID
	b.Status = status
	b.TriggeredBy = "Anonymous" //TODO fill in with logged-in user
	b.StartedAt = time.Now()
	b.Team = repo.Team
	b.Branch = branch
	b.LogType = LogType_File
	b.CloneURL = repo.CloneURL

	// Build file path - (for NFS)
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/log
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/reports
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/archive.zip
	// cache must be mounted as /elasticshift to containers
	// b.Log = filepath.Join(repo.Team, repo.Identifier, repo.Name, branch, b.ID.Hex(), logfile)

	err = r.store.Save(&b)
	if err != nil {
		return nil, fmt.Errorf("Failed to save build details: %v", err)
	}

	// Pass the build data to builder
	r.BuildQueue <- b

	return b, err
}

func (r *resolver) SLog(id interface{}, log string) error {
	return r.Log(id, types.Log{Time: time.Now(), Data: log})
}

func (r *resolver) Log(id interface{}, log types.Log) error {
	return r.store.UpdateId(id, bson.M{"$push": bson.M{"log": log}})
}

func (r *resolver) FetchBuild(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	repository_id, _ := params.Args["repository_id"].(string)
	branch, _ := params.Args["branch"].(string)
	status, _ := params.Args["status"].(int)

	result := types.BuildList{}
	res, err := r.store.FetchBuild(team, repository_id, branch, types.BuildStatus(status))
	if err != nil {
		return result, fmt.Errorf("Failed to fetch the build : %v", err)
	}

	result.Nodes = res
	result.Count = len(res)

	return result, nil
}

func (r *resolver) CancelBuild(params graphql.ResolveParams) (interface{}, error) {

	res, err := r.FetchBuildByID(params)
	if err != nil && !strings.EqualFold("not found", err.Error()) {
		return nil, fmt.Errorf("Failed to cancel the build : %v", err)
	}

	b := res.(types.Build)
	if b.ID == "" {
		return nil, fmt.Errorf("Build id not found")
	}

	if types.BS_CANCELLED == b.Status || types.BS_FAILED == b.Status || types.BS_SUCCESS == b.Status {
		return fmt.Sprintf("Cancelling the build is not possible, because it seems that it was already %s", b.Status.String()), nil
	}

	// TODO trigger the cancel build, only if the current status is RUNNING | WAITING | STUCK

	err = r.store.UpdateBuildStatus(b.ID, types.BS_CANCELLED)
	if err != nil {
		return nil, fmt.Errorf("Failed to cancel the build: %v", err)
	}
	return nil, nil
}

func (r *resolver) FetchBuildByID(params graphql.ResolveParams) (interface{}, error) {

	id, _ := params.Args["id"].(string)
	if id == "" {
		return nil, errIDCantBeEmpty
	}

	res, err := r.store.FetchBuildByID(id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *resolver) findImageName(b types.Build) (string, error) {

	var file []byte
	// https://github.com/nshahm/hybrid.test.runner/raw/master/Shiftfile
	if strings.Contains(b.CloneURL, "github.com") {

		// repoUrl := fmt.Sprintf(RAW_GUTHUB_URL, )
		repoUrl := strings.TrimRight(b.CloneURL, ".git")
		repoUrl += "/raw/" + b.Branch + "/Shiftfile"

		fmt.Println("Repo URL: " + repoUrl)

		resp, err := http.Get(repoUrl)
		if err != nil {
			return "", err
		}

		fmt.Println("Raw URL:" + resp.Request.URL.String())
		resp, err = http.Get(resp.Request.URL.String())
		if err != nil {
			return "", err
		}

		// read the response body
		file, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
	}

	// r.logger.Infoln("Response = ", string(file[:]))

	// write shift file

	sf, err := parser.AST(file)
	if err != nil {
		r.SLog(b.ID, fmt.Sprintf("Failed to parse shift file: %v", err))
	}

	return sf.Image()[keys.NAME].(string), nil
}
