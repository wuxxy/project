package handlers

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/wuxxy/project/echo/ws/session"
	"github.com/wuxxy/project/echo/ws/utils"
)

type helloResponse struct {
	HeartbeatTime uint16 `json:"heartbeat_time"`
	SessionId     string `json:"session_id"`
}

func Hello(conn *websocket.Conn, data map[string]any) error {

	if data == nil {
		utils.CloseConnection(conn, websocket.StatusUnsupportedData, "Invalid data format")
		return nil
	}
	if data["token"] == nil {
		utils.CloseConnection(conn, websocket.StatusUnsupportedData, "Invalid data format")
	}
	newSession := session.Connection{
		Conn:          conn,
		UserId:        data["token"].(string),
		HeartbeatTime: 5000,
		LastPing:      time.Now(),
	}
	newSessionID := uuid.New().String()
	session.ConnectionSessions[newSessionID] = newSession
	conn.Write(context.Background(), websocket.MessageText, []byte(`{"op":1, "data":{"heartbeat_time":5000,"session_id":"`+newSessionID+`"}}`))
	return nil
}
