package main

import (
	"aetrix/observer/services"
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
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

// command image_name or command -options
// example pull mongodb:latest
// example Ilist -2
const example_usage string = `The following commands are available:
pull <image-name>          - Pull a Docker image
Iremove <image-name>      - Remove a Docker image
Ilist -<level>            - List all Docker images
Ifind <image-name>        - Find a specific Docker image \n
start <container-name>    - Start a Docker container
stop <container-name>     - Stop a Docker container
Cremove <container-name>  - Remove a Docker container
Clist   -<level>            - List all Docker containers
Cfind <container-name>    - Find a specific Docker container
monitor <container-name>  - Monitor a Docker container
`

const (
	PLACEHOLDER = "<xyz>"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Initialize Docker client
	Cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer Cli.Close()
	imageService := services.NewImageService(Cli)
	containerService := services.NewContainerService(Cli)
	emailService := services.NewEmailService("smtp.example.com", 587, "user", "pass")

	// Use the services
	_ = imageService
	_ = containerService
	_ = emailService

	// implement a CliTool to interact with these services

	// when the cli starte we should check for config.json or create one with user inputs

	ctx := context.Background()
	for {
		timt := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] Aetrix> ", timt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		// Handle commands command look like  pull:image-name
		// TODO: repalce with a better command parser
		parts := strings.Split(input, " ")
		if len(parts) < 1 {
			// add empty input handling can ignore la
			fmt.Println("Invalid command")
			fmt.Print(example_usage)
			continue
		}

		// Ensure parts has at least 2 elements for commands without parameters error free handling
		if len(parts) < 2 {
			parts = append(parts, PLACEHOLDER)
		}
		command := parts[0]
		Options := parts[1]
		switch command {
		case PULL_IMAGE:
			handlePullImage(ctx, imageService, Options)
			continue
		case REMOVE_IMAGE:
			handleRemoveImage(ctx, imageService, Options)
			continue
		case LIST_IMAGES:
			handleListImages(ctx, imageService, Options)
			continue
		case FIND_IMAGE:
			handleFindImage(ctx, imageService, Options)
			continue
		case START_CONTAINER:
			handleStartContainer(ctx, containerService , Options)
			continue
		case STOP_CONTAINER:
			// Implement stop container
			continue
		case REMOVE_CONTAINER:
			// Implement remove container
			continue
		case LIST_CONTAINERS:
			handleListContainers(ctx , containerService)
			continue
			// Implement list containers
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

func handlePullImage(ctx context.Context, imageService *services.ImageService, imageName string) {
	if imageName == PLACEHOLDER {
		fmt.Println("Please provide an image name to pull.")
		return
	}
	fmt.Printf("Pulling image: %s\n", imageName)
	err := imageService.PullImage(ctx, imageName)
	if err != nil {
		fmt.Printf("Error pulling image: %v\n", err)
		return
	}
	fmt.Println("Image pulled successfully")
}

func handleRemoveImage(ctx context.Context, imageService *services.ImageService, imageName string) {
	if imageName == PLACEHOLDER {
		fmt.Println("Please provide an image name to pull.")
		return
	}
	fmt.Printf("Removing image: %s\n", imageName)
	err := imageService.RemoveImage(ctx, imageName)
	if err != nil {
		fmt.Printf("Error removing image: %v\n", err)
	}

}

// TODO: add addiions options to print (differtent l- list images , detailed list etc)
func handleListImages(ctx context.Context, imageService *services.ImageService, Options string) {
	fmt.Println("Listing images with options:", Options)
	if Options == PLACEHOLDER {
		Options = "-1"
	}
	fmt.Println("Listing images:")
	opt := Options[1:]
	level, err := strconv.Atoi(opt)
	if err != nil {
		fmt.Println("Invalid level option, defaulting to level 1")
		level = 1
	}
	imageService.ListImages(ctx, level)
}

func handleFindImage(ctx context.Context, 	imageService *services.ImageService, imageName string) {
	if imageName == PLACEHOLDER {
		fmt.Println("Please provide an image name to pull.")
		return
	}
	fmt.Println("Finding image:")
	imageService.FindImage(ctx, imageName)
}

func handleStartContainer(ctx context.Context, containerService *services.ContainerService , containerName string) {
	fmt.Println("Starting container:")
	containerService.StartContainer(ctx , containerName)
}


func handleListContainers(ctx context.Context, containerService *services.ContainerService) {
	fmt.Println("Listing containers:")
	containerService.ListContainers(ctx)

}

func handleFindContainer(ctx context.Context){}
