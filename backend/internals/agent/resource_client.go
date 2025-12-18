package agent

import (
	"aetrix/observer/internals/lib"
	"context"
	"log"
	"time"
)

func StartResourceMonitor(
	ctx context.Context,
	ws *WSClient,
	machineID string,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Resource monitor stopped")
			return

		case <-ticker.C:
			metrics := map[string]any{
				"cpu":    12.3, // placeholder
				"memory": 512,  // placeholder
			}

			msg := lib.NewWsMessage(
				lib.TypeStream,
				machineID,
				lib.PayloadType{Data: metrics},
			)

			if err := ws.Send(msg); err != nil {
				log.Println("metrics send error:", err)
			}
		}
	}
}
