package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/coder/websocket"
	"github.com/wuxxy/project/echo/ws/handlers"
	"github.com/wuxxy/project/echo/ws/utils"
)

type Request struct {
	Data      map[string]any `json:"data"`
	Op        uint8          `json:"op"`
	SessionId *string        `json:"sid"`
}
type OpCode uint8
type HandlerFunction func(*websocket.Conn, map[string]any) error

var handlersRegistry = map[OpCode]HandlerFunction{}

func ServeWS(port string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Println("WebSocket accept error:", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "connection closed gracefully")
		conn.SetReadLimit(1024)
		log.Println("WebSocket connection established")

		// Heartbeat tracking
		const heartbeatTimeout = 10 * time.Second // allow some buffer
		lastHeartbeat := time.Now()

		done := make(chan struct{})

		// Start heartbeat monitor
		go func() {
			ticker := time.NewTicker(heartbeatTimeout)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if time.Since(lastHeartbeat) > heartbeatTimeout {
						log.Println("Heartbeat timeout. Closing connection.")
						conn.Close(websocket.StatusPolicyViolation, "missed heartbeat")
						return
					}
				case <-done:
					return
				}
			}
		}()

		// Read loop
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			_, data, err := conn.Read(ctx)
			cancel()

			if err != nil {
				log.Println("WebSocket read error:", err)
				break
			}

			var incoming Request
			if err := json.Unmarshal(data, &incoming); err != nil {
				utils.CloseConnection(conn, websocket.StatusUnsupportedData, "Invalid JSON")
				break
			}

			if incoming.Data == nil {
				utils.CloseConnection(conn, websocket.StatusUnsupportedData, "Missing data")
				break
			}

			// Check for heartbeat op: 2
			if incoming.Op == 2 {
				lastHeartbeat = time.Now()
				continue // no need to call a handler
			}

			handler := handlersRegistry[OpCode(incoming.Op)]
			if handler == nil {
				utils.CloseConnection(conn, websocket.StatusUnsupportedData, "Unknown op")
				break
			}

			if err := handler(conn, incoming.Data); err != nil {
				log.Println("Handler error:", err)
				break
			}
		}

		close(done)
	})

	handlersRegistry[1] = handlers.Hello

	log.Printf("Starting WebSocket server on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
