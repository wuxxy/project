package session

import (
	"time"

	"github.com/coder/websocket"
)

type Connection struct {
	HeartbeatTime uint16
	Conn          *websocket.Conn
	UserId        string
	LastPing      time.Time
}

var ConnectionSessions = make(map[string]Connection)
