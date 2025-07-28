package ws

import (
	"context"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

type AuthRequestData struct {
	Token string `json:"access_token"`
}
type Request[T AuthRequestData] struct {
	Data      T       `json:"data"`
	Op        uint8   `json:"op"`
	SessionId *string `json:"sid"`
}
type OpCode uint8
type Handler func(*websocket.Conn, any) error

var handlers = map[OpCode]Handler{}

func ServeWS(port string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade to WebSocket
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true, // skip Origin check (you can tighten this later)
		})
		if err != nil {
			log.Println("Accept error:", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "goodbye")

		log.Println("WebSocket connection established")

		// Echo loop
		for {
			// Read message
			msgType, data, err := conn.Read(context.Background())
			if err != nil {
				log.Println("Read error:", err)
				break
			}

			log.Printf("Received: %s\n", data)

			// Write back
			if err := conn.Write(context.Background(), msgType, data); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	})

}
func RegisterOP(op uint8, handler func(*websocket.Conn, []byte) error) {

}
