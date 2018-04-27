/*
Copyright 2018 The Elasticshift Authors.
*/
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"context"
)

const (
	log_embedded = "embedded"
	log_file     = "file"
)

const (
// DIR_LOGS = "/Users/ghazni/elasticshift/elasticshift/logs"
// DIR_LOGS = "/shift/logs"
)

type Logr struct {
	opts options

	logger *log.Logger

	file *os.File

	Writer io.Writer
}

type options struct {
	// file based logger, the default one.
	logfile string

	// Minio
	minio_url        string
	minio_accesscode string
	minio_key        string
	minio            bool

	// GCE logger

	// Amazon block

	// NFS

	timeOut time.Duration
}

var defaultLoggerOptions = options{
	timeOut: 120 * time.Minute,
}

// A LoggerOption let you provide options to where to log the build
// logs such as local file or GCE etc.
type LoggerOption func(*options)

func FileLogger(path string) LoggerOption {
	return func(o *options) {
		o.logfile = path
	}
}

func MinioLogger(url, accesscode, key string) LoggerOption {
	return func(o *options) {
		o.minio_url = url
		o.minio_accesscode = accesscode
		o.minio_key = key
		o.minio = true
	}
}

func New(ctx context.Context, buildId string, opt ...LoggerOption) (*Logr, error) {

	opts := defaultLoggerOptions

	for _, o := range opt {
		o(&opts)
	}

	l := &Logr{
		opts: opts,
	}

	writers := []io.Writer{os.Stdout}

	var f *os.File
	var err error
	if opts.logfile != "" {

		f, err = os.OpenFile(opts.logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}

		l.file = f
		writers = append(writers, f)
	}

	if opts.minio {

		// way to connect minio
	}

	l.Writer = io.MultiWriter(writers...)
	log.SetOutput(l.Writer)

	return l, nil
}

func (l *Logr) Close() {
	l.file.Close()
	// close minio and other logger providers
}
