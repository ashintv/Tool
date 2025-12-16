package main

import (
	"aetrix/observer/internals/config"
	"aetrix/observer/internals/db"
	"aetrix/observer/internals/handlers"
	"aetrix/observer/internals/middlewares"
	"aetrix/observer/internals/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	logger := services.GetLogger()
	cnf := config.LoadConfig(logger)
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // allow ALL origins
		},
		AllowMethods:     cnf.AllowMethods,
		AllowHeaders:     cnf.AllowHeaders,
		ExposeHeaders:    cnf.ExposeHeaders,
		AllowCredentials: cnf.AllowCredentials,
	}))
	wsService := handlers.NewWebsocketService()
	CommandHandler := handlers.NewCommandHandler(wsService)
	database := db.InitializeDB()
	userHandler := handlers.NewUserHandler(database, cnf)
	MachineHandler := handlers.NewMachineHandler(database)

	userMiddleware := middlewares.NewUserMiddleware(cnf)
	machineMiddleware := middlewares.NewMachineMiddleware(cnf)
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("user/login", userHandler.Login)
			auth.POST("user/signup", userHandler.Signup)

		}
		// user := api.Group("/user", userMiddleware.UserMiddleware())
		_ = userMiddleware
		user := api.Group("/user")
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
		_ = machineMiddleware
		agent := r.Group("/agent")
		{
			agent.GET("/:machine_id", wsService.Wss)
		}
	}
	err := r.Run(cnf.Port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
