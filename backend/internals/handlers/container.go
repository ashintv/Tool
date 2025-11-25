package handlers

import (
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

type RequestStartNewContainer struct {
	ImageName string `json:"image_name" binding:"required"`
	MachineID string `json:"machine_id" binding:"required"`
}

// common for delete , start , restart etc
type RequestConatiner struct {
	ContainerID string `json:"container_id" binding:"required"`
	MachineID string `json:"machine_id" binding:"required"`
}
type RequestListContainers struct {
	MachineID string `json:"machine_id" binding:"required"`

}

type Reponse struct{
	ContainerId string `json:"container_id"`
	MachineID string `json:"machine_id"`
	Error error `json:"error"`
	Data string `json:"data"`
}
func NewContainerHandler(ws *WebsocketService) *ContainerHandler {
	return &ContainerHandler{
		ws: ws,
	}
}
const (
	LIST_CONTAINER = "listContainers"

)
var upgrader_2 = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins (or customize)
    },
}

func (h *ContainerHandler) LstContainers(ctx *gin.Context) {
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

	err := h.ws.SendCommandToMachine(req.MachineID, LIST_CONTAINER)
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


func (h *ContainerHandler) DeleteContainer(ctx *gin.Context){
	var res Reponse
	var req RequestConatiner
	res.ContainerId = req.ContainerID
	res.MachineID = req.MachineID

	if err:= ctx.ShouldBindJSON(&req);err!=nil{
		res.Error = err
		ctx.JSON(404 , res)
	}


}



func (h *ContainerHandler) StartNewContainer(ctx *gin.Context){
	var req RequestStartNewContainer
	if err:= ctx.ShouldBindJSON(&req);err!=nil{
		ctx.JSON(404 , gin.H{"error": err.Error()})
		return
	}
	conn, err := upgrader_2.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to upgrade to websocket"})
		return
	}
	defer conn.Close()
	var wg sync.WaitGroup
	// Add subscriber
	wg.Add(1)
	h.ws.AddSubscriber(req.MachineID, conn, "username" , &wg) // Replace "username" with actual username
	h.ws.SendCommandToMachine(req.MachineID ,req)
	wg.Wait()
}

func (h *ContainerHandler) RestartContainer(ctx *gin.Context){
}

func (h *ContainerHandler) StartContainer(ctx *gin.Context){
}

func (h *ContainerHandler) PauseContainer(ctx *gin.Context){
}
