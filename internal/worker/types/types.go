/*
Copyright 2018 The Elasticshift Authors.
*/
package types

import (
	"context"
	"io"

	"gitlab.com/conspico/elasticshift/api"
)

type Context struct {
	Client      api.ShiftClient
	Context     context.Context
	Config      Config
	ContainerID string
	Writer      io.Writer
	Logdir      string
}

type Config struct {
	GRPC string //worker port

	Host string // shift server host
	Port string // shift server port

	ShiftDir           string
	Timeout            string
	BuildID            string
	TeamID             string
	RepoBasedShiftFile bool
}
