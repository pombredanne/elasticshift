/*
Copyright 2018 The Elasticshift Authors.
*/
package worker

import (
	"path/filepath"
	"time"

	homedir "github.com/minio/go-homedir"
)

var (

	// Default builder timeout
	DEFAULT_TIMEOUT, _ = time.ParseDuration("120m")

	// SSH directory name
	DIR_SSH, _ = homedir.Expand("~/.ssh")
)

const (

	// Default GRPC port where the worker is listening for commands from shift server
	DEFAULT_GRPC_PORT = "5053"

	// Bit size used when generating RSA keys.
	DEFAULT_BIT_SIZE = 2048

	PRIV_KEY_NAME = "shift.privatekey"
	PUB_KEY_NAME  = "shift.publickey"
)

var (
	// Default private key filepath
	PRIV_KEY_PATH = filepath.Join(DIR_SSH, PRIV_KEY_NAME)

	// Default public key filepath
	PUB_KEY_PATH = filepath.Join(DIR_SSH, PUB_KEY_NAME)
)

func GetSSHDir() (string, error) {

	expanded, err := homedir.Expand(DIR_SSH)
	if err != nil {
		return "", err
	}
	return expanded, err
}
