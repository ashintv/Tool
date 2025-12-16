package agents

import (
	"aetrix/observer/internals/lib"

	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
)

func TestHandleListContainers(t *testing.T) {

	// test cases with mock data
	tests := []struct {
		name          string
		mockClient    *MockDockerClient
		req           lib.Command
		expectError   bool
		expectedCount int
	}{
		{
			name: "Successful listing of containers",
			mockClient: &MockDockerClient{
				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
					return []container.Summary{
						// returns more filled container summaries
						{ID: "container1"},
						{ID: "container2"},
					}, nil
				},
			},
			req: lib.Command{
				CMD:       lib.LIST_CONTAINERS,
				MachineID: "test-machine",
			},
			expectError:   false,
			expectedCount: 2,
		},

		{
			name: "Error listing containers",
			mockClient: &MockDockerClient{
				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
					return nil, fmt.Errorf("Failed to list")
				},
			},
			req: lib.Command{
				CMD:       lib.LIST_CONTAINERS,
				MachineID: "test-machine",
			},
			expectError:   true,
			expectedCount: 0,
		},

		{
			name: "Test Options with All:true",
			mockClient: &MockDockerClient{
				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
					if !options.All {
						return nil, fmt.Errorf(
							"docker.ContainerList called with All=false; expected All=true from command params",
						)
					}
					return []container.Summary{
						{ID: "container", State: container.StateExited},
						{ID: "container1", State: container.StateRunning},
						{ID: "container2", State: container.StateRunning},
					}, nil
				},
			},
			req: lib.Command{
				CMD:       lib.LIST_CONTAINERS,
				MachineID: "test-machine",
				Params: lib.Params{
					ListContainersConfig: &lib.ListContainersConfig{
						All: true,
					},
				},
			},
			expectError:   false,
			expectedCount: 3,
		},

		{
			name: "Test Options with size:true",
			mockClient: &MockDockerClient{
				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
					if !options.Size {
						return nil, fmt.Errorf(
							"docker.ContainerList called with Size=false; expected Size=true from command params",
						)
					}
					return []container.Summary{
						{ID: "container", State: container.StateExited, SizeRootFs: 1233},
						{ID: "container1", State: container.StateRunning, SizeRootFs: 1233},
						{ID: "container2", State: container.StateRunning, SizeRootFs: 1233},
					}, nil
				},
			},
			req: lib.Command{
				CMD:       lib.LIST_CONTAINERS,
				MachineID: "test-machine",
				Params: lib.Params{
					ListContainersConfig: &lib.ListContainersConfig{
						Size: true,
					},
				},
			},
			expectError:   false,
			expectedCount: 3,
		},
	}

	// run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			handler := NewHandler(tt.mockClient)
			resp := handler.HandleListContainers(ctx, tt.req)

			if tt.expectError {
				if resp.Payload.Error == "" {
					t.Fatal("Error expected got no error")
				}
				return
			}

			if !tt.expectError {
				if resp.Payload.Error != "" {
					t.Fatal("Error unexpected got no error")
				}
			}

			containers, ok := resp.Payload.Data.([]container.Summary)
			if !ok {
				t.Fatalf("payload data is not []container.Summary")
			}
			if len(containers) != tt.expectedCount {
				t.Fatalf("expected %d containers, got %d", tt.expectedCount, len(containers))
			}

		})

	}

}

func TestHandleStartNewContainer(t *testing.T) {

}

func TestHandleDeleteContainer(t *testing.T) {}

func TestHandleStopContainer(t *testing.T) {}

func TestHandleRestartContainer(t *testing.T) {}

func TestHandleStartContainer(t *testing.T) {}
