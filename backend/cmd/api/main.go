package main

import (
	"aetrix/observer/internals/db"
	"aetrix/observer/internals/handlers"
	"aetrix/observer/internals/middlewares"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	wsService := handlers.NewWebsocketService()
	CommandHandler := handlers.NewCommandHandler(wsService)
	database := db.InitializeDB()
	userHandler := handlers.NewUserHandler(database)
	
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("user/login", userHandler.Login)
			auth.POST("user/signup", userHandler.Signup)

		}
		user := api.Group("/user", middlewares.UserMiddleware())
		{
			user.GET("/profile", userHandler.GetUser)
			user.PUT("/update", userHandler.UpdateUserDetailas)
			user.PUT("/change-password", userHandler.ChangePassword)
			user.POST("/cmd", CommandHandler.SendCommand)
			user.GET("/cmd", CommandHandler.SendCommandWs)

		}

		agent := r.Group("/agent", middlewares.MachineMiddleware())
		{
			// WebSocket endpoint for machines to connect
			agent.GET("/:machine_id", wsService.Wss)
		}
	}
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
