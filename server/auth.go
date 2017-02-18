// Package server ..
// Author Ghazni Nattarshah
// Date: 2/11/17
package server

import (
	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/pb"
	"gitlab.com/conspico/elasticshift/store"
)

type authServer struct {
	store  store.AuthStore
	cache  *Cache
	logger logrus.FieldLogger
}

// NewAuthServer ..
// Implementation of pb.UserServer
func NewAuthServer(s *Server) pb.AuthServer {
	return &authServer{
		store:  store.NewAuthStore(s.Store),
		cache:  s.Cache,
		logger: s.Logger,
	}
}

// Authorize ..
func (s *authServer) Authorize(ctx context.Context, req *pb.AuthorizeReq) (*pb.AuthorizeRes, error) {

	return nil, nil
}
