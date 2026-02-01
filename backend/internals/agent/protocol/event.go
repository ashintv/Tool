package protocol

import "context"

type Event struct {
	Data    any
	Error   error
	Message string
}

type ErrorInfo struct {
	Code    string
	Message string
}

type EventOption func(*Event)

func NewEvent(opts ...EventOption) Event {
	e := Event{}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

func WithMessage(msg string) EventOption {
	return func(e *Event) {
		e.Message = msg
	}
}

func WithError(e error) EventOption {
	return func(ev *Event) {
		ev.Error = e
	}
}

func WithData(data any) EventOption {
	return func(e *Event) {
		e.Data = data
	}
}

func (e Event) Send(ctx context.Context, ch chan<- Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- e:
		return nil
	}
}

func (e Event) TrySend(ctx context.Context, ch chan<- Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case ch <- e:
		return nil
	default:
	}
    return nil
}
