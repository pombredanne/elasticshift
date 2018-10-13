/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/internal/worker/builder"
	"gitlab.com/conspico/elasticshift/internal/worker/logger"
	"gitlab.com/conspico/elasticshift/internal/worker/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultTimeout = "120m"
)

// Run ..
func Run() error {

	log.Print("S:~0.1:Environment Setup:~\n")

	bctx := context.Background()
	cfg := types.Config{}

	// os.Setenv("SHIFT_HOST", "127.0.0.1")
	// os.Setenv("SHIFT_PORT", "9101")
	// os.Setenv("SHIFT_BUILDID", "5ba20b3016986f00017f5b00")
	// os.Setenv("SHIFT_DIR", "/Users/ghazni/.elasticshift/storage")
	// os.Setenv("WORKER_PORT", "9200")
	// os.Setenv("SHIFT_TEAMID", "5a3a41f08011e098fb86b41f")
	// os.Setenv("SHIFT_REPOFILE", "true")

	// logLevel := os.Getenv("SHIFT_LOG_LEVEL")
	// if logLevel == "" {
	// 	log.Fatalln("SHIFT_LOG_LEVEL must be passed through environment variable.")
	// }

	// logFormat := os.Getenv("SHIFT_LOG_FORMAT")
	// if logFormat == "" {
	// 	log.Fatalln("SHIFT_LOG_FORMAT must be passed through environment variable.")
	// }

	// loggr, err := shiftlogger.New(logLevel, logFormat)
	// if err != nil {
	// 	return fmt.Errorf("invalid config: %v", err)
	// }
	// slog := loggr.GetLogger("Worker")
	// log.Printf("SHIFT_LOG_LEVEL=%s\n", logLevel)
	// log.Printf("SHIFT_LOG_FORMAT=%s\n", logFormat)

	shiftDir := os.Getenv("SHIFT_DIR")
	if shiftDir == "" {
		log.Println("SHIFT_DIR must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_DIR=%s\n", shiftDir)
	}
	cfg.ShiftDir = shiftDir

	// cfg.ShiftDir = shiftDir
	// opts := []logger.LoggerOption{
	// 	logger.FileLogger(shiftDir),
	// }

	dir, writers, err := logger.Initialize()
	if err != nil {
		log.Fatalf("Cannot initialize logger: %v", err)
	}

	buildID := os.Getenv("SHIFT_BUILDID")
	if buildID == "" {
		log.Println("SHIFT_BUILDID must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_BUILDID=%s\n", buildID)
	}
	cfg.BuildID = buildID

	teamID := os.Getenv("SHIFT_TEAMID")
	if teamID == "" {
		log.Println("SHIFT_TEAMID  must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_TEAMID=%s\n", teamID)
	}
	cfg.TeamID = teamID

	// logr, err := logger.New(bctx, buildID, teamID)
	// defer logr.Close()
	// if err != nil {
	// 	log.Printf("Initializing logger failed : %v", err)
	// 	return fmt.Errorf("Error initializing logger: %v", err)
	// }

	var isError bool
	host := os.Getenv("SHIFT_HOST")
	if host == "" {
		log.Println("SHIFT_HOST must be passed through environment variable.")
		isError = true
	} else {
		log.Printf("SHIFT_HOST=%s\n", host)
	}
	cfg.Host = host

	port := os.Getenv("SHIFT_PORT")
	if port == "" {
		log.Println("SHIFT_PORT must be passed through environment variable")
		isError = true
	} else {
		log.Printf("SHIFT_PORT=%s\n", port)
	}
	cfg.Port = port

	workerPort := os.Getenv("WORKER_PORT")
	if workerPort == "" {
		log.Println("WORKER_PORT must be passed though environment variable.")
		isError = true
	} else {
		log.Printf("WORKER_PORT=%s\n", workerPort)
	}
	cfg.GRPC = workerPort

	if isError {
		log.Println("One or more arguments required to start the worker has not passed through environment variables.")
		return errors.New("Failed to start the worker")
	}

	cfg.Timeout = os.Getenv("SHIFT_TIMEOUT")
	if cfg.Timeout == "" {
		log.Println("SHIFT_TIMEOUT defaulted to 120m")
		cfg.Timeout = defaultTimeout
	} else {
		log.Println("SHIFT_TIMEOUT=" + cfg.Timeout)
	}

	repoBasedShiftFile := os.Getenv("SHIFT_REPOFILE")
	cfg.RepoBasedShiftFile, _ = strconv.ParseBool(repoBasedShiftFile)

	ctx := types.Context{}
	ctx.Context = bctx
	ctx.Config = cfg
	ctx.Writer = writers
	ctx.Logdir = dir

	// Start the worker
	err = Start(ctx)
	if err != nil {
		return fmt.Errorf("Failed to start the worker: %+v", err)
	}

	return nil
}

// W holds the worker related values
type W struct {
	Config types.Config

	GRPCServer *grpc.Server
	errch      chan error
	done       chan int

	// logger      logshipper.Logger
	ShiftServer *grpc.ClientConn
	Context     types.Context

	privKeyPath string
	pubKeyPath  string
}

// Start is the launch point of the worker that accepts the environment
// variables as config, to kick start the worker
func Start(ctx types.Context) error {

	// initialize the logger.

	w := &W{}
	w.Config = ctx.Config
	w.Context = ctx
	w.errch = make(chan error)
	w.done = make(chan int)

	var timeout time.Duration
	timeout, _ = time.ParseDuration(ctx.Config.Timeout)
	log.Println("Idle Timeout :" + timeout.String())

	var err error
	go func() {
		// // Start the log shipper
		// err = w.StartLogShipper()
		// if err != nil {
		// 	w.errch <- err
		// 	return
		// }

		// Connects to shift server
		err = w.ConnectToShiftServer()
		if err != nil {
			w.errch <- err
			return
		}

		// Generate RSA keys, used to SSH
		err = w.GenerateRSAKeys()
		if err != nil {
			w.errch <- err
			return
		}

		// Register the worker to shift server
		err = w.RegisterWorker()
		if err != nil {
			w.errch <- err
			return
		}

		// Listener on worker to receive command from shift server.
		w.StartGRPCServer()

		// Kick start the builder
		err = w.StartBuilder()
		if err != nil {
			w.errch <- err
			return
		}
	}()

	// close the log file & shift server connection
	// defer w.Halt()

	// Stops when receive the fatal error or when worker timeout
	select {
	case err := <-w.errch:
		w.UpdateShiftServer(statusFailed, "")
		return err
	case <-time.After(timeout):
		w.Halt()
		msg := fmt.Sprintf("Worker has been timed-out after running for about %s minutes, and all the process have been halted", ctx.Config.Timeout)
		w.UpdateShiftServer(statusFailed, "")
		log.Println(msg)
	case <-w.done:
	}

	return nil
}

// ConnectToShiftServer establish the connection to elasticshift GRPC server
// Worker -> shift server communication channel (thru GRPC)
func (w *W) ConnectToShiftServer() error {

	log.Println("Connecting to shift server..")

	// TODO connect ssl
	cp := keepalive.ClientParameters{}
	cp.Time, _ = time.ParseDuration("1m")
	cp.Timeout, _ = time.ParseDuration("120m")

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 120 * time.Second,
		}),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	grpc.EnableTracing = true
	conn, err := grpc.Dial(w.Config.Host+":"+w.Config.Port, opts...)
	if err != nil {
		return fmt.Errorf("Failed to connect to shift server %s:%s, %v", w.Config.Host, w.Config.Port, err)
	}
	w.ShiftServer = conn
	log.Printf("Connection state: %v", w.ShiftServer.GetState())

	w.Context.Client = api.NewShiftClient(conn)

	log.Printf("Connection state: %v", w.ShiftServer.GetState())

	return nil
}

// Start the log shipper which is used to send
// the logs generated out of build to different logging system
// By default embedded logger is set to work
// func (w *W) StartLogShipper() error {

// 	log.Println("Starting the logshipper..")

// 	// Start logshipper
// 	f, err := logshipper.New(w.Context, w.Config.BuildID)
// 	if err != nil {
// 		return fmt.Errorf("Failed to start log shipper for '%s' type : %v ", w.Config.ShiftDir, err)
// 	}
// 	w.logfile = f
// 	log.Println("Started.")
// 	return nil
// }

// Register the worker to elasticshift server
// This is to let elasticshift know where the worker
// is running for further communication
func (w *W) RegisterWorker() error {

	log.Println("Registering the worker..")

	// Loads the generate private key
	key, err := w.ReadPrivateKey(PRIV_KEY_PATH)
	if err != nil {
		return fmt.Errorf("Failed to read the key : %v", err)
	}

	// perform registration
	req := &api.RegisterReq{}
	req.BuildId = w.Config.BuildID
	req.Privatekey = key

	log.Printf("Connection state: %v", w.ShiftServer.GetState())

	res, err := w.Context.Client.Register(w.Context.Context, req)
	if err != nil {
		log.Printf("registration response: %s\n", res)
		return fmt.Errorf("Worker registration failed: %v", err)
	}

	if res != nil && res.Registered {
		log.Println("Registration Successful.")
	} else {
		return fmt.Errorf("Registration failed.")
	}
	return nil
}

// Start the GRPC server to listen for commands from elasticshift server
func (w *W) StartGRPCServer() {

	if w.Config.GRPC == "" {
		w.Config.GRPC = DEFAULT_GRPC_PORT
	}

	var grpcServer *grpc.Server

	go func() {

		//start grpc
		w.errch <- func() error {

			log.Println("Starting listener to obey shift server commands on " + w.Config.GRPC)
			listen, err := net.Listen("tcp", ":"+w.Config.GRPC)
			if err != nil {
				return fmt.Errorf("Failed to start GRPC server on %s : %v", w.Config.GRPC, err)
			}
			grpcOpts := []grpc.ServerOption{}
			grpcServer = grpc.NewServer(grpcOpts...)
			w.GRPCServer = grpcServer

			log.Println("Exposing GRPC services on " + w.Config.GRPC)

			// register the grpc services
			api.RegisterWorkServer(grpcServer, NewServer(w.Context))

			err = grpcServer.Serve(listen)
			return fmt.Errorf("Listening on %s failed : %v", w.Config.GRPC, err)
		}()
	}()
}

// Halt stops all the process running by the worker
// This would run if the build is timed out or executed successfully
func (w *W) Halt() {

	// Stop the grpc server
	w.GRPCServer.GracefulStop()

	// Close the shift server connection
	w.ShiftServer.Close()
}

// start the builder where the real execution happens.
func (w *W) StartBuilder() error {
	return builder.New(w.Context, w.ShiftServer, w.Context.Writer, w.done)
}

// Post the log to error channel, that denotes the startup of the worker is failed
func (w *W) Fatal(err error) {
	log.Println(err)
	w.errch <- err
}

var (
	statusFailed = "F"
)

func (w *W) UpdateShiftServer(status, checkpoint string) {

	req := &api.UpdateBuildStatusReq{}
	req.BuildId = w.Context.Config.BuildID
	req.Status = statusFailed

	if w.Context.Client != nil {
		_, err := w.Context.Client.UpdateBuildStatus(w.Context.Context, req)
		if err != nil {
			log.Printf("Failed to update buld graph: %v", err)
		}
	}
}
