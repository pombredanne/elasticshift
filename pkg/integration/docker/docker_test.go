package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"gitlab.com/conspico/elasticshift/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

const (
	LogType_Embedded = "Embedded"
	LogType_NFS      = "NFS"
)

func TestDocker(t *testing.T) {

	buildId := bson.NewObjectId().Hex()
	team := "elasticshift"

	ctx := context.Background()
	opts := &ClientOptions{}
	opts.Host = DefaultHost
	opts.Ctx = ctx

	cli, err := NewClient(opts)
	if err != nil {
		t.Logf("%v", err)
	}

	env := []string{
		"SHIFT_HOST=127.0.0.1",
		"SHIFT_PORT=5050",
		"SHIFT_LOGGER=" + LogType_Embedded,
		"SHIFT_BUILDID=" + buildId,
	}

	imgName := "openjdk:7"
	// imgName := "alpine:latest"
	storage, err := filepath.Abs("/Users/ghazni/elasticshift")
	if err != nil {
		t.Logf("%v", err)
	}

	// volumes := make(map[string]struct{}, 1)
	// volumes[filepath.Join(storage.Path, "team", b.Team, "code")] = struct{ Code string }{}{"/code"}
	// volumes[filepath.Join(storage.Path, "plugins")] = "/plugins"
	// volumes[filepath.Join(storage.Path, "worker")] = "/worker"

	// volumes := map[string]struct{}{
	// 	filepath.Join(storage, team, "code") + ":/code": struct{}{},
	// 	filepath.Join(storage, "plugins") + ":/plugins": struct{}{},
	// 	filepath.Join(storage, "worker") + ":/worker":   struct{}{},
	// }

	// fmt.Println(volumes)

	err = utils.Mkdir(filepath.Join(storage, "code", team))
	if err != nil {
		t.Log(err)
	}

	hc := &container.HostConfig{}
	hc.Binds = []string{
		filepath.Join(storage, "code", team) + ":/code",
		filepath.Join(storage, "plugins") + ":/plugins",
		filepath.Join(storage, "worker") + ":/worker",
	}

	c := &container.Config{
		Image: imgName,
		// Cmd:   []string{"tail", "-f", "/dev/null"},
		// Entrypoint: strslice.StrSlice{"/worker/worker"},
		// Volumes: volumes,
		// Cmd:          []string{"ls", "-a", "/plugins"},
		Entrypoint:   strslice.StrSlice{"/worker/worker"},
		Tty:          true,
		Env:          env,
		AttachStdout: true,
		AttachStderr: true,
	}

	containerID, err := cli.CreateContainer(c, hc, buildId)
	if err != nil {
		str := fmt.Sprintf("Unable to create the container %v", err)
		t.Log(str)
	}

	err = cli.StartContainer(containerID)
	if err != nil {
		t.Log(err)
	}

	// statusCh, errCh := cli.CLI().ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	// select {
	// case err := <-errCh:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// case <-statusCh:
	// }

	options := dtypes.ContainerLogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := cli.CLI().ContainerLogs(ctx, containerID, options)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	fmt.Println("Container ID =", containerID)
}
