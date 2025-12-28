package websocket

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// wsConn abstracts the websocket connection methods for easier testing and mocking.
type Conn interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

// WSClient manages WebSocket connections for the agent, handling connection lifecycle and message routing.
type WSClient struct {
	host       string
	path       string
	clientName string
	logger     *zerolog.Logger
	conn       Conn
	//TODO Need a better type
	Dispatch func(ctx context.Context, cmd lib.Command) <-chan interface{}
}

// TODO: Need better struct initialzer
// NewWSClient creates and returns a new WSClient instance with the specified connection parameters.
func NewWSClient(host, path, clientName string, logger *zerolog.Logger , dispatch func(ctx context.Context, cmd lib.Command) <-chan interface{}) *WSClient {
	return &WSClient{
		host:       host,
		path:       path,
		clientName: clientName,
		logger:     logger,
		Dispatch: dispatch,
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
	c.logger.Info().Msgf("Connection Established %s/%s", c.host, c.path)
	return nil
}

// Receive continuously reads messages from the WebSocket connection and invokes the provided handler function.
// It blocks until the connection is closed or an error occurs.
func (ws *WSClient) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, data, err := ws.conn.ReadMessage()
			if err != nil {
				ws.logger.Err(err).Msg("Read Error")
				continue
			}

			var cmd lib.Command
			if err := json.Unmarshal(data, &cmd); err != nil {
				ws.logger.Err(err).Msg("Command Error")
				continue
			}

			ws.logger.Info().Msgf("Request Recieved %s", cmd.CMD)
			// Start processing
			stream := ws.Dispatch(ctx, cmd)

			// handle stream in new go routine
			go func() {
				for msg := range stream {
					//
					str, err := json.Marshal(msg)
					if err != nil {
						ws.logger.Err(err).Msg("Marshal Error")
					}
					ws.conn.WriteMessage(websocket.TextMessage, str)
				}
			}()
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
