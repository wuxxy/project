package auth

import (
	"github.com/kataras/iris/v12"
)

type DeleteSessionBody struct {
	SessionID string `json:"session_id" `
}

func DeleteSession(c iris.Context) {
	var req DeleteSessionBody

	// Parse input as request body
	if err := c.ReadJSON(&req); err != nil {
		c.StatusCode(iris.StatusBadRequest)
		_ = c.JSON(iris.Map{"error": "Invalid request body"})
		return
	}
}
