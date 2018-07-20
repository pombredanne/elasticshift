/*
Copyright 2018 The Elasticshift Authors.
*/
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"context"

	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
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
	shiftDir string

	buildID string
	teamID  string

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
		o.shiftDir = path
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

func New(ctx context.Context, buildID, teamID string, opt ...LoggerOption) (*Logr, error) {

	opts := defaultLoggerOptions

	for _, o := range opt {
		o(&opts)
	}

	l := &Logr{
		opts: opts,
	}

	opts.buildID = buildID
	opts.teamID = teamID

	// writers := []io.Writer{}
	// writers := []io.Writer{os.Stdout}

	var f *os.File
	var err error
	if opts.shiftDir != "" {

		// construct the logfile path
		// <storage-path>/teamid/buildid/log
		p := filepath.Join(opts.shiftDir, "logs", teamID, buildID)

		var exist bool
		exist, err = utils.PathExist(p)
		if err != nil {
			return nil, fmt.Errorf("Error Initializing logger: %v", err)
		}

		if !exist {

			err = utils.Mkdir(p)
			if err != nil {
				return nil, fmt.Errorf("Error mkdir (%s) : %v", p, err)
			}
		}

		f, err = os.OpenFile(filepath.Join(p, "log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}

		l.file = f
		// writers = append(writers, bufio.NewWriter(f))
		// l.Writer = bufio.NewWriter(f)
		l.Writer = f
	}

	if opts.minio {

		// way to connect minio
	}

	// l.Writer = io.MultiWriter(writers...)

	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC)
	// log.SetOutput(l.Writer)

	return l, nil
}

func (l *Logr) Log(message string) {
	l.Writer.Write([]byte(message + "\n"))
}

func (l *Logr) Close() {
	l.file.Close()
	// close minio and other logger providers
}
