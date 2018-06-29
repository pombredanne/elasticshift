/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"fmt"
	"log"
	"net"
	"time"

	"gitlab.com/conspico/elasticshift/api"
	"gitlab.com/conspico/elasticshift/pkg/builder"
	"gitlab.com/conspico/elasticshift/pkg/worker/logger"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// W holds the worker related values
type W struct {
	Config types.Config

	GRPCServer *grpc.Server
	errch      chan error

	// logger      logshipper.Logger
	ShiftServer *grpc.ClientConn
	Context     types.Context

	Logr *logger.Logr
}

// Start is the launch point of the worker that accepts the environment
// variables as config, to kick start the worker
func Start(ctx types.Context, logr *logger.Logr) error {

	// initialize the logger.

	w := &W{}
	w.Config = ctx.Config
	w.Context = ctx
	w.Logr = logr
	w.errch = make(chan error)

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
		return err
	case <-time.After(timeout):
		w.Halt()
		msg := fmt.Sprintf("Worker has been timed-out after running for about %s minutes, and all the process have been halted", ctx.Config.Timeout)
		log.Println(msg)
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

	log.Println(w.Context.Client)

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

	// log.Println(w.Context.Client)

	log.Printf("Connection state: %v", w.ShiftServer.GetState())
	// if w.ShiftServer.GetState().String() == "TRANSIENT_FAILURE" {
	// 	w.ConnectToShiftServer()
	// }

	res, err := w.Context.Client.Register(w.Context.Context, req)
	if err != nil {
		fmt.Printf("registration response: %s\n", res)
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
	return builder.New(w.Context, w.ShiftServer, w.Logr)
}

// Post the log to error channel, that denotes the startup of the worker is failed
func (w *W) Fatal(err error) {
	log.Println(err)
	w.errch <- err
}
