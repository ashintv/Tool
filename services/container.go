package services

import (
	"context"

	"github.com/docker/docker/client"
)

type ContainerService struct {
	Name string
	ctx  context.Context
	cli  *client.Client
}

func NewContainerService(name string, ctx context.Context, cli *client.Client) *ContainerService {
	return &ContainerService{
		Name: name,
		ctx:  ctx,
		cli:  cli,
	}
}

func (c *ContainerService) StartContainer() {
	
}

func (c *ContainerService) StopContainer() {

}

func (c *ContainerService) RemoveContainer() {

}

func ListContainers() {

}

func FindContainer() {

}


/*
MonitorContainer
Monitors the status and performance of a container

*/
func MonitorContainer() {

}
