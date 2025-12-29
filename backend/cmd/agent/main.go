package main

import (
	"aetrix/observer/internals/agent/runtime"
	"aetrix/observer/internals/agent/websocket"
	"aetrix/observer/internals/services"
	"context"

	"github.com/docker/docker/client"
)

func main() {

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	dockerRuntime := runtime.NewDockerRuntime(cli)
	dispatcher := websocket.NewDispatcher(dockerRuntime)

	logger := services.GetLogger()
	wsCli := websocket.NewWSClient(
		"localhost:8080",
		"agent",
		"agent-001",
		logger,
		dispatcher.Dispatch,
	)

	err = wsCli.Connect()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect WebSocket client")
		return
	}

	wsCli.Start(ctx)

}
