/*
Copyright 2018 The Elasticshift Authors.
*/
package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	dclient "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	itypes "github.com/elasticshift/elasticshift/internal/shiftserver/integration/types"
)

type dockerClient struct {
	cli  *dclient.Client
	ctx  context.Context
	opts *ConnectOptions
}

// ConnectDocker ...
func ConnectDocker(logger *logrus.Entry, opts *ConnectOptions) (ContainerEngineInterface, error) {
	var httpClient *http.Client
	c, err := dclient.NewClient(opts.Host, opts.Version, httpClient, nil)
	if err != nil {
		return &dockerClient{}, fmt.Errorf("Failed to create Docker Client for host '%s:%s, :%v", opts.Host, opts.Version, err)
	}

	dc := &dockerClient{c, opts.Ctx, opts}
	return dc, nil
}

func (c *dockerClient) CreateContainer(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error) {

	_, err := c.cli.ImagePull(c.ctx, opts.Image, dtypes.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("Image pull error if any: %v \n ", err)
	}
	//defer r.Close()

	// _, err = io.Copy(&StreamWriter{w: os.Stdout}, r)
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to grab image pull progress %v", err)
	// }

	var envs = []string{}
	for _, env := range opts.Environment {
		envs = append(envs, env.Key+"="+env.Value)
	}

	if c.opts.Storage.Kind == 4 { //NFS

	} else if c.opts.Storage.Kind == 1 { // Minio

		var bucketName string
		if c.opts.Storage.Minio.BucketName != "" {
			bucketName = c.opts.Storage.Minio.BucketName
		} else if c.opts.Storage.Name != "" {
			bucketName = c.opts.Storage.Name
		} else {
			bucketName = "elasticshift"
		}

		m := c.opts.Storage.StorageSource.Minio
		envs = append(envs, KEY_SHIFTDIR+"=/tmp")
		envs = append(envs, KEY_WORKER_URL+"="+m.Host+"/"+filepath.Join(bucketName, c.opts.Storage.WorkerPath))
		envs = append(envs, KEY_BUCKET+"="+bucketName)
	}

	hc := &container.HostConfig{}
	// hc.Binds = []string{
	// 	filepath.Join(storage, "code", team) + ":/code",
	// 	filepath.Join(storage, "plugins") + ":/plugins",
	// 	filepath.Join(storage, "worker") + ":/worker",
	// }

	cfg := &container.Config{
		Image: opts.Image,
		//Cmd:   []string{opts.Command},
		//Entrypoint: strslice.StrSlice{opts.Command},
		// Volumes: volumes,
		Entrypoint:   strslice.StrSlice{"/bin/bash", "-c", opts.Command},
		Tty:          true,
		Env:          envs,
		AttachStdout: true,
		AttachStderr: true,
	}

	name := opts.BuildID + "-" + opts.SubBuildID
	containerResult, err := c.cli.ContainerCreate(c.ctx, cfg, hc, nil, name)
	if err != nil {
		return nil, fmt.Errorf("Failed to create container %s:%v", opts.Image, err)
	}

	cinfo := &itypes.ContainerInfo{
		Name:              name,
		CreationTimestamp: time.Now().String(),
		Image:             opts.Image,
		ImageVersion:      opts.ImageVersion,
		// ClusterName:  md.GetClusterName(),
		UID: containerResult.ID,
		// ShiftID: shiftID.String(),
	}

	err = c.cli.ContainerStart(c.ctx, containerResult.ID, dtypes.ContainerStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to start the container: %v", err)
	}

	// go func() {
	// 	statusCh, err := c.cli.ContainerWait(c.ctx, containerResult.ID)
	// 	if err != nil {
	// 		cinfo.Status = err.Error()
	// 	}
	// 	stcode, _ := strconv.ParseInt(statusCh, 10, 64)
	// }()

	return cinfo, nil
}

func (c *dockerClient) CreateContainerWithVolume(opts *itypes.CreateContainerOptions) (*itypes.ContainerInfo, error) {
	return nil, nil
}

func (c *dockerClient) CreatePersistentVolume(opts *itypes.CreatePersistentVolumeOptions) (*itypes.PersistentVolumeInfo, error) {
	return nil, nil
}

func (c *dockerClient) DeleteContainer(id string) error {
	err := c.cli.ContainerStop(c.ctx, id, nil)
	if err != nil {
		return fmt.Errorf("Failed to stop the container(%s) :%v", id, err)
	}

	err = c.cli.ContainerRemove(c.ctx, id, dtypes.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("Failed to remote the container(%s): %v", id, err)
	}
	return nil
}

func (c *dockerClient) StreamLog(opts *itypes.StreamLogOptions) (io.ReadCloser, error) {

	options := dtypes.ContainerLogsOptions{ShowStdout: true}
	out, err := c.cli.ContainerLogs(c.ctx, opts.ContainerID, options)
	if err != nil {
		return nil, fmt.Errorf("Failed to stream logs for container %s: %v", opts.ContainerID, err)
	}
	return out, nil
}
