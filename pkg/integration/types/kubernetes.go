/*
Copyright 2018 The Elasticshift Authors.
*/
package types

type KubernetesClientOptions struct {
	KubeConfig []byte
	Namespace  string
}

//go:generate stringer -type=PersistentVolumeProvider
type PersistentVolumeProvider int

const (
	NetworkFileShare PersistentVolumeProvider = iota
	GoogleCloudStorage
	HostLocalDirectory
)

type CreatePersistentVolumeOption struct {
	Server       string
	Path         string
	Name         string
	MountOptions []string
	provider     PersistentVolumeProvider
	Capacity     string // Specfic to kubernetes
}
