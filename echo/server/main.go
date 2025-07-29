package main

import (
	"sync"

	"github.com/wuxxy/project/echo/http"
	"github.com/wuxxy/project/echo/ws"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		http.StartServer(":6000")
	}()

	go func() {
		defer wg.Done()
		ws.ServeWS(":6001")
	}()

	wg.Wait()
}
