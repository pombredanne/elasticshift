// Package server ..
// Author Ghazni Nattarshah
// Date: 2/11/17
package server

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/pb"
	"gitlab.com/conspico/elasticshift/store"
)

var (
	errNoClient           = errors.New("No client provided")
	errNoClientIDProvided = errors.New("No client id provided")

	errDeleteClient = errors.New("Failed to delete the client")
)

type clientServer struct {
	store  store.ClientStore
	logger logrus.FieldLogger
}

// NewClientServer ..
// Implementation of pb.UserServer
func NewClientServer(s *Server) pb.ClientServer {
	return &clientServer{
		store:  store.NewClientStore(s.Store),
		logger: s.Logger,
	}
}

// Create ..
func (s *clientServer) Create(ctx context.Context, req *pb.CreateReq) (*pb.CreateRes, error) {

	if req.Client == nil {
		return nil, errNoClient
	}

	c := store.Client{}
	c.ID = store.NewID()
	c.Secret = store.NewID() + store.NewID()
	c.Name = req.Client.Name
	c.RedirectURIs = req.Client.RedirectUris
	c.Public = req.Client.Public
	c.TrustedPeers = req.Client.TrustedPeers

	exist, err := s.store.Exist(c.Name)
	if err != nil {
		s.logger.Errorf("failed to insert client : %v", err)
		return nil, errors.New("Failed to create client")
	}

	res := &pb.CreateRes{
		Exists: exist,
	}

	if exist {
		return res, nil
	}

	err = s.store.Insert(&c)
	if err != nil {
		s.logger.Errorf("failed to insert client : %v", err)
		return nil, errDeleteClient
	}

	req.Client.Id = c.ID
	req.Client.Secret = c.Secret
	res.Client = req.Client

	return res, nil
}

// Delete ..
func (s *clientServer) Delete(ctx context.Context, req *pb.DeleteReq) (*pb.DeleteRes, error) {

	if req.ClientId == "" {
		return nil, errNoClientIDProvided
	}

	err := s.store.Delete(req.ClientId)
	if err != nil {
		s.logger.Errorf("Failed to delete the client %s : %v", req.ClientId, err)
		return nil, errDeleteClient
	}

	return &pb.DeleteRes{}, nil
}
