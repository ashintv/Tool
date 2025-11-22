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

func (c *ContainerService) StartContainerByContainerId(ctx context.Context, ContainerId string) error {
	err := c.cli.ContainerStart(ctx, ContainerId, container.StartOptions{})
	if err != nil {
		fmt.Println("Error starting container", err)
		return err
	}
	fmt.Println("Container started successfully with ID:", ContainerId)
	return nil
}

// start a container by image name
// if not image pull image ,
// start a new one
func (c *ContainerService) StartContainerByImage(ctx context.Context, imageName string) error {
	resp, err := c.cli.ImagePull(ctx, imageName, image.PullOptions{})

	if err != nil {
		fmt.Println("Error pulling image", err)
		return err
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
	// this part is commend becoz whever imagename give we have to start a new container not an old one
	// for running old one we have toCreate a new service call StartbyContainerId

	// conts, err := c.cli.ContainerList(ctx, container.ListOptions{All: true})
	// if err != nil {
	// 	fmt.Println("Error listing containers", err)
	// 	return
	// }
	// for _, cont := range conts {
	// 	for _, name := range cont.Names {
	// 		if name == "/"+c.Name {
	// 			fmt.Println("Container with name", c.Name, "already exists. Starting it.")
	// 			err := c.cli.ContainerStart(ctx, cont.ID, container.StartOptions{})
	// 			if err != nil {
	// 				fmt.Println("Error starting existing container", err)
	// 				return
	// 			}
	// 			fmt.Println("Container started successfully")
	// 			return
	// 		}
	// 	}
	// }
	// create container
	res, err := c.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, c.Name)
	if err != nil {
		fmt.Println("Error creating container", err)
		return err
	}
	for {
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
		return err
	}
	fmt.Println("Container started successfully with ID:", res.ID)
	return nil
}

// Pause a conatiner without removing the process
func (c *ContainerService) PauseContainer(ctx context.Context, ContainerId string) error {
	err := c.cli.ContainerPause(ctx, ContainerId)
	if err != nil {
		fmt.Println("failed to Pause Container - ", err)
		return err
	}
	return nil
}

// remove a container fully from the system
func (c *ContainerService) RemoveContainer(ctx context.Context, ContainerId string) error {
	err := c.cli.ContainerRemove(ctx, ContainerId, container.RemoveOptions{})
	if err != nil {
		fmt.Println("failed to remove Container - ", err)
		return err
	}
	return nil
}

// Rename a container using id
func (c *ContainerService) RenameContainer(ctx context.Context, ContainerId string, newName string) error {
	err := c.cli.ContainerRename(ctx, ContainerId, newName)
	if err != nil {
		fmt.Println("failed to rename Container - ", err)
		return err
	}
	return nil
}
