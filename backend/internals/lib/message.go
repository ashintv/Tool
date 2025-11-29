package lib

// MessageType defines the type of websocket message
type MessageType string

const (
	TypeEvent    MessageType = "event"
	TypeResponse MessageType = "response"
	TypeStream   MessageType = "stream"
	TypeStreamEnd MessageType = "stream_end"
)
type PayloadType struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}
type WsMessage struct {
	Type      MessageType       `json:"type" binding:"required"`
	MachineID string            `json:"machine_id" binding:"required"`
	Payload   PayloadType       `json:"payload"`
}

func NewWsMessage(messageType MessageType, machineID string, payload PayloadType) WsMessage {
	return WsMessage{
		Type:      messageType,
		MachineID: machineID,
		Payload:   payload,
	}
}
