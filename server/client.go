// Package server ..
// Author Ghazni Nattarshah
// Date: 2/11/17
package server

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/store"
)

var (
	errNoClient           = errors.New("No client provided")
	errNoClientIDProvided = errors.New("No client id provided")
	errFailedCreateClient = errors.New("Failed to create the client")
	errDeleteClient       = errors.New("Failed to delete the client")
)

type clientServer struct {
	store  store.ClientStore
	logger logrus.FieldLogger
	dex    dex.DexClient
}

// NewClientServer ..
// Implementation of api.UserServer
func NewClientServer(s *Server) api.ClientServer {
	return &clientServer{
		store:  store.NewClientStore(s.Store),
		dex:    s.Dex,
		logger: s.Logger,
	}
}

// Create ..
func (s *clientServer) Create(ctx context.Context, req *api.CreateClientReq) (*api.CreateClientRes, error) {

	in := &dex.CreateClientReq{}
	in.Client = &dex.Client{}
	in.Client.Name = req.Name
	in.Client.RedirectUris = req.RedirectUris
	in.Client.Public = req.Public
	in.Client.TrustedPeers = req.TrustedPeers
	in.Client.LogoUrl = req.LogoUrl

	out, err := s.dex.CreateClient(ctx, in)
	if err != nil {
		return nil, errFailedCreateClient
	}

	res := &api.CreateClientRes{}
	res.Id = out.Client.Id
	res.Secret = out.Client.Secret
	res.Name = out.Client.Name
	res.Public = out.Client.Public
	res.TrustedPeers = out.Client.TrustedPeers
	res.LogoUrl = out.Client.LogoUrl
	res.RedirectUris = out.Client.RedirectUris

	return res, nil
}

// Delete ..
func (s *clientServer) Delete(ctx context.Context, req *api.DeleteClientReq) (*api.DeleteClientRes, error) {

	if req.ClientId == "" {
		return nil, errNoClientIDProvided
	}

	in := &dex.DeleteClientReq{}
	in.Id = req.ClientId

	out, err := s.dex.DeleteClient(ctx, in)
	if err != nil {
		return nil, errDeleteClient
	}

	return &api.DeleteClientRes{NotFound: out.NotFound}, nil
}
