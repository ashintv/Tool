package services

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type imageService struct {
	cli *client.Client
}

func NewImageService(cli *client.Client) *imageService {
	return &imageService{
		cli: cli,
	}
}

type Client interface {
	PullImage(ctx context.Context, imageName string) (bool, error)
	RemoveImage(ctx context.Context, imageName string) (bool, error)
	ListImages(ctx context.Context, level int) ([]image.Summary, error)
	FindImage(ctx context.Context, imageName string) (image.InspectResponse, error)
}

func (I *imageService) PullImage(ctx context.Context, imageName string)(bool, error) {
	res, err := I.cli.ImagePull(ctx, imageName, image.PullOptions{})

	if err != nil {
		fmt.Println("Error pulling image", err)
		return false, err
	}

	defer res.Close()

	buff := make([]byte, 1024)
	for {
		n, err := res.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading image pull response", err)
			return false, err
		}
		fmt.Print(string(buff[:n]))
	}

	return true, nil
}

func (I *imageService) RemoveImage(ctx context.Context, imageName string) (bool, error) {
	_, err := I.cli.ImageRemove(ctx, imageName, image.RemoveOptions{})
	if err != nil {
		fmt.Println("Error removing image", err)
		return false, err
	}
	return true, nil
}

func (I *imageService) ListImages(ctx context.Context, level int) ([]image.Summary ,error){
	images, err := I.cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		fmt.Println("Error listing images", err)
		return nil, err
	}
	return images, nil
}

func (I *imageService) FindImage(ctx context.Context, imageName string) (image.InspectResponse, error) {
	// Use ImageInspectWithRaw to find image details
	img, err := I.cli.ImageInspect(ctx, imageName)
	if err != nil {
		fmt.Println("Error finding image", err)
		return 	image.InspectResponse{}, err
	}
	fmt.Println("Image found:", img.ID)
	return img, nil
}
