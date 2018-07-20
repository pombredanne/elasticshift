/*
Copyright 2018 The Elasticshift Authors.
*/
package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	logLevels  = []string{"debug", "info", "error"}
	logFormats = []string{"json", "text"}
)

type utcFormatter struct {
	f logrus.Formatter
}

func (f *utcFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return f.f.Format(e)
}

func New(level string, format string) (logrus.Logger, error) {
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "", "info":
		logLevel = logrus.InfoLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		return logrus.Logger{}, fmt.Errorf("log level is not one of the supported values (%s): %s", strings.Join(logLevels, ", "), level)
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "", "text":
		formatter.f = &logrus.TextFormatter{DisableColors: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		return logrus.Logger{}, fmt.Errorf("log format is not one of the supported values (%s): %s", strings.Join(logFormats, ", "), format)
	}

	return logrus.Logger{
		Out:       os.Stderr,
		Formatter: &formatter,
		Level:     logLevel,
	}, nil
}
