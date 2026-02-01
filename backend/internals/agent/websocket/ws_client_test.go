package websocket

import (
	"aetrix/observer/internals/agent/protocol"
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"errors"

	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

type MockConn struct {
	ReadMessageFn  func() (messageType int, p []byte, err error)
	WriteMessageFn func(messageType int, data []byte) error
	CloseFn        func() error
}

func (m *MockConn) ReadMessage() (messageType int, p []byte, err error) {
	if m.ReadMessageFn == nil {
		cmd := lib.Command{
			ContainerID: "test",
			MachineID:   "test",
			CMD:         lib.CREATE_CONTAINER,
		}

		srt, _ := json.Marshal(cmd)
		return 1, srt, nil
	}
	return m.ReadMessageFn()
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	if m.WriteMessageFn == nil {
		_ = data
		return nil
	}
	return m.WriteMessageFn(messageType, data)
}

func (m *MockConn) Close() error {
	return nil
}

func testLogger() *zerolog.Logger {
	l := zerolog.New(io.Discard)
	return &l
}

func validCommandBytes() []byte {
	cmd := lib.Command{
		ContainerID: "test",
		MachineID:   "test",
		CMD:         lib.CREATE_CONTAINER,
	}
	b, _ := json.Marshal(cmd)
	return b
}

func TestStart(t *testing.T) {
	t.Run("Data Flow", func(t *testing.T) {
		writeCalled := make(chan struct{}, 1)

		ws := &WSClient{
			logger: testLogger(),
			conn: &MockConn{
				ReadMessageFn: func() (int, []byte, error) {
					return 1, validCommandBytes(), nil
				},
				WriteMessageFn: func(int, []byte) error {
					writeCalled <- struct{}{}
					return nil
				},
			},

			Dispatch: func(ctx context.Context, cmd lib.Command) <-chan protocol.Event {
				out := make(chan protocol.Event)
				go func() {
					defer close(out)
					out <- protocol.NewEvent(
						protocol.WithMessage("test message"),
					)
				}()
				return out
			},
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ws.Start(ctx)

		select {
		case <-writeCalled:
		case <-time.After(time.Second):
			t.Fatal("expected WriteMessage to be called")
		}
	})

	t.Run("read error exits loop", func(t *testing.T) {
		readCalled := make(chan struct{}, 1)

		ws := &WSClient{
			logger: testLogger(),
			conn: &MockConn{
				ReadMessageFn: func() (int, []byte, error) {
					readCalled <- struct{}{}
					return 0, nil, errors.New("read failed")
				},
			},
			Dispatch: func(ctx context.Context, cmd lib.Command) <-chan protocol.Event {
				t.Fatal("dispatcher should not be called")
				return nil
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ws.Start(ctx)

		select {
		case <-readCalled:
		case <-time.After(time.Second):
			t.Fatal("expected ReadMessage to be called")
		}
	})



	t.Run("marsha error exits loop", func(t *testing.T) {
		readCalled := make(chan struct{}, 1)

		ws := &WSClient{
			logger: testLogger(),
			conn: &MockConn{
				ReadMessageFn: func() (int, []byte, error) {
					readCalled <- struct{}{}
					return 0, nil, errors.New("read failed")
				},
			},
			Dispatch: func(ctx context.Context, cmd lib.Command) <-chan protocol.Event {
				t.Fatal("dispatcher should not be called")
				return nil
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ws.Start(ctx)

		select {
		case <-readCalled:
		case <-time.After(time.Second):
			t.Fatal("expected ReadMessage to be called")
		}
	})

}
