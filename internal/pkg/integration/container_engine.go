/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	itypes "gitlab.com/conspico/elasticshift/internal/pkg/integration/types"
)

//container engine
const (
	Kubernetes int = iota + 1
	DockerSwarm
	DCOS
)

type containerEngine struct {
	provider int
	logger   logrus.Logger
	i        types.ContainerEngine
}

type ContainerEngineInterface interface {
	CreateContainer(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error)
	CreateContainerWithVolume(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error)
	CreatePersistentVolume(opts *itypes.CreatePersistentVolumeOptions) (*itypes.PersistentVolumeInfo, error)
}

func NewContainerEngine(logger logrus.Logger, i types.ContainerEngine) (ContainerEngineInterface, error) {

	switch i.Kind {
	case DCOS:
	case DockerSwarm:
	case Kubernetes:
		opts := &ConnectOptions{}
		opts.Host = i.Host
		opts.ServerCertificate = i.Certificate
		opts.Token = i.Token
		opts.InsecureSkipTLSVerify = true
		return ConnectKubernetes(logger, opts)
	}

	return nil, fmt.Errorf("No container engine to connect.")
}
