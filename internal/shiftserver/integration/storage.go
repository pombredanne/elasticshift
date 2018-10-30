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
	GetObject(bucketName, objectName string) (io.ReadCloser, error)
	PutFObject(bucketName, objectName, filepath, contentType string) (int64, error)
	GetFObject(bucketName, objectName, filepath string) error
}

// NewStorage ..
func NewStorage(logger *logrus.Entry, s types.Storage) (StorageInterface, error) {

	switch s.Kind {
	case AmazonS3:
	case GoogleCloudStorage:
	case Minio:
		return ConnectMinio(logger, s)
	}

	return nil, fmt.Errorf("No storage provider to connect.")
}
