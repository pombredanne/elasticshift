/*
Copyright 2018 The Elasticshift Authors.
*/
package logger

import (
	"fmt"
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

type loggr struct {
	logLevel  logrus.Level
	formatter utcFormatter
}

// Loggr ..
type Loggr interface {
	GetLogger(prefix string) *logrus.Entry
}

// GetLogger ..
func (l *loggr) GetLogger(component string) *logrus.Entry {
	logger := logrus.New()
	logger.Formatter = &l.formatter
	logger.Level = l.logLevel
	return logger.WithField("component", component)
}

// New ..
// Creates a new logger
func New(level string, format string) (Loggr, error) {

	l := &loggr{}
	var logLevel logrus.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = logrus.DebugLevel
	case "", "info":
		logLevel = logrus.InfoLevel
	case "error":
		logLevel = logrus.ErrorLevel
	default:
		return l, fmt.Errorf("log level is not one of the supported values (%s): %s", strings.Join(logLevels, ", "), level)
	}

	var formatter utcFormatter
	switch strings.ToLower(format) {
	case "", "text":
		formatter.f = &logrus.TextFormatter{FullTimestamp: true}
	case "json":
		formatter.f = &logrus.JSONFormatter{}
	default:
		return l, fmt.Errorf("log format is not one of the supported values (%s): %s", strings.Join(logFormats, ", "), format)
	}

	/*return logrus.Logger{
		Out:       os.Stderr,
		Formatter: &formatter,
		Level:     logLevel,
	}, nil*/

	l.logLevel = logLevel
	l.formatter = formatter
	return l, nil
}
