/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"io"
	"strings"

	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
)

type minioClient struct {
	opts   types.Storage
	cli    *minio.Client
	logger *logrus.Entry
}

func ConnectMinio(logger *logrus.Entry, opts types.Storage) (StorageInterface, error) {

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
	// return m.cli.PutObjectStreaming(bucketName, objectName, r)
	return 0, nil
}

func (m minioClient) PutObject(bucketName, objectName string, r io.Reader, contentType string) (int64, error) {
	return m.cli.PutObject(bucketName, objectName, r, -1, minio.PutObjectOptions{ContentType: contentType})
}

func (m minioClient) GetObject(bucketName, objectName string) (io.ReadCloser, error) {
	return m.cli.GetObject(bucketName, objectName, minio.GetObjectOptions{})
}

func (m minioClient) PutFObject(bucketName, objectName, filepath, contentType string) (int64, error) {
	return m.cli.FPutObject(bucketName, objectName, filepath, minio.PutObjectOptions{ContentType: contentType})
}

func (m minioClient) GetFObject(bucketName, objectName, filepath string) error {
	return m.cli.FGetObject(bucketName, objectName, filepath, minio.GetObjectOptions{})
}
