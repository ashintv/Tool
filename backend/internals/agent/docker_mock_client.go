package agents

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// MockDockerClient is a fake Docker client used for unit tests.
// It does NOT talk to real Docker.
// It simply returns fake values so your business logic can be tested.
type MockDockerClient struct{}

// ImagePull pretends to pull a Docker image successfully.
func (m *MockDockerClient) ImagePull(
	ctx context.Context,
	ref string,
	opts image.PullOptions,
) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("fake image pull")), nil
}

// ContainerCreate pretends to create a container and returns a fake container ID.
func (m *MockDockerClient) ContainerCreate(
	ctx context.Context,
	config *container.Config,
	hostConfig *container.HostConfig,
	networkingConfig *network.NetworkingConfig,
	platform *ocispec.Platform,
	containerName string,
) (container.CreateResponse, error) {
	return container.CreateResponse{
		ID: "fake_container_id",
	}, nil
}

// ContainerStart pretends to start a container successfully.
func (m *MockDockerClient) ContainerStart(
	ctx context.Context,
	containerID string,
	options container.StartOptions,
) error {
	return nil
}

// ContainerList pretends that one container exists.
func (m *MockDockerClient) ContainerList(
	ctx context.Context,
	options container.ListOptions,
) ([]container.Summary, error) {
	return []container.Summary{
		{
			ID:    "fake_container_id",
			Image: "redis",
		},
	}, nil
}

// ContainerRemove pretends to remove a container successfully.
func (m *MockDockerClient) ContainerRemove(
	ctx context.Context,
	containerID string,
	options container.RemoveOptions,
) error {
	return nil
}

// ContainerStop pretends to stop a container successfully.
func (m *MockDockerClient) ContainerStop(
	ctx context.Context,
	containerID string,
	options container.StopOptions,
) error {
	return nil
}

// ContainerRestart pretends to restart a container successfully.
func (m *MockDockerClient) ContainerRestart(
	ctx context.Context,
	containerID string,
	options container.StopOptions,
) error {
	return nil
}

// Compile-time check to ensure MockDockerClient implements DockerClient.
var _ DockerClient = (*MockDockerClient)(nil)
