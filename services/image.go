package services

import (
	"context"

	"github.com/docker/docker/client"
)

type ImageService struct {
	ctx context.Context
	cli *client.Client
}

func NewImageService(ctx context.Context, cli *client.Client) *ImageService {
	return &ImageService{
		ctx:  ctx,
		cli:  cli,
	}
}

func (I *ImageService) PullImage(imageName string){

}
func (I *ImageService) RemoveImage(imageName string){
}

func (I *ImageService) ListImages(){

}

func (I *ImageService) FindImage(imageName string){
}

