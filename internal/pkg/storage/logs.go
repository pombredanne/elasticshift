/*
Copyright 2018 The Elasticshift Authors.
*/
package storage

import (
	"io"
	"strings"

	"gitlab.com/conspico/elasticshift/api/types"
)

func (s *ShiftStorage) GetLog(name string) (io.ReadCloser, error) {
	return s.GetLogWithMetadata(name, s.metadata)
}

func (s *ShiftStorage) GetLogWithMetadata(name string, m *types.StorageMetadata) (io.ReadCloser, error) {

	objectName := GetObjectName(m, name, logDir)
	return s.stor.GetObject(s.bucketName, objectName)
}

func (s *ShiftStorage) PutLog(name string, r io.Reader) error {
	return s.PutLogWithMetadata(name, r, s.metadata)
}

func (s *ShiftStorage) PutLogWithMetadata(name string, r io.Reader, m *types.StorageMetadata) error {

	objectName := GetObjectName(m, name, logDir)
	_, err := s.stor.PutObject(s.bucketName, objectName, r, logContentType)
	return err
}

func GetObjectName(m *types.StorageMetadata, name, typee string) string {

	var b strings.Builder
	b.WriteString(GetPath(m, typee))
	b.WriteString(objectSeparator)
	b.WriteString(name)

	return b.String()
}

func GetPath(m *types.StorageMetadata, typee string) string {

	var b strings.Builder
	b.WriteString(typee)
	b.WriteString(objectSeparator)
	b.WriteString(m.TeamID)
	b.WriteString(objectSeparator)
	b.WriteString(m.RepositoryID)
	b.WriteString(objectSeparator)
	b.WriteString(m.BuildID)
	b.WriteString(objectSeparator)
	b.WriteString(m.SubBuildID)

	return b.String()
}
