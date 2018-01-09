/*
Copyright 2017 The Elasticshift Authors.
*/
package kubernetes

import "time"

type Env struct {
	Key   string
	Value string
}

type CreateContainerOptions struct {
	Image                string
	ImageVersion         string
	Command              string
	Environment          []Env
	PersistentVolumeName string // Used internally
	privileged           bool
	RestartPolicy        string
}

type ContainerInfo struct {
	StartedAt time.Time
	//StoppedAt    time.Time
	Status       string
	Image        string
	ImageVersion string

	ClusterName       string
	CreationTimestamp string
	Uid               string
	Namespace         string
	Name              string
}

type KubernetesClientOptions struct {
	KubeConfigFile string
	Namespace      string
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
