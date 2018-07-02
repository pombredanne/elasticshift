/*
Copyright 2017 The Elasticshift Authors.
*/
package build

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/keys"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/parser"
	"gitlab.com/conspico/elasticshift/internal/pkg/vcs"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty           = errors.New("Build ID cannot be empty")
	errRepositoryIDCantBeEmpty = errors.New("Repository ID is must in order to trigger the build")
	errInvalidRepositoryID     = errors.New("Please provide the valid repository ID")

	logfile = "logfile"
)

// Resolver ...
type Resolver interface {
	TriggerBuild(params graphql.ResolveParams) (interface{}, error)
	FetchBuild(params graphql.ResolveParams) (interface{}, error)
	CancelBuild(params graphql.ResolveParams) (interface{}, error)
	FetchBuildByID(params graphql.ResolveParams) (interface{}, error)

	SLog(id interface{}, log string) error
	Log(id interface{}, log types.Log) error
}

type resolver struct {
	store            store.Build
	repositoryStore  store.Repository
	sysconfStore     store.Sysconf
	teamStore        store.Team
	integrationStore store.Integration
	defaultStore     store.Defaults
	shiftfileStore   store.Shiftfile
	logger           logrus.Logger
	Ctx              context.Context
	BuildQueue       chan types.Build
}

// NewResolver ...
func NewResolver(ctx context.Context, logger logrus.Logger, s store.Shift) (Resolver, error) {

	r := &resolver{
		store:            s.Build,
		repositoryStore:  s.Repository,
		sysconfStore:     s.Sysconf,
		teamStore:        s.Team,
		integrationStore: s.Integration,
		defaultStore:     s.Defaults,
		shiftfileStore:   s.Shiftfile,
		logger:           logger,
		Ctx:              ctx,
		BuildQueue:       make(chan types.Build),
	}

	// Launch a background process to launch container after build trigger.
	go r.ContainerLauncher()

	return r, nil
}

func (r *resolver) TriggerBuild(params graphql.ResolveParams) (interface{}, error) {

	repositoryID, _ := params.Args["repository_id"].(string)
	if repositoryID == "" {
		return nil, errRepositoryIDCantBeEmpty
	}

	repo, err := r.repositoryStore.GetRepositoryByID(repositoryID)
	if err != nil {
		return nil, errInvalidRepositoryID
	}

	branch, _ := params.Args["branch"].(string)
	if branch == "" {
		branch = repo.DefaultBranch
	}

	// Check if default container engine is set
	def, err := r.defaultStore.FindByReferenceId(repo.Team)
	if err != nil {
		return nil, err
	}

	if def.ContainerEngineID == "" {
		return nil, fmt.Errorf("No default container engine found, please configure it.")
	}

	status := types.BS_RUNNING
	rb, err := r.store.FetchBuild(repo.Team, repositoryID, branch, types.BS_RUNNING)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate if there are any build running", err)
	}

	if len(rb) > 0 {
		status = types.BS_WAITING
	}

	b := types.Build{}
	b.ID = bson.NewObjectId()
	b.RepositoryID = repositoryID
	b.ContainerEngineID = def.ContainerEngineID
	b.VcsID = repo.VcsID
	b.Status = status
	b.TriggeredBy = "Anonymous" //TODO fill in with logged-in user
	b.StartedAt = time.Now()
	b.Team = repo.Team
	b.Branch = branch
	b.StorageID = def.StorageID
	b.CloneURL = repo.CloneURL
	b.Language = repo.Language

	// Build file path - (for NFS)
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/log
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/reports
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/archive.zip
	// cache must be mounted as /elasticshift to containers
	b.StoragePath = filepath.Join(repo.Team, repo.Identifier, repo.Name, branch, b.ID.Hex())

	err = r.store.Save(&b)
	if err != nil {
		return nil, fmt.Errorf("Failed to save build details: %v", err)
	}

	if b.Status == types.BS_RUNNING {

		// Pass the build data to builder
		r.BuildQueue <- b
	}

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

	f, err := r.repoImageName(b)
	if err != nil && f == nil {

		// falling back to see if there is a team's default configured
		defs, err := r.defaultStore.FindByReferenceId(b.Team)
		if err != nil {
			return "", fmt.Errorf("Failed to fetch defaults by reference id: %v", err)
		}

		var file types.Shiftfile
		err = r.shiftfileStore.FindByID(defs.Languages[b.Language], &file)
		if err != nil {
			return "", fmt.Errorf("Failed to fetch the default shiftfile for language: %v", err)
		}
		f = file.File
	}
	// write shift file

	sf, err := parser.AST([]byte(f))
	if err != nil {
		r.SLog(b.ID, fmt.Sprintf("Failed to parse shift file: %v", err))
	}

	return sf.Image()[keys.NAME].(string), nil
}

func (r *resolver) repoImageName(b types.Build) ([]byte, error) {

	var source string
	if strings.Contains(b.CloneURL, vcs.GITHUB_DOT_COM) {
		source = vcs.GITHUB_DOT_COM
	}

	f, err := vcs.GetShiftFile(source, b.CloneURL, b.Branch)
	if err != nil {
		return nil, err
	}

	return f, nil
}
