// Package server ..
// Author Ghazni Nattarshah
// Date: 1/10/17
package server

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/store"
	"golang.org/x/net/context"

	mgo "gopkg.in/mgo.v2"
)

// Server ..
type Server struct {
	Logger logrus.FieldLogger
	Store  store.Store
	Mux    *runtime.ServeMux
	Dex    dex.DexClient
}

// Config ..
type Config struct {
	Store   store.Config
	Logger  logrus.FieldLogger
	Session *mgo.Session
	Dex     dex.DexClient
}

// NewServer ..
// Creates a new server
func NewServer(ctx context.Context, c Config) (*Server, error) {

	s := &Server{}

	if c.Logger == nil {
		return nil, fmt.Errorf("No logger found")
	}
	s.Logger = c.Logger
	s.Dex = c.Dex
	s.Store = store.NewStore(c.Store.Name, c.Session)

	return s, nil
}
