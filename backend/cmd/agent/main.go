package main

import (
	"aetrix/observer/internals/agent/websocket"
	"aetrix/observer/internals/services"
	"context"
)

func main() {

	ctx := context.Background()
	dispatcher := websocket.NewDispatcher()
	logger  := services.GetLogger()
	client := websocket.NewWSClient(
		"localhost:8080",
		"agent",
		"agent-001",
		logger,
		dispatcher.Dispatch,
	)

	err := client.Connect()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect WebSocket client")
		return
	}

	client.Start(ctx)


}
