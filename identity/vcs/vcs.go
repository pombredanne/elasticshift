/*
Copyright 2017 The Elasticshift Authors.
*/
package vcs

import (
	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/core/store"
	"gitlab.com/conspico/elasticshift/identity/team"
	"golang.org/x/net/context"
)

// VCSServer ..
//type VCSServer interface {
//authorize(w http.ResponseWriter, r *http.Request)
//handleAuthorizeCallback(w http.ResponseWriter, r *http.Request)
//GetVCS(teamID string) (GetVCSResponse, error)
//SyncVCS(r SyncVCSRequest) (bool, error)
//}

type server struct {
	store  team.Store
	logger logrus.FieldLogger
}

// NewVCSServer ..
// Implementation of api.VCSServer
func NewVCSServer(s store.Store, logger logrus.FieldLogger) api.VCSServer {
	return &server{
		store:  team.NewStore(s),
		logger: logger,
	}
}

func (s server) GetVCS(ctx context.Context, req *api.GetVCSReq) (*api.GetVCSRes, error) {
	return nil, nil
}

func (s server) Sync(ctx context.Context, req *api.SyncVCSReq) (*api.SyncVCSRes, error) {
	return nil, nil
}
