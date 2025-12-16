package main

import (
	agents "aetrix/observer/internals/agent"
	"aetrix/observer/internals/lib"
	"context"
	"log"
	"time"

	"github.com/docker/docker/client"
)



func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	ws := agents.NewWSClient("0.0.0.0", "/agent", "agent-1")
	handler := agents.NewHandler(cli)
	err = ws.Connect()
	if err != nil {
		panic(err)
	}

	defer ws.Close()

	// Start resource monitoring in a separate goroutine
	go agents.StartResourceMonitor(ctx, ws, "agent-1", time.Second*3)

	// Listen for incoming messages and route them to the appropriate handler
	go ws.Receive(ctx, func(cmd lib.Command) {
		response := agents.MessageRouter(ctx, cmd, handler)
		if err := ws.Send(response); err != nil {
			log.Println("send error:", err)
		}
	})

}
