package agent

import (
	"aetrix/observer/internals/lib"
	"context"
)

// handleMessages routes incoming commands to the appropriate handler methods.
// It is passed to the WSClient receive method as a callback function and returns a WsMessage with the operation result.
func MessageRouter(ctx context.Context, req lib.Command, h *Handler) lib.WsMessage {
	wsMessage := lib.NewWsMessage(lib.TypeResponse, req.MachineID, lib.PayloadType{})

	switch req.CMD {
	case lib.LIST_CONTAINERS:
		wsMessage = h.HandleListContainers(ctx, req)

	case lib.DELETE_CONTAINER:
		wsMessage = h.HandleDeleteContainer(ctx, req)

	case lib.STOP_CONTAINER:
		wsMessage = h.HandleStopContainer(ctx, req)

	case lib.START_NEW_CONTAINER:
		wsMessage = h.HandleStartNewContainer(ctx, req)

	case lib.RESTART_CONTAINER:
		wsMessage = h.HandleRestartContainer(ctx, req)

	case lib.START_CONTAINER:
		wsMessage = h.HandleStartContainer(ctx, req)

	default:
		wsMessage.Payload.Error = "unknown command: " + string(req.CMD)
	}

	return wsMessage
}
