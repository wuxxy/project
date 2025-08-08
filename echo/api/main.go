package main

import (
	"github.com/wuxxy/project/echo/database"
	"github.com/wuxxy/project/echo/ipc"
	"github.com/wuxxy/project/echo/server"
)

func main() {
	database.ConnectToDb()
	ipc.ConnectToNats()
	server.StartServer(":6000")
}
