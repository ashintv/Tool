package websocket

import (
	"aetrix/observer/internals/agent/protocol"
	"aetrix/observer/internals/agent/runtime"
	"aetrix/observer/internals/lib"
	"context"
	"fmt"
	"time"
)

type HandlerConfig struct {
	HandlerFn func(ctx context.Context, out chan<- protocol.Event , cmd lib.Command)
	Timeout    time.Duration
}
type Dispatcher struct{
	handlers map[string]HandlerConfig
	invalidHandler func(req lib.Command) HandlerConfig
	DockerRuntime runtime.DockerRuntimeInterface
}

func NewDispatcher(dockerRuntime runtime.DockerRuntimeInterface) *Dispatcher {
	d :=  Dispatcher{
		handlers: make(map[string]HandlerConfig),
		invalidHandler: InvalidCommandConfig,
		DockerRuntime: dockerRuntime,
	}
	d.RegisterHandler()
	return &d
}


// creates a map of command to handler function
func (d *Dispatcher) RegisterHandler() {
	d.handlers[lib.START_NEW_CONTAINER] = HandlerConfig{
		HandlerFn: d.DockerRuntime.StartNewContainer,
		Timeout: 10 * time.Minute,
	}
}

var handler = map[string]HandlerConfig{

}

func InvalidCommandConfig(req lib.Command) HandlerConfig {
	return HandlerConfig{
		HandlerFn: func(ctx context.Context, out chan<- protocol.Event  , cmd lib.Command) {
			protocol.NewEvent(
				protocol.WithError(fmt.Errorf("invalid command: %s", cmd.CMD)),
			).Send(ctx, out)
		},
		Timeout: 2 * time.Second,
	}
}


func (d *Dispatcher) Dispatch(ctx context.Context, req lib.Command) <-chan protocol.Event {

	cfg, ok := d.handlers[string(req.CMD)]
	if !ok {
		cfg = d.invalidHandler(req)
	}

	Stream := make(chan protocol.Event , 8)

	go func() {
		defer close(Stream)
		timeoutCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
		cfg.HandlerFn(timeoutCtx, Stream, req)
	}()
	return Stream
}
