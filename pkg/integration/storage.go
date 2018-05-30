/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
)

//storage
const (
	Minio int = iota + 1
	AmazonS3
	GoogleCloudStorage
)

type StorageInterface interface {
	CreateBucket(name, region string) error
	PutObjectStreaming(bucketName, objectName string, reader io.Reader) (int64, error)
	PutObject(bucketName, objectName string, r io.Reader, contentType string) (int64, error)
}

/*
Copyright 2018 The Elasticshift Authors.
*/
func NewStorage(logger logrus.Logger, s types.Storage) (StorageInterface, error) {

	switch s.Kind {
	case AmazonS3:
	case GoogleCloudStorage:
	case Minio:
		return ConnectMinio(logger, s)
	}

	return nil, fmt.Errorf("No storage provider to connect.")
}
