/*
Copyright 2018 The Elasticshift Authors.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"context"

	"gitlab.com/conspico/elasticshift/pkg/worker"
	"gitlab.com/conspico/elasticshift/pkg/worker/logger"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
)

const (
	DEFAULT_TIMEOUT = "120m"
)

func main() {

	bctx := context.Background()
	cfg := types.Config{}

	os.Setenv("SHIFT_HOST", "127.0.0.1")
	os.Setenv("SHIFT_PORT", "5051")
	os.Setenv("SHIFT_BUILDID", "5b0393b3dc294a2d45fa2232")
	os.Setenv("SHIFT_DIR", "/Users/ghazni/.elasticshift/storage")
	os.Setenv("WORKER_PORT", "6060")
	os.Setenv("SHIFT_TEAMID", "5a3a41f08011e098fb86b41f")

	buildID := os.Getenv("SHIFT_BUILDID")
	if buildID == "" {
		panic("SHIFT_BUILDID must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_BUILDID=%s\n", buildID)
	}
	cfg.BuildID = buildID

	teamID := os.Getenv("SHIFT_TEAMID")
	if teamID == "" {
		panic("SHIFT_TEAMID  must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_TEAMID=%s\n", teamID)
	}
	cfg.TeamID = teamID

	shiftDir := os.Getenv("SHIFT_DIR")
	if shiftDir == "" {
		panic("SHIFT_DIR must be passed through environment variable.")
	} else {
		log.Printf("SHIFT_DIR=%s\n", shiftDir)
	}
	cfg.ShiftDir = shiftDir

	opts := []logger.LoggerOption{
		logger.FileLogger(shiftDir),
	}

	logr, err := logger.New(bctx, buildID, teamID, opts...)
	defer logr.Close()
	if err != nil {
		panic(err)
	}

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
	err = worker.Start(ctx, logr)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to start the worker %v", err))
	}
}
