package agents

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// WSClientInterface defines the contract for websocket client interactions.
type WSClientInterface interface {
	Connect() error
	Send(msg lib.WsMessage) error
	Receive(ctx context.Context, onMessage func(lib.Command))
	Close() error
}

// WSClient manages WebSocket connections for the agent, handling connection lifecycle and message routing.
type WSClient struct {
	host       string
	path       string
	clientName string
	conn       *websocket.Conn
}

// NewWSClient creates and returns a new WSClient instance with the specified connection parameters.
func NewWSClient(host, path, clientName string) *WSClient {
	return &WSClient{
		host:       host,
		path:       path,
		clientName: clientName,
	}
}

// Connect establishes a WebSocket connection to the configured server.
// It returns an error if the connection fails.
func (c *WSClient) Connect() error {
	u := url.URL{
		Scheme: "ws",
		Host:   c.host,
		Path:   c.path + "/" + c.clientName,
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	log.Println("WebSocket connected:", u.String())
	return nil
}

// Send marshals and transmits a WsMessage over the WebSocket connection.
// It returns an error if marshaling or sending fails.
func (c *WSClient) Send(msg lib.WsMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Receive continuously reads messages from the WebSocket connection and invokes the provided handler function.
// It blocks until the connection is closed or an error occurs.
func (c *WSClient) Receive(ctx context.Context, onMessage func(lib.Command)) {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			// Happens when socket is closed or network error
			return
		}

		var cmd lib.Command
		if err := json.Unmarshal(data, &cmd); err != nil {
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
			onMessage(cmd) // Assuming NewHandler initializes a Handler with necessary dependencies

		}
	}
}

// Close gracefully closes the WebSocket connection if it exists.
// It returns nil if the connection is already closed or an error if closing fails.
func (c *WSClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
