package handlers

import (
	"aetrix/observer/lib"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ContainerHandler struct {
	ws *WebsocketService
}

// genral request structure for container commands like list containers etc
// request structure for containers for specific container commands like delete container etc
type RequestType struct {
	MachineID   string          `json:"machine_id" binding:"required"`
	ContainerId string          `json:"container_id"`
	Command     lib.CommandType `json:"command_type" binding:"required"`
}

type MachineReponse struct {
	HTTPResponse struct {
		Message     string `json:"message"`
		ContainerId string `json:"container_id"`
		MachineID   string `json:"machine_id"`
		Error       error  `json:"error"`
		Data        string `json:"data"`
	}
	Status int // status code need top passed
	Mu     sync.Mutex
}

func NewContainerHandler(ws *WebsocketService) *ContainerHandler {
	return &ContainerHandler{
		ws: ws,
	}
}

var upgrader_2 = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins (or customize)
	},
}

// request handler to send command via websocket and wait for response
// post endpoint to send command to machine and wait for response
func (h *ContainerHandler) SendCommand(ctx *gin.Context) {
	var req RequestType
	var res MachineReponse
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//TODO: add task to database or in-memory store
	//start listening for response from machine
	ctxTime, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()
	var mu sync.Mutex
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	// start
	go func() {
		defer waitgroup.Done()
		Res_, err := h.ws.WaitforRespose(req.MachineID, ctxTime)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			res.HTTPResponse.Error = err
			res.Status = 500 // TODO: replace with proper one
			res.HTTPResponse.Message = "agent Failed"
			return
		}

		res.HTTPResponse.Data = Res_
		res.Status = 200
		res.HTTPResponse.Message = "Containers fetched successfully"
		res.HTTPResponse.MachineID = req.MachineID
	}()
	// send command to machine
	command := lib.GetCommand(req.ContainerId, req.MachineID, req.Command)
	err := h.ws.SendCommandToMachine(command)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// wait for response or timeout
	waitgroup.Wait()
	mu.Lock()
	defer mu.Unlock()
	ctx.JSON(res.Status, res.HTTPResponse)
}

// request handler to send command via websocket and keep the connection open to receive response
// get endpoint no body should be from parameters
func (h *ContainerHandler) SendCommandWs(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	machineID := params.Get("machine_id")
	containerID := params.Get("container_id")
	commandType := params.Get("command_type")
	if machineID == "" || commandType == "" {
		ctx.JSON(400, gin.H{"error": "invalid  parameters"})
		return
	}

	command := lib.GetCommand(containerID, machineID, lib.CommandType(commandType))
	conn, err := upgrader_2.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to upgrade to websocket"})
		return
	}
	defer conn.Close()
	var wg sync.WaitGroup
	//TODO: should  pass a timout
	// Add subscriber
	wg.Add(1)
	h.ws.AddSubscriber(machineID, conn, "username", &wg) // Replace "username" with actual username
	h.ws.SendCommandToMachine(command)
	wg.Wait()
	h.ws.Unsubscribe(machineID, "username") // Replace "username" with actual username
}
