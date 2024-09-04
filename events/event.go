package events

import (
	"context"

	"github.com/ezeportela/go-cqrs/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	OnCreateFeed(f func(CreatedFeedMessage)) error
}

var eventStore EventStore

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

func SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

func OnCreateFeed(f func(CreatedFeedMessage)) error {
	return eventStore.OnCreateFeed(f)
}
