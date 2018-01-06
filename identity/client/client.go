/*
Copyright 2017 The Elasticshift Authors.
*/
package client

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/dex"
	core "gitlab.com/conspico/elasticshift/core/store"
)

var (
	errNoClient           = errors.New("No client provided")
	errNoClientIDProvided = errors.New("No client id provided")
	errFailedCreateClient = errors.New("Failed to create the client")
	errDeleteClient       = errors.New("Failed to delete the client")
)

type server struct {
	store  Store
	logger logrus.Logger
	dex    dex.DexClient
}

// NewServer ..
// Implementation of api.UserServer
func NewServer(s core.Store, logger logrus.Logger, dex dex.DexClient) api.ClientServer {
	return &server{
		store:  NewStore(s),
		dex:    dex,
		logger: logger,
	}
}

// Create ..
func (s *server) Create(ctx context.Context, req *api.InsertClientReq) (*api.InsertClientRes, error) {

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

	res := &api.InsertClientRes{}
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
func (s *server) Delete(ctx context.Context, req *api.RemoveClientReq) (*api.RemoveClientRes, error) {

	if req.ClientId == "" {
		return nil, errNoClientIDProvided
	}

	in := &dex.DeleteClientReq{}
	in.Id = req.ClientId

	out, err := s.dex.DeleteClient(ctx, in)
	if err != nil {
		return nil, errDeleteClient
	}

	return &api.RemoveClientRes{NotFound: out.NotFound}, nil
}
