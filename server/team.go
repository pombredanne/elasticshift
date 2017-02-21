// Package server ..
// Teamor Ghazni Nattarshah
// Date: 2/11/17
package server

import (
	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/store"
)

type teamServer struct {
	store  store.TeamStore
	logger logrus.FieldLogger
}

// NewTeamServer ..
// Implementation of api.UserServer
func NewTeamServer(s *Server) api.TeamServer {
	return &teamServer{
		store:  store.NewTeamStore(s.Store),
		logger: s.Logger,
	}
}

// Teamorize ..
func (s *teamServer) Create(context.Context, *api.CreateTeamReq) (*api.CreateTeamRes, error) {

	return nil, nil
}
