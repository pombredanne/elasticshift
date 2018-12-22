/*
Copyright 2018 The Elasticshift Authors.
*/
package storage

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/elasticshift/elasticshift/api"
	"github.com/elasticshift/elasticshift/api/types"
	"github.com/elasticshift/elasticshift/internal/shiftserver/integration"
)

var (

	// bucket
	defaultBucketName = "elasticshift"

	// specific dirs
	logDir     = "logs"
	cacheDir   = "cache"
	archiveDir = "archive"

	// content type
	logContentType   = "text/plain"
	cacheContentType = "application/octet-stream"

	objectSeparator = "/"
)

type ShiftStorage struct {
	stor       integration.StorageInterface
	logger     *logrus.Entry
	stortype   *types.Storage
	bucketName string
	metadata   *types.StorageMetadata
}

// New ..
// Storage interaction
func New(logger *logrus.Entry, s *types.Storage) (*ShiftStorage, error) {
	ss := &ShiftStorage{
		logger:   logger,
		stortype: s,
	}

	bucketName := os.Getenv("SHIFT_STORAGE_BUCKET")
	if bucketName == "" {
		bucketName = defaultBucketName
	}
	ss.bucketName = bucketName

	err := ss.connect()
	if err != nil {
		return nil, err
	}

	return ss, nil
}

// NewWithMetadata ..
// Storage interaction
func NewWithMetadata(logger *logrus.Entry, s *types.Storage, m *types.StorageMetadata) (*ShiftStorage, error) {

	ss, err := New(logger, s)
	if err != nil {
		return nil, err
	}
	ss.metadata = m
	return ss, nil
}

// Connect ..
func (s *ShiftStorage) connect() error {

	st, err := integration.NewStorage(s.logger, *s.stortype)
	if err != nil {
		return fmt.Errorf("Failed to connect to storage : %v", err)
	}
	s.stor = st

	return nil
}

func Convert(s *api.Storage) *types.Storage {

	var stor types.Storage
	if int(s.Kind) == integration.Minio {

		stor.Kind = integration.Minio
		ms := s.Minio

		m := &types.MinioStorage{}
		m.Host = ms.Host
		m.Certificate = ms.Certificate
		m.AccessKey = ms.AccessKey
		m.SecretKey = ms.SecretKey

		stor.Minio = m
	}

	return &stor
}
