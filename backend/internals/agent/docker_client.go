package agent

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// DockerClient defines the methods our application uses from the Docker SDK.
// This interface allows for easier testing and mocking of Docker interactions.
type DockerClient interface {
	ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)

	ContainerCreate(
		ctx context.Context,
		config *container.Config,
		hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig,
		platform *specs.Platform,
		containerName string,
	) (container.CreateResponse, error)

	ContainerStart(
		ctx context.Context,
		containerID string,
		options container.StartOptions,
	) error

	ContainerList(
		ctx context.Context,
		options container.ListOptions,
	) ([]container.Summary, error)

	ContainerRemove(
		ctx context.Context,
		containerID string,
		options container.RemoveOptions,
	) error

	ContainerStop(
		ctx context.Context,
		containerID string,
		options container.StopOptions,
	) error

	ContainerRestart(
		ctx context.Context,
		containerID string,
		options container.StopOptions,
	) error
}
