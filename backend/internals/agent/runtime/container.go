package runtime

import (
	"aetrix/observer/internals/agent/protocol"
	"aetrix/observer/internals/lib"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

type DockerRuntimeInterface interface {
	StartNewContainer(ctx context.Context, ResWriter chan<- protocol.Event, req lib.Command)
}

type DockerRuntime struct {
	dockerClient DockerClient
}

// NewHandler creates and returns a new Handler instance with the provided DockerClient.
func NewDockerRuntime(dockerClient DockerClient) *DockerRuntime {
	return &DockerRuntime{
		dockerClient: dockerClient,
	}
}

// starts a new container by
// pull image ,creates container , start
func (R *DockerRuntime) StartNewContainer(ctx context.Context, ResWriter chan<- protocol.Event, req lib.Command) {
	StartConfig := req.Params.StartNewContainerConfig
	fmt.Println("rec", req)
	// 1. Pull the Docker image
	imageName := StartConfig.Image

	reader, err := R.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})

	if err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx, ResWriter)
		return
	}
	buf := make([]byte, 1024)
	for {
		n, readErr := reader.Read(buf)
		if n > 0 {
			protocol.NewEvent(
				protocol.WithMessage(string(buf[:n])),
			).TrySend(ctx, ResWriter)
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
		}
	}
	// Close the reader to free resources
	reader.Close()
	protocol.NewEvent(
		protocol.WithMessage("Image pulled successfully: "+imageName),
	).TrySend(ctx, ResWriter)

	if StartConfig.Protocol == "" {
		StartConfig.Protocol = "tcp"
	}
	if StartConfig.HostIP == "" {
		StartConfig.HostIP = "0.0.0.0"
	}

	// 3. Create the container
	ContainerPort := nat.Port(fmt.Sprintf("%d/%s", StartConfig.ContainerPort, StartConfig.Protocol))

	cont_config := &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{ContainerPort: struct{}{}},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			ContainerPort: []nat.PortBinding{
				{
					HostIP:   StartConfig.HostIP,
					HostPort: fmt.Sprintf("%d", StartConfig.HostPort),
				},
			},
		},
	}

	resp, err := R.dockerClient.ContainerCreate(
		ctx,
		cont_config,
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		StartConfig.Name,
	)

	if err != nil {
		errEvent := protocol.NewEvent(protocol.WithError(err))
		errEvent.Send(ctx, ResWriter)
		return
	}

	protocol.NewEvent(
		protocol.WithMessage("Container created successfully with ID: "+resp.ID),
		protocol.WithData(resp),
	).TrySend(ctx, ResWriter)
	// 4. Start the container

	if err := R.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx, ResWriter)
		return
	}

	protocol.NewEvent(
		protocol.WithMessage("Container started successfully with ID: "+resp.ID),
	).Send(ctx, ResWriter)

}

// HandleListContainers retrieves and returns a list of Docker containers.
// The context handles cancellation and timeout, and the request may include optional filters.
// It returns a WsMessage containing the container list or an error message.
func (R *DockerRuntime) ListContainers(ctx context.Context, ResWriter chan<- protocol.Event, req lib.Command) {
	options := container.ListOptions{}
	if req.Params.ListContainersConfig != nil {
		options.All = req.Params.ListContainersConfig.All
		options.Size = req.Params.ListContainersConfig.Size
	}

	containers, err := R.dockerClient.ContainerList(ctx, options)
	if err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx, ResWriter)
		return
	}
	protocol.NewEvent(
		protocol.WithData(containers),
		protocol.WithMessage("Listed containers successfully"),
	).Send(ctx, ResWriter)
}

// HandleDeleteContainer removes a Docker container by its ID.
// The context handles cancellation and timeout, and the request contains the container ID and removal options.
func (R *DockerRuntime) DeleteContainer(ctx context.Context, ResWriter chan<- protocol.Event, req lib.Command) {
	configs := req.Params.DeleteContainerConfig
	options := container.RemoveOptions{}
	if configs != nil {
		options.Force = configs.Force
		options.RemoveVolumes = configs.Volume
		options.RemoveLinks = configs.Links
	} else {
		options.Force = true
	}
	err := R.dockerClient.ContainerRemove(ctx, req.ContainerID, options)
	if err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx, ResWriter)
		return
	}

	protocol.NewEvent(
		protocol.WithMessage("Cotainer deleted successfully"),
		protocol.WithData(struct {
			ID string
		}{
			ID: req.ContainerID,
		}),
	).Send(ctx, ResWriter)
}

func (R *DockerRuntime) StopContainer(ctx context.Context, ResWriter chan<- protocol.Event, req lib.Command) {
	// Using a timeout of 5 seconds for stopping the container
	// overrides the default 10 seconds
	timeout := 5 // seconds
	err := R.dockerClient.ContainerStop(ctx, req.ContainerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx, ResWriter)
		return
	}

	data := struct {
		Id string
	}{
		Id: req.ContainerID,
	}

	protocol.NewEvent(
		protocol.WithMessage("Container stoped successfully"),
		protocol.WithData(data),
	).Send(ctx, ResWriter)
}


// timeout of one second (ie it take one second btw starting and stoping)
func (R * DockerRuntime ) RestartContainer(ctx context.Context,resWriter chan<- protocol.Event  , req lib.Command){

	timeout := 1 // immediate restart
	err := R.dockerClient.ContainerRestart(ctx, req.ContainerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx , resWriter)

		return
	}
	data:= struct{
		Id string
	}{
		Id : req.ContainerID,
	}
	protocol.NewEvent(
		protocol.WithData(data),
		protocol.WithMessage("Container restarted successfully"),
	).Send(ctx , resWriter)
}


func (R *DockerRuntime) StartContainer(ctx context.Context, resWirter chan<- protocol.Event, req lib.Command){
	err := R.dockerClient.ContainerStart(ctx, req.ContainerID, container.StartOptions{})
	if err != nil {
		protocol.NewEvent(
			protocol.WithError(err),
		).Send(ctx , resWirter)
		return
	}
	data:= struct{
		Id string
	}{
		Id : req.ContainerID,
	}
	protocol.NewEvent(
		protocol.WithData(data),
		protocol.WithMessage("Container started suceesfully"),
	).Send(ctx , resWirter)
}
