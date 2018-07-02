/*
Copyright 2018 The Elasticshift Authors.
*/
package types

type KubernetesClientOptions struct {
	KubeConfig []byte
	Namespace  string
}

type CreatePersistentVolumeOptions struct {
	Name     string
	Capacity string // Specfic to kubernetes
	
	//NFS
	Server       string
	Path         string
	MountOptions []string
	
	// minio
	Url         string
	AccessKey   string
	AccessToken string
}

type PersistentVolumeInfo struct {

	Name string
}
