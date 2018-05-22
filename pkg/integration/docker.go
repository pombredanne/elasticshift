/*
Copyright 2017 The Elasticshift Authors.
*/
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dclient "github.com/docker/docker/client"
)

const DefaultHost = "unix:///var/run/docker.sock"

type ClientOptions struct {
	Host        string
	Certificate string
	Key         string
	CACert      string
	Ctx         context.Context
	Version     string
}

type ContainerInfo struct {
	ID        string
	Name      string
	StartedAt time.Time
	StoppedAt time.Time
	LifeTime  time.Time
}

type dockerClient struct {
	cli *dclient.Client
	ctx context.Context
}

type ImagePullStatus struct {
	ID       string `json:"id,omitempty"`
	Status   string `json:"status,omitempty"`
	Progress string `json:"progress,omitempty"`
}

type DockerClient interface {
	CreateContainer(opts *container.Config, hostConfig *container.HostConfig, name string) (string, error)
	DeleteContainer(id string) error
	StartContainer(id string) error
	StopContainer(id string) error
	CLI() *dclient.Client
}

type StreamWriter struct {
	w io.Writer
}

func (d *dockerClient) CLI() *dclient.Client {
	return d.cli
}

func NewClient(opts *ClientOptions) (DockerClient, error) {

	var httpClient *http.Client
	c, err := dclient.NewClient(opts.Host, opts.Version, httpClient, nil)
	if err != nil {
		return &dockerClient{}, fmt.Errorf("Failed to create Docker Client for host '%s:%s, :%v", opts.Host, opts.Version, err)
	}

	dc := &dockerClient{c, opts.Ctx}
	return dc, nil
}

func (sw StreamWriter) Write(b []byte) (int, error) {

	l := len(b)
	stats := strings.Split(string(b), "\n")
	for _, stat := range stats {

		if stat != "" {

			ps := ImagePullStatus{}
			err := json.Unmarshal([]byte(stat), &ps)
			if err != nil {
				return l, err
			}
			sw.w.Write([]byte(ps.String()))
		}
	}
	return l, nil
}

func (p ImagePullStatus) String() string {

	var buf bytes.Buffer
	if p.ID != "" {
		buf.WriteString(p.ID)
		buf.WriteString(": ")
	}
	buf.WriteString(p.Status)

	if p.Progress != "" {

		idx := strings.Index(p.Progress, "]")
		buf.WriteString(" ")
		buf.WriteString(strings.TrimSpace(p.Progress[idx+1:]))
	}
	buf.WriteString("\n")
	return buf.String()
}

func (dc *dockerClient) CreateContainer(opts *container.Config, hostConfig *container.HostConfig, name string) (string, error) {

	r, err := dc.cli.ImagePull(dc.ctx, opts.Image, types.ImagePullOptions{})
	if err != nil {
		fmt.Printf("Image pull error if any: %v\n", err)
		return "", err
	}
	defer r.Close()

	_, err = io.Copy(&StreamWriter{w: os.Stdout}, r)
	if err != nil {
		return "", fmt.Errorf("Failed to grab image pull progress %v", err)
	}

	container, err := dc.cli.ContainerCreate(dc.ctx, opts, hostConfig, nil, name)
	if err != nil {
		return "", fmt.Errorf("Failed to create container %s:%v", opts.Image, err)
	}

	return container.ID, nil
}

func (dc *dockerClient) DeleteContainer(id string) error {
	return dc.cli.ContainerRemove(dc.ctx, id, types.ContainerRemoveOptions{})
}

func (dc *dockerClient) StartContainer(id string) error {
	return dc.cli.ContainerStart(dc.ctx, id, types.ContainerStartOptions{})
}

func (dc *dockerClient) StopContainer(id string) error {
	fmt.Print("Stopping container ", id[:10], "... ")
	return dc.cli.ContainerStop(dc.ctx, id, nil)
}
