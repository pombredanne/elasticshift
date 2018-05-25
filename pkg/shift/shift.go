/*
Copyright 2018 The Elasticshift Authors.
*/
package shift

import (
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/build"
	"gitlab.com/conspico/elasticshift/pkg/secret"
	"gitlab.com/conspico/elasticshift/pkg/vcs/repository"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

type shift struct {
	logger          logrus.Logger
	Ctx             context.Context
	buildStore      build.Store
	repositoryStore repository.Store
	vault           secret.Vault
}

func NewServer(logger logrus.Logger, ctx context.Context, buildStore build.Store, repositoryStore repository.Store, vault secret.Vault) api.ShiftServer {
	return &shift{logger, ctx, buildStore, repositoryStore, vault}
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

func (s *shift) LogShip(reqStream api.Shift_LogShipServer) error {

	for {

		in, err := reqStream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		logTime, err := ptypes.Timestamp(in.GetTime())
		if err != nil {
			return err
		}

		log := types.Log{
			Time: logTime,
			Data: in.GetLog(),
		}

		err = s.buildStore.UpdateId(in.GetBuildId(), bson.M{"$push": bson.M{"log": log}})
		if err != nil {
			return err
		}
	}
}

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

	return res, nil
}
