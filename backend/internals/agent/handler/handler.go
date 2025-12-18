package handler

import (
	"aetrix/observer/internals/agent/docker"
	"aetrix/observer/internals/lib"
	"context"

	"fmt"
	"log"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types/container"
)

// Handler manages Docker operations and provides methods for container lifecycle management.
type Handler struct {
	dockerClient docker.DockerClient
}

// NewHandler creates and returns a new Handler instance with the provided DockerClient.
func NewHandler(dockerClient docker.DockerClient) *Handler {
	return &Handler{
		dockerClient: dockerClient,
	}
}

// HandleStartNewContainer pulls a Docker image, creates and starts a new container with the specified configuration.
// It returns a WsMessage containing the container ID on success or an error message on failure.
func (h *Handler) HandleStartNewContainer(ctx context.Context, req lib.Command) lib.WsMessage {
	StartConfig := req.Params.StartNewContainerConfig
	imageName := StartConfig.Image
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	wsMessage.Type = lib.TypeResponse

	_, err := h.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}
	if StartConfig.Protocol == "" {
		StartConfig.Protocol = "tcp"
	}
	if StartConfig.HostIP == "" {
		StartConfig.HostIP = "0.0.0.0"
	}

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

	resp, err := h.dockerClient.ContainerCreate(
		ctx,
		cont_config,
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		StartConfig.Name,
	)

	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}

	// 4. Start the container
	if err := h.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {

		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}

	wsMessage.Payload.Data = "Container started successfully with ID: " + resp.ID
	return wsMessage

}

// HandleListContainers retrieves and returns a list of Docker containers.
// The context handles cancellation and timeout, and the request may include optional filters.
// It returns a WsMessage containing the container list or an error message.
func (h *Handler) HandleListContainers(ctx context.Context, req lib.Command) lib.WsMessage {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})

	options := container.ListOptions{}
	if req.Params.ListContainersConfig != nil {
		options.All = req.Params.ListContainersConfig.All
		options.Size = req.Params.ListContainersConfig.Size
	}

	containers, err := h.dockerClient.ContainerList(ctx, options)
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}
	wsMessage.Type = lib.TypeResponse
	wsMessage.Payload.Data = containers
	return wsMessage

}

// HandleDeleteContainer removes a Docker container by its ID.
// The context handles cancellation and timeout, and the request contains the container ID and removal options.
// It returns a WsMessage with a success confirmation or an error message.
func (h *Handler) HandleDeleteContainer(ctx context.Context, req lib.Command) lib.WsMessage {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	configs := req.Params.DeleteContainerConfig

	options := container.RemoveOptions{}
	if configs != nil {
		options.Force = configs.Force
		options.RemoveVolumes = configs.Volume
		options.RemoveLinks = configs.Links
	} else {
		options.Force = true
	}
	err := h.dockerClient.ContainerRemove(ctx, req.ContainerID, options)
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}
	wsMessage.Payload.Data = "Container " + req.ContainerID + " deleted successfully"
	return wsMessage
}

// HandleStopContainer stops a running Docker container by its ID.
// The context handles cancellation and timeout, and a 5-second timeout is used for graceful shutdown.
// It returns a WsMessage with a success confirmation or an error message.
func (h *Handler) HandleStopContainer(ctx context.Context, req lib.Command) lib.WsMessage {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	fmt.Println("Processing STOP_CONTAINER command for ContainerID:", req.ContainerID)

	// Using a timeout of 5 seconds for stopping the container
	// overrides the default 10 seconds
	timeout := 5 // seconds
	err := h.dockerClient.ContainerStop(ctx, req.ContainerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}
	wsMessage.Payload.Data = "Container " + req.ContainerID + " stopped successfully"
	return wsMessage
}

// HandleRestartContainer restarts a Docker container by its ID.
// The context handles cancellation and timeout, and an immediate restart (timeout=0) is performed.
// It returns a WsMessage with a success confirmation or an error message.
func (h *Handler) HandleRestartContainer(ctx context.Context, req lib.Command) lib.WsMessage {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	log.Println("Processing RESTART_CONTAINER command for ContainerID:", req.ContainerID)
	timeout := 0 // immediate restart
	err := h.dockerClient.ContainerRestart(ctx, req.ContainerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}
	wsMessage.Payload.Data = "Container " + req.ContainerID + " restarted successfully"

	return wsMessage
}

// HandleStartContainer starts a stopped Docker container by its ID.
// The context handles cancellation and timeout, and it starts an existing container (not creating a new one).
// It returns a WsMessage with a success confirmation or an error message.
func (h *Handler) HandleStartContainer(ctx context.Context, req lib.Command) lib.WsMessage {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
	err := h.dockerClient.ContainerStart(ctx, req.ContainerID, container.StartOptions{})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		return wsMessage
	}
	wsMessage.Payload.Data = "Container " + req.ContainerID + " started successfully"
	return wsMessage
}
