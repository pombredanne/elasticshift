/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"context"

	"github.com/elasticshift/elasticshift/api"
	"github.com/elasticshift/elasticshift/internal/worker/types"
)

type server struct {
	ctx types.Context
}

func NewServer(ctx types.Context) api.WorkServer {
	return &server{ctx}
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
