/*
Copyright 2018 The Elasticshift Authors.
*/
package logwriter

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/internal/pkg/logger"
)

const (
	logdir = "/tmp/shiftlogs/"
)

type LogWriter interface {
	GetLogger(nodeid string) (*logrus.Entry, error)
	LogFile(nodeid string) (*os.File, error)
}

type logw struct {
	loggr       logger.Loggr
	nodewriters map[string]NodeWriter
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
	lw.nodewriters = make(map[string]NodeWriter)

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
	lw.nodewriters[nodeid] = nw

	return lw.loggr.GetLoggerWithField("node_id", nodeid, nw), nil
}

func (lw *logw) LogFile(nodeid string) (*os.File, error) {

	var err error
	var f *os.File
	nw := lw.nodewriters[nodeid]
	if nw != nil {
		f = nw.File()
		if f == nil {
			err = fmt.Errorf("No logfile availe for %s \n ", nodeid)
		}
	}
	return f, err
}
