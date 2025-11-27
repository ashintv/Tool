package main

import (
	"aetrix/observer/internals/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	wsService := handlers.NewWebsocketService()
	CommandHandler := handlers.NewCommandHandler(wsService)

	api := r.Group("/api")
	{
		command := api.Group("/command")
		{
			command.POST("/", CommandHandler.SendCommand)
			command.GET("/ws", CommandHandler.SendCommandWs)
		}
	}

	agent := r.Group("/agent-ws")
	{
		agent.GET("/:machine_id", wsService.Wss)
	}

	r.Run(":8080")
}
