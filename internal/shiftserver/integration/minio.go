/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"io"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/minio/minio-go"
	"gitlab.com/conspico/elasticshift/api/types"
)

type minioClient struct {
	opts   types.Storage
	cli    *minio.Client
	logger logrus.Logger
}

func ConnectMinio(logger logrus.Logger, opts types.Storage) (StorageInterface, error) {

	mc := minioClient{
		opts:   opts,
		logger: logger,
	}

	var err error
	mc.cli, err = minio.New(opts.Minio.Host, opts.Minio.AccessKey, opts.Minio.SecretKey, strings.HasPrefix(opts.Minio.Host, "https"))

	return mc, err
}

func (m minioClient) CreateBucket(name, region string) error {
	return m.cli.MakeBucket(name, region)
}

func (m minioClient) PutObjectStreaming(bucketName, objectName string, r io.Reader) (int64, error) {
	return m.cli.PutObjectStreaming(bucketName, objectName, r)
}

func (m minioClient) PutObject(bucketName, objectName string, r io.Reader, contentType string) (int64, error) {
	return m.cli.PutObject(bucketName, objectName, r, contentType)
}
