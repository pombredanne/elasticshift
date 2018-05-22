/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
	itypes "gitlab.com/conspico/elasticshift/pkg/integration/types"
)

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
}

func New(logger logrus.Logger, i types.ContainerEngine) (ContainerEngineInterface, error) {

	// ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	// Name        string        `json:"name" bson:"name"`
	// Provider    int           `json:"provider" bson:"provider"`
	// Kind        int           `json:"kind" bson:"kind"`
	// Host        string        `json:"host" bson:"host"`
	// Certificate string        `json:"certificate" bson:"certificate"`
	// Token       string        `json:"token" bson:"token"`
	// Team        string        `json:"team" bson:"team"`

	switch i.Kind {
	case Kubernetes:
		opts := &ConnectOptions{}
		opts.Host = i.Host
		opts.ServerCertificate = i.Certificate
		opts.Token = i.Token
		opts.InsecureSkipTLSVerify = true
		return ConnectKubernetes(logger, opts)
	case DockerSwarm:
	case DCOS:
	}

	return nil, fmt.Errorf("Failed to connect to default container engine.")
}
