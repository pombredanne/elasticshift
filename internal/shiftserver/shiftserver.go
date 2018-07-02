/*
Copyright 2017 The Elasticshift Authors.
*/
package shiftserver

import (
	"encoding/base64"
	"net/http/pprof"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/plugin"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/secret"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/vcs"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/schema"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/shift"
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
	DB        store.Database
	Router    *mux.Router
	Dex       dex.DexClient
	Providers providers.Providers
	Ctx       context.Context

	Shift store.Shift

	Vault secret.Vault

	BuildType *graphql.Object
}

// ServerConfig ..
type ServerConfig struct {
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
func New(ctx context.Context, c ServerConfig) (*Server, error) {

	s := &Server{}
	s.Ctx = ctx
	s.Logger = c.Logger

	s.DB = store.NewDatabase(c.Store.Name, c.Session)
	s.Shift = s.DB.InitShiftStore()

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

	s.Providers = providers.New(s.Logger, s.Shift)

	// initialize graphql based services
	err := s.registerGraphQLServices()
	if err != nil {
		return nil, err
	}

	// initialize oauth2 providers
	s.registerEndpointServices()

	s.registerWebSocketServices()

	err = s.bootstrap()
	if err != nil {
		return nil, err
	}

	// err := NewAuthServer(ctx, r, c)
	// if err != nil {
	// 	return nil, err
	// }

	return s, nil
}

func (s *Server) registerEndpointServices() {

	// VCS service to link repositories.
	vcsServ := vcs.NewService(s.Logger, s.DB, s.Providers, s.Shift, s.Vault)

	// Oauth2 providers
	s.Router.HandleFunc("/api/{team}/link/{provider}", vcsServ.Authorize)
	s.Router.HandleFunc("/api/link/{provider}/callback", vcsServ.Authorized)

	// TODO remove kubeconfig
	// Sysconf Upload kube file
	// sysconfServ := sysconf.NewService(s.Logger, s.DB, s.TeamStore)
	// s.Router.HandleFunc("/sysconf/upload", sysconfServ.UploadKubeConfigFile)

	// Plugin bundle push
	pluginServ := plugin.NewService(s.Logger, s.DB, s.Shift)
	s.Router.HandleFunc("/api/plugin/push", pluginServ.PushPlugin)
}

func (s *Server) registerGraphQLServices() error {

	logger := s.Logger
	r := s.Router

	vault := secret.NewVault(s.Shift, logger, s.Ctx)
	s.Vault = vault

	schm, err := schema.Construct(s.Ctx, s.Logger, s.Providers, s.Shift, s.Vault)
	if err != nil {
		return err
	}

	// initialize graphql
	h := handler.New(&handler.Config{
		Schema:   schm,
		Pretty:   true,
		GraphiQL: true,
	})
	r.Handle("/graphql", h)
	r.Handle("/graphql/", h)

	return nil
}

func (s *Server) registerWebSocketServices() {

	//r := s.Router
	//r.HandleFunc("/api/ws/")

}

// Registers the GRPC services
func RegisterGRPCServices(grpcServer *grpc.Server, s *Server) {
	api.RegisterShiftServer(grpcServer, shift.NewServer(s.Logger, s.Ctx, s.Shift, s.Vault))
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
