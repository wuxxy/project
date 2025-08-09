package main

import (
	"github.com/joho/godotenv"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/ipc"
	"github.com/wuxxy/project/main/server"
)

func main() {
	_ = godotenv.Load()

	database.ConnectToDb()
	ipc.ConnectToNats()

	server.StartServer(":5000")
}
