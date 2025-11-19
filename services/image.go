package services

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type ImageService struct {
	ctx context.Context
	cli *client.Client
}

func NewImageService(ctx context.Context, cli *client.Client) *ImageService {
	return &ImageService{
		ctx: ctx,
		cli: cli,
	}
}

func (I *ImageService) PullImage(imageName string) error {
	res, err := I.cli.ImagePull(I.ctx, imageName, image.PullOptions{})
	if err != nil {
		fmt.Println("Error pulling image", err)
		return err
	}
	defer res.Close()
	io.Copy(io.Discard, res)
	return nil
}

func (I *ImageService) RemoveImage(imageName string) error {
	_, err := I.cli.ImageRemove(I.ctx, imageName, image.RemoveOptions{})
	if err != nil {
		fmt.Println("Error removing image", err)
		return err
	}
	return nil
}

func (I *ImageService) ListImages() error {
	images, err := I.cli.ImageList(I.ctx, image.ListOptions{})
	if err != nil {
		fmt.Println("Error listing images", err)
		return err
	}
	for i, img := range images {
		fmt.Println(i+1, " : Image ID:", img.ID)
	}
	return nil
}

func (I *ImageService) FindImage(imageName string) error {
	// Use ImageInspectWithRaw to find image details
	img, err := I.cli.ImageInspect(I.ctx, imageName)
	if err != nil {
		fmt.Println("Error finding image", err)
		return err
	}
	fmt.Println("Image found:", img.ID)
	return nil
}
