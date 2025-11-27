package lib

// MessageType defines the type of websocket message
type MessageType string

const (
	TypeEvent    MessageType = "event"
	TypeResponse MessageType = "response"
	TypeStream   MessageType = "stream"
)
type WsMessage struct {
	Type      MessageType       `json:"type" binding:"required"`
	MachineID string            `json:"machine_id" binding:"required"`
	Payload   string            `json:"payload"`
}