/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"path/filepath"
	"time"
)

var (

	// Default builder timeout
	DEFAULT_TIMEOUT, _ = time.ParseDuration("120m")
)

const (

	// Default GRPC port where the worker is listening for commands from shift server
	DEFAULT_GRPC_PORT = "5053"

	// Bit size used when generating RSA keys.
	DEFAULT_BIT_SIZE = 2048

	// SSH directory
	DIR_SSH = "~/.ssh"
)

var (
	// Default private key filepath
	PRIV_KEY_PATH = filepath.Join(DIR_SSH, "shift.privatekey")

	// Default public key filepath
	PUB_KEY_PATH = filepath.Join(DIR_SSH, "shift.publickey")
)
