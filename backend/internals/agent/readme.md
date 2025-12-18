# Agent Package

The `agent` package provides a Docker container management agent that communicates with a central server via WebSocket. It enables remote Docker operations, resource monitoring, and real-time command execution.

## Overview

This package implements an agent service that:
- Connects to a WebSocket server for bidirectional communication
- Executes Docker container operations remotely
- Monitors and reports system resource metrics
- Handles container lifecycle management (start, stop, restart, delete)

## Architecture

The agent package is organized into the following subpackages:

```
internals/agent/
├── docker/                 # Docker client interface and mocking
│   ├── docker_client.go    # DockerClient interface definition
│   └── docker_mock_client.go # Mock client for testing
├── handler/                # Container operation handlers
│   ├── handler.go          # Docker container lifecycle handlers
│   └── handler_test.go     # Comprehensive handler tests
├── monitor/                # System resource monitoring
│   └── resource.go         # Resource monitoring implementation
└── websocket/              # WebSocket communication
    ├── ws_client.go        # WebSocket client implementation
    ├── ws_client_test.go   # WebSocket client tests
    └── router.go           # Command routing logic
```

### Package Structure

```
┌─────────────┐         WebSocket         ┌──────────────────┐
│   Server    │ ◄─────────────────────────► │ websocket.Client │
└─────────────┘                            └──────────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │ Router       │
                                           │ (router.go)  │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │   Handler    │
                                           │ (handler.go) │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │DockerClient  │
                                           │ (interface)  │
                                           └──────────────┘
```

## Components

### Docker Package (`docker/`)

Defines the Docker client interface for easy testing and mocking.

#### `docker.DockerClient` (`docker_client.go`)

#### `docker.DockerClient` (`docker_client.go`)

Defines the `DockerClient` interface for Docker SDK interactions, enabling:
- Easy testing and mocking
- Abstraction of Docker API complexity
- Dependency injection for testability

**Interface Methods:**
- `ImagePull` - Pull Docker images
- `ContainerCreate` - Create new containers
- `ContainerStart` - Start containers
- `ContainerList` - List containers
- `ContainerRemove` - Remove containers
- `ContainerStop` - Stop containers
- `ContainerRestart` - Restart containers

#### `docker.MockDockerClient` (`docker_mock_client.go`)

A `MockDockerClient` is provided for testing handler methods without requiring an actual Docker daemon. Each method can be customized via function fields for precise test control.

### Handler Package (`handler/`)

Manages all Docker operations and provides methods for container lifecycle management.

#### `handler.Handler` (`handler.go`)

The `Handler` struct uses a `docker.DockerClient` to perform Docker operations.

#### Methods

- **`HandleStartNewContainer`** - Pulls a Docker image and creates/starts a new container with specified configuration
- **`HandleListContainers`** - Retrieves a list of Docker containers with optional filters
- **`HandleDeleteContainer`** - Removes a Docker container by ID with configurable options
- **`HandleStopContainer`** - Stops a running container with a 5-second graceful shutdown timeout
- **`HandleRestartContainer`** - Restarts a container immediately
- **`HandleStartContainer`** - Starts a stopped container

#### Tests (`handler_test.go`)

Comprehensive unit tests for all handler methods using `docker.MockDockerClient`.

### WebSocket Package (`websocket/`)

Manages WebSocket connections and communication with the central server.

#### `websocket.Client` (`ws_client.go`)

The `Client` manages WebSocket connections for the agent.

#### Features

- Automatic connection management
- JSON message serialization/deserialization
- Continuous message receiving with context cancellation
- Graceful connection closure
**Methods:**
- **`Connect()`** - Establishes WebSocket connection to the server
- **`Send(msg lib.WsMessage)`** - Sends a message over the WebSocket
- **`Receive(ctx, onMessage, onError)`** - Continuously reads and processes incoming messages
- **`Close()`** - Gracefully closes the WebSocket connection

#### Message Router (`router.go`)

Routes incoming commands to appropriate handler methods based on the command type.

#### Tests (`ws_client_test.go`)

Unit tests for WebSocket client functionality.

### Monitor Package (`monitor/`)

Provides system resource monitoring capabilities.

#### `monitor.ResourceMonitor` (`resource.go`)

The `ResourceMonitor` provides periodic system resource monitoring.

**Features:**
- Configurable monitoring interval
- CPU and memory metrics reporting (placeholders for actual implementation)
- Context-based cancellation support
- Automatic metric streaming via WebSocket

**Methods:**
- **`NewResourceMonitor(ws, machineID, interval)`** - Creates a new monitor instance
- **`Start(ctx)`** - Begins the resource monitoring loop

## Usage Example

```go
package main

import (
    "aetrix/observer/internals/agent/docker"
    "aetrix/observer/internals/agent/handler"
    "aetrix/observer/internals/agent/monitor"
    "aetrix/observer/internals/agent/websocket"
    "aetrix/observer/internals/lib"
    "context"
    "log"
    "time"
    dockerclient "github.com/docker/docker/client"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Initialize Docker client
    cli, err := dockerclient.NewClientWithOpts(
        dockerclient.FromEnv,
        dockerclient.WithAPIVersionNegotiation(),
    )
    if err != nil {
        panic(err)
    }

    // Create WebSocket client
    ws := websocket.NewClient("0.0.0.0:8080", "/agent", "agent-1")

    // Create handler with Docker client
    h := handler.NewHandler(cli)

    // Connect to server
    err = ws.Connect()
    if err != nil {
        panic(err)
    }
    defer ws.Close()

    // Start resource monitoring
    rm := monitor.NewResourceMonitor(ws, "agent-1", time.Second*3)
    go rm.Start(ctx)

    // Listen for incoming commands and route them
    go ws.Receive(ctx,
        func(cmd lib.Command) {
            response := websocket.RouteCommand(ctx, cmd, h)
            ws.Send(response)
        },
        func(err error) {
            log.Println("receive error:", err)
        },
    )

    // Wait for shutdown signal
    // ... (signal handling code)
}
```

## Usage Examples

### Starting a New Container

The agent receives a command via WebSocket with the following structure:

```json
{
  "cmd": "START_NEW_CONTAINER",
  "machine_id": "agent-1",
  "params": {
    "start_new_container_config": {
      "name": "my-nginx",
      "image": "nginx:latest",
      "host_port": 8080,
      "container_port": 80,
      "protocol": "tcp",
      "host_ip": "0.0.0.0"
    }
  }
}
```

### Listing Containers

```json
{
## Testing

The package includes comprehensive unit tests with mock implementations.

### Running Tests

```bash
# Test all packages
go test ./internals/agent/...

# Test specific package
go test ./internals/agent/handler
go test ./internals/agent/websocket
go test ./internals/agent/docker
```

### Mock Implementations

- **`docker.MockDockerClient`** - Mock Docker client for testing without a Docker daemon
- WebSocket client tests in `ws_client_test.go`
- Handler tests in `handler_test.go` with comprehensive coverage
## Dependencies

- `github.com/docker/docker/client` - Docker Engine API client
- `github.com/gorilla/websocket` - WebSocket implementation
- `aetrix/observer/internals/lib` - Internal message types and utilities

## Package Benefits

### Modularity
- Clear separation of concerns across subpackages
- Easy to test individual components
- Simple to extend with new features

### Testability
- Interface-based design for easy mocking
- Comprehensive test coverage
- Isolated unit tests for each component

### Maintainability
- Well-documented code
- Consistent naming conventions
- Clear package boundaries

## Package Structure Summary

| Package | Files | Purpose |
|---------|-------|---------|
| `docker` | `docker_client.go`, `docker_mock_client.go` | Docker client interface and mock implementation |
| `handler` | `handler.go`, `handler_test.go` | Container lifecycle management and tests |
| `monitor` | `resource.go` | System resource monitoring |
| `websocket` | `ws_client.go`, `ws_client_test.go`, `router.go` | WebSocket communication and command routing |

## Dependencies

- `github.com/docker/docker/client` - Docker Engine API client
- `github.com/gorilla/websocket` - WebSocket implementation
- `aetrix/observer/internals/lib` - Internal message types and utilities

## Configuration

The agent requires:
- WebSocket server host and port
- WebSocket endpoint path
- Unique agent/machine identifier
- Docker daemon access (typically via `/var/run/docker.sock`)

## Error Handling

All handler methods return `lib.WsMessage` with either:
- **Success**: `Payload.Data` contains the result
- **Error**: `Payload.Error` contains the error message

## Future Enhancements

- Actual CPU and memory metrics collection (currently placeholders)
- Container logs streaming
- Docker network management
- Volume management operations
- Docker Compose support

## License

Internal package for Aetrix Observer project.
