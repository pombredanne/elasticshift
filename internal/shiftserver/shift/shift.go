/*
Copyright 2018 The Elasticshift Authors.
*/
package shift

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/integration"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/pubsub"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/secret"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

type shift struct {
	logger           logrus.Logger
	Ctx              context.Context
	buildStore       store.Build
	containerStore   store.Container
	repositoryStore  store.Repository
	defaultStore     store.Defaults
	integrationStore store.Integration
	vault            secret.Vault
	ps               pubsub.Engine
}

func NewServer(logger logrus.Logger, ctx context.Context, s store.Shift, vault secret.Vault, ps pubsub.Engine) api.ShiftServer {
	return &shift{logger, ctx, s.Build, s.Container, s.Repository, s.Defaults, s.Integration, vault, ps}
}

func (s *shift) Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterRes, error) {

	s.logger.Println("Registration request for build " + req.GetBuildId())
	if req.GetBuildId() == "" {
		return nil, fmt.Errorf("Registration failed: Build ID cannot be empty.")
	}

	if req.GetPrivatekey() == "" {
		return nil, fmt.Errorf("Registration failed: No key provided")
	}

	// TODO store the secret key id in build and the actual key in secret store
	buildId := bson.ObjectIdHex(req.GetBuildId())
	err := s.buildStore.UpdateId(buildId, bson.M{"$push": bson.M{"private_key": req.GetPrivatekey()}})
	if err != nil {
		return nil, fmt.Errorf("Registration failed: Due to internal server error %v", err)
	}

	res := &api.RegisterRes{}
	res.Registered = true

	return res, nil
}

func (s *shift) UpdateBuildStatus(ctx context.Context, req *api.UpdateBuildStatusReq) (*api.UpdateBuildStatusRes, error) {

	if req == nil {
		return nil, fmt.Errorf("UpdateBuildGraphReq cannot be nil")
	}

	if req.GetBuildId() == "" {
		return nil, fmt.Errorf("BuildID is empty")
	}

	res := &api.UpdateBuildStatusRes{}

	var b types.Build
	var err error
	err = s.buildStore.FindByID(req.GetBuildId(), &b)
	if err != nil {
		return res, fmt.Errorf("Failed to fetch build by id : %v", err)
	}

	b.Graph = req.GetGraph()
	status := req.GetStatus()
	cp := req.GetCheckpoint()

	var stopContainer bool
	if status != "" {

		if status == "F" {
			b.Status = types.BS_FAILED
			stopContainer = true
		} else if status == "S" && cp == "END" {
			b.Status = types.BS_SUCCESS
			stopContainer = true
		} else if status == "C" {
			b.Status = types.BS_CANCELLED
		}
		b.EndedAt = time.Now()
	}

	err = s.buildStore.UpdateId(b.ID, &b)
	if err != nil {
		return res, fmt.Errorf("Failed to update the graph : %v", err)
	}

	// publish pubsub to fetch latest update to subscribers
	s.ps.Publish(pubsub.SubscribeBuildUpdate, req.GetBuildId())

	if stopContainer {

		fmt.Println("-------------------------------------------------------------")
		fmt.Println("Stopping the container..... ")
		fmt.Println("Status = ", status)
		fmt.Println("Checkpoint = ", cp)
		fmt.Println("-------------------------------------------------------------")
		// request container engine to stop the live container
		ce, err := s.getContainerEngine(b.Team)
		if err != nil {
			return res, errors.Wrap(err, "Failed to get the default container engine: %v")
		}

		err = ce.DeleteContainer(req.GetBuildId())
		if err != nil {
			return res, fmt.Errorf("Failed to stop the container: %v", err)
		}
	}

	return res, nil
}

func (s *shift) getContainerEngine(team string) (integration.ContainerEngineInterface, error) {

	// Get the default container engine id based on team
	// dce, err := s.defaultStore.GetDefaultContainerEngine(team)
	def, err := s.defaultStore.FindByReferenceId(team)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get default container engine:")
	}

	// Get the details of the integration
	var i types.ContainerEngine
	err = s.integrationStore.FindByID(def.ContainerEngineID, &i)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get default integration:")
	}

	// Get the details of the storeage
	var stor types.Storage
	err = s.integrationStore.FindByID(def.StorageID, &stor)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get default storage:")
	}

	// connect to container engine cluster
	return integration.NewContainerEngine(s.logger, i, stor)
}

// func (s *shift) LogShip(reqStream api.Shift_LogShipServer) error {

// 	for {

// 		in, err := reqStream.Recv()
// 		if err == io.EOF {
// 			return nil
// 		}

// 		if err != nil {
// 			return err
// 		}

// 		logTime, err := ptypes.Timestamp(in.GetTime())
// 		if err != nil {
// 			return err
// 		}

// 		log := types.Log{
// 			Time: logTime,
// 			Data: in.GetLog(),
// 		}

// 		err = s.buildStore.UpdateId(in.GetBuildId(), bson.M{"$push": bson.M{"log": log}})
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

func (s *shift) GetProject(ctx context.Context, req *api.GetProjectReq) (*api.GetProjectRes, error) {

	if req == nil {
		return nil, fmt.Errorf("GetProjectReq cannot be nil")
	}

	if req.BuildId == "" {
		return nil, fmt.Errorf("BuildID is empty")
	}

	var b types.Build
	err := s.buildStore.FindByID(req.BuildId, &b)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch : %v", err)
	}

	r, err := s.repositoryStore.GetRepositoryByID(b.RepositoryID)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch the repository: %v", err)
	}

	res := &api.GetProjectRes{}
	res.VcsId = b.VcsID
	res.Branch = b.Branch
	res.CloneUrl = r.CloneURL
	res.Language = r.Language
	res.Name = r.Name
	res.StoragePath = b.StoragePath
	res.Source = b.Source

	return res, nil
}
