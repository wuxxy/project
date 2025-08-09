package ipc

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/wuxxy/project/main/tokens"
)
import "github.com/vmihailenco/msgpack/v5"

type VerifyTokenResponse struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Error     string `json:"error"`
}

func InitHandler() {
	NC.Subscribe("auth.verify_token", func(msg *nats.Msg) {
		var response VerifyTokenResponse
		if len(msg.Data) == 0 {
			response.Error = "Missing data"
			errorResponse, _ := msgpack.Marshal(response)
			_ = msg.Respond([]byte(errorResponse))
			return
		}

		parts := strings.Split(string(msg.Data), " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			_ = msg.Respond([]byte("Invalid Authorization header format"))
			return
		}
		sessionID, userID, err := tokens.VerifyAccessToken(parts[1])
		if err != nil {
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
	})
}
