package ipc

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

var NC *nats.Conn

func ConnectToNats() {
	var err error
	user := os.Getenv("NATSUSERNAME")
	pass := os.Getenv("NATSPASS")

	url := "nats://" + user + ":" + pass + "@localhost:4222"

	NC, err = nats.Connect(url)
	if err != nil {
		log.Fatalf("❌ Failed to connect to NATS: %v", err)
	}

	defer NC.Close()

}
