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

	VolumeMounts []Volume
	BuildID      string

	FailureFunc func(string, string, time.Time)
}

type Volume struct {
	Name      string
	MountPath string
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
	Name     string
	Capacity string
}
