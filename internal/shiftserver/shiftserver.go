/*
Copyright 2017 The Elasticshift Authors.
*/
package shiftserver

import (
	"encoding/base64"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/identity/oauth2/providers"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/integration"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/plugin"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/pubsub"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/schema"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/secret"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/shift"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/vcs"
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
	Config ServerConfig

	Loggr     logger.Loggr
	Logger    *logrus.Entry
	DB        store.Database
	Router    *mux.Router
	Dex       dex.DexClient
	Providers providers.Providers
	Ctx       context.Context

	NSQ pubsub.NSQConfig

	Shift store.Shift

	Vault secret.Vault

	Pubsub pubsub.Engine
}

// ServerConfig ..
type ServerConfig struct {
	Store    store.Config
	Logger   logger.Loggr
	NSQ      NSQ
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
	s.Config = c
	s.Ctx = ctx
	s.Loggr = c.Logger
	s.Logger = s.Loggr.GetLogger("shiftserver")
	s.NSQ.Consumer.Address = c.NSQ.ConsumerAddress
	s.NSQ.Producer.Address = c.NSQ.ProducerAddress

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

	s.Providers = providers.New(s.Loggr, s.Shift)

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
	vcsServ := vcs.NewService(s.Loggr, s.DB, s.Providers, s.Shift, s.Vault)

	// Oauth2 providers
	s.Router.HandleFunc("/api/{team}/link/{provider}", vcsServ.Authorize)
	s.Router.HandleFunc("/api/link/{provider}/callback", vcsServ.Authorized)

	// TODO the directory is only applicable for dev testing
	// s.Router.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("/Users/ghazni/.elasticshift/cloud"))))

	s.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

	// TODO remove kubeconfig
	// Sysconf Upload kube file
	integrationServ := integration.NewService(s.Loggr, s.DB, s.Shift)
	s.Router.HandleFunc("/api/integration/kubernetes", integrationServ.UploadKubeConfigFile)

	// Plugin bundle push
	pluginServ := plugin.NewService(s.Loggr, s.DB, s.Shift)
	s.Router.HandleFunc("/api/plugin/push", pluginServ.PushPlugin)
}

func (s *Server) registerGraphQLServices() error {

	r := s.Router

	vault := secret.NewVault(s.Shift, s.Loggr, s.Ctx)
	s.Vault = vault

	nsqc := pubsub.NSQConfig{}
	nsqc.Consumer.Address = s.Config.NSQ.ConsumerAddress
	nsqc.Producer.Address = s.Config.NSQ.ProducerAddress
	s.NSQ = nsqc

	cons := pubsub.NewConsumers(nsqc, s.Loggr)

	// subscription handler through websocket
	sh := pubsub.NewSubscriptionHandler(s.Loggr, cons)

	eng := pubsub.NewEngine(s.Loggr, sh, nsqc, cons)
	s.Pubsub = eng

	schm, err := schema.Construct(s.Ctx, s.Loggr, s.Providers, s.Shift, s.Vault, s.Pubsub)
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

	// This is really important to validate the schema
	// during subscription, and also used when pushing
	// the results to consumers
	s.Pubsub.Schema(schm)

	// Graphql endpoint works with websocket only for subscription
	psh := pubsub.NewGraphqlWSHandler(s.Pubsub, s.Loggr)
	r.Handle("/subscription", psh)

	return nil
}

func (s *Server) registerWebSocketServices() {

	//r := s.Router
	//r.HandleFunc("/api/ws/")

}

// Registers the GRPC services ...
func RegisterGRPCServices(grpcServer *grpc.Server, s *Server) {
	api.RegisterShiftServer(grpcServer, shift.NewServer(s.Loggr, s.Ctx, s.Shift, s.Vault, s.Pubsub))
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
