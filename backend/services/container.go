package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type ContainerService struct {
	Name string
	cli  *client.Client
}

func NewContainerService(cli *client.Client) *ContainerService {
	return &ContainerService{
		cli: cli,
	}
}

func (c *ContainerService) ListContainers(ctx context.Context) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error listing containers", err)
		return
	}
	for i, cont := range containers {
		fmt.Printf("%d Container: %s - %s\n", i+1, cont.Names, cont.ID)
	}
}

func (c *ContainerService) StartContainer(ctx context.Context, imageName string) {
	resp, err := c.cli.ImagePull(ctx, imageName, image.PullOptions{})

	if err != nil {
		fmt.Println("Error pulling image", err)
		return
	}
	// log response , resp
	defer resp.Close()
	buff := make([]byte, 1024)
	for {
		n, err := resp.Read(buff)
		if err != nil {
			break
		}
		fmt.Println(string(buff[:n]))
	}

	// check if container with c.Name

	conts, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		fmt.Println("Error listing containers", err)
		return
	}
	for _, cont := range conts {
		for _, name := range cont.Names {
			if name == "/"+c.Name {
				fmt.Println("Container with name", c.Name, "already exists. Starting it.")
				err := c.cli.ContainerStart(ctx, cont.ID, container.StartOptions{})
				if err != nil {
					fmt.Println("Error starting existing container", err)
					return
				}
				fmt.Println("Container started successfully")
				return
			}
		}
	}
	// create container
	res, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, c.Name)
	if err != nil {
		fmt.Println("Error creating container", err)
		return
	}
    for{
		n, err := resp.Read(buff)
		if err != nil {
			break
		}
		fmt.Println(string(buff[:n]))
	}

	// start container
	err = c.cli.ContainerStart(ctx, res.ID, container.StartOptions{})
	if err != nil {
		fmt.Println("Error starting container", err)
		return
	}
	fmt.Println("Container started successfully with ID:", res.ID)
}

func (c *ContainerService) StopContainer(ctx context.Context) {

}

func (c *ContainerService) RemoveContainer(ctx context.Context) {

}

func (c *ContainerService) CreateContainer(ctx context.Context, imageName string) {
	res, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, c.Name)
	fmt.Println("Container created successfully with ID:", res.ID)
	if err != nil {
		fmt.Println("Error creating container", err)
		return
	}
}

func FindContainer(ctx context.Context) {
}
