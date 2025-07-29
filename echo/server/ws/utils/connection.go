package utils

import (
	"log"

	"github.com/coder/websocket"
)

func CloseConnection(conn *websocket.Conn, code websocket.StatusCode, reason string) {
	if err := conn.Close(code, reason); err != nil {
		log.Println("Error closing WebSocket connection:", err)
	}
}
