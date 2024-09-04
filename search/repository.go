package search

import (
	"context"

	"github.com/ezeportela/go-cqrs/models"
)

type SearchRepository interface {
	Close()
	IndexFeed(ctx context.Context, feed *models.Feed) error
	SearchFeed(ctx context.Context, query string) ([]*models.Feed, error)
}

var repository SearchRepository

func SetSearchRepository(r SearchRepository) {
	repository = r
}

func Close() {
	repository.Close()
}

func IndexFeed(ctx context.Context, feed *models.Feed) error {
	return repository.IndexFeed(ctx, feed)
}

func SearchFeed(ctx context.Context, query string) ([]*models.Feed, error) {
	return repository.SearchFeed(ctx, query)
}
