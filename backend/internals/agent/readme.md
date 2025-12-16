# Agent Package

The `agent` package provides a Docker container management agent that communicates with a central server via WebSocket. It enables remote Docker operations, resource monitoring, and real-time command execution.

## Overview

This package implements an agent service that:
- Connects to a WebSocket server for bidirectional communication
- Executes Docker container operations remotely
- Monitors and reports system resource metrics
- Handles container lifecycle management (start, stop, restart, delete)

## Components

### Handler (`handler.go`)

The `Handler` manages Docker operations and provides methods for container lifecycle management.

#### Methods

- **`HandleStartNewContainer`** - Pulls a Docker image and creates/starts a new container with specified configuration
- **`HandleListContainers`** - Retrieves a list of Docker containers with optional filters
- **`HandleDeleteContainer`** - Removes a Docker container by ID with configurable options
- **`HandleStopContainer`** - Stops a running container with a 5-second graceful shutdown timeout
- **`HandleRestartContainer`** - Restarts a container immediately
- **`HandleStartContainer`** - Starts a stopped container

### Router (`router.go`)

The `MessageRouter` function routes incoming commands to the appropriate handler methods based on command type.

#### Supported Commands

- `LIST_CONTAINERS` - List all containers
- `DELETE_CONTAINER` - Remove a container
- `STOP_CONTAINER` - Stop a running container
- `START_NEW_CONTAINER` - Create and start a new container
- `RESTART_CONTAINER` - Restart a container
- `START_CONTAINER` - Start an existing stopped container

### Docker Client (`docker_client.go`)

Defines the `DockerClient` interface for Docker SDK interactions, enabling:
- Easy testing and mocking
- Abstraction of Docker API complexity
- Dependency injection for testability

#### Interface Methods

- `ImagePull` - Pull Docker images
- `ContainerCreate` - Create new containers
- `ContainerStart` - Start containers
- `ContainerList` - List containers
- `ContainerRemove` - Remove containers
- `ContainerStop` - Stop containers
- `ContainerRestart` - Restart containers

### WebSocket Client (`ws_client.go`)

The `WSClient` manages WebSocket connections for the agent.

#### Features

- Automatic connection management
- JSON message serialization/deserialization
- Continuous message receiving with context cancellation
- Graceful connection closure

#### Methods

- **`Connect()`** - Establishes WebSocket connection to the server
- **`Send(msg lib.WsMessage)`** - Sends a message over the WebSocket
- **`Receive(ctx, onMessage)`** - Continuously reads and processes incoming messages
- **`Close()`** - Gracefully closes the WebSocket connection

### Resource Monitor (`resource_client.go`)

The `StartResourceMonitor` function provides periodic system resource monitoring.

#### Features

- Configurable monitoring interval
- CPU and memory metrics reporting (placeholders for actual implementation)
- Context-based cancellation support
- Automatic metric streaming via WebSocket

## Usage

### Basic Setup

```go
package main

import (
    agents "aetrix/observer/internals/agent"
    "context"
    "time"
    "github.com/docker/docker/client"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Initialize Docker client
    cli, err := client.NewClientWithOpts(
        client.FromEnv,
        client.WithAPIVersionNegotiation(),
    )
    if err != nil {
        panic(err)
    }

    // Create WebSocket client
    ws := agents.NewWSClient("0.0.0.0:8080", "/agent", "agent-1")

    // Create handler
    handler := agents.NewHandler(cli)

    // Connect to server
    err = ws.Connect()
    if err != nil {
        panic(err)
    }
    defer ws.Close()

    // Start resource monitoring
    go agents.StartResourceMonitor(ctx, ws, "agent-1", time.Second*3)

    // Listen for incoming commands
    go ws.Receive(ctx, func(cmd lib.Command) {
        response := agents.MessageRouter(ctx, cmd, handler)
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
  "cmd": "LIST_CONTAINERS",
  "machine_id": "agent-1",
  "params": {
    "list_containers_config": {
      "all": true,
      "size": true
    }
  }
}
```

## Testing

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
                                           │ MessageRouter │
                                           └──────────────┘
                                                    │
                                                    ▼
                                           ┌──────────────┐
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
