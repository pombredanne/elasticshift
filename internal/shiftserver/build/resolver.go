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
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/ast"
	"gitlab.com/conspico/elasticshift/internal/pkg/shiftfile/parser"
	"gitlab.com/conspico/elasticshift/internal/pkg/vcs"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/pubsub"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty           = errors.New("Build ID cannot be empty")
	errRepositoryIDCantBeEmpty = errors.New("Repository ID is must in order to trigger the build")
	errInvalidRepositoryID     = errors.New("Please provide the valid repository ID")

	logfile      = "logfile"
	defaultGraph = `[
  {
    "node": {
      "name": "START",
      "description": "Start node of the execution/graph",
      "status": "N"
    }
  },
  {
    "node": {
      "name": "END",
      "description": "End node of the execution/graph",
      "status": "N"
    }
  }
]`
)

// Resolver ...
type Resolver interface {
	TriggerBuild(params graphql.ResolveParams) (interface{}, error)
	FetchBuild(params graphql.ResolveParams) (interface{}, error)
	CancelBuild(params graphql.ResolveParams) (interface{}, error)
	FetchBuildByID(params graphql.ResolveParams) (interface{}, error)

	SLog(id interface{}, log string) error
	Log(id interface{}, log types.Log) error
	TriggerNextIfAny(teamID, repositoryID, branch string)
}

type resolver struct {
	store            store.Build
	repositoryStore  store.Repository
	sysconfStore     store.Sysconf
	teamStore        store.Team
	integrationStore store.Integration
	defaultStore     store.Defaults
	shiftfileStore   store.Shiftfile
	logger           *logrus.Entry
	loggr            logger.Loggr
	Ctx              context.Context
	BuildQueue       chan types.Build
	ps               pubsub.Engine
}

// NewResolver ...
func NewResolver(ctx context.Context, loggr logger.Loggr, s store.Shift, ps pubsub.Engine) (Resolver, error) {

	r := &resolver{
		store:            s.Build,
		repositoryStore:  s.Repository,
		sysconfStore:     s.Sysconf,
		teamStore:        s.Team,
		integrationStore: s.Integration,
		defaultStore:     s.Defaults,
		shiftfileStore:   s.Shiftfile,
		logger:           loggr.GetLogger("graphql/build"),
		loggr:            loggr,
		Ctx:              ctx,
		BuildQueue:       make(chan types.Build),
		ps:               ps,
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
		return nil, errors.New("No default container engine found, please configure it.")
	}

	status := types.BuildStatusPreparing
	rb, err := r.store.FetchBuild(repo.Team, repositoryID, branch, "", []string{types.BuildStatusPreparing, types.BuildStatusRunning})
	if err != nil {
		return nil, fmt.Errorf("Failed to validate if there are any build running: %v", err)
	}

	if len(rb) > 0 {
		status = types.BuildStatusWaiting
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
	b.Source = repo.Source
	b.Graph = defaultGraph

	buildID := b.ID.Hex()

	// Build file path - (for NFS)
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/log
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/reports
	// <cache>/team-id/vcs-id/repository-id/branch-name/build-id/archive.zip
	// cache must be mounted as /elasticshift to containers
	b.StoragePath = filepath.Join(repo.Team, repo.Identifier, repo.Name, branch, buildID)

	err = r.store.Save(&b)
	if err != nil {
		return nil, fmt.Errorf("Failed to save build details: %v", err)
	}

	if b.Status == types.BuildStatusPreparing {
		r.pushToQueue(b)
	}

	return b, err
}

func (r *resolver) TriggerNextIfAny(teamID, repositoryID, branch string) {

	query := bson.M{
		"team":          teamID,
		"repository_id": repositoryID,
		"branch":        branch,
		"status":        types.BuildStatusWaiting,
	}

	var b types.Build
	var err error
	r.store.Execute(func(c *mgo.Collection) {
		err = c.Find(query).Sort("-started_at").Limit(1).One(&b)
	})

	if err != nil && err.Error() != "not found" {
		r.logger.Errorf("Trigger next build failed. [TeamID: %s, RepositoryID = %s, Branch = %s]\n", teamID, repositoryID, branch)
	} else if err != nil && err.Error() == "not found" {
		return
	} else {

		if b.ID != "" {

			// update the status to preparing
			r.store.UpdateBuildStatus(b.ID, types.BuildStatusPreparing)

			// post to build queue
			r.pushToQueue(b)
		}
	}
}

func (r *resolver) pushToQueue(b types.Build) {

	// publish to topic to push build updates to subscribers
	r.ps.Publish(pubsub.SubscribeBuildUpdate, b.ID.Hex())

	// Pass the build data to builder
	r.BuildQueue <- b
}

func (r *resolver) UpdateBuildStatusAsFailed(id, reason string, endedAt time.Time) {
	r.UpdateBuildStatus(id, reason, types.BuildStatusFailed, endedAt)
}

func (r *resolver) UpdateBuildStatus(id, reason, status string, endedAt time.Time) {

	// should be the container startup log, if startup failed
	// should be the err log if build failed.
	var b types.Build
	err := r.store.FindByID(id, &b)
	if err != nil {
		// TODO handler error
	}

	b.Reason = reason
	b.Status = status
	b.EndedAt = time.Now()

	err = r.store.UpdateId(b.ID, b)
	if err != nil {
		r.logger.Errorf("failed to update build status: %v", err)
	}
}

func (r *resolver) SLog(id interface{}, log string) error {
	return r.Log(id, types.Log{Time: time.Now(), Data: log})
}

func (r *resolver) UpdateReason(id interface{}, reason string) error {
	return r.store.UpdateId(id, bson.M{"$set": bson.M{"reason": reason}})
}

func (r *resolver) Log(id interface{}, log types.Log) error {
	return r.store.UpdateId(id, bson.M{"$push": bson.M{"log": log}})
}

func (r *resolver) FetchBuild(params graphql.ResolveParams) (interface{}, error) {

	team, _ := params.Args["team"].(string)
	repository_id, _ := params.Args["repository_id"].(string)
	branch, _ := params.Args["branch"].(string)
	id, _ := params.Args["id"].(string)
	statusParam, _ := params.Args["status"].(int)

	statusArr := []string{}
	if statusParam > 0 {

		var status string
		switch statusParam {
		case 1: // waiting
			status = types.BuildStatusWaiting
		case 2: // preparing
			status = types.BuildStatusPreparing
		case 3: // running
			status = types.BuildStatusRunning
		case 4: // success
			status = types.BuildStatusSuccess
		case 5: // failed
			status = types.BuildStatusFailed
		case 6: // cancelled
			status = types.BuildStatusCancel
		case 7: // stuck
			status = types.BuildStatusStuck
		}

		statusArr = append(statusArr, status)
	}

	result := types.BuildList{}
	res, err := r.store.FetchBuild(team, repository_id, branch, id, statusArr)
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

	if types.BuildStatusCancel == b.Status || types.BuildStatusFailed == b.Status || types.BuildStatusSuccess == b.Status {
		return fmt.Sprintf("Cancelling the build is not possible, because it seems that it was already %s", b.Status), nil
	}

	// TODO trigger the cancel build, only if the current status is RUNNING | WAITING | STUCK

	err = r.store.UpdateBuildStatus(b.ID, types.BuildStatusCancel)
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

func (r *resolver) GetShiftfile(b types.Build) (*ast.File, bool, error) {

	var repoFile bool
	repoFile = true
	f, err := r.repoImageName(b)
	if err != nil && f == nil {

		repoFile = false

		// falling back to see if there is a team's default configured
		defs, err := r.defaultStore.FindByReferenceId(b.Team)
		if err != nil {
			return nil, repoFile, fmt.Errorf("Failed to fetch defaults by reference id: %v", err)
		}

		fileId := defs.Languages[b.Language]
		if fileId == "" {
			return nil, repoFile, fmt.Errorf("No default shiftfile configured for language [%s].", b.Language)
		}

		var file types.Shiftfile
		err = r.shiftfileStore.FindByID(fileId, &file)
		if err != nil {
			return nil, repoFile, fmt.Errorf("Failed to fetch the default shiftfile for language: %v", err)
		}
		f = file.File
	}
	// write shift file

	sf, err := parser.AST([]byte(f))
	if err != nil {
		r.SLog(b.ID, fmt.Sprintf("Failed to parse shift file: %v", err))
	}

	return sf, repoFile, nil
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
