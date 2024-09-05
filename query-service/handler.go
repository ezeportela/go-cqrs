package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ezeportela/go-cqrs/events"
	"github.com/ezeportela/go-cqrs/models"
	"github.com/ezeportela/go-cqrs/repositories"
	"github.com/ezeportela/go-cqrs/search"
)

func onCreateFeed(m events.CreatedFeedMessage) {
	feed := &models.Feed{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
	if err := search.IndexFeed(context.Background(), feed); err != nil {
		log.Printf("failed to insert feed into search: %v", err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	feeds, err := repositories.ListFeeds(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}
	feeds, err := search.SearchFeed(ctx, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feeds)
}
