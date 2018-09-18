/*
Copyright 2017 The Elasticshift Authors.
*/
package shiftserver

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
	"gitlab.com/conspico/elasticshift/internal/shiftserver/store"
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
	NSQ      NSQ
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

// NSQ ..
type NSQ struct {
	ConsumerAddress string `json:"consumer_address"`
	ProducerAddress string `json:"producer_address"`
}

// Web ..
// Holds the web server configuration
type Web struct {
	HTTP string `json:"http"`
	GRPC string `json:"grpc"`
}

// Dex ..
// Holds the web server configuration
// type Identity struct {
// 	GRPC        string `json:"grpc"`
// 	Issuer      string `json:"issuer"`
// 	ID          string `json:"id"`
// 	Secret      string `json:"secret"`
// 	RedirectURI string `json:"redirect_uri"`
// }

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

// Run ...
func Run() error {

	cfgFile := "/elasticshift/config.yaml"
	cfgData, err := ioutil.ReadFile(cfgFile)

	var c Config
	switch {
	case os.IsNotExist(err):
		c = Config{

			Store: Store{

				//Server: "127.0.0.1",
				Server:    "10.10.7.152",
				Name:      "elasticshift",
				Username:  "elasticshift",
				Password:  "3l@$t1c$h1ft",
				Monotonic: true,
				Timeout:   "10s",
				Retry:     "10s",
			},

			Logger: Logger{

				Level:  "debug",
				Format: "text",
			},

			NSQ: NSQ{
				ConsumerAddress: "127.0.0.1:4161",
				ProducerAddress: "127.0.0.1:4150",
			},

			Web: Web{
				HTTP: "0.0.0.0:9100",
				GRPC: "0.0.0.0:9101",
			},

			Identity: Identity{
				HostAndPort: "127.0.0.1:5557",
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

	sc := ServerConfig{}
	// logger
	loggr, err := logger.New(c.Logger.Level, c.Logger.Format)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	if c.Logger.Level != "" {
		log.Println(fmt.Printf("config using log level: %s", c.Logger.Level))
	}
	sc.Logger = loggr

	logger := loggr.GetLogger("main")

	// NSQ config
	// override nsq consumer address
	if consumerAddress := os.Getenv("NSQ_CONSUMER_ADDRESS"); consumerAddress != "" {
		sc.NSQ.ConsumerAddress = consumerAddress
	} else {
		sc.NSQ.ConsumerAddress = c.NSQ.ConsumerAddress
	}

	if sc.NSQ.ConsumerAddress == "" {
		return fmt.Errorf("No NSQ consumer address provided.")
	}

	// override nsq producer address
	if producerAddress := os.Getenv("NSQ_PRODUCER_ADDRESS"); producerAddress != "" {
		sc.NSQ.ProducerAddress = producerAddress
	} else {
		sc.NSQ.ProducerAddress = c.NSQ.ProducerAddress
	}

	if sc.NSQ.ProducerAddress == "" {
		return fmt.Errorf("No NSQ producer address provided.")
	}

	// parse db config
	sc.Store.Timeout, err = time.ParseDuration(c.Store.Timeout)
	if err != nil {
		return fmt.Errorf("Failed to parse database timeout duration %s :%v", c.Store.Timeout, err)
	}

	sc.Store.RetryIn, err = time.ParseDuration(c.Store.Retry)
	if err != nil {
		return fmt.Errorf("Failed to parse database retryin duration %s :%v", c.Store.Retry, err)
	}
	sc.Store.AutoReconnect = true

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

	s, err := New(ctx, sc)
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
			s.registerGRPCServices(grpcServer)

			err = grpcServer.Serve(listen)

			return fmt.Errorf("Listening on %s failed : %v", c.Web.GRPC, err)
		}()
	}()

	var serv *http.Server

	// start http
	go func() {

		errch <- func() error {

			dialopts := []grpc.DialOption{grpc.WithInsecure()}

			router := mux.NewRouter()

			err := RegisterHTTPServices(ctx, router, c.Web.GRPC, dialopts)
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
