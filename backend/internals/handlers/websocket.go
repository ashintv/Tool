package handlers

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	// "github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Subscriber struct {
	username string
	conn     *websocket.Conn
	wt       *sync.WaitGroup
}
type WebsocketService struct {
	mu                      sync.Mutex                 // docker client
	Machines                map[string]*websocket.Conn // storing machine info for connections
	pendingResponseChannels map[string]chan string
	Subscribers             map[string][]Subscriber // map of machineID to list of subscriber channels
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		Machines:                make(map[string]*websocket.Conn),
		pendingResponseChannels: make(map[string]chan string),
		Subscribers:             make(map[string][]Subscriber),
	}
}

type MessageType string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins (or customize)
	},
}

func (s *WebsocketService) Wss(ctx *gin.Context) {
	// Upgrade initial GET request to a websocket
	machineID := ctx.Param("machine_id")

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	s.Machines[machineID] = conn
	fmt.Printf("Machine %s registered.\n", machineID)

	defer conn.Close()
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var message lib.WsMessage
		res := json.Unmarshal(msg, &message)
		if res != nil {
			fmt.Printf("Invalid message from machine %s: %v\n", machineID, res)
			continue
		}
		// Check if there's a pending response channel for this machine
		switch message.Type {
		case lib.TypeResponse:
			s.mu.Lock()
			responseChan, exists := s.pendingResponseChannels[message.MachineID]
			s.mu.Unlock()
			if exists {
				responseChan <- message.Payload
			}
		case lib.TypeEvent:
			// Notify all subscribers about the event
			fmt.Printf("Event from machine %s: %s\n", message.MachineID, message.Payload)
			// send mail 
			// send restart request
			// handle other event types
		case lib.TypeStream:
			// Notify all subscribers about the stream data
			err := s.SendToSubscribers(message.MachineID, message.Payload)
			if err != nil {
				fmt.Printf("Error notifying subscribers for machine %s: %v\n", message.MachineID, err)
			}
		default:
			fmt.Printf("Unknown message type from machine %s: %s\n", message.MachineID, message.Type)
		}
	}
}

func (s *WebsocketService) WaitforRespose(machineID string, ctx context.Context) (string, error) {
	// clode if timout
	responseChan, exists := s.pendingResponseChannels[machineID]
	if exists {
		return "", fmt.Errorf("waiting for Existing command  %s", machineID)
	}
	responseChan = make(chan string)
	s.pendingResponseChannels[machineID] = responseChan

	defer func() {
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

func (s *WebsocketService) SendCommandToMachine(command lib.Command) error {
	stringMessage, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %v", err)
	}
	conn, exists := s.Machines[command.MachineID]
	if !exists {
		return fmt.Errorf("machine %s not connected", command.MachineID)
	}
	conn.WriteMessage(websocket.TextMessage, []byte(stringMessage))
	return nil
}

func (s *WebsocketService) AddSubscriber(machineID string, conn *websocket.Conn, username string, wt *sync.WaitGroup) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Subscribers[machineID] = append(s.Subscribers[machineID], Subscriber{
		username: username,
		conn:     conn,
		wt:       wt,
	})

	go func() {
		time.Sleep(30 * time.Second) // Example timeout duration
		wt.Done()
	}()
}

func (s *WebsocketService) SendToSubscribers(machineID string, event string) error {
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
		// TODO: only if stream end message (DONE) is received
		if event == "DONE" {
			sub.wt.Done()
		}
	}
	return nil
}

func (s *WebsocketService) Unsubscribe(machineID string, username string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	subscribers, exists := s.Subscribers[machineID]
	if !exists {
		return fmt.Errorf("no subscribers for machine %s", machineID)
	}
	for i, sub := range subscribers {
		if username == sub.username {
			s.Subscribers[machineID] = append(subscribers[:i], subscribers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("subscriber %s not found for machine %s", username, machineID)
}
