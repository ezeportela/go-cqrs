package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ezeportela/go-cqrs/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	hub := NewHub()

	natsAddress := fmt.Sprintf("nats://%s", cfg.NatsAddress)
	n, err := events.NewNatsEventStore(natsAddress)
	if err != nil {
		log.Fatal(err)
	}

	err = n.OnCreateFeed(func(m events.CreatedFeedMessage) {
		hub.Broadcast(
			newCreatedFeedMessage(m.ID, m.Title, m.Description, m.CreatedAt),
			nil,
		)
	})
	events.SetEventStore(n)

	defer events.Close()

	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWebSocket)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("%v", err)
	}
}
