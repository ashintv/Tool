package runtime

import (
	"aetrix/observer/internals/agent/protocol"
	"aetrix/observer/internals/lib"
	"bytes"
	"io"

	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Test struct {
	name          string
	mockClient    *MockDockerClient
	req           lib.Command
	expectError   bool
	expectedCount int
}

// func TestHandleListContainers(t *testing.T) {

// 	// test cases with mock data
// 	tests := []Test{
// 		{
// 			name: "Successful listing of containers",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
// 					return []container.Summary{
// 						// returns more filled container summaries
// 						{ID: "container1"},
// 						{ID: "container2"},
// 					}, nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:       lib.LIST_CONTAINERS,
// 				MachineID: "test-machine",
// 			},
// 			expectError:   false,
// 			expectedCount: 2,
// 		},

// 		{
// 			name: "Error listing containers",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
// 					return nil, fmt.Errorf("Failed to list")
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:       lib.LIST_CONTAINERS,
// 				MachineID: "test-machine",
// 			},
// 			expectError:   true,
// 			expectedCount: 0,
// 		},

// 		{
// 			name: "Test Options with All:true",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
// 					if !options.All {
// 						return nil, fmt.Errorf(
// 							"docker.ContainerList called with All=false; expected All=true from command params",
// 						)
// 					}
// 					return []container.Summary{
// 						{ID: "container", State: container.StateExited},
// 						{ID: "container1", State: container.StateRunning},
// 						{ID: "container2", State: container.StateRunning},
// 					}, nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:       lib.LIST_CONTAINERS,
// 				MachineID: "test-machine",
// 				Params: lib.Params{
// 					ListContainersConfig: &lib.ListContainersConfig{
// 						All: true,
// 					},
// 				},
// 			},
// 			expectError:   false,
// 			expectedCount: 3,
// 		},

// 		{
// 			name: "Test Options with size:true",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerListFn: func(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
// 					if !options.Size {
// 						return nil, fmt.Errorf(
// 							"docker.ContainerList called with Size=false; expected Size=true from command params",
// 						)
// 					}
// 					return []container.Summary{
// 						{ID: "container", State: container.StateExited, SizeRootFs: 1233},
// 						{ID: "container1", State: container.StateRunning, SizeRootFs: 1233},
// 						{ID: "container2", State: container.StateRunning, SizeRootFs: 1233},
// 					}, nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:       lib.LIST_CONTAINERS,
// 				MachineID: "test-machine",
// 				Params: lib.Params{
// 					ListContainersConfig: &lib.ListContainersConfig{
// 						Size: true,
// 					},
// 				},
// 			},
// 			expectError:   false,
// 			expectedCount: 3,
// 		},
// 	}

// 	// run tests
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 			defer cancel()

// 			handler := NewHandler(tt.mockClient)
// 			resp := handler.HandleListContainers(ctx, tt.req)

// 			if tt.expectError {
// 				if resp.Payload.Error == nil {
// 					t.Fatal("Error expected got no error")
// 				}
// 				return
// 			}

// 			if !tt.expectError {
// 				if resp.Payload.Error != nil {
// 					t.Fatal("Error unexpected got no error")
// 				}
// 			}

// 			containers, ok := resp.Payload.Data.([]container.Summary)
// 			if !ok {
// 				t.Fatalf("payload data is not []container.Summary")
// 			}
// 			if len(containers) != tt.expectedCount {
// 				t.Fatalf("expected %d containers, got %d", tt.expectedCount, len(containers))
// 			}

// 		})

// 	}

// }

func TestHandleStartNewContainer(t *testing.T) {
	tests := []Test{
		{
			name:        "Default",
			mockClient:  &MockDockerClient{},
			expectError: false,
			req: lib.Command{
				MachineID: "test-machine",
				CMD:       lib.START_NEW_CONTAINER,
				Params: lib.Params{
					StartNewContainerConfig: &lib.StartNewContainerConfig{
						Image:         "test-image",
						ContainerPort: 8080,
						HostPort:      8080,
						Name:          "hello",
					},
				},
			},
		},
		{
			name: "Parameter testing",
			mockClient: &MockDockerClient{
				ContainerStartFn: func(
					ctx context.Context,
					containerID string,
					options container.StartOptions,
				) error {
					if containerID != "test_container" {
						return fmt.Errorf(
							`expected container id "test_container", got %s`,
							containerID,
						)
					}
					return nil
				},
				ImagePullFn: func(
					ctx context.Context,
					ref string,
					opts image.PullOptions,
				) (io.ReadCloser, error) {
					if ref != "test_image" {
						return nil, fmt.Errorf(
							`expected image "test_image", got %s`,
							ref,
						)
					}
					return io.NopCloser(bytes.NewBuffer(nil)), nil
				},
				ContainerCreateFn: func(
					ctx context.Context,
					config *container.Config,
					hostConfig *container.HostConfig,
					networkingConfig *network.NetworkingConfig,
					platform *v1.Platform,
					containerName string,
				) (container.CreateResponse, error) {

					if containerName != "test_container" {
						return container.CreateResponse{}, fmt.Errorf(
							`expected container name "test_container", got %s`,
							containerName,
						)
					}

					containerPort := nat.Port("8080/tcp")
					expected := nat.PortMap{
						containerPort: []nat.PortBinding{
							{
								HostIP:   "tcp",
								HostPort: "8080",
							},
						},
					}

					actual := hostConfig.PortBindings

					if actual[containerPort][0].HostIP != expected[containerPort][0].HostIP ||
						actual[containerPort][0].HostPort != expected[containerPort][0].HostPort {
						return container.CreateResponse{}, fmt.Errorf(
							"port bindings mismatch, expected %v, got %v",
							expected,
							actual,
						)
					}

					return container.CreateResponse{ID: "test_container"}, nil
				},
			},
			expectError: false,
			req: lib.Command{
				MachineID: "test-machine",
				CMD:       lib.START_NEW_CONTAINER,
				Params: lib.Params{
					StartNewContainerConfig: &lib.StartNewContainerConfig{
						Image:         "test_image",
						ContainerPort: 8080,
						HostPort:      8080,
						Protocol:      "tcp",
						HostIP:        "tcp",
						Name:          "test_container",
					},
				},
			},
		},
		{
			name: "Image pull failure",
			mockClient: &MockDockerClient{
				ImagePullFn: func(
					ctx context.Context,
					ref string,
					opts image.PullOptions,
				) (io.ReadCloser, error) {
					return nil, fmt.Errorf("failed to pull image")
				},
			},
			expectError: true,
			req: lib.Command{
				MachineID: "test-machine",
				CMD:       lib.START_NEW_CONTAINER,
				Params: lib.Params{
					StartNewContainerConfig: &lib.StartNewContainerConfig{
						Image:         "test-image",
						ContainerPort: 8080,
						HostPort:      8080,
						Name:          "hello",
					},
				},
			},
		},
		{
			name: "Container creation failure",
			mockClient: &MockDockerClient{
				ContainerCreateFn: func(
					ctx context.Context,
					config *container.Config,
					hostConfig *container.HostConfig,
					networkingConfig *network.NetworkingConfig,
					platform *v1.Platform,
					containerName string,
				) (container.CreateResponse, error) {
					return container.CreateResponse{}, fmt.Errorf("failed to create container")
				},
			},
			expectError: true,
			req: lib.Command{
				MachineID: "test-machine",
				CMD:       lib.START_NEW_CONTAINER,
				Params: lib.Params{
					StartNewContainerConfig: &lib.StartNewContainerConfig{
						Image:         "test-image",
						ContainerPort: 8080,
						HostPort:      8080,
						Name:          "hello",
					},
				},
			},
		},
		{
			name: "Container start failure",
			mockClient: &MockDockerClient{
				ContainerStartFn: func(
					ctx context.Context,
					containerID string,
					options container.StartOptions,
				) error {
					return fmt.Errorf("failed to start container")
				},
			},
			expectError: true,
			req: lib.Command{
				MachineID: "test-machine",
				CMD:       lib.START_NEW_CONTAINER,
				Params: lib.Params{
					StartNewContainerConfig: &lib.StartNewContainerConfig{
						Image:         "test-image",
						ContainerPort: 8080,
						HostPort:      8080,
						Name:          "hello",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			var err error
			var data interface{}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			respChan := make(chan protocol.Event, 8)

			runtime := NewDockerRuntime(tt.mockClient)

			//closing actually inside Dispatcher
			go func ()  {
				runtime.StartNewContainer(ctx, respChan, tt.req)
				defer close(respChan)
			}()

			for resp := range respChan {
				if resp.Error != nil {
					err = resp.Error
				}
				if resp.Data != nil {
					data = resp.Data
				}
			}

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if data == nil {
				t.Fatal("expected data, got nil")
			}
		})
	}
}

// func TestHandleDeleteContainer(t *testing.T) {
// 	tests := []Test{
// 		{
// 			name: "Successful deletion of container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerRemoveFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.RemoveOptions,
// 				) error {
// 					if containerID != "test_container" {
// 						return fmt.Errorf(`expected container Id to be: "test_container" , got %s `, containerID)
// 					}
// 					return nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.DELETE_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "Error deleting container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerRemoveFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.RemoveOptions,
// 				) error {
// 					return fmt.Errorf("Failed to delete")
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.DELETE_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: true,
// 		},
// 		{
// 			name: "Delete with Custom options",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerRemoveFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.RemoveOptions,
// 				) error {
// 					if containerID != "test_container" {
// 						return fmt.Errorf(`expected container Id to be: "test_container" , got %s `, containerID)
// 					}
// 					if !options.Force || !options.RemoveVolumes || !options.RemoveLinks {
// 						return fmt.Errorf(
// 							`expected default options Force:true, Volume:true, Links:true but got Force:%v, Volume:%v, Links:%v`,
// 							options.Force, options.RemoveVolumes, options.RemoveLinks,
// 						)
// 					}
// 					return nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.DELETE_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 				Params: lib.Params{
// 					DeleteContainerConfig: &lib.DeleteContainerConfig{
// 						Force:  true,
// 						Volume: true,
// 						Links:  true,
// 					},
// 				},
// 			},
// 			expectError: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 			defer cancel()

// 			handler := NewHandler(tt.mockClient)
// 			resp := handler.HandleDeleteContainer(ctx, tt.req)

// 			if tt.expectError {
// 				if resp.Payload.Error == nil {
// 					t.Fatal("Error expected got no error")
// 				}
// 				return
// 			}

// 			if !tt.expectError {
// 				if resp.Payload.Error != nil {
// 					t.Fatalf("Error unexpected got error: %s", resp.Payload.Error)
// 				}
// 				return
// 			}

// 			data, ok := resp.Payload.Data.(string)
// 			if !ok {
// 				t.Fatalf("payload data is not string")
// 				return
// 			}
// 			expectedMessage := "Container deleted successfully"
// 			if data != expectedMessage {
// 				t.Fatalf("expected payload data: %s, got: %s", expectedMessage, data)
// 			}

// 		})
// 	}
// }

// func TestHandleStopContainer(t *testing.T) {
// 	tests := []Test{
// 		{
// 			name: "Successful stopping of container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerStopFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.StopOptions,
// 				) error {
// 					if containerID != "test_container" {
// 						return fmt.Errorf(`expected container Id to be: "test_container" , got %s `, containerID)
// 					}
// 					return nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.STOP_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "Error stopping container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerStopFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.StopOptions,
// 				) error {
// 					return fmt.Errorf("Failed to stop")
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.STOP_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 			defer cancel()

// 			handler := NewHandler(tt.mockClient)
// 			resp := handler.HandleStopContainer(ctx, tt.req)

// 			if tt.expectError {
// 				if resp.Payload.Error == nil {
// 					t.Fatal("Error expected got no error")
// 				}
// 				return
// 			}

// 			if !tt.expectError {
// 				if resp.Payload.Error != nil {
// 					t.Fatalf("Error unexpected got error: %s", resp.Payload.Error)
// 				}
// 				return
// 			}

// 			data, ok := resp.Payload.Data.(string)
// 			if !ok {
// 				t.Fatalf("payload data is not string")
// 				return
// 			}
// 			expectedMessage := "Container stopped successfully"
// 			if data != expectedMessage {
// 				t.Fatalf("expected payload data: %s, got: %s", expectedMessage, data)
// 			}

// 		})
// 	}
// }

// func TestHandleRestartContainer(t *testing.T) {
// 	tests := []Test{
// 		{
// 			name: "Successful restarting of container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerRestartFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.StopOptions,
// 				) error {
// 					if containerID != "test_container" {
// 						return fmt.Errorf(`expected container Id to be: "test_container" , got %s `, containerID)
// 					}
// 					return nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.RESTART_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "Error restarting container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerRestartFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.StopOptions,
// 				) error {
// 					return fmt.Errorf("Failed to restart")
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.RESTART_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 			defer cancel()

// 			handler := NewHandler(tt.mockClient)
// 			resp := handler.HandleRestartContainer(ctx, tt.req)

// 			if tt.expectError {
// 				if resp.Payload.Error == nil {
// 					t.Fatal("Error expected got no error")
// 				}
// 				return
// 			}

// 			if !tt.expectError {
// 				if resp.Payload.Error != nil {
// 					t.Fatalf("Error unexpected got error: %s", resp.Payload.Error)
// 				}
// 				return
// 			}

// 			data, ok := resp.Payload.Data.(string)
// 			if !ok {
// 				t.Fatalf("payload data is not string")
// 				return
// 			}
// 			expectedMessage := "Container restarted successfully"
// 			if data != expectedMessage {
// 				t.Fatalf("expected payload data: %s, got: %s", expectedMessage, data)
// 			}

// 		})
// 	}
// }

// func TestHandleStartContainer(t *testing.T) {
// 	tests := []Test{
// 		{
// 			name: "Successful starting of container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerStartFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.StartOptions,
// 				) error {
// 					if containerID != "test_container" {
// 						return fmt.Errorf(`expected container Id to be: "test_container" , got %s `, containerID)
// 					}
// 					return nil
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.START_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: false,
// 		},
// 		{
// 			name: "Error starting container",
// 			mockClient: &docker.MockDockerClient{
// 				ContainerStartFn: func(
// 					ctx context.Context,
// 					containerID string,
// 					options container.StartOptions,
// 				) error {
// 					return fmt.Errorf("Failed to start")
// 				},
// 			},
// 			req: lib.Command{
// 				CMD:         lib.START_CONTAINER,
// 				MachineID:   "test-machine",
// 				ContainerID: "test_container",
// 			},
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 			defer cancel()

// 			handler := NewHandler(tt.mockClient)
// 			resp := handler.HandleStartContainer(ctx, tt.req)

// 			if tt.expectError {
// 				if resp.Payload.Error == nil {
// 					t.Fatal("Error expected got no error")
// 				}
// 				return
// 			}

// 			if !tt.expectError {
// 				if resp.Payload.Error != nil {
// 					t.Fatalf("Error unexpected got error: %s", resp.Payload.Error)
// 				}
// 				return
// 			}

// 			data, ok := resp.Payload.Data.(string)
// 			if !ok {
// 				t.Fatalf("payload data is not string")
// 				return
// 			}
// 			expectedMessage := "Container started successfully"
// 			if data != expectedMessage {
// 				t.Fatalf("expected payload data: %s, got: %s", expectedMessage, data)
// 			}

// 		})
// 	}
// }
