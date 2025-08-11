package router

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/echo/router/user"
)

func Init(app *iris.Application) {
	apiRouter := app.Party("/api", AuthMiddleware)
	apiRouter.Get("/user/me", user.Me)
}
