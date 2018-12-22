/*
Copyright 2018 The Elasticshift Authors.
*/
package types

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api"
	"github.com/elasticshift/elasticshift/internal/pkg/utils"
	"github.com/elasticshift/elasticshift/internal/worker/logwriter"
)

type Context struct {
	Client      api.ShiftClient
	Context     context.Context
	Config      Config
	ContainerID string
	Writer      io.Writer
	Logdir      string

	LogWriter logwriter.LogWriter
	EnvLogger *logrus.Entry
	EnvTimer  utils.Timer
}

type Config struct {
	GRPC string //worker port

	Host string // shift server host
	Port string // shift server port

	ShiftDir           string
	Timeout            string
	BuildID            string
	SubBuildID         string
	TeamID             string
	RepoBasedShiftFile bool
}
