package services

import (
	"context"

	"github.com/docker/docker/client"
)

type imageService struct {
	Name string
	ctx context.Context
	cli *client.Client
}

func NewImageService(name string, ctx context.Context, cli *client.Client) *imageService {
	return &imageService{
		Name: name,
		ctx:  ctx,
		cli:  cli,
	}
}

func (I *imageService) PullImage(){

}
func (I *imageService) RemoveImage(){
}

func ListImages(){

}

func FindImage(){

}

