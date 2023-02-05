package dockerhelper

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/dapr/cli/pkg/print"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/katasec/ark/utils"
)

var (
	arkSpinner = utils.NewArkSpinner()
)

type DockerHelper struct {
	cli *client.Client
	ctx context.Context
}

func NewDockerHelper() *DockerHelper {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &DockerHelper{
		cli: cli,
		ctx: ctx,
	}

}

func (d *DockerHelper) StartContainerUI(imageName string, envVars []string, port string, containerName string, cmd []string, volumemounts ...string) {

	running, _, _ := d.IsRunning(imageName, containerName)

	if !running {

		// Force container remove
		d.ContainerRemove("/" + containerName)

		// Pull Image
		note := "Download docker image for " + containerName + ": " + imageName
		arkSpinner.Start(note)
		err := d.Pull(imageName)
		arkSpinner.Stop(err, note)

		// Run Image
		note = "Running " + containerName + ": " + imageName
		arkSpinner.Start(note)

		// Pass volumemounts if any and run container
		if volumemounts == nil {
			d.RunContainer(imageName, envVars, port, containerName, cmd)
		} else {
			d.RunContainer(imageName, envVars, port, containerName, cmd, volumemounts...)
		}

		arkSpinner.Stop(err, note)
	} else {
		note := fmt.Sprintf("%s %s is already running !", containerName, imageName)
		arkSpinner.InfoStatusEvent(note)
	}
}

func (d *DockerHelper) StopContainerUI(imageName string, containerName string) {
	note := "Stopping " + containerName + ": " + imageName
	arkSpinner.Start(note)

	err := d.ContainerRemove("/" + containerName)
	arkSpinner.Stop(err, note)
}

func (d *DockerHelper) ListContainers() (containers []types.Container, err error) {
	// Get list of containers
	containers, err = d.cli.ContainerList(d.ctx, types.ContainerListOptions{
		All: true,
	})

	if err != nil {
		arkSpinner.ErrorStatusEvent("Plesae start docker for running Ark dev, exitting !")
		os.Exit(1)
	}

	return containers, err
}

func (d *DockerHelper) Pull(imageName string) (err error) {

	// Start image pull
	out, err := d.cli.ImagePull(d.ctx, imageName, types.ImagePullOptions{})
	//defer out.Close()

	if err != nil {
		print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Could not pull image, error: %s\n", err.Error()))
		return err
	}

	// Waiting for pull to complete
	if _, err := ioutil.ReadAll(out); err != nil {
		panic(err)
	}
	return err
}

func (d *DockerHelper) IsRunning(imageName string, name string, status ...string) (running bool, state string, id string) {

	var checkingFor string
	if len(status) == 0 {
		checkingFor = "running"
	} else {
		checkingFor = strings.ToLower(status[0])
	}

	// init locals
	running = false
	id = ""
	state = "none"

	// Get list of containers
	containers, _ := d.ListContainers()

	// Iterate to find container with imageName
	for _, container := range containers {
		// Check & return container state
		if container.Image == imageName && container.State == checkingFor && container.Names[0] == "/"+name {
			//fmt.Println("The status was:" + container.State)
			id = container.ID[:12]
			running = true
			state = container.State
		} else if container.Image == imageName {
			state = container.State
		}
	}

	return running, state, id
}

func (d *DockerHelper) RunContainer(imageName string, envvars []string, port string, containerName string, cmd []string, volumemounts ...string) (err error) {

	// Create port spec for e.g "tcp/80"
	portSpec, _ := nat.NewPort("tcp", port)

	// Define container config
	containerConfig := &container.Config{
		Image: imageName,
		Env:   envvars,
		ExposedPorts: nat.PortSet{
			portSpec: struct{}{},
		},
		Cmd: cmd,
	}

	// Genrate bind mounts
	mounts := []mount.Mount{}

	if len(volumemounts) > 0 {
		for _, vmount := range volumemounts {
			source := strings.Split(vmount, "|")[0]
			target := strings.Split(vmount, "|")[1]
			myMount := mount.Mount{
				Type:   mount.TypeBind,
				Source: source,
				Target: target,
			}
			mounts = append(mounts, myMount)
		}
	} else {
		mounts = nil
	}

	// Define container->host port map
	var portBindings nat.PortMap
	if port != "0" {
		portBindings = nat.PortMap{
			portSpec: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		}
	} else {
		portBindings = nil
	}

	// Setup ports, mounts etc
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Mounts:       mounts,
	}

	// Create container
	resp, err := d.cli.ContainerCreate(d.ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())

	}

	// Start the container
	containerID := resp.ID
	if err := d.cli.ContainerStart(d.ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// Loop until container starts
	//var state string
	var running bool
	for i := 1; i <= 10; i++ {
		//running, state, _ = d.IsRunning(imageName, containerName)
		running, _, _ = d.IsRunning(imageName, containerName)
		if running {
			break
		}

		//fmt.Printf("Container state is :%v\n", state)
		time.Sleep(5 * time.Second)
	}

	return err
}

// ContainerRemove kills and removes a container from the docker host.
func (d *DockerHelper) ContainerRemove(containerID string) error {
	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}
	err := d.cli.ContainerRemove(d.ctx, "/"+containerID, options)
	if err != nil && !strings.Contains(err.Error(), "No such container") {
		fmt.Println(err.Error())
	}
	return err
}
