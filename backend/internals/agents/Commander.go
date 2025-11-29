package agents

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/docker/docker/api/types/image"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

type Commander struct {
	dockerClient *client.Client
	Config       CommanderConfig
}

func NewCommander(dockerClient *client.Client , Config *CommanderConfig) *Commander {
	return &Commander{
		dockerClient: dockerClient,
		Config:       *Config,
	}
}

func (Commader *Commander) Run(ctx context.Context) {
	u := url.URL{Scheme: "ws", Host: Commader.Config.WsServerHOST, Path: Commader.Config.Path}

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
			}else{
				wsMessage.Type = lib.TypeResponse
				wsMessage.Payload.Data = containers
				SendMessage(conn, wsMessage)
			}
			continue
		case lib.CREATE_CONTAINER:
			log.Println("Processing CREATE_CONTAINER command")
			// Add logic to create a new container
		case lib.DELETE_CONTAINER:
			log.Println("Processing DELETE_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to delete the specified container
		case lib.STOP_CONTAINER:
			log.Println("Processing STOP_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to stop the specified container
		case lib.START_CONTAINER:
			log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to start the specified container
			imageName := req.Params[0]
			wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})
			log.Println("Pulling image:", imageName)
			stream, err := Commader.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				wsMessage.Payload.Error = err.Error()
				SendMessage(conn, wsMessage)
				continue
			}
			if req.Stream {
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
			wsMessage.Payload.Data = "STREAM_END"
			wsMessage.Type = lib.TypeStreamEnd
			SendMessage(conn, wsMessage)
			} else {
				wsMessage.Type = lib.TypeResponse
				wsMessage.Payload.Data = "Image pulled successfully"
				SendMessage(conn, wsMessage)
			}
		case lib.RESTART_CONTAINER:
			log.Println("Processing RESTART_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to restart the specified container
		default:
			log.Println("Unknown command:", req.CMD)
		}

	}
}


func SendMessage(conn *websocket.Conn, message lib.WsMessage)  {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		log.Println("write error:", err)
	}
}

