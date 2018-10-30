/*
Copyright 2018 The Elasticshift Authors.
*/
package storage

import (
	"mime/multipart"

	"gitlab.com/conspico/elasticshift/api/types"
)

var (
	DIR_PLUGIN_BUNDLE = "plugin-bundles"
	DIR_PLUGIN        = "plugins"

	TARBALL_EXT = ".tar.gz"
	BUNDLE_NAME = "bundle" + TARBALL_EXT
)

const (
	MINIO = iota + 1
	AmazonS3
	GoogleCloudStorage
	NFS
)

func WritePluginBundle(stor types.Storage, f multipart.File, destPath string) error {

	var err error
	switch stor.Kind {
	//writeMinio(stor, f)
	case MINIO:
	case AmazonS3:
	case GoogleCloudStorage:
	case NFS:
		err = writeNFS(stor, f, destPath)
	}
	return err
}
