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

	"github.com/elasticshift/elasticshift/api"
	"github.com/elasticshift/elasticshift/internal/pkg/utils"
	"github.com/elasticshift/elasticshift/internal/worker/builder"
	"github.com/elasticshift/elasticshift/internal/worker/logwriter"
	"github.com/elasticshift/elasticshift/internal/worker/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultTimeout = "120m"
)

// Run ..
func Run() error {

	timer := utils.NewTimer()
	timer.Start()
	bctx := context.Background()
	cfg := types.Config{}

	logLevel := os.Getenv("SHIFT_LOG_LEVEL")
	if logLevel == "" {
		log.Fatalln("SHIFT_LOG_LEVEL must be passed through environment variable.")
	}

	logFormat := os.Getenv("SHIFT_LOG_FORMAT")
	if logFormat == "" {
		log.Fatalln("SHIFT_LOG_FORMAT must be passed through environment variable.")
	}

	lw, err := logwriter.New(logLevel, logFormat)
	if err != nil {
		return fmt.Errorf("invalid config: %v\n", err)
	}

	log1, err := lw.GetLogger("1")
	if err != nil {
		return fmt.Errorf("Failed to initialize logger: %v\n", err)
	}

	// slog := loggr.GetLogger("Worker")
	// log.Printf("SHIFT_LOG_LEVEL=%s\n", logLevel)
	// log.Printf("SHIFT_LOG_FORMAT=%s\n", logFormat)

	shiftDir := os.Getenv("SHIFT_DIR")
	if shiftDir == "" {
		log1.Print("SHIFT_DIR must be passed through environment variable. \n")
	} else {
		log1.Printf("SHIFT_DIR=%s\n", shiftDir)
	}
	cfg.ShiftDir = shiftDir

	// cfg.ShiftDir = shiftDir
	// opts := []logger.LoggerOption{
	// 	logger.FileLogger(shiftDir),
	// }

	// dir, writers, err := logger.Initialize()
	// if err != nil {
	// 	log1.Fatalf("Cannot initialize logger: %v", err)
	// }

	buildID := os.Getenv("SHIFT_BUILDID")
	if buildID == "" {
		log1.Print("SHIFT_BUILDID must be passed through environment variable.\n")
	} else {
		log1.Printf("SHIFT_BUILDID=%s\n", buildID)
	}
	cfg.BuildID = buildID

	subBuildID := os.Getenv("SHIFT_SUBBUILDID")
	if subBuildID == "" {
		log1.Print("SHIFT_SUBBUILDID must be passed through environment variable.\n")
	} else {
		log1.Printf("SHIFT_SUBBUILDID=%s\n", subBuildID)
	}
	cfg.SubBuildID = subBuildID

	teamID := os.Getenv("SHIFT_TEAMID")
	if teamID == "" {
		log1.Print("SHIFT_TEAMID  must be passed through environment variable.\n")
	} else {
		log1.Printf("SHIFT_TEAMID=%s\n", teamID)
	}
	cfg.TeamID = teamID

	var isError bool
	host := os.Getenv("SHIFT_HOST")
	if host == "" {
		log1.Print("SHIFT_HOST must be passed through environment variable.\n")
		isError = true
	} else {
		log1.Printf("SHIFT_HOST=%s\n", host)
	}
	cfg.Host = host

	port := os.Getenv("SHIFT_PORT")
	if port == "" {
		log1.Print("SHIFT_PORT must be passed through environment variable. \n")
		isError = true
	} else {
		log1.Printf("SHIFT_PORT=%s\n", port)
	}
	cfg.Port = port

	workerPort := os.Getenv("WORKER_PORT")
	if workerPort == "" {
		log1.Print("WORKER_PORT must be passed though environment variable. \n")
		isError = true
	} else {
		log1.Printf("WORKER_PORT=%s\n", workerPort)
	}
	cfg.GRPC = workerPort

	if isError {
		log1.Print("One or more arguments required to start the worker has not passed through environment variables. \n")
		return errors.New("Failed to start the worker")
	}

	cfg.Timeout = os.Getenv("SHIFT_TIMEOUT")
	if cfg.Timeout == "" {
		log1.Print("SHIFT_TIMEOUT defaulted to 120m \n")
		cfg.Timeout = defaultTimeout
	} else {
		log1.Printf("SHIFT_TIMEOUT= %s \n", cfg.Timeout)
	}

	repoBasedShiftFile := os.Getenv("SHIFT_REPOFILE")
	cfg.RepoBasedShiftFile, _ = strconv.ParseBool(repoBasedShiftFile)

	ctx := types.Context{}
	ctx.Context = bctx
	ctx.Config = cfg
	//ctx.Writer = writers
	//ctx.Logdir = dir
	ctx.LogWriter = lw
	ctx.EnvLogger = log1
	ctx.EnvTimer = timer

	// Start the worker
	err = Start(ctx)
	if err != nil {
		return fmt.Errorf("Failed to start the worker: %+v \n ", err)
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

	log1 := ctx.EnvLogger

	var timeout time.Duration
	timeout, _ = time.ParseDuration(ctx.Config.Timeout)
	log1.Printf("Idle Timeout : %s \n", timeout.String())

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
		log.Print(msg + "\n")
	case <-w.done:
	}

	return nil
}

// ConnectToShiftServer establish the connection to elasticshift GRPC server
// Worker -> shift server communication channel (thru GRPC)
func (w *W) ConnectToShiftServer() error {

	log1 := w.Context.EnvLogger

	log1.Print("Connecting to shift server.. \n")

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
	log1.Printf("Connection state: %v\n", w.ShiftServer.GetState())

	w.Context.Client = api.NewShiftClient(conn)

	log1.Printf("Connection state: %v\n", w.ShiftServer.GetState())

	return nil
}

// Start the log shipper which is used to send
// the logs generated out of build to different logging system
// By default embedded logger is set to work
// func (w *W) StartLogShipper() error {

// 	log.Print("Starting the logshipper.. \n")

// 	// Start logshipper
// 	f, err := logshipper.New(w.Context, w.Config.BuildID)
// 	if err != nil {
// 		return fmt.Errorf("Failed to start log shipper for '%s' type : %v ", w.Config.ShiftDir, err)
// 	}
// 	w.logfile = f
// 	log.Print("Started. \n")
// 	return nil
// }

// Register the worker to elasticshift server
// This is to let elasticshift know where the worker
// is running for further communication
func (w *W) RegisterWorker() error {

	log1 := w.Context.EnvLogger

	log1.Print("Registering the worker..\n")

	// Loads the generate private key
	key, err := w.ReadPrivateKey(PRIV_KEY_PATH)
	if err != nil {
		return fmt.Errorf("Failed to read the key : %v", err)
	}

	// perform registration
	req := &api.RegisterReq{}
	req.BuildId = w.Config.BuildID
	req.Privatekey = key

	log1.Printf("Connection state: %v\n", w.ShiftServer.GetState())

	res, err := w.Context.Client.Register(w.Context.Context, req)
	if err != nil {
		log1.Printf("registration response: %s \n ", res)
		return fmt.Errorf("Worker registration failed: %v", err)
	}

	if res != nil && res.Registered {
		log1.Print("Registration Successful.\n ")
	} else {
		return fmt.Errorf("Registration failed. \n ")
	}
	return nil
}

// StartGRPCServer ..
// Start the GRPC server to listen for commands from elasticshift server
func (w *W) StartGRPCServer() {

	log1 := w.Context.EnvLogger

	if w.Config.GRPC == "" {
		w.Config.GRPC = DEFAULT_GRPC_PORT
	}

	var grpcServer *grpc.Server

	go func() {

		//start grpc
		w.errch <- func() error {

			log1.Printf("Starting listener to obey shift server commands on %s \n", w.Config.GRPC)
			listen, err := net.Listen("tcp", ":"+w.Config.GRPC)
			if err != nil {
				return fmt.Errorf("Failed to start GRPC server on %s : %v", w.Config.GRPC, err)
			}
			grpcOpts := []grpc.ServerOption{}
			grpcServer = grpc.NewServer(grpcOpts...)
			w.GRPCServer = grpcServer

			log1.Printf("Exposing GRPC services on %s \n", w.Config.GRPC)

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

// StartBuilder ..
// start the builder where the real execution happens.
func (w *W) StartBuilder() error {
	return builder.New(w.Context, w.ShiftServer, w.Context.Writer, w.done)
}

// Fatal ..
// Post the log to error channel, that denotes the startup of the worker is failed
func (w *W) Fatal(err error) {
	w.Context.EnvLogger.Printf("%v \n", err)
	w.errch <- err
}

var (
	statusFailed = "F"
)

// UpdateShiftServer ..
func (w *W) UpdateShiftServer(status, checkpoint string) {

	req := &api.UpdateBuildStatusReq{}
	req.TeamId = w.Config.BuildID
	req.BuildId = w.Context.Config.BuildID
	req.SubBuildId = w.Config.SubBuildID
	req.Status = statusFailed

	if w.Context.Client != nil {
		_, err := w.Context.Client.UpdateBuildStatus(w.Context.Context, req)
		if err != nil {
			log.Printf("Failed to update buld graph: %v\n", err)
		}
	}
}
