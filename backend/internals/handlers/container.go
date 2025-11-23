package handlers

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ContainerHandler struct {
	ws *WebsocketService
}

type RequestStartNewContainer struct {
	ImageName string `json:"image_name" binding:"required"`
	MachineID string `json:"machine_id" binding:"required"`
}

type RequestListContainers struct {
	MachineID string `json:"machine_id" binding:"required"`
}

func NewContainerHandler(ws *WebsocketService) *ContainerHandler {
	return &ContainerHandler{
		ws: ws,
	}
}

func (h *ContainerHandler) ListContainers(ctx *gin.Context) {
	var req RequestListContainers
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//TODO: add task to database or in-memory store

	//start listening for response from machine
	ctxTime, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()

	response := make(map[string]interface{})
	var mu sync.Mutex
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	// start
	go func() {
		defer waitgroup.Done()
		res, err := h.ws.WaitforRespose(req.MachineID, ctxTime)

		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			response["error"] = err.Error()
			response["status"] = 504 //TODO: define error codes

			return
		}

		response["data"] = res
		response["status"] = 200
	}()
	// send command to machine
	err := h.ws.SendCommandToMachine(req.MachineID, req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// wait for response or timeout
	waitgroup.Wait()

	mu.Lock()
	defer mu.Unlock()

	if status, ok := response["status"]; ok {
		ctx.JSON(status.(int), response)
	} else {
		ctx.JSON(504, gin.H{
			"error":  "machine timeout or no response",
			"status": 504,
		})
	}

}
