// Package server ..
// Author Ghazni Nattarshah
// Date: 1/10/17
package server

import (
	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gitlab.com/conspico/elasticshift/store"
	"golang.org/x/net/context"

	"time"

	mgo "gopkg.in/mgo.v2"
)

// Server ..
type Server struct {
	Issuer *url.URL
	Logger logrus.FieldLogger
	Store  store.Store
	Mux    *runtime.ServeMux
	Cache  *Cache
}

// Config ..
type Config struct {
	Issuer             string
	Store              store.Config
	Logger             logrus.FieldLogger
	Session            *mgo.Session
	SignerKeysLifeSpan time.Duration
	IDTokensLifeSpan   time.Duration
}

// NewServer ..
// Creates a new server
func NewServer(ctx context.Context, c Config) (*Server, error) {

	s := &Server{}

	iss, err := url.Parse(c.Issuer)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse the issuer URL %s: %v", c.Issuer, err)
	}
	s.Issuer = iss

	if c.Logger == nil {
		return nil, fmt.Errorf("No logger found")
	}
	s.Logger = c.Logger

	s.Store = store.NewStore(c.Store.Name, c.Session)

	s.Cache = newCache(store.NewKeyStore(s.Store))

	s.StartKeySpinner(ctx, c)
	return s, nil
}
