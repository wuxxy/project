package server

import "github.com/kataras/iris/v12"

func StartServer(addr string) {
	// Initialize the Iris application
	app := iris.New()

	// Set up routes and middleware here
	// For example:
	// app.Get("/", func(ctx iris.Context) {
	// 	ctx.WriteString("Hello, World!")
	// })

	// Start the server
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
