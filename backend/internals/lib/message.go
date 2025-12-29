package lib

// MessageType defines the type of websocket message
type MessageType string
const (
	TypeEvent     MessageType = "event"
	TypeResponse  MessageType = "response"
	TypeStream    MessageType = "stream"
	TypeStreamEnd MessageType = "stream_end"
)

type EventTypes string

const (
	UNEXPECTED_STOP EventTypes = "UNEXPECTED"
	PING            EventTypes = "PING"
)

type PayloadType struct {
	Error string       `json:"error"`
	Data  interface{} `json:"data"`
}

type WsMessage struct {
	Type      MessageType `json:"type" binding:"required"`
	MachineID string      `json:"machine_id" binding:"required"`
	Payload   PayloadType `json:"payload"`
}

type optFunc func(*WsMessage)

func defaultOpts() WsMessage {
	return WsMessage{
		Type:      TypeResponse,
		MachineID: "",
		Payload: PayloadType{
			Error: "",
			Data:  nil,
		},
	}
}

func WithMessageType(t MessageType) optFunc {
	return func(wsm *WsMessage) {
		wsm.Type = t
	}
}

func WithMachineID(id string) optFunc {
	return func(wsm *WsMessage) {
		wsm.MachineID = id
	}
}
func WithPayloadData(data interface{}) optFunc {
	return func(wsm *WsMessage) {
		wsm.Payload.Data = data
	}
}

func WithPayloadError(err error) optFunc {
	return func(wsm *WsMessage) {
		if err != nil {
			wsm.Payload.Error = err.Error()
		} else {
			wsm.Payload.Error = ""
		}
	}
}

func NewWsMessage(opt ...optFunc) WsMessage {
	WsMessage := defaultOpts()
	for _, o := range opt {
		o(&WsMessage)
	}
	return WsMessage
}

// SteamMEssage


