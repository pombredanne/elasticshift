/*
Copyright 2018 The Elasticshift Authors.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/conspico/elasticshift/pkg/worker"
	"gitlab.com/conspico/elasticshift/pkg/worker/logger"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
	"golang.org/x/net/context"
)

const (
	DEFAULT_TIMEOUT = "120m"
)

func main() {

	ctx := context.Background()

	os.Setenv("SHIFT_HOST", "127.0.0.1")
	os.Setenv("SHIFT_PORT", "5051")
	os.Setenv("SHIFT_LOGGER", "embedded")
	os.Setenv("SHIFT_BUILDID", "5ad21b9ddc294a6133cdd77d")
	os.Setenv("WORKER_PORT", "6060")

	logType := os.Getenv("SHIFT_LOGGER")
	if logType == "" {
		panic("SHIFT_LOGGER must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_LOGGER=%s\n", logType)
	}

	buildID := os.Getenv("SHIFT_BUILDID")
	if buildID == "" {
		panic("SHIFT_BUILDID must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_BUILDID=%s\n", buildID)
	}

	logfile := os.Getenv("FILELOGR_PATH")
	if logfile == "" {
		panic("FILELOGR_PATH must be passed through environment variable.")
	} else {
		log.Printf("FILELOGR_PATH=%s\n", logfile)
	}

	opts := []logger.LoggerOptions{
		logger.FileLogger(logfile),
	}

	logr := logger.New(ctx, buildID, opts...)
	defer logr.Close()

	cfg := types.Config{}
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
		log.Fatal("Halting the worker.")
		os.Exit(1)
	}

	cfg.Timeout = os.Getenv("SHIFT_TIMEOUT")
	if cfg.Timeout == "" {
		log.Println("SHIFT_TIMEOUT defaulted to 120m")
		cfg.Timeout = DEFAULT_TIMEOUT
	} else {
		log.Println("SHIFT_TIMEOUT=" + cfg.Timeout)
	}

	ctx := types.Context{}
	ctx.Context = context.Background()
	ctx.Config = cfg

	// Start the worker
	err := worker.Start(ctx, logr)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to start the worker %v", err))
	}
}
