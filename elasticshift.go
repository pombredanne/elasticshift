// Package armor
// Author Ghazni Nattarshah
// Date: 1/10/17
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"log"

	"github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/server"
	"gitlab.com/conspico/elasticshift/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Config ..
type Config struct {
	Store  Store  `json:"store"`
	Web    Web    `json:"web"`
	Logger Logger `json:"logger"`
	Dex    Dex    `json:"dex"`
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
type Dex struct {
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
				Name:      "esh",
				Username:  "esh",
				Password:  "eshpazz",
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

			Dex: Dex{
				GRPC:        "127.0.0.1:5557",
				Issuer:      "http://127.0.0.1:5556/dex",
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

	// parse Dex config
	if c.Dex.Issuer != "" {
		sc.Dex.Issuer = c.Dex.Issuer
	}

	if c.Dex.ID != "" {
		sc.Dex.ID = c.Dex.ID
	}

	if c.Dex.Secret != "" {
		sc.Dex.Secret = c.Dex.Secret
	}

	if c.Dex.RedirectURI != "" {
		sc.Dex.RedirectURI = c.Dex.RedirectURI
	}

	sc.Dex.HostAndPort = c.Dex.GRPC

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

	s, err := server.NewServer(ctx, sc)
	if err != nil {
		logger.Fatalln(fmt.Errorf("Failed to initialize the server [%v]", err))
	}

	errch := make(chan error, 2)

	//start grpc
	go func() {

		errch <- func() error {

			listen, err := net.Listen("tcp", c.Web.GRPC)
			if err != nil {
				return fmt.Errorf("Listening on %s : %v", c.Web.GRPC, err)
			}
			grpcOpts := []grpc.ServerOption{}
			grpcServer := grpc.NewServer(grpcOpts...)

			api.RegisterUserServer(grpcServer, server.NewUserServer(s))
			api.RegisterTeamServer(grpcServer, server.NewTeamServer(s))
			api.RegisterClientServer(grpcServer, server.NewClientServer(s))

			err = grpcServer.Serve(listen)
			return fmt.Errorf("Listening on %s : %v", c.Web.GRPC, err)
		}()
	}()

	// start http
	go func() {

		errch <- func() error {

			dialopts := []grpc.DialOption{grpc.WithInsecure()}

			// user
			userMux := runtime.NewServeMux()
			api.RegisterUserHandlerFromEndpoint(ctx, userMux, c.Web.GRPC, dialopts)
			s.Mux.Handle("/user", userMux)

			// team
			teamMux := runtime.NewServeMux()
			api.RegisterTeamHandlerFromEndpoint(ctx, teamMux, c.Web.GRPC, dialopts)
			s.Mux.Handle("/team", teamMux)

			// client
			clientMux := runtime.NewServeMux()
			api.RegisterClientHandlerFromEndpoint(ctx, clientMux, c.Web.GRPC, dialopts)
			s.Mux.Handle("/client", clientMux)

			err := http.ListenAndServe(c.Web.HTTP, s.Mux)
			return fmt.Errorf("Listing on %s failed with : %v", c.Web.HTTP, err)
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

func newLogger(level string, format string) (logrus.FieldLogger, error) {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "", "info":
		logLevel = logrus.InfoLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		return nil, fmt.Errorf("log level is not one of the supported values (%s): %s", strings.Join(logLevels, ", "), level)
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "", "text":
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		return nil, fmt.Errorf("log format is not one of the supported values (%s): %s", strings.Join(logFormats, ", "), format)
	}

	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: &formatter,
		Level:     logLevel,
	}, nil
}
