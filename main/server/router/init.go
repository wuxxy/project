package router

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/middleware"
	"github.com/wuxxy/project/main/router/auth"
)

func Init(app *iris.Application) {
	app.Post("/auth/register", auth.Register)
	app.Post("/auth/login", auth.Login)
	app.Post("/auth/token", auth.Token)
	authRouter := app.Party("/api", middleware.AuthMiddleware)
	{
		authRouter.Get("/me", auth.Me)
		authRouter.Post("/logout", auth.Logout)
	}
}
