package handlers

import (
	"log"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuxxy/project/main/tokens"
)

type VerifyTokenResponse struct {
	Error     string `msgpack:"error"`
	SessionID string `msgpack:"session_id,omitempty"`
	UserID    string `msgpack:"user_id,omitempty"`
}

func VerifyToken(msg *nats.Msg) {
	var response VerifyTokenResponse
	if len(msg.Data) == 0 {
		response.Error = "Missing data"
		errorResponse, _ := msgpack.Marshal(response)
		_ = msg.Respond([]byte(errorResponse))
		return
	}

	parts := strings.Split(string(msg.Data), " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		response.Error = "Invalid Authorization header format"
		errorResponse, _ := msgpack.Marshal(response)
		_ = msg.Respond([]byte(errorResponse))
		return
	}
	sessionID, userID, err := tokens.VerifyAccessToken(parts[1])
	if err != nil {
		log.Print("this shit expired")
		response.Error = "Invalid or expired token"
		
		errorResponse, _ := msgpack.Marshal(response)
		_ = msg.Respond([]byte(errorResponse))
		return
	}
	response.Error = ""
	response.SessionID = sessionID
	response.UserID = userID
	errorResponse, _ := msgpack.Marshal(response)
	_ = msg.Respond([]byte(errorResponse))
	return
}
