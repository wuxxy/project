package router

import (
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/middleware"
	"github.com/wuxxy/project/main/router/admin/services"
	"github.com/wuxxy/project/main/router/admin/users"
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
	adminRouter := app.Party("/admin", middleware.AdminMiddleware)
	{
		adminRouter.Get("/services", services.ServicesReadAll)
		adminRouter.Post("/services", services.ServicesCreate)
		adminRouter.Get("/services/struct", services.ServicesStruct)
		adminRouter.Delete("/services/:id", services.ServicesDelete)
		adminRouter.Put("/services/:id", services.ServicesUpdate)

		adminRouter.Get("/users", users.UsersReadAll)
		adminRouter.Post("/users", users.UsersCreate)
		adminRouter.Get("/users/struct", users.UsersStruct)
		adminRouter.Delete("/users/:id", users.UsersDelete)
		adminRouter.Put("/users/:id", users.UsersUpdate)
	}
}
