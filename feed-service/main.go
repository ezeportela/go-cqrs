package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ezeportela/go-cqrs/database"
	"github.com/ezeportela/go-cqrs/events"
	"github.com/ezeportela/go-cqrs/repositories"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgreDB       string `envconfig:"POSTGRES_DB"`
	PostgreUser     string `envconfig:"POSTGRES_USER"`
	PostgrePassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress     string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgreUser, cfg.PostgrePassword, cfg.PostgreDB)
	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		log.Fatal(err)
	}

	repositories.SetRepository(repo)

	natsAddress := fmt.Sprintf("nats://%s", cfg.NatsAddress)
	n, err := events.NewNatsEventStore(natsAddress)
	if err != nil {
		log.Fatal(err)
	}
	events.SetEventStore(n)

	defer events.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("%v", err)
	}
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/feeds", createFeedHandler).Methods(http.MethodPost)
	router.HandleFunc("/feeds", listFeedHandler).Methods(http.MethodGet)
	return
}
