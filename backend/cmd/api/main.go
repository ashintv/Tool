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
	MachineHandler := handlers.NewMachineHandler(database)

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

			user.GET("/machine/:machine_id", MachineHandler.GetMachine)
			user.GET("/machine/owned", MachineHandler.ListMachinesOfUser)
			user.GET("/machine/usable", MachineHandler.ListUsableMachine)
			user.POST("/machine", MachineHandler.RegisterMachine)
			user.PUT("/machine", MachineHandler.UpdateMachine)
			user.PUT("/machine/add/user", MachineHandler.AddUser)
			user.PUT("/machine/remove/user", MachineHandler.RemoveUser)
			user.DELETE("/machine/:machine_id", MachineHandler.DeleteMachine)

		}

		agent := r.Group("/agent", middlewares.MachineMiddleware())
		{
			agent.GET("/:machine_id", wsService.Wss)
		}
	}
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}
