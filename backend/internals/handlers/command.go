package handlers

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type CommandHandler struct {
	ws *WebsocketService
}

// general request structure for container commands like list containers etc
// request structure for containers for specific container commands like delete container etc
type RequestType struct {
	MachineID   string          `json:"machine_id" binding:"required"`
	ContainerId string          `json:"container_id"`
	Command     lib.CommandType `json:"command_type" binding:"required"`
	Params      lib.Params      `json:"params"` // parameters for starting container etc
}

type MachineResponse struct {
	HTTPResponse struct {
		Message     string          `json:"message"`
		ContainerId string          `json:"container_id"`
		MachineID   string          `json:"machine_id"`
		Error       error           `json:"error"`
		Data        lib.PayloadType `json:"data"`
	}
	Status int // status code need top passed
	Mu     sync.Mutex
}

func NewCommandHandler(ws *WebsocketService) *CommandHandler {
	return &CommandHandler{
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
func (h *CommandHandler) SendCommand(ctx *gin.Context) {
	var req RequestType
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//TODO: add task to database for command history feature

	ctxTime, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()
	ResponseChan := make(chan lib.PayloadType)
	ErrorChan := make(chan error)

	go func() {
		res, err := h.ws.WaitForResponse(req.MachineID, ctxTime)
		if err != nil {
			ErrorChan <- err
			return
		}
		ResponseChan <- res
	}()

	command := lib.GetCommand(req.ContainerId, req.MachineID, req.Command, false, req.Params)
	err := h.ws.SendCommandToMachine(command)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	select {
	case Res := <-ResponseChan:
		ctx.JSON(200, Res)
		return
	case err := <-ErrorChan:
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

// request handler to send command via websocket and keep the connection open to receive response
// get endpoint no body should be from parameters
func (h *CommandHandler) SendCommandWs(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	machineID := params.Get("machine_id")
	containerID := params.Get("container_id")
	commandType := params.Get("command_type")
	Params := params.Get("params")
	ParamsStruct := lib.Params{}
	if Params != "" {
		err := json.Unmarshal([]byte(Params), &ParamsStruct)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "invalid params"})
			return
		}
	}
	if machineID == "" || commandType == "" {
		ctx.JSON(400, gin.H{"error": "invalid  parameters"})
		return
	}

	command := lib.GetCommand(containerID, machineID, lib.CommandType(commandType), true, ParamsStruct)
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
