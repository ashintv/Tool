package services

import (
	"context"

	"github.com/docker/docker/client"
)

type ContainerService struct {
	Name string
	cli  *client.Client
}

func NewContainerService(cli *client.Client) *ContainerService {
	return &ContainerService{
		cli:  cli,
	}
}

func (c *ContainerService) StartContainer(ctx context.Context , ) {

}

func (c *ContainerService) StopContainer(ctx context.Context) {

}

func (c *ContainerService) RemoveContainer(ctx context.Context) {

}

func ListContainers(ctx context.Context) {

}

func FindContainer(ctx context.Context) {
}


