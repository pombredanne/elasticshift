/*
Copyright 2018 The Elasticshift Authors.
*/
package logwriter

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
)

const (
	logdir = "/tmp/shiftlogs/"
)

type LogWriter interface {
	GetLogger(nodeid string) (*logrus.Entry, error)
}

type logw struct {
	loggr logger.Loggr
}

/*
{
	"level": "info" | "error",
	"node_id" : "1",
	"content": "this is sample log"
}
*/

func New(logLevel, logFormat string) (LogWriter, error) {

	lw := &logw{}

	loggr, err := logger.New(logLevel, logFormat)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}
	lw.loggr = loggr

	return lw, nil
}

func (lw *logw) GetLogger(nodeid string) (*logrus.Entry, error) {

	nw, err := newNodeWriter(nodeid)
	if err != nil {
		return nil, err
	}

	return lw.loggr.GetLoggerWithField("node_id", nodeid, nw), nil
}
