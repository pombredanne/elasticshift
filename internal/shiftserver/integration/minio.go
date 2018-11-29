/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/conspico/elasticshift/api/types"
)

var (
	policyStr = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "*"
        ]
      },
      "Action": [
        "s3:GetBucketLocation"
      ],
      "Resource": [
        "arn:aws:s3:::%s"
      ]
    },
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "*"
        ]
      },
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::%s"
      ],
      "Condition": {
        "StringEquals": {
          "s3:prefix": [
            "sys"
          ]
        }
      }
    },
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": [
          "*"
        ]
      },
      "Action": [
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::%s/sys*"
      ]
    }
  ]
}`
	KeyHttps = "https://"
	KeyHttp  = "http://"

	KeyHttpsSize = len(KeyHttps)
	KeyHttpSize  = len(KeyHttp)
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

	var secure bool
	var host string
	if strings.HasPrefix(opts.Minio.Host, KeyHttps) {
		idx := strings.Index(opts.Minio.Host, KeyHttps)
		host = opts.Minio.Host[idx+KeyHttpsSize:]
		secure = true
	} else if strings.HasPrefix(opts.Minio.Host, KeyHttp) {
		idx := strings.Index(opts.Minio.Host, KeyHttp)
		host = opts.Minio.Host[idx+KeyHttpSize:]
	} else {
		host = opts.Minio.Host
	}

	var err error
	mc.cli, err = minio.New(host, opts.Minio.AccessKey, opts.Minio.SecretKey, secure)

	return mc, err
}

func (m minioClient) SetupStorage(bucketName, workerURL string) (string, error) {

	err := m.CreateBucket(bucketName, "")
	if err != nil {
		return "", fmt.Errorf("Failed to create bucket: %s: %v", bucketName, err)
	}

	idx := strings.LastIndex(workerURL, "/")
	name := workerURL[idx+1:]
	objectName := filepath.Join(defaultSysObject, name)

	r, err := http.Get(workerURL)
	if err != nil {
		return "", fmt.Errorf("Failed to download file from (%s) :%v", workerURL, err)
	}
	defer r.Body.Close()

	_, err = m.PutObject(bucketName, objectName, r.Body, defaultWorkerContentType)
	if err != nil {
		return "", fmt.Errorf("Failed to upload worker to storage : %v", err)
	}

	// fmt.Println("Bucket name=", bucketName)

	// privObj := filepath.Join(bucketName, defaultSysObject)
	// fmt.Println("Priv=", privObj)

	//	set bucket policy to download
	// var p = policy.BucketAccessPolicy{Version: "2012-10-2017"}
	// p.Statements = policy.SetPolicy(p.Statements, policy.BucketPolicyReadOnly, bucketName, "sys")

	// fpolicy, err := json.Marshal(p)
	// if err != nil {
	// 	return "", fmt.Errorf("Failed to set construct download policy: %v", err)
	// }

	err = m.cli.SetBucketPolicy(bucketName, fmt.Sprintf(policyStr, bucketName, bucketName, bucketName))
	if err != nil {
		return "", fmt.Errorf("Failed to set download privilege. :%v", err)
	}

	return objectName, nil
}

func (m minioClient) CreateBucket(name, region string) error {

	var err error
	exist, err := m.cli.BucketExists(name)
	if err != nil {
		return err
	}

	if !exist {
		err = m.cli.MakeBucket(name, region)
	}

	return err
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
