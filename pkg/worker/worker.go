/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"fmt"
	"net"
	"time"

	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/pkg/builder"
	"gitlab.com/conspico/elasticshift/pkg/worker/logshipper"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
	"google.golang.org/grpc"
)

// W holds the worker related values
type W struct {
	Config types.Config

	GRPCServer *grpc.Server
	errch      chan error

	logger      logshipper.Logger
	ShiftServer *grpc.ClientConn
	Context     types.Context
}

// Start is the launch point of the worker that accepts the environment
// variables as config, to kick start the worker
func Start(ctx types.Context) error {

	w := &W{}
	w.Config = ctx.Config
	w.Context = ctx
	w.errch = make(chan error)

	// Connects to shift server
	w.ConnectToShiftServer()
	defer w.ShiftServer.Close()

	// Start the log shipper
	w.StartLogShipper()

	// Generate RSA keys, used to SSH
	w.GenerateRSAKeys()

	// Register the worker to shift server
	w.RegisterWorker()

	// Listener on worker to receive command from shift server.
	w.StartGRPCServer()

	// Kick start the builder
	w.StartBuilder()

	var timeout time.Duration
	if ctx.Config.Timeout == "" {
		timeout = DEFAULT_TIMEOUT
	} else {
		timeout, _ = time.ParseDuration(ctx.Config.Timeout)
	}

	// Stops when receive the fatal error or when worker timeout
	select {
	case err := <-w.errch:
		return err
	case <-time.After(timeout):
		w.Halt()
		msg := fmt.Sprintf("Worker has been timed-out after running for about %s minutes, and all the process have been halted", ctx.Config.Timeout)
		return fmt.Errorf(msg)
	}
}

// ConnectToShiftServer establish the connection to elasticshift GRPC server
// Worker -> shift server communication channel (thru GRPC)
func (w *W) ConnectToShiftServer() {

	// TODO connect ssl
	conn, err := grpc.Dial(w.Config.Host+":"+w.Config.Port, grpc.WithInsecure())
	if err != nil {
		w.Fatal(fmt.Errorf("Failed to connect to shift server %s:%s, %v", w.Config.Host, w.Config.Port, err))
	}

	w.ShiftServer = conn
	w.Context.Client = api.NewShiftClient(conn)
}

// Start the log shipper which is used to send
// the logs generated out of build to different logging system
// By default embedded logger is set to work
func (w *W) StartLogShipper() {

	// Start logshipper
	l, err := logshipper.New(w.Context)
	if err != nil {
		w.Fatal(fmt.Errorf("Failed to start log shipper for '%s' type : %v ", w.Config.LogType, err))
	}
	w.logger = l

	l.Info("Connected to elasticshift server.")
	l.Info("Started the logshipper")
}

// Register the worker to elasticshift server
// This is to let elasticshift know where the worker
// is running for further communication
func (w *W) RegisterWorker() {
	// Loads the generate private key
	key, err := w.ReadPrivateKey(PRIV_KEY_PATH)
	if err != nil {
		w.logger.Error(err)
	}

	// perform registration
	req := &api.RegisterReq{}
	req.BuildId = w.Config.BuildID
	req.Privatekey = key

	res, err := w.Context.Client.Register(w.Context.Context, req)
	if err != nil {
		w.logger.Error(fmt.Errorf("Worker registration failed: %v", err))
	}

	if res.Registered {
		w.logger.Info("Registration Successful.")
	}
}

// Start the GRPC server to listen for commands from elasticshift server
func (w *W) StartGRPCServer() {

	if w.Config.GRPC == "" {
		w.Config.GRPC = DEFAULT_GRPC_PORT
	}

	var grpcServer *grpc.Server

	//start grpc
	go func(config types.Config, logger logshipper.Logger) {

		listen, err := net.Listen("tcp", config.GRPC)
		if err != nil {
			w.logger.Error(fmt.Errorf("Failed to start GRPC server on %s : %v", config.GRPC, err))
		}
		grpcOpts := []grpc.ServerOption{}
		grpcServer = grpc.NewServer(grpcOpts...)

		w.logger.Info("Exposing GRPC services on " + config.GRPC)

		// register the grpc services
		api.RegisterWorkServer(grpcServer, NewServer(w.Context, w.logger))

		err = grpcServer.Serve(listen)
		if err != nil {
			w.logger.Error(fmt.Errorf("Failed to start GRPC server on %s : %v", config.GRPC, err))
		}

	}(w.Config, w.logger)

	w.GRPCServer = grpcServer
}

// Halt stops all the process running by the worker
// This would run if the build is timed out or executed successfully
func (w *W) Halt() {

	// Stop the grpc server
	w.GRPCServer.GracefulStop()

	// Stop the log shipper
	w.logger.Halt()

	// Close the shift server connection
	w.ShiftServer.Close()
}

// start the builder where the real execution happens.
func (w *W) StartBuilder() {

	err := builder.Start(w.logger)
	w.errch <- err
}

// Post the log to error channel, that denotes the startup of the worker is failed
func (w *W) Fatal(err error) {
	w.errch <- err
}
