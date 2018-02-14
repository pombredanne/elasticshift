package shift

import (
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/ptypes"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/build"
	"gitlab.com/conspico/elasticshift/pkg/vcs/repository"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2/bson"
)

type shift struct {
	logger          logrus.Logger
	Ctx             context.Context
	buildStore      build.Store
	repositoryStore repository.Store
}

func NewServer(logger logrus.Logger, ctx context.Context, buildStore build.Store, repositoryStore repository.Store) api.ShiftServer {
	return &shift{logger, ctx, buildStore, repositoryStore}
}

func (s *shift) Register(ctx context.Context, req *api.RegisterReq) (*api.RegisterRes, error) {
	return nil, nil
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
