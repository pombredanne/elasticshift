/*
Copyright 2018 The Elasticshift Authors.
*/
package types

import (
	"context"

	"gitlab.com/conspico/elasticshift/api"
)

type Context struct {
	Client      api.ShiftClient
	Context     context.Context
	Config      Config
	ContainerID string
}

type Config struct {
	GRPC string

	Host    string
	Port    string
	LogType string
	Timeout string
	BuildID string
}
