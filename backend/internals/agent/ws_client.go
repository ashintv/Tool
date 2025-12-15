// interface for websocket client interactions
package agents

import (
	"aetrix/observer/internals/lib"
	"encoding/json"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type WSClientInterface interface {
	Connect() error
	Send(msg lib.WsMessage) error
	Receive(handler func(lib.Command))
	Close() error
}

type WSClient struct {
	host       string
	path       string
	clientName string
	conn       *websocket.Conn
}

func NewWSClient(host, path, clientName string) *WSClient {
	return &WSClient{
		host:       host,
		path:       path,
		clientName: clientName,
	}
}

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

func (c *WSClient) Send(msg lib.WsMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *WSClient) Receive(handler func(lib.Command)) {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var cmd lib.Command
		if err := json.Unmarshal(data, &cmd); err != nil {
			continue
		}

		handler(cmd)
	}
}

func (c *WSClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
