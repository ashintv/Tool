package main

import (
	"aetrix/observer/internals/handlers"
	"aetrix/observer/internals/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	wsService := handlers.NewWebsocketService()
	CommandHandler := handlers.NewCommandHandler(wsService)

	api := r.Group("/api", middlewares.UserMiddleware())
	{
		command := api.Group("/command")
		{
			command.POST("/", CommandHandler.SendCommand)
			command.GET("/ws", CommandHandler.SendCommandWs)
		}
	}

	agent := r.Group("/agent-ws", middlewares.MachineMiddleware())
	{
		agent.GET("/:machine_id", wsService.Wss)
	}

	r.Run(":8080")
}
