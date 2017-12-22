/*
Copyright 2017 The Elasticshift Authors.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"log"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gitlab.com/conspico/elasticshift/core/server"
	"gitlab.com/conspico/elasticshift/core/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

// Config ..
type Config struct {
	Store    Store    `json:"store"`
	Web      Web      `json:"web"`
	Logger   Logger   `json:"logger"`
	Identity Identity `json:"dex"`
}

// Store ..
type Store struct {
	Server    string
	Name      string
	Username  string
	Password  string
	Timeout   string
	Monotonic bool
	Retry     string // duration

	// old info
	IdleConnection int
	MaxConnection  int
	Log            bool
}

// Web ..
// Holds the web server configuration
type Web struct {
	HTTP string `json:"http"`
	GRPC string `json:"grpc"`
}

// Dex ..
// Holds the web server configuration
type Identity struct {
	GRPC        string `json:"grpc"`
	Issuer      string `json:"issuer"`
	ID          string `json:"id"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
}

// Logger ..x``
type Logger struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// issuer: http://elasticshift.com/armor

// store:
//     server: 127.0.0.1:27017
//     name: armor
//     username: armor
//     password: armorpazz
//     monotonic: true
//     timeout: "10s"
//     retry: "10s"

// logger:
//   level: "debug"
//   format: "json"

// web:
//   #static: dist
//   http: 127.0.0.1:5050
//   grpc: 127.0.0.1:5051

func main() {

	err := elasticshift()
	log.Fatalln(err)
	os.Exit(-1)
}

func elasticshift() error {

	cfgFile := "/etc/conspico/elasticshift/config.yaml"
	cfgData, err := ioutil.ReadFile(cfgFile)

	var c Config
	switch {
	case os.IsNotExist(err):
		c = Config{

			Store: Store{

				Server:    "127.0.0.1",
				Name:      "elasticshift",
				Username:  "elasticshift",
				Password:  "3l@$t1c$h1ft",
				Monotonic: true,
				Timeout:   "10s",
				Retry:     "10s",
			},

			Logger: Logger{

				Level:  "debug",
				Format: "json",
			},

			Web: Web{
				HTTP: "127.0.0.1:5050",
				GRPC: "127.0.0.1:5051",
			},

			Identity: Identity{
				GRPC:        "127.0.0.1:5557",
				Issuer:      "http://127.0.0.1:5556/Identity",
				ID:          "yyjw66rn2hso6wriuzlic62jiy",
				Secret:      "l77r6wixjjtgmo4iym2kmk3jcuuxetj3afnqaw5w3rnl5nu5hehu",
				RedirectURI: "http://127.0.0.1:5050/login/callback",
			},
		}
	case err == nil:
		if err := yaml.Unmarshal(cfgData, &c); err != nil {
			return fmt.Errorf("Failed to parse config file %s: %v", cfgFile, err)
		}
	default:
		log.Println(err)
	}

	// override configuration with environment variables.
	if storeServer := os.Getenv("STORE_SERVER"); storeServer != "" {
		c.Store.Server = storeServer
	}

	sc := server.Config{}
	// logger
	logger, err := newLogger(c.Logger.Level, c.Logger.Format)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	if c.Logger.Level != "" {
		log.Println(fmt.Printf("config using log level: %s", c.Logger.Level))
	}
	sc.Logger = logger

	// parse db config
	sc.Store.Timeout, err = time.ParseDuration(c.Store.Timeout)
	if err != nil {
		return fmt.Errorf("Failed to parse database timeout duration %s :%v", c.Store.Timeout, err)
	}

	sc.Store.RetryIn, err = time.ParseDuration(c.Store.Retry)
	if err != nil {
		return fmt.Errorf("Failed to parse database retryin duration %s :%v", c.Store.Retry, err)
	}

	// parse identity config
	if c.Identity.Issuer != "" {
		sc.Identity.Issuer = c.Identity.Issuer
	}

	if c.Identity.ID != "" {
		sc.Identity.ID = c.Identity.ID
	}

	if c.Identity.Secret != "" {
		sc.Identity.Secret = c.Identity.Secret
	}

	if c.Identity.RedirectURI != "" {
		sc.Identity.RedirectURI = c.Identity.RedirectURI
	}

	sc.Identity.HostAndPort = c.Identity.GRPC

	ctx := context.Background()

	// set rest of databse properties to server config
	sc.Store.Name = c.Store.Name
	sc.Store.Username = c.Store.Username
	sc.Store.Password = c.Store.Password
	sc.Store.Server = c.Store.Server
	sc.Store.Monotonic = c.Store.Monotonic

	// open the db connection & retries
	session, err := store.Connect(sc.Logger, sc.Store)
	if err != nil {
		logger.Fatalln(fmt.Errorf("Failed to initalize store (database) connection : %v", err))
	}
	sc.Session = session
	defer session.Close()

	corsOpts := handlers.AllowedOrigins([]string{"*"})
	corsHandler := handlers.CORS(corsOpts)
	recoveryHandler := handlers.RecoveryHandler()

	publicChain := alice.New(recoveryHandler, corsHandler)

	//extractHandler := handlers.ExtractHandler(ctx.Context, ctx.Router)
	//ctx.PublicChain = commonChain.Extend(alice.New(extractHandler))

	//secureHandler := handlers.SecurityHandler(ctx.Context, ctx.Logger, ctx.Signer, ctx.Verifier)
	//ctx.SecureChain = commonChain.Extend(alice.New(secureHandler, extractHandler))

	s, err := server.New(ctx, sc)
	if err != nil {
		logger.Fatalln(fmt.Errorf("Failed to initialize the server [%v]", err))
	}

	errch := make(chan error, 2)

	var grpcServer *grpc.Server

	//start grpc
	go func() {

		errch <- func() error {

			listen, err := net.Listen("tcp", c.Web.GRPC)
			if err != nil {
				return fmt.Errorf("Listening on %s : %v", c.Web.GRPC, err)
			}
			grpcOpts := []grpc.ServerOption{}
			grpcServer = grpc.NewServer(grpcOpts...)

			logger.Info("Exposing GRPC services on ", c.Web.GRPC)
			server.RegisterGRPCServices(grpcServer, s)

			err = grpcServer.Serve(listen)

			return fmt.Errorf("Listening on %s : %v", c.Web.GRPC, err)
		}()
	}()

	var serv *http.Server

	// start http
	go func() {

		errch <- func() error {

			dialopts := []grpc.DialOption{grpc.WithInsecure()}

			router := mux.NewRouter()

			err := server.RegisterHTTPServices(ctx, router, c.Web.GRPC, dialopts)
			if err != nil {
				return fmt.Errorf("Error when registering services.. : %v", err)
			}

			s.Router.Handle("/", router)

			serv = &http.Server{Addr: c.Web.HTTP, Handler: publicChain.Then(s.Router)}

			logger.Info("Exposing HTTP services on ", c.Web.HTTP)
			err = serv.ListenAndServe()

			return fmt.Errorf("Listing on %s failed with : %v", c.Web.HTTP, err)
		}()
	}()

	//graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {

		errch <- func() error {

			<-sigs

			logger.Infoln("Stopping GRPC Server..")
			grpcServer.GracefulStop()

			logger.Infoln("Stopping HTTP(S) Server..")
			serv.Shutdown(ctx)

			return fmt.Errorf("Server gracefully stopped ")
		}()
	}()

	return <-errch
}

var (
	logLevels  = []string{"debug", "info", "error"}
	logFormats = []string{"json", "text"}
)

type utcFormatter struct {
	f logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

func newLogger(level string, format string) (logrus.Logger, error) {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "", "info":
		logLevel = logrus.InfoLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		return logrus.Logger{}, fmt.Errorf("log level is not one of the supported values (%s): %s", strings.Join(logLevels, ", "), level)
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "", "text":
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		return logrus.Logger{}, fmt.Errorf("log format is not one of the supported values (%s): %s", strings.Join(logFormats, ", "), format)
	}

	return logrus.Logger{
		Out:       os.Stderr,
		Formatter: &formatter,
		Level:     logLevel,
	}, nil
}
