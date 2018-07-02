package types

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

	BuildID string
}

type ContainerInfo struct {
	StartedAt    time.Time
	StoppedAt    time.Time
	Status       string
	Image        string
	ImageVersion string

	ClusterName       string
	CreationTimestamp string
	UID               string
	ShiftID           string
	Namespace         string
	Name              string
}

type PersistentVolumeClaimOptions struct {
	Name string
	Capacity string
}
