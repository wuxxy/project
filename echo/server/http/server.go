package http

import (
	"log"

	"github.com/kataras/iris/v12"
)

func StartServer(port string) {
	// Initialize the Echo server
	app := iris.New()

	// Start the server
	if err := app.Listen(port); err != nil {
		log.Fatal(err)
	}
}
