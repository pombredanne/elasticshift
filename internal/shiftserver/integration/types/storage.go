/*
Copyright 2018 The Elasticshift Authors.
*/
package types

import "io"

type CreateBucketOptions struct {
	Name   string
	Region string
}

type PutObjectStreamingOptions struct {
	BucketName string
	ObjectName string
	Reader     io.Reader
}
