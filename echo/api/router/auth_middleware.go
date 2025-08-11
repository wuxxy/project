package router

import (
	"log"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuxxy/project/echo/ipc"
)

type VerifyTokenResponse struct {
	Error     string `msgpack:"error"`
	SessionID string `msgpack:"session_id,omitempty"`
	UserID    string `msgpack:"user_id,omitempty"`
}

func AuthMiddleware(ctx iris.Context) {
	msgPackResponse, err := ipc.NC.Request("auth.verify_token", []byte(ctx.GetHeader("Authorization")), 2*time.Second)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		log.Print(err.Error())
		_ = ctx.JSON(iris.Map{"error": "Invalid or expired token"})
		return
	}
	var res VerifyTokenResponse
	err = msgpack.Unmarshal(msgPackResponse.Data, &res)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_ = ctx.JSON(iris.Map{"error": "Internal server error"})
		return
	}
	if res.Error != "" {

		ctx.StatusCode(iris.StatusUnauthorized)
		_ = ctx.JSON(iris.Map{"error": res.Error})
		return
	}
	log.Print(res.SessionID)
	ctx.Values().Set("user_id", res.UserID)
	ctx.Values().Set("session_id", res.SessionID)
	ctx.Next()
}
