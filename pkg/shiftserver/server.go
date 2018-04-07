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
	"gitlab.com/conspico/elasticshift/pkg/build"
	"gitlab.com/conspico/elasticshift/pkg/container"
	"gitlab.com/conspico/elasticshift/pkg/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/pkg/identity/team"
	"gitlab.com/conspico/elasticshift/pkg/plugin"
	"gitlab.com/conspico/elasticshift/pkg/shift"
	"gitlab.com/conspico/elasticshift/pkg/store"
	stypes "gitlab.com/conspico/elasticshift/pkg/store/types"
	"gitlab.com/conspico/elasticshift/pkg/sysconf"
	"gitlab.com/conspico/elasticshift/pkg/vcs"
	"gitlab.com/conspico/elasticshift/pkg/vcs/repository"
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
	DB        stypes.Database
	Router    *mux.Router
	Dex       dex.DexClient
	Providers providers.Providers
	Ctx       context.Context

	TeamStore       team.Store
	VCSStore        vcs.Store
	SysConfStore    sysconf.Store
	BuildStore      build.Store
	RepositoryStore repository.Store
	PluginStore     plugin.Store
	ContainerStore  container.Store
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
	s.Ctx = ctx
	s.Logger = c.Logger

	s.DB = store.NewDatabase(c.Store.Name, c.Session)

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

	s.Providers = providers.New(s.Logger, sysconf.NewStore(s.DB))

	// initialize oauth2 providers
	s.registerEndpointServices()

	// initialize graphql based services
	s.registerGraphQLServices()

	// err := NewAuthServer(ctx, r, c)
	// if err != nil {
	// 	return nil, err
	// }
	return s, nil
}

func (s Server) registerEndpointServices() {

	teamStore := team.NewStore(s.DB)
	vcsServ := vcs.NewService(s.Logger, s.DB, s.Providers, teamStore)

	// Oauth2 providers
	s.Router.HandleFunc("/{team}/link/{provider}", vcsServ.Authorize)
	s.Router.HandleFunc("/link/{provider}/callback", vcsServ.Authorized)

	// Sysconf Upload kube file
	sysconfServ := sysconf.NewService(s.Logger, s.DB, teamStore)
	s.Router.HandleFunc("/sysconf/upload", sysconfServ.UploadKubeConfigFile)
}

func (s Server) registerGraphQLServices() {

	logger := s.Logger
	r := s.Router

	// initialize schema
	queries := graphql.Fields{}
	mutations := graphql.Fields{}

	// data store
	teamStore := team.NewStore(s.DB)
	s.TeamStore = teamStore

	vcsStore := vcs.NewStore(s.DB)
	s.VCSStore = vcsStore

	sysconfStore := sysconf.NewStore(s.DB)
	s.SysConfStore = sysconfStore

	repositoryStore := repository.NewStore(s.DB)
	s.RepositoryStore = repositoryStore

	buildStore := build.NewStore(s.DB)
	s.BuildStore = buildStore

	pluginStore := plugin.NewStore(s.DB)
	s.PluginStore = pluginStore

	containerStore := container.NewStore(s.DB)
	s.ContainerStore = containerStore

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

	// build fields
	buildQ, buildM := build.InitSchema(logger, s.Ctx, buildStore, repositoryStore, sysconfStore)
	appendFields(queries, buildQ)
	appendFields(mutations, buildM)

	// app fields
	pluginQ, pluginM := plugin.InitSchema(logger, s.Ctx, pluginStore)
	appendFields(queries, pluginQ)
	appendFields(mutations, pluginM)

	// container fields
	containerQ, containerM := container.InitSchema(logger, s.Ctx, containerStore)
	appendFields(queries, containerQ)
	appendFields(mutations, containerM)

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

// Registers the GRPC services
func RegisterGRPCServices(grpcServer *grpc.Server, s *Server) {
	api.RegisterShiftServer(grpcServer, shift.NewServer(s.Logger, s.Ctx, s.BuildStore, s.RepositoryStore))
}

// Registers the exposed http services
func RegisterHTTPServices(ctx context.Context, router *mux.Router, grpcAddress string, dialopts []grpc.DialOption) error {
	return nil
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
