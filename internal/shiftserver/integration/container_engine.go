/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"context"
	"fmt"
	"io"

	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/pkg/logger"
	itypes "github.com/elasticshift/elasticshift/internal/shiftserver/integration/types"
)

//container engine
const (
	Kubernetes int = iota + 1
	Docker
	DCOS
)

type containerEngine struct {
	provider int
	i        types.ContainerEngine
}

type ContainerEngineInterface interface {
	CreateContainer(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error)
	CreateContainerWithVolume(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error)
	CreatePersistentVolume(opts *itypes.CreatePersistentVolumeOptions) (*itypes.PersistentVolumeInfo, error)
	DeleteContainer(id string) error
	StreamLog(opts *itypes.StreamLogOptions) (io.ReadCloser, error)
}

func NewContainerEngine(loggr logger.Loggr, i types.ContainerEngine, s types.Storage) (ContainerEngineInterface, error) {

	switch i.Kind {
	case DCOS:
	case Docker:
		opts := &ConnectOptions{}
		opts.Storage = s
		opts.Host = i.Host
		opts.Version = i.Version
		opts.InsecureSkipTLSVerify = true
		opts.ServerCertificate = i.Certificate
		opts.Ctx = context.Background()
		return ConnectDocker(loggr.GetLogger("engine/docker"), opts)
	case Kubernetes:
		opts := &ConnectOptions{}
		opts.Storage = s
		opts.UseConfig = len(i.KubeFile) > 0
		if opts.UseConfig {
			opts.Config = i.KubeFile
		} else {
			opts.Host = i.Host
			opts.ServerCertificate = i.Certificate
			opts.Token = i.Token
		}
		opts.InsecureSkipTLSVerify = true
		return ConnectKubernetes(loggr.GetLogger("engine/kubernetes"), opts)
	}

	return nil, fmt.Errorf("No container engine to connect.")
}
