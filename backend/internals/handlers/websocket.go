package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Subscriber struct {
	username string
	conn *websocket.Conn
	wt *sync.WaitGroup
}
type WebsocketService struct{
	mu sync.Mutex
	cli *client.Client // docker client
	Machines map[string]*websocket.Conn // storing machine info for connections
	pendingResponseChannels map[string]chan string
	Subscribers map[string][]Subscriber // map of machineID to list of subscriber channels
}


func NewWebsocketService(cli *client.Client) *WebsocketService {
	return &WebsocketService{
		cli: cli,
		Machines: make(map[string]*websocket.Conn),
		pendingResponseChannels: make(map[string]chan string),
		Subscribers: make(map[string][]Subscriber),
	}
}
type MessageType string
const (
	TypeRegister MessageType = "register"
	TypeResponse MessageType = "response"
	TypeEvent    MessageType = "event"
	TypeAck      MessageType = "ack"
)

type WsMessage struct {
	MachineID string `json:"machine_id" binding:"required"`
	Type    MessageType   `json:"type" binding:"required"`
	Data   string        `json:"data" binding:"required"`
}
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // allow all origins (or customize)
    },
}

func (s *WebsocketService) Wss(ctx *gin.Context) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        fmt.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }

		// Handle incoming message
        var MsgJson WsMessage
		if err:= json.Unmarshal(msg , &MsgJson); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}
		fmt.Printf("Received message from %s: %s\n", conn.RemoteAddr().String(), MsgJson)

		// Process based on message type
		switch MessageType(MsgJson.Type) {
		case TypeRegister:
			s.HanldeRegisterRequest(MsgJson, conn)
		case TypeEvent:
			s.HandleEventMessage(MsgJson , conn)
		default:
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Unknown message type: %s", MsgJson.Type)))
		}
        // echo back

    }
}

func (s *WebsocketService) HanldeRegisterRequest(req WsMessage, conn *websocket.Conn) error {
	s.Machines[req.MachineID] = conn
	fmt.Printf("Machine %s registered.\n", req.MachineID)
	return nil
}

func (s *WebsocketService) WaitforRespose(machineID string , ctx context.Context) (string, error) {
	// clode if timout
	responseChan, exists := s.pendingResponseChannels[machineID]
	if exists {
		return "", fmt.Errorf("waiting for Existing command  %s", machineID)
	}
	responseChan = make(chan string)
	s.pendingResponseChannels[machineID] = responseChan

	defer func(){
		close(responseChan)
		delete(s.pendingResponseChannels, machineID) // Clean up
	}()

	// Wait for the response (this will block)
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("timeout waiting for response from machine %s", machineID)
	case response := <-responseChan:
		fmt.Printf("Received response for machine %s: %s\n", machineID, response)
		return response, nil
	}
}

func (s *WebsocketService) SendCommandToMachine(machineID string, command interface{}) error {
	stringMessage, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %v", err)
	}
	conn, exists := s.Machines[machineID]
	if !exists {
		return fmt.Errorf("machine %s not connected", machineID)
	}
	conn.WriteMessage(websocket.TextMessage, []byte(stringMessage))
	return nil
}


func (s *WebsocketService) HandleEventMessage(msg WsMessage, conn *websocket.Conn) error {
	// Broadcast to all subscribers
	// replace with notify function

	// save event to db or log
	// retry mechanism
	// notify subscribers

	fmt.Printf("Event from machine %s: %s\n", msg.MachineID, msg.Data)
	return nil
}


func (s *WebsocketService) AddSubscriber(machineID string, conn *websocket.Conn , username string , wt *sync.WaitGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Subscribers[machineID] = append(s.Subscribers[machineID] , Subscriber{
		username: username,
		conn: conn,
		wt: wt,
	})
	return nil
}

func (s* WebsocketService) NotifySubscribers(machineID string, event string) error {
	s.mu.Lock()
	subscribers, exists := s.Subscribers[machineID]
	s.mu.Unlock()
	if !exists {
		return fmt.Errorf("no subscribers for machine %s", machineID)
	}
	for _, sub := range subscribers {
		err := sub.conn.WriteMessage(websocket.TextMessage, []byte(event))
		if err != nil {
			fmt.Printf("Error notifying subscriber %s: %v\n", sub.username, err)
			continue
		}
		sub.wt.Done()
	}
	return nil
}

func (s *WebsocketService) Unsubscribe(machineID string,  username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	subscribers, exists := s.Subscribers[machineID]
	if !exists {
		return fmt.Errorf("no subscribers for machine %s", machineID)
	}
	for i, sub := range subscribers {
		if username== sub.username {
			s.Subscribers[machineID] = append(subscribers[:i], subscribers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("subscriber %s not found for machine %s", username, machineID)
}


func SendAndWait(){}
