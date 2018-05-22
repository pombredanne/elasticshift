/*
Copyright 2018 The Elasticshift Authors.
*/
package build

import (
	"fmt"

	"gitlab.com/conspico/elasticshift/api/types"
	"gitlab.com/conspico/elasticshift/pkg/integration"
	itypes "gitlab.com/conspico/elasticshift/pkg/integration/types"
)

func (r *resolver) GetContainerEngine(team string) (integration.ContainerEngineInterface, error) {

	// Get the default container engine id based on team
	dce, err := r.defaultStore.GetDefaultContainerEngine(team)
	if err != nil {
		return nil, err
	}

	// Get the details of the integration
	var i types.ContainerEngine
	err = r.integrationStore.FindByID(dce, &i)

	// connect to container engine cluster
	return integration.New(r.logger, i)
}

func (r *resolver) ContainerLauncher() {

	defer r.recoverErrorIfAny()

	// TODO handle panic
	for b := range r.BuildQueue {

		go func(b types.Build) {

			// start the container
			// TODO select the default orchestration, by config
			// opts := &docker.ClientOptions{}
			// opts.Host = docker.DefaultHost
			// opts.Ctx = r.Ctx

			// cli, err := docker.NewClient(opts)
			// if err != nil {
			// 	r.SLog(b.ID, fmt.Sprintf("Failed to connect to docker daemon: %v", err))
			// }

			imgName, err := r.findImageName(b)
			if err != nil {
				r.SLog(b.ID, fmt.Sprintf("Unable to find the build image from Shiftfile", b.CloneURL))
				return
			}
			fmt.Println("Image name: " + imgName)

			// Identify the default orchestration based integration
			// such as docker swarm or kubernetes etc
			engine, err := r.GetContainerEngine(b.Team)
			if err != nil {
				//udpate the build log and set the status to failed
				r.logger.Errorf("Failed to connect container engine: %v", err)
			}

			// find the system storage
			// storage, err := r.sysconfStore.GetDefaultStorage()
			// if err != nil {
			// 	r.SLog(b.ID, "Failed to fetch the default storage: "+err.Error())
			// 	return
			// }

			// err = utils.Mkdir(filepath.Join(storage.Path, "code", b.Team))
			// if err != nil {
			// 	r.SLog(b.ID, "Unable to create directory for cloning the project:"+err.Error())
			// }

			// hostIp := utils.GetIP()
			// if hostIp == "" {
			// 	hostIp = "127.0.0.1"
			// }

			// env := []string{
			// 	"SHIFT_HOST=shiftserver",
			// 	"SHIFT_PORT=5051",
			// 	"SHIFT_LOGGER=" + LogType_File,
			// 	"SHIFT_BUILDID=" + b.ID.Hex(),
			// 	"SHIFT_TIMEOUT=120m",
			// 	"WORKER_PORT=" + "6060",
			// }

			// filepath.Join(storage.Path, b.Team, DIR_CODE)

			// hc := &container.HostConfig{}
			// hc.Binds = []string{
			// 	filepath.Join(storage.Path, b.Team, DIR_CODE) + ":" + VOL_CODE,
			// 	filepath.Join(storage.Path, b.Team, DIR_LOGS) + ":" + VOL_LOGS,
			// 	filepath.Join(storage.Path, DIR_PLUGINS) + ":" + VOL_PLUGINS,
			// 	filepath.Join(storage.Path, DIR_WORKER) + ":" + VOL_SHIFT,
			// }

			// workerPort, _ := nat.NewPort("tcp", "6060")
			// serverPort, _ := nat.NewPort("tcp", "5051")

			// exposedPorts := map[nat.Port]struct{}{
			// 	serverPort: struct{}{},
			// 	workerPort: struct{}{},
			// }

			// c := &container.Config{
			// 	Image:        imgName,
			// 	Entrypoint:   strslice.StrSlice{"./shift/worker"},
			// 	Env:          env,
			// 	AttachStdout: true,
			// 	ExposedPorts: exposedPorts,
			// }

			envs := []itypes.Env{
				itypes.Env{"SHIFT_HOST", "shiftserver"},
				itypes.Env{"SHIFT_PORT", "5051"},
				itypes.Env{"SHIFT_LOGGER", LogType_File},
				itypes.Env{"SHIFT_BUILDID", b.ID.Hex()},
				itypes.Env{"SHIFT_TIMEOUT", "120m"},
				itypes.Env{"WORKER_PORT", "6060"},
			}

			opts := &itypes.CreateContainerOptions{}
			opts.Image = imgName
			opts.Command = "./shift/worker"
			opts.Environment = envs
			opts.BuildID = b.ID.Hex()

			res, err := engine.CreateContainer(opts)
			if err != nil {
				str := fmt.Sprintf("Unable to create the container %v", err)
				r.SLog(b.ID, str)
				return
			}

			fmt.Println("Container ID =", res.UID)
			err = r.store.UpdateContainerID(b.ID, res.UID)
			if err != nil {
				r.logger.Errorln("Failed to update the container id: ", res.UID)
			}

			// err = cli.StartContainer(containerID)
			// if err != nil {
			// 	r.logger.Errorln("Failed to start the container: %v", err)
			// }
		}(b)
	}
}

func (r *resolver) recoverErrorIfAny() {

	err := recover()
	fmt.Println("recovered : %v", err)
}
