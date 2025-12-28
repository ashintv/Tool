package websocket

import (
	"aetrix/observer/internals/lib"
	"context"
	"fmt"
	"time"
)



type Dispatcher struct{
	handlers map[string]HandlerConfig
	invalidHandler func(req lib.Command) HandlerConfig
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: handler,
		invalidHandler: InvalidCommandConfig,
	}
}

type HandlerConfig struct {
	HandlerFn func(ctx context.Context, c chan<- interface{})
	Timeout   time.Duration
}

var handler = map[string]HandlerConfig{
	lib.LIST_CONTAINERS: {
		HandlerFn: func(ctx context.Context, c chan<- interface{}) {},
		Timeout:   500 * time.Millisecond,
	},
	lib.CREATE_CONTAINER: {
		HandlerFn: func(ctx context.Context, c chan<- interface{}) {},
		Timeout:   500 * time.Millisecond,
	},
	lib.DELETE_CONTAINER: {
		HandlerFn: func(ctx context.Context, c chan<- interface{}) {},
		Timeout:   500 * time.Millisecond,
	},
	lib.START_CONTAINER: {
		HandlerFn: func(ctx context.Context, c chan<- interface{}) {},
		Timeout:   500 * time.Millisecond,
	},
	lib.STOP_CONTAINER: {
		HandlerFn: func(ctx context.Context, c chan<- interface{}) {},
		Timeout:   500 * time.Millisecond,
	},
}

func InvalidCommandConfig(req lib.Command) HandlerConfig {
	return HandlerConfig{
		HandlerFn: func(ctx context.Context, out chan<- interface{}) {
			select {
			case <-ctx.Done():
				return
			case out <- fmt.Errorf("invalid command: %s", req.CMD):
			}
		},
		Timeout: 2 * time.Second,
	}
}


func (d *Dispatcher) Dispatch(ctx context.Context, req lib.Command) <-chan interface{} {

	cfg, ok := d.handlers[string(req.CMD)]
	if !ok {
		cfg = d.invalidHandler(req)
	}

	Stream := make(chan interface{})
	go func() {
		defer close(Stream)
		timeoutCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
		cfg.HandlerFn(timeoutCtx, Stream)
	}()
	return Stream
}
