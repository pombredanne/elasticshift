/*
Copyright 2017 The Elasticshift Authors.
*/
package build

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/container"
	"github.com/graphql-go/graphql"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/cloudprovider/docker"
	"gitlab.com/conspico/elasticshift/pkg/sysconf"
	"gitlab.com/conspico/elasticshift/pkg/vcs/repository"
	"gopkg.in/mgo.v2/bson"
)

var (
	errIDCantBeEmpty           = errors.New("Build ID cannot be empty")
	errRepositoryIDCantBeEmpty = errors.New("Repository ID is must in order to trigger the build")
	errInvalidRepositoryID     = errors.New("Please provide the valid repository ID")

	logfile = "logfile"
)

const (
	LogType_Embedded = "Embedded"
	LogType_NFS      = "NFS"
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

			env := []string{
				"SHIFT_HOST=127.0.0.1",
				"SHIFT_PORT=5050",
				"SHIFT_LOGGER=" + LogType_Embedded,
				"SHIFT_BUILDID=" + b.ID.Hex(),
			}

			c := &container.Config{
				Image: "alpine",

				// Cmd:   []string{"/bin/sh"},
				// Entrypoint: strslice.StrSlice{"/bin/sh"},
				Env: env,
			}
			containerID, err := cli.CreateContainer(c, b.ID.Hex())
			if err != nil {
				str := fmt.Sprintf("Unable to create the container %v", err)
				r.SLog(b.ID, str)
			}

			fmt.Println("Container ID =", containerID)
			err = r.store.UpdateContainerID(b.ID, containerID)
			if err != nil {
				r.logger.Errorln("Failed to update the container id: ", containerID)
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

	status := types.BuildStatus_Running
	rb, err := r.store.FetchBuild(repo.Team, repository_id, branch, types.BuildStatus_Running)
	if err != nil {
		return nil, fmt.Errorf("Failed to validate if there are any build running", err)
	}

	if len(rb) > 0 {
		status = types.BuildStatus_Waiting
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
	b.LogType = LogType_Embedded

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

	if types.BuildStatus_Cancelled == b.Status || types.BuildStatus_Failed == b.Status || types.BuildStatus_Success == b.Status {
		return fmt.Sprintf("Cancelling the build is not possible, because it seems that it was already %s", b.Status.String()), nil
	}

	// TODO trigger the cancel build, only if the current status is RUNNING | WAITING | STUCK

	err = r.store.UpdateBuildStatus(b.ID, types.BuildStatus_Cancelled)
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
