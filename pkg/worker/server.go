/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"context"

	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/pkg/worker/logshipper"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
)

type server struct {
	ctx    types.Context
	logger logshipper.Logger
}

func NewServer(ctx types.Context, logger logshipper.Logger) api.WorkServer {
	return &server{ctx, logger}
}

func (s *server) Top(req *api.TopReq, stream api.Work_TopServer) error {
	return nil
}

func (s *server) KillTask(ctx context.Context, req *api.KillTaskReq) (*api.KillTaskRes, error) {
	return nil, nil
}

func (s *server) StopBuild(ctx context.Context, req *api.StopBuildReq) (*api.StopBuildReq, error) {
	return nil, nil
}
