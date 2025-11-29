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

func NewCommander(dockerClient *client.Client, Config *CommanderConfig) *Commander {
	return &Commander{
		dockerClient: dockerClient,
		Config:       *Config,
	}
}

func (Commader *Commander) Run(ctx context.Context) {
	Path := Commader.Config.Path + "/" + Commader.Config.Servername
	u := url.URL{Scheme: "ws", Host: Commader.Config.WsServerHOST, Path: Path}

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
			wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
			log.Println("Processing LIST_CONTAINERS command")
			containers, err := Commader.dockerClient.ContainerList(ctx, container.ListOptions{})
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
			continue
		case lib.DELETE_CONTAINER:
			log.Println("Processing DELETE_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to delete the specified container
		case lib.STOP_CONTAINER:
			log.Println("Processing STOP_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to stop the specified container
		case lib.START_CONTAINER:
			log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
			if req.Stream {
				Commader.HandleStartNewContainerStream(ctx, req, conn)
				continue
			}
			Commader.HandleStartNewContainer(ctx, req, conn)
		case lib.RESTART_CONTAINER:
			log.Println("Processing RESTART_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to restart the specified container
		default:
			log.Println("Unknown command:", req.CMD)
		}

	}
}

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

//TODO: all harder coded config should passes as params
func (Cmdr *Commander) HandleStartNewContainerStream(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	Params := req.Params
	imageName := Params.StartParams.Image
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	wsMessage.Payload.Data = "Starting container with image: " + imageName
	wsMessage.Type = lib.TypeStream
	SendMessage(conn, wsMessage)
	log.Println("Pulling image:", imageName)
	stream, err := Cmdr.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		wsMessage.Type = lib.TypeStreamEnd
		SendMessage(conn, wsMessage)
		return
	}
	buffer := make([]byte, 1024)
	for {
		n, err := stream.Read(buffer)
		if err != nil {
			break
		}
		wsMessage.Type = lib.TypeStream
		wsMessage.Payload.Data = string(buffer[:n])
		SendMessage(conn, wsMessage)
	}
	wsMessage.Payload.Data = "Image pulled successfully"
	SendMessage(conn, wsMessage)

	ExposedPorts := nat.PortSet{}
	if Params.StartParams.ExposedPorts != 0 {
		port := nat.Port(fmt.Sprintf("%d/%s", Params.StartParams.ExposedPorts, Params.StartParams.Protocol))
		ExposedPorts[port] = struct{}{}
	}
	config := &container.Config{
		Image: imageName,
		ExposedPorts: ExposedPorts,
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"6379/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", Params.StartParams.ExposedPorts),
				},
			},
		},
	}
	wsMessage.Payload.Data = "Config setted successfully"
	SendMessage(conn, wsMessage)
	resp, err := Cmdr.dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		"my-redis-container",
	)
	wsMessage.Payload.Data = "Container created successfully"
	SendMessage(conn, wsMessage)
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
		wsMessage.Type = lib.TypeStreamEnd
		SendMessage(conn, wsMessage)
		return
	}

	wsMessage.Payload.Data = "Container started successfully with ID: " + resp.ID
	wsMessage.Type = lib.TypeStreamEnd
	SendMessage(conn, wsMessage)

}


func (Cmdr *Commander) HandleStartNewContainer(ctx context.Context, req lib.Command, conn *websocket.Conn) {
	Params := req.Params
	imageName := Params.StartParams.Image
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
	wsMessage.Type = lib.TypeResponse
	log.Println("Pulling image:", imageName)
	_, err := Cmdr.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		wsMessage.Payload.Error = err.Error()
		SendMessage(conn, wsMessage)
		return
	}

	ExposedPorts := nat.PortSet{}
	if Params.StartParams.ExposedPorts != 0 {
		port := nat.Port(fmt.Sprintf("%d/%s", Params.StartParams.ExposedPorts, Params.StartParams.Protocol))
		ExposedPorts[port] = struct{}{}
	}

	config := &container.Config{
		Image: imageName,
		ExposedPorts: ExposedPorts,
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"6379/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", Params.StartParams.ExposedPorts),
				},
			},
		},
	}
	resp, err := Cmdr.dockerClient.ContainerCreate(
		ctx,
		config,
		hostConfig,
		&network.NetworkingConfig{},
		nil,
		Params.StartParams.Name,
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


