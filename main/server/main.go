package main

import (
	"os"

	"github.com/iris-contrib/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/wuxxy/auth/database"
	"github.com/wuxxy/auth/middleware"
	"github.com/wuxxy/auth/router/auth"
)

func main() {
	_ = godotenv.Load()
	app := iris.New()
	database.ConnectToDb()
	app.Logger().SetLevel("debug")
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5602", "http://localhost:5603"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)
	app.Use(func(ctx iris.Context) {
		// Prevent content sniffing
		ctx.Header("X-Content-Type-Options", "nosniff")

		// Prevent iframe embedding (clickjacking)
		ctx.Header("X-Frame-Options", "DENY")

		// Basic XSS protection in older browsers (ignored in modern ones)
		ctx.Header("X-XSS-Protection", "1; mode=block")

		// Prevent referrer leakage
		ctx.Header("Referrer-Policy", "no-referrer")
		csrfHeader := ctx.GetHeader("X-CSRF-Token")
		if csrfHeader != os.Getenv("CSRF_SECRET") {
			ctx.StatusCode(iris.StatusForbidden)
			_ = ctx.JSON(iris.Map{"error": "Invalid request"})
			return
		}

		// Continue to next handler
		ctx.Next()
	})
	app.Post("/auth/register", auth.Register)
	app.Post("/auth/login", auth.Login)
	app.Post("/auth/token", auth.Token)
	authRouter := app.Party("/api", middleware.AuthMiddleware)
	{
		authRouter.Get("/me", auth.Me)
		authRouter.Post("/logout", auth.Logout)
	}
	err := app.Listen(":5000")
	if err != nil {
		return
	}
}
