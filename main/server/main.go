package main

import (
	"github.com/joho/godotenv"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/server"
)

func main() {
	_ = godotenv.Load()

	database.ConnectToDb()
	server.StartServer(":5000")
}
