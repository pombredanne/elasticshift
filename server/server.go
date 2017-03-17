// Package server ..
// Author Ghazni Nattarshah
// Date: 1/10/17
package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/pprof"
	"regexp"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/dex"
	"gitlab.com/conspico/elasticshift/store"
	"golang.org/x/net/context"

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
	Logger logrus.FieldLogger
	Store  store.Store
	Router *http.ServeMux
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

	//r := mux.NewRouter()
	r := http.NewServeMux()

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

// isAlphaNumericOnly ..
// Check to see if the given text is alpha-numeric only
func isAlphaNumericOnly(str string) bool {
	matched, _ := regexp.MatchString("^[A-Za-z0-9]*$", str)
	return matched
}
