// Package server ..
// Author Ghazni Nattarshah
// Date: 1/10/17
package server

import (
	"fmt"
	"net/http/pprof"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/store"
	"golang.org/x/net/context"

	mgo "gopkg.in/mgo.v2"
)

// Server ..
type Server struct {
	Logger logrus.FieldLogger
	Store  store.Store
	Router *mux.Router
	Dex    dex.DexClient
}

// Config ..
type Config struct {
	Store   store.Config
	Logger  logrus.FieldLogger
	Session *mgo.Session
	Dex     Dex
}

// Dex ..
type Dex struct {
	Issuer      string
	HostAndPort string
	caPath      string
	ID          string
	Secret      string
	RedirectURI string
}

// NewServer ..
// Creates a new server
func NewServer(ctx context.Context, c Config) (*Server, error) {

	s := &Server{}

	if c.Logger == nil {
		return nil, fmt.Errorf("No logger found")
	}
	s.Logger = c.Logger

	s.Store = store.NewStore(c.Store.Name, c.Session)

	d, err := newDexClient(ctx, c.Dex)
	if err != nil {
		return nil, err
	}
	s.Dex = d

	r := mux.NewRouter()

	// pprof
	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/symbol", pprof.Symbol)
	r.HandleFunc("/debug/profile", pprof.Profile)
	r.Handle("/debug/heap", pprof.Handler("heap"))
	r.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/block", pprof.Handler("block"))

	s.Router = r

	err = NewAuthServer(ctx, r, c)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func newDexClient(ctx context.Context, c Dex) (dex.DexClient, error) {
	// creds, err := credentials.NewClientTLSFromFile(caPath, "")
	// if err != nil {
	//     return nil, fmt.Errorf("load dex cert: %v", err)
	// }

	//conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(c.HostAndPort, grpc.WithInsecure())
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", c.HostAndPort, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Printf("Failed to close conn to %s: %v", c.HostAndPort, cerr)
			}
		}()
	}()
	return dex.NewDexClient(conn), nil
}
