package router

import (
	"log"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/wuxxy/project/echo/ipc"
)

type VerifyTokenResponse struct {
	SessionID string
	UserID    string
	Error     string
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
	log.Print(res.UserID)
	ctx.Next()
}
