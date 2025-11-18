package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"aetrix/observer/services"
)



func main(){
	fmt.Println("Hello, CLI!")
	// Initialize Docker client
	Cli , err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer Cli.Close()
	ctx := context.Background()
	imageService := services.NewImageService("my-image", ctx, Cli)
	containerService := services.NewContainerService("my-container", ctx, Cli)
	emailService := services.NewEmailService("smtp.example.com", 587, "user", "pass")

	// Use the services
	_ = imageService
	_ = containerService
	_ = emailService

	// implement a CliTool to interact with these services

	// when the cli starte we should check for config.json or create one with user inputs
	
}
