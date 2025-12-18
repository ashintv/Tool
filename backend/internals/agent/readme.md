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
├── docker/          # Docker container management
│   ├── client.go    # Docker client interface
│   ├── handler.go   # Container operation handlers
│   └── mock_client.go # Mock client for testing
├── websocket/       # WebSocket communication
│   └── client.go    # WebSocket client implementation
├── monitoring/      # System resource monitoring
│   └── resource.go  # Resource monitoring implementation
└── router.go        # Command routing logic
```

### Package Structure

```
┌─────────────┐         WebSocket         ┌──────────────────┐
│   Server    │ ◄─────────────────────────► │ websocket.Client │
└─────────────┘                            └──────────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │ MessageRouter │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │docker.Handler│
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │docker.Client │
                                           └──────────────┘
```

## Components

### Docker Package (`docker/`)

Handles all Docker-related operations including container lifecycle management.

#### `docker.Handler` (`handler.go`)

The `Handler` manages Docker operations and provides methods for container lifecycle management.

#### Methods

- **`HandleStartNewContainer`** - Pulls a Docker image and creates/starts a new container with specified configuration
- **`HandleListContainers`** - Retrieves a list of Docker containers with optional filters
- **`HandleDeleteContainer`** - Removes a Docker container by ID with configurable options
- **`HandleStopContainer`** - Stops a running container with a 5-second graceful shutdown timeout
- **`HandleRestartContainer`** - Restarts a container immediately
- **`HandleStartContainer`** - Starts a stopped container

#### `docker.Client` (`client.go`)

Defines the `Client` interface for Docker SDK interactions, enabling:
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

#### `docker.MockClient` (`mock_client.go`)

A `MockClient` is provided for testing handler methods without requiring an actual Docker daemon.

### WebSocket Package (`websocket/`)

Manages WebSocket connections and communication with the central server.

#### `websocket.Client` (`client.go`)

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

### Monitoring Package (`monitoring/`)

Provides system resource monitoring capabilities.

#### `monitoring.ResourceMonitor` (`resource.go`)

The `ResourceMonitor` provides periodic system resource monitoring.

**Features:**
- Configurable monitoring interval
- CPU and memory metrics reporting (placeholders for actual implementation)
- Context-based cancellation support
- Automatic metric streaming via WebSocket

**Methods:**
- **`NewResourceMonitor(ws, machineID, interval)`** - Creates a new monitor instance
- **`Start(ctx)`** - Begins the resource monitoring loop
```go
package main

import (
    "aetrix/observer/internals/agent"
    "aetrix/observer/internals/agent/docker"
    "aetrix/observer/internals/agent/monitoring"
    "aetrix/observer/internals/agent/websocket"
    "aetrix/observer/internals/lib"
    "context"
    "time"
    dockerclient "github.com/docker/docker/client"
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

    // Create handler
    handler := docker.NewHandler(cli)

    // Connect to server
    err = ws.Connect()
    if err != nil {
        panic(err)
    }
    defer ws.Close()

    // Start resource monitoring
    go monitoring.StartResourceMonitor(ctx, ws, "agent-1", time.Second*3)

    // Listen for incoming commands
    go ws.Receive(ctx,
        func(cmd lib.Command) {
            response := agent.MessageRouter(ctx, cmd, handler)
            ws.Send(response)
        },
        func(err error) {
            log.Println("receive error:", err)
        },
    )
    // Connect to server
    err = ws.Connect()
    if err != nil {
        panic(err)
    }
    defer ws.Close()

    // Start resource monitoring
    go agents.StartResourceMonitor(ctx, ws, "agent-1", time.Second*3)
}
```

## Usage Examples

### Starting a New ContaineressageRouter(ctx, cmd, handler)
        ws.Send(response)
    })

    // Wait for shutdown signal
    // ... (signal handling code)
}
```

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
go test ./internals/agent/docker
go test ./internals/agent/websocket
go test ./internals/agent/monitoring
```

### Mock Implementations

- **`docker.MockClient`** - Mock Docker client for testing without a Docker daemon
- Mock WebSocket connections for integration testing
The package includes comprehensive unit tests with mock implementations.

### Running Tests

```bash
go test ./internals/agent/...
```

### Mock Client (`docker_mock_client.go`)

A `MockDockerClient` is provided for testing handler methods without requiring an actual Docker daemon.

## Architecture

```
┌─────────────┐         WebSocket         ┌──────────────┐
│   Server    │ ◄─────────────────────────► │  WSClient    │
└─────────────┘                            └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
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
## Error Handling

All handler methods return `lib.WsMessage` with either:
- **Success**: `Payload.Data` contains the result
- **Error**: `Payload.Error` contains the error message

## Migration Guide

If you're upgrading from the previous flat structure:

### Old Import
```go
import "aetrix/observer/internals/agent"
ws := agent.NewWSClient(...)
handler := agent.NewHandler(cli)
agent.StartResourceMonitor(...)
```

### New Import
```go
import (
    "aetrix/observer/internals/agent"
    "aetrix/observer/internals/agent/docker"
    "aetrix/observer/internals/agent/websocket"
    "aetrix/observer/internals/agent/monitoring"
)
ws := websocket.NewClient(...)
handler := docker.NewHandler(cli)
monitoring.StartResourceMonitor(...)
// Router stays in agent package
agent.MessageRouter(...)
```
- Well-documented code
- Consistent naming conventions
- Clear package boundaries
                                           │   Handler    │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
                                           │ DockerClient │
                                           └──────────────┘
```

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
