/*
Copyright 2018 The Elasticshift Authors.
*/
package main

import (
	"fmt"
	"log"
	"os"

	"gitlab.com/conspico/elasticshift/pkg/worker"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
	"golang.org/x/net/context"
)

func main() {

	cfg := types.Config{}

	host := os.Getenv("SHIFT_HOST")
	if host == "" {
		log.Println("SHIFT_HOST must be passed through environment variable when starting the container")
	}
	cfg.Host = host

	port := os.Getenv("SHIFT_PORT")
	if port == "" {
		log.Println("SHIFT_PORT must be passed through environment variable when starting the container")
	}
	cfg.Port = port

	logType := os.Getenv("SHIFT_LOGGER")
	if logType == "" {
		log.Println("SHIFT_LOGGER must be passed through environment variable when starting the container")
	}
	cfg.LogType = logType

	buildID := os.Getenv("SHIFT_BUILDID")
	if buildID == "" {
		log.Println("SHIFT_BUILDID must be passed through environment variable when starting the container")
	}
	cfg.BuildID = buildID

	cfg.GRPC = os.Getenv("WORKER_PORT")
	cfg.Timeout = os.Getenv("SHIFT_TIMEOUT")

	ctx := types.Context{}
	ctx.Context = context.Background()
	ctx.Config = cfg

	// Start the worker
	err := worker.Start(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to start the worker %v", err))
	}
}
