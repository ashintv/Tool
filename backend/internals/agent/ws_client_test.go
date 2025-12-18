package agents

// Test file for ws_client.go
// Focuses on testing the WSClient methods:  Send, Receive,

import (
	"aetrix/observer/internals/lib"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// mock for wsConn
type mockWsConn struct {
	writeMessageFn func(messageType int, data []byte) error
	readMessageFn  func() (int, []byte, error)
	closeFn        func() error
}

func (m *mockWsConn) WriteMessage(messageType int, data []byte) error {
	if m.writeMessageFn == nil {
		return nil
	}
	return m.writeMessageFn(messageType, data)
}
func (m *mockWsConn) ReadMessage() (int, []byte, error) {
	return m.readMessageFn()
}
func (m *mockWsConn) Close() error {
	return m.closeFn()
}

type test struct {
	name        string
	expectError bool
	mockWsConn  *mockWsConn
	onMessage   func(lib.Command)
}

// TESTS
func TestSend(t *testing.T) {
	tests := []test{
		{
			name: "Successful Send",
			mockWsConn: &mockWsConn{
				writeMessageFn: func(messageType int, data []byte) error {
					return nil
				},
			},
			expectError: false,
		},
		{
			name: "Failed Send",
			mockWsConn: &mockWsConn{
				writeMessageFn: func(messageType int, data []byte) error {
					return fmt.Errorf("send error")
				},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &WSClient{
				conn: tc.mockWsConn,
			}

			err := client.Send(lib.WsMessage{Type: lib.TypeResponse, Payload: lib.PayloadType{}})
			if tc.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}

		})
	}
}

func TestReceive(t *testing.T) {
	tests := []struct {
		name        string
		mockWsConn  *mockWsConn
		expectError bool
	}{
		{
			name: "Successful Receive",
			mockWsConn: func() *mockWsConn {
				readOnce := true
				return &mockWsConn{
					readMessageFn: func() (int, []byte, error) {
						if readOnce {
							readOnce = false
							cmd := lib.Command{CMD: "TEST_CMD"}
							data, _ := json.Marshal(cmd)
							return 1, data, nil
						}
						return 0, nil, fmt.Errorf("stop")
					},
				}
			}(),
			expectError: false,
		},
		{
			name: "Receive Error",
			mockWsConn: &mockWsConn{
				readMessageFn: func() (int, []byte, error) {
					return 0, nil, fmt.Errorf("read error")
				},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			success := make(chan struct{}, 1)
			errorCh := make(chan error, 1)

			client := &WSClient{
				conn: tc.mockWsConn,
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go client.Receive(
				ctx,
				func(cmd lib.Command) {
					if cmd.CMD != "TEST_CMD" {
						t.Errorf("unexpected cmd: %s", cmd.CMD)
					}
					success <- struct{}{}
					cancel() 
				},
				func(err error) {
					errorCh <- err
				},
			)

			select {
			case <-success:
				if tc.expectError {
					t.Fatal("expected error but got success")
				}
			case err := <-errorCh:
				if !tc.expectError {
					t.Fatalf("did not expect error but got: %v", err)
				}
			case <-time.After(time.Second):
				t.Fatal("test timed out")
			}
		})
	}
}
