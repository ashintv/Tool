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
type MockDockerClient struct {
	ImagePullFn       func(ctx context.Context, ref string, opts image.PullOptions) (io.ReadCloser, error)
	ContainerCreateFn func(ctx context.Context, config *container.Config, hostConfig *container.HostConfig,
		networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string,
	) (container.CreateResponse, error)
	
	ContainerStartFn   func(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerListFn    func(ctx context.Context, options container.ListOptions) ([]container.Summary, error)
	ContainerRemoveFn  func(ctx context.Context, containerID string, options container.RemoveOptions) error
	ContainerStopFn    func(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRestartFn func(ctx context.Context, containerID string, options container.StopOptions) error
}

// ImagePull pretends to pull a Docker image successfully.
func (m *MockDockerClient) ImagePull(
	ctx context.Context, ref string, opts image.PullOptions,
) (io.ReadCloser, error) {
	if m.ImagePullFn == nil {
		return io.NopCloser(strings.NewReader("")), nil
	}
	return m.ImagePullFn(ctx, ref, opts)
}

func (m *MockDockerClient) ContainerCreate(
	ctx context.Context,
	config *container.Config,
	hostConfig *container.HostConfig,
	networkingConfig *network.NetworkingConfig,
	platform *ocispec.Platform,
	containerName string,
) (container.CreateResponse, error) {
	if m.ContainerCreateFn == nil {
		return container.CreateResponse{ID: "fake_container_id"}, nil
	}
	return m.ContainerCreateFn(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func (m *MockDockerClient) ContainerStart(
	ctx context.Context, containerID string, options container.StartOptions,
) error {
	if m.ContainerStartFn == nil {
		return nil
	}
	return m.ContainerStartFn(ctx, containerID, options)
}

func (m *MockDockerClient) ContainerList(
	ctx context.Context, options container.ListOptions,
) ([]container.Summary, error) {
	if m.ContainerListFn == nil {
		return nil, nil
	}
	return m.ContainerListFn(ctx, options)
}

func (m *MockDockerClient) ContainerRemove(
	ctx context.Context, containerID string, options container.RemoveOptions,
) error {
	if m.ContainerRemoveFn == nil {
		return nil
	}
	return m.ContainerRemoveFn(ctx, containerID, options)
}

func (m *MockDockerClient) ContainerStop(
	ctx context.Context, containerID string, options container.StopOptions,
) error {
	if m.ContainerStopFn == nil {
		return nil
	}
	return m.ContainerStopFn(ctx, containerID, options)
}

func (m *MockDockerClient) ContainerRestart(
	ctx context.Context, containerID string, options container.StopOptions,
) error {
	if m.ContainerRestartFn == nil {
		return nil
	}
	return m.ContainerRestartFn(ctx, containerID, options)
}

