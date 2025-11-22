package handlers

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)


type ContainerHandler struct {
	cli client.APIClient
}

type RequestStartNewContainer struct {
	ImageName  string `json:"image_name" binding:"required"`
	MachineID string `json:"machine_id" binding:"required"`
}

func NewContainerHandler(cli client.APIClient) *ContainerHandler {
	return &ContainerHandler{
		cli: cli,
	}
}


func (h *ContainerHandler) StartNewContainer(ctx *gin.Context) {
	var req RequestStartNewContainer
	if err := ctx.ShouldBindJSON(&req);  err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}


}
