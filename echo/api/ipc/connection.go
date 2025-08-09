package ipc

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var NC *nats.Conn

func ConnectToNats() {
	var err error
	url := os.Getenv("NATS_URL")
	NC, err = nats.Connect(url,
		nats.Name("api-gateway"),
		nats.MaxReconnects(-1),            // keep trying
		nats.ReconnectWait(2*time.Second), // backoff
		nats.PingInterval(10*time.Second),
		nats.Timeout(3*time.Second), // dial timeout
		nats.DisconnectErrHandler(func(nc *nats.Conn, e error) {
			log.Printf("NATS disconnected: %v", e)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("NATS closed: %v", nc.LastError())
		}),
	)
	if err != nil {
		log.Fatalf("NATS connect failed: %v", err)
	} else {
		log.Println("NATS Connected")
	}

}
