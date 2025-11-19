package main

import (
	"aetrix/observer/services"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

const (
	PULL_IMAGE   = "pull"
	REMOVE_IMAGE = "Iremove"
	LIST_IMAGES  = "Ilist"
	FIND_IMAGE   = "Ifind"

	START_CONTAINER   = "start"
	STOP_CONTAINER    = "stop"
	REMOVE_CONTAINER  = "Cremove"
	LIST_CONTAINERS   = "Clist"
	FIND_CONTAINER    = "Cfind"
	MONITOR_CONTAINER = "monitor"
)

const example_usage string = `The following commands are available:
pull:<image-name>          - Pull a Docker image
Iremove:<image-name>      - Remove a Docker image
Ilist                     - List all Docker images
Ifind:<image-name>        - Find a specific Docker image \n
start:<container-name>    - Start a Docker container
stop:<container-name>     - Stop a Docker container
Cremove:<container-name>  - Remove a Docker container
Clist                     - List all Docker containers
Cfind:<container-name>    - Find a specific Docker container
monitor:<container-name>  - Monitor a Docker container
`

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Initialize Docker client
	Cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer Cli.Close()
	ctx := context.Background()
	imageService := services.NewImageService(ctx, Cli)
	containerService := services.NewContainerService("my-container", ctx, Cli)
	emailService := services.NewEmailService("smtp.example.com", 587, "user", "pass")

	// Use the services
	_ = imageService
	_ = containerService
	_ = emailService

	// implement a CliTool to interact with these services

	// when the cli starte we should check for config.json or create one with user inputs

	for {
		timt := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] Aetrix> ", timt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// Handle commands command look like  pull:image-name
		parts := strings.Split(input, ":")
		if len(parts) < 1 {
			fmt.Println("Invalid command")
			fmt.Print(example_usage)
			continue
		}
		switch parts[0] {
		case PULL_IMAGE:
			handlePullImage(imageService, parts[1])
			continue
		case REMOVE_IMAGE:
			handleRemoveImage(imageService, parts[1])
			continue
		case LIST_IMAGES:
			handleListImages(imageService)
			continue
		case FIND_IMAGE:
			handleFindImage(imageService, parts[1])
			continue
		case START_CONTAINER:
			handleStartContainer(containerService)
			continue
		case STOP_CONTAINER:
			// Implement stop container
			continue
		case REMOVE_CONTAINER:
			// Implement remove container
			continue
		case LIST_CONTAINERS:
			// Implement list containers
			continue
		case FIND_CONTAINER:
			// Implement find container
			continue
		case MONITOR_CONTAINER:
			// Implement monitor container
			continue
		default:
			fmt.Println("Unknown command")
			fmt.Print(example_usage)
			continue
		}
	}
}

func handlePullImage(imageService *services.ImageService, imageName string) {
	fmt.Printf("Pulling image: %s\n", imageName)
	err := imageService.PullImage(imageName)
	if err != nil {
		fmt.Printf("Error pulling image: %v\n", err)
		return
	}
	fmt.Println("Image pulled successfully")
}

func handleRemoveImage(imageService *services.ImageService, imageName string) {
	fmt.Printf("Removing image: %s\n", imageName)
	err := imageService.RemoveImage(imageName)
	if err != nil {
		fmt.Printf("Error removing image: %v\n", err)
	}
	
}

func handleListImages(imageService *services.ImageService) {
	fmt.Println("Listing images:")
	imageService.ListImages()
}

func handleFindImage(imageService *services.ImageService, imageName string) {
	fmt.Println("Finding image:")
	imageService.FindImage(imageName)
}

func handleStartContainer(containerService *services.ContainerService) {
	fmt.Println("Starting container:")
	containerService.StartContainer()
}
