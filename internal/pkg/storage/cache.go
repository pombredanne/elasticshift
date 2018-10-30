/*
Copyright 2018 The Elasticshift Authors.
*/
package storage

import (
	"gitlab.com/conspico/elasticshift/api/types"
)

// PutCacheFile ..
func (s *ShiftStorage) PutCacheFile(name, path string) (int64, error) {
	return s.PutCacheFileWithMetadata(name, path, s.metadata)
}

// PutCacheFileWithMetadata ..
func (s *ShiftStorage) PutCacheFileWithMetadata(name, path string, m *types.StorageMetadata) (int64, error) {

	objectName := GetObjectName(m, name, cacheDir)
	return s.stor.PutFObject(s.bucketName, objectName, path, cacheContentType)
}

// GetCacheFile ..
func (s *ShiftStorage) GetCacheFile(name, path string) error {
	return s.GetCacheFileWithMetadata(name, path, s.metadata)
}

// GetCacheFileWithMetadata ..
func (s *ShiftStorage) GetCacheFileWithMetadata(name, path string, m *types.StorageMetadata) error {

	objectName := GetObjectName(m, name, cacheDir)
	return s.stor.GetFObject(s.bucketName, objectName, path)
}
