package repositories

import (
	"context"

	"github.com/ezeportela/go-cqrs/models"
)

type Repository interface {
	Close()
	InsertFeed(ctx context.Context, feed *models.Feed)
	ListFeeds(ctx context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepository(r Repository) {
	repository = r
}

func Close() {
	repository.Close()
}

func InsertFeed(ctx context.Context, feed *models.Feed) {
	repository.InsertFeed(ctx, feed)
}

func ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	return repository.ListFeeds(ctx)
}
