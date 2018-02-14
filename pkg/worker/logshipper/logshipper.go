/*
Copyright 2018 The Elasticshift Authors.
*/
package logshipper

import (
	"strings"
	"time"

	stypes "gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/worker/types"
)

const (
	log_embedded = "embedded"
	log_file     = "file"
)

type Logger interface {

	// Send the log to embedded logger (by default)
	// Configurable to send logs to different source
	Send(log stypes.Log)

	Log(log string)
	Error(err error)
	Info(msg string)

	Halt() error
}

func New(ctx types.Context) (Logger, error) {

	var logger Logger
	var err error
	if strings.EqualFold(log_embedded, ctx.Config.LogType) {
		logger, err = newEmbeddedLogger(ctx)
	}

	// } else if strings.EqualFold(log_file, logType) {
	// 	err = newFileLogger(ctx)
	// }

	if err != nil {
		return nil, err
	}

	return logger, nil
}

func constructLog(log string) stypes.Log {
	return stypes.Log{Data: log, Time: time.Now()}
}
