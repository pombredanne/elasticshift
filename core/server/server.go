/*
Copyright 2017 The Elasticshift Authors.
*/
package server

import (
	"encoding/base64"
	"log"
	"net/http/pprof"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/build"
	"gitlab.com/conspico/elasticshift/core/store"
	"gitlab.com/conspico/elasticshift/identity/client"
	"gitlab.com/conspico/elasticshift/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/identity/team"
	"gitlab.com/conspico/elasticshift/identity/user"
	"gitlab.com/conspico/elasticshift/sysconf"
	"gitlab.com/conspico/elasticshift/vcs"
	"gitlab.com/conspico/elasticshift/vcs/repository"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	mgo "gopkg.in/mgo.v2"
)

// Constants for performing encode decode
const (
	EQUAL        = "="
	DOUBLEEQUALS = "=="
	DOT0         = ".0"
	DOT1         = ".1"
	DOT2         = ".2"
	SLASH        = "/"
	SEMICOLON    = ";"
)

// Server ..
type Server struct {
	Logger    logrus.Logger
	Store     store.Store
	Router    *mux.Router
	Dex       dex.DexClient
	Providers providers.Providers
}

// Config ..
type Config struct {
	Store    store.Config
	Logger   logrus.Logger
	Session  *mgo.Session
	Identity Identity
}

// Identity ..
type Identity struct {
	Issuer      string
	HostAndPort string
	caPath      string
	ID          string
	Secret      string
	RedirectURI string
}

// New ..
// Creates a new server
func New(ctx context.Context, c Config) (*Server, error) {

	s := &Server{}

	s.Logger = c.Logger
	s.Store = store.New(c.Store.Name, c.Session)

	// d, err := newDexClient(ctx, c.Identity)
	// if err != nil {
	// 	return nil, err
	// }
	// s.Dex = d

	r := mux.NewRouter()
	s.Router = r

	// pprof
	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/symbol", pprof.Symbol)
	r.HandleFunc("/debug/profile", pprof.Profile)
	r.Handle("/debug/heap", pprof.Handler("heap"))
	r.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/block", pprof.Handler("block"))

	s.Providers = providers.New(s.Logger, sysconf.NewStore(s.Store))

	// initialize oauth2 providers
	s.registerOauth2Providers()

	// initialize graphql based services
	s.registerGraphQLServices()

	// err := NewAuthServer(ctx, r, c)
	// if err != nil {
	// 	return nil, err
	// }
	return s, nil
}

func (s Server) registerOauth2Providers() {

	teamStore := team.NewStore(s.Store)
	vcsServ := vcs.NewService(s.Logger, s.Store, s.Providers, teamStore)

	s.Router.HandleFunc("/{team}/link/{provider}", vcsServ.Authorize)
	s.Router.HandleFunc("/link/{provider}/callback", vcsServ.Authorized)
}

func (s Server) registerGraphQLServices() {

	logger := s.Logger
	r := s.Router

	// initialize schema
	queries := graphql.Fields{}
	mutations := graphql.Fields{}

	// data store
	teamStore := team.NewStore(s.Store)
	vcsStore := vcs.NewStore(s.Store)
	sysconfStore := sysconf.NewStore(s.Store)
	repositoryStore := repository.NewStore(s.Store)
	buildStore := build.NewStore(s.Store)

	// team fields
	teamQ, teamM := team.InitSchema(logger, teamStore)
	appendFields(queries, teamQ)
	appendFields(mutations, teamM)

	// vcs fields
	vcsQ, vcsM := vcs.InitSchema(logger, s.Providers, vcsStore, teamStore, repositoryStore)
	appendFields(queries, vcsQ)
	appendFields(mutations, vcsM)

	// sysconf fields
	sysconfQ, sysconfM := sysconf.InitSchema(logger, sysconfStore)
	appendFields(queries, sysconfQ)
	appendFields(mutations, sysconfM)

	// sysconf fields
	buildQ, buildM := build.InitSchema(logger, buildStore, repositoryStore, sysconfStore)
	appendFields(queries, buildQ)
	appendFields(mutations, buildM)

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queries}
	rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: mutations}

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create schema due to errors :v", err)
	}

	// initialize graphql
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	r.Handle("/graphql", h)
}

// Utility method to append fields
func appendFields(fields graphql.Fields, input graphql.Fields) {

	for k, v := range input {
		fields[k] = v
	}
}

func RegisterGRPCServices(grpcServer *grpc.Server, s *Server) {

	api.RegisterUserServer(grpcServer, user.NewServer(s.Logger, s.Dex))
	api.RegisterClientServer(grpcServer, client.NewServer(s.Store, s.Logger, s.Dex))
}

// Registers the exposed http services
func RegisterHTTPServices(ctx context.Context, router *mux.Router, grpcAddress string, dialopts []grpc.DialOption) error {

	return nil
	// err := api.RegisterUserHandlerFromEndpoint(ctx, router, grpcAddress, dialopts)
	// if err != nil {
	// 	return fmt.Errorf("Registering User handler failed : %v", err)
	// }

	// err = api.RegisterClientHandlerFromEndpoint(ctx, router, grpcAddress, dialopts)
	// if err != nil {
	// 	return fmt.Errorf("Registering Client handler failed : %v", err)
	// }

	// return err
}

//func newDexClient(ctx context.Context, c Dex) (dex.DexClient, error) {
//	// creds, err := credentials.NewClientTLSFromFile(caPath, "")
//	// if err != nil {
//	//     return nil, fmt.Errorf("load dex cert: %v", err)
//	// }

//	//conn, err := grpc.Dial(hostAndPort, grpc.WithTransportCredentials(creds))

//	conn, err := grpc.Dial(c.HostAndPort, grpc.WithInsecure())
//	defer func() {
//		if err != nil {
//			if cerr := conn.Close(); cerr != nil {
//				grpclog.Printf("Failed to close conn to %s: %v", c.HostAndPort, cerr)
//			}
//			return
//		}
//		go func() {
//			<-ctx.Done()
//			if cerr := conn.Close(); cerr != nil {
//				grpclog.Printf("Failed to close conn to %s: %v", c.HostAndPort, cerr)
//			}
//		}()
//	}()
//	return dex.NewDexClient(conn), nil
//}

func encode(id string) string {

	eid := base64.URLEncoding.EncodeToString([]byte(id))
	if strings.Contains(eid, DOUBLEEQUALS) {
		eid = strings.TrimRight(eid, DOUBLEEQUALS) + DOT2
	} else if strings.Contains(eid, EQUAL) {
		eid = strings.TrimRight(eid, EQUAL) + DOT1
	} else {
		eid = eid + DOT0
	}
	return eid
}

func decode(id string) string {

	if strings.Contains(id, DOT2) {
		id = strings.TrimRight(id, DOT2) + DOUBLEEQUALS
	} else if strings.Contains(id, DOT1) {
		id = strings.TrimRight(id, DOT1) + EQUAL
	} else {
		id = strings.TrimRight(id, DOT0)
	}
	did, _ := base64.URLEncoding.DecodeString(id)
	return string(did[:])
}