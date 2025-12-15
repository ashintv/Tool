package agents

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

type Commander struct {
	dockerClient *client.Client
	Config       CommanderConfig
}

// NewCommander creates a new Commander instance with the provided Docker client and configuration.
// Parameters:
//   - dockerClient: Docker client instance for container operations
//   - Config: Commander configuration containing connection and server settings
//
// Returns:
//   - *Commander: A new Commander instance
func NewCommander(dockerClient *client.Client, Config *CommanderConfig) *Commander {
	return &Commander{
		dockerClient: dockerClient,
		Config:       *Config,
	}
}

// Run starts the Commander's main execution loop, establishing a WebSocket connection
// and continuously listening for commands from the server.
// Parameters:
//   - ctx: Context for cancellation and timeout handling
func (Cmdr *Commander) Run(ctx context.Context) {
	Path := Cmdr.Config.Path + "/" + Cmdr.Config.Servername
	u := url.URL{Scheme: "ws", Host: Cmdr.Config.WsServerHOST, Path: Path}

	log.Println("Connecting to", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	// Send a message
	err = conn.WriteMessage(websocket.TextMessage, []byte("hello from client"))
	if err != nil {
		log.Println("write error:", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		log.Println("received:", string(msg))

		var req lib.Command
		err = json.Unmarshal(msg, &req)
		if err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		log.Printf("Received command: %+v\n", req)
		// Process the command here
		switch req.CMD {
		case lib.LIST_CONTAINERS:
			log.Println("Processing LIST_CONTAINERS command")
			Cmdr.HandleListContainers(ctx, req, conn)
			continue
		case lib.DELETE_CONTAINER:
			log.Println("Processing DELETE_CONTAINER command for ContainerID:", req.ContainerID)
			Cmdr.HandleDeleteContainer(ctx, req, conn)
		case lib.STOP_CONTAINER:
			log.Println("Processing STOP_CONTAINER command for ContainerID:", req.ContainerID)
			Cmdr.HandleStopContainer(ctx, req, conn)
			// Add logic to stop the specified container
		case lib.START_NEW_CONTAINER:
			log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
			if req.Stream {
				// Cmdr.HandleStartNewContainerStream(ctx, req, conn)
				continue
			}
			Cmdr.HandleStartNewContainer(ctx, req, conn)
			continue
		case lib.RESTART_CONTAINER:
			log.Println("Processing RESTART_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to restart the specified containee :
			continue
		case lib.START_CONTAINER:
			log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
			Cmdr.HandleStartContainer(ctx, req, conn)
			continue

		default:
			log.Println("Unknown command:", req.CMD)
		}

	}
}

// SendMessage sends a WebSocket message to the connected client.
// It marshals the message to JSON and writes it to the WebSocket connection.
// Parameters:
//   - conn: WebSocket connection to send the message through
//   - message: WsMessage to be sent
func SendMessage(conn *websocket.Conn, message lib.WsMessage) {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	log.Println("Sending message:", string(msgBytes), conn.RemoteAddr())
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		log.Println("write error:", err)
	}
}

// HandleStartNewContainerStream handles starting a new Docker container with streaming output.
// It pulls the specified image, creates and starts the container, sending progress updates
// through the WebSocket connection as a stream.
// TODO: all harder coded config should passes as params
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request containing container start parameters
//   - conn: WebSocket connection for streaming progress updates
// func (Cmdr *Commander) HandleStartNewContainerStream(ctx context.Context, req lib.Command, conn *websocket.Conn) {
// 	configs := req.Params.StartNewContainerConfig
// 	imageName := configs.Image
// 	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
// 	wsMessage.Payload.Data = "Starting container with image: " + imageName
// 	wsMessage.Type = lib.TypeStream
// 	SendMessage(conn, wsMessage)
// 	log.Println("Pulling image:", imageName)
// 	stream, err := Cmdr.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
// 	if err != nil {
// 		wsMessage.Payload.Error = err.Error()
// 		wsMessage.Type = lib.TypeStreamEnd
// 		SendMessage(conn, wsMessage)
// 		return
// 	}
// 	buffer := make([]byte, 1024)
// 	for {
// 		n, err := stream.Read(buffer)
// 		if err != nil {
// 			break
// 		}
// 		wsMessage.Type = lib.TypeStream
// 		wsMessage.Payload.Data = string(buffer[:n])
// 		SendMessage(conn, wsMessage)
// 	}
// 	wsMessage.Payload.Data = "Image pulled successfully"
// 	SendMessage(conn, wsMessage)

// 	ExposedPorts := nat.PortSet{}
// 	if configs.ExposedPorts != 0 {
// 		port := nat.Port(fmt.Sprintf("%d/%s", configs.ExposedPorts, configs.Protocol))
// 		ExposedPorts[port] = struct{}{}
// 	}
// 	config := &container.Config{
// 		Image:        imageName,
// 		ExposedPorts: ExposedPorts,
// 	}
// 	hostConfig := &container.HostConfig{
// 		PortBindings: nat.PortMap{
// 			nat.Port(configs.HostPort): []nat.PortBinding{
// 				{
// 					HostIP:   "0.0.0.0",
// 					HostPort: fmt.Sprintf("%d", configs.ExposedPorts),
// 				},
// 			},
// 		},
// 	}
// 	wsMessage.Payload.Data = "Config set successfully"
// 	SendMessage(conn, wsMessage)
// 	resp, err := Cmdr.dockerClient.ContainerCreate(
// 		ctx,
// 		config,
// 		hostConfig,
// 		&network.NetworkingConfig{},
// 		nil,
// 		configs.Name,
// 	)
// 	wsMessage.Payload.Data = "Container created successfully"
// 	SendMessage(conn, wsMessage)
// 	if err != nil {
// 		log.Println("Container creation error:", err)
// 		wsMessage.Payload.Error = err.Error()
// 		wsMessage.Type = lib.TypeStreamEnd
// 		SendMessage(conn, wsMessage)
// 		return
// 	}

// 	// 4. Start the container
// 	if err := Cmdr.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
// 		log.Println("Container start error:", err)
// 		wsMessage.Payload.Error = err.Error()
// 		wsMessage.Type = lib.TypeStreamEnd
// 		SendMessage(conn, wsMessage)
// 		return
// 	}

// 	wsMessage.Payload.Data = "Container started successfully with ID: " + resp.ID
// 	wsMessage.Type = lib.TypeStreamEnd
// 	SendMessage(conn, wsMessage)

// }

// HandleStartNewContainer handles starting a new Docker container without streaming.
// It pulls the specified image, creates and starts the container, sending a single
// response message upon completion.
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request containing container start parameters
//   - conn: WebSocket connection for sending response
func (Cmdr *Commander) HandleStartNewContainer(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	StartConfig := req.Params.StartNewContainerConfig
	imageName := StartConfig.Image
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	wsMessage.Type = lib.TypeResponse
	log.Println("Pulling image:", imageName)
	_, err := Cmdr.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		SendMessage(conn, wsMessage)
		return
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

	resp, err := Cmdr.dockerClient.ContainerCreate(
		ctx,
		cont_config,
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		StartConfig.Name,
	)

	if err != nil {
		log.Println("Container creation error:", err)
		wsMessage.Payload.Error = err.Error()
		wsMessage.Type = lib.TypeStreamEnd
		SendMessage(conn, wsMessage)
		return
	}

	// 4. Start the container
	if err := Cmdr.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Println("Container start error:", err)
		wsMessage.Payload.Error = err.Error()
		SendMessage(conn, wsMessage)
		return
	}

	wsMessage.Payload.Data = "Container started successfully with ID: " + resp.ID
	SendMessage(conn, wsMessage)

}

// HandleListContainers retrieves and returns a list of Docker containers.
// It can operate in streaming mode (sending containers one by one) or
// standard mode (sending all containers in a single response).
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request with optional streaming flag
//   - conn: WebSocket connection for sending container list
func (cmdr *Commander) HandleListContainers(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	log.Println("Processing LIST_CONTAINERS command")

	options := container.ListOptions{}
	if req.Params.ListContainersConfig != nil {
		options.All = req.Params.ListContainersConfig.All
		options.Size = req.Params.ListContainersConfig.Size
	}

	containers, err := cmdr.dockerClient.ContainerList(ctx, options)
	if err != nil {
		wsMessage.Payload.Error = err.Error()
	}
	if req.Stream {
		wsMessage.Type = lib.TypeStream
		for _, container := range containers {
			wsMessage.Payload.Data = container
			SendMessage(conn, wsMessage)
		}
		wsMessage.Payload.Data = "STREAM_END"
		wsMessage.Type = lib.TypeStreamEnd
		SendMessage(conn, wsMessage)
	} else {
		wsMessage.Type = lib.TypeResponse
		wsMessage.Payload.Data = containers
		SendMessage(conn, wsMessage)
	}
}

// HandleDeleteContainer removes a Docker container by its ID.
// It forcefully removes the container and sends a confirmation message
// through the WebSocket connection.
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request containing the container ID to delete
//   - conn: WebSocket connection for sending deletion response
func (cmdr *Commander) HandleDeleteContainer(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	configs := req.Params.DeleteContainerConfig
	log.Println("Processing DELETE_CONTAINER command for ContainerID:", req.ContainerID)
	options := container.RemoveOptions{}
	if configs != nil {
		options.Force = configs.Force
		options.RemoveVolumes = configs.Volume
		options.RemoveLinks = configs.Links
	} else {
		options.Force = true
	}
	err := cmdr.dockerClient.ContainerRemove(ctx, req.ContainerID, options)
	if err != nil {
		wsMessage.Payload.Error = err.Error()
	} else {
		wsMessage.Payload.Data = "Container " + req.ContainerID + " deleted successfully"
	}
	if req.Stream {
		wsMessage.Type = lib.TypeStream
	}
	SendMessage(conn, wsMessage)
}

// HandleStopContainer stops a running Docker container by its ID.
// It uses a 10-second timeout for graceful shutdown and sends a confirmation
// message through the WebSocket connection.
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request containing the container ID to stop
//   - conn: WebSocket connection for sending stop response
func (Cmdr *Commander) HandleStopContainer(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	fmt.Println("Processing STOP_CONTAINER command for ContainerID:", req.ContainerID)

	// Using a timeout of 5 seconds for stopping the container
	// overrides the default 10 seconds
	timeout := 5 // seconds
	err := Cmdr.dockerClient.ContainerStop(ctx, req.ContainerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
	} else {
		wsMessage.Payload.Data = "Container " + req.ContainerID + " stopped successfully"
	}
	if req.Stream {
		wsMessage.Type = lib.TypeStream
	}
	SendMessage(conn, wsMessage)
}

// HandleRestartContainer restarts a Docker container by its ID.
// It performs an immediate restart (timeout=0) and sends a confirmation
// message through the WebSocket connection.
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request containing the container ID to restart
//   - conn: WebSocket connection for sending restart response
func (Cmdr *Commander) HandleRestartContainer(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	log.Println("Processing RESTART_CONTAINER command for ContainerID:", req.ContainerID)
	timeout := 0 // immediate restart
	err := Cmdr.dockerClient.ContainerRestart(ctx, req.ContainerID, container.StopOptions{Timeout: &timeout})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
	} else {
		wsMessage.Payload.Data = "Container " + req.ContainerID + " restarted successfully"
	}
	if req.Stream {
		wsMessage.Type = lib.TypeStream
	}
	SendMessage(conn, wsMessage)
}

// HandleStartContainer starts a stopped Docker container by its ID.
// It starts an existing container (not creating a new one) and sends a confirmation
// message through the WebSocket connection.
// Parameters:
//   - ctx: Context for cancellation and timeout handling
//   - req: Command request containing the container ID to start
//   - conn: WebSocket connection for sending start response
func (Cmdr *Commander) HandleStartContainer(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
	err := Cmdr.dockerClient.ContainerStart(ctx, req.ContainerID, container.StartOptions{})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
	} else {
		wsMessage.Payload.Data = "Container " + req.ContainerID + " started successfully"
	}
	if req.Stream {
		wsMessage.Type = lib.TypeStream
	}
	SendMessage(conn, wsMessage)
}
