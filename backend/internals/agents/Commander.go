package agents

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

type Commander struct {
	dockerClient *client.Client
}

func (Doc *Commander) Run(ctx context.Context) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/agent-ws"}

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
		var WsMessage lib.WsMessage
		WsMessage.MachineID = req.MachineID
		// Process the command here
		switch req.CMD {
		case lib.LIST_CONTAINERS:
			log.Println("Processing LIST_CONTAINERS command")
			containers, err := Doc.dockerClient.ContainerList(ctx, container.ListOptions{})
			if err != nil {
				
			} else {
				log.Printf("Containers: %+v\n", containers)
			}
		case lib.CREATE_CONTAINER:
			log.Println("Processing CREATE_CONTAINER command")
			// Add logic to create a container
		case lib.DELETE_CONTAINER:
			log.Println("Processing DELETE_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to delete the specified container
		case lib.STOP_CONTAINER:
			log.Println("Processing STOP_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to stop the specified container
		case lib.START_CONTAINER:
			log.Println("Processing START_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to start the specified container
		case lib.RESTART_CONTAINER:
			log.Println("Processing RESTART_CONTAINER command for ContainerID:", req.ContainerID)
			// Add logic to restart the specified container
		default:
			log.Println("Unknown command:", req.CMD)
		}

	}
}


