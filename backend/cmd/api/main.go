package main

import (
	"aetrix/observer/internals/db"
	"aetrix/observer/internals/handlers"
	"aetrix/observer/internals/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	wsService := handlers.NewWebsocketService()
	CommandHandler := handlers.NewCommandHandler(wsService)
	db := db.InitializeDB()
	_ = db // to avoid unused variable error
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
