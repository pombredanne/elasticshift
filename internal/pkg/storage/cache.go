/*
Copyright 2018 The Elasticshift Authors.
*/
package storage

import (
	"strings"

	"github.com/elasticshift/elasticshift/api/types"
)

var (
	errKeyDoesNotExist = "The specified key does not exist."
)

// PutCacheFile ..
func (s *ShiftStorage) PutCacheFile(name, path string) (int64, error) {
	return s.PutCacheFileWithMetadata(name, path, s.metadata)
}

// PutCacheFileWithMetadata ..
func (s *ShiftStorage) PutCacheFileWithMetadata(name, path string, m *types.StorageMetadata) (int64, error) {

	objectName := GetCacheObjectName(m, name)
	return s.stor.PutFObject(s.bucketName, objectName, path, cacheContentType)
}

// GetCacheFile ..
func (s *ShiftStorage) GetCacheFile(name, path string) error {
	return s.GetCacheFileWithMetadata(name, path, s.metadata)
}

// GetCacheFileWithMetadata ..
func (s *ShiftStorage) GetCacheFileWithMetadata(name, path string, m *types.StorageMetadata) error {

	objectName := GetCacheObjectName(m, name)
	err := s.stor.GetFObject(s.bucketName, objectName, path)
	if err != nil && err.Error() != errKeyDoesNotExist {
		return err
	}

	return nil
}

func GetCacheObjectName(m *types.StorageMetadata, name string) string {

	var b strings.Builder
	b.WriteString(GetCachePath(m, cacheDir))
	b.WriteString(objectSeparator)
	b.WriteString(name)

	return b.String()
}

func GetCachePath(m *types.StorageMetadata, typee string) string {

	var b strings.Builder
	b.WriteString(typee)
	b.WriteString(objectSeparator)
	b.WriteString(m.Path)

	return b.String()
}
