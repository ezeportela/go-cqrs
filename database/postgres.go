package database

import (
	"context"
	"database/sql"

	"github.com/ezeportela/go-cqrs/models"
	_ "github.com/lib/pq"
)

type PostgreRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgreRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgreRepository{db}, nil
}

func (r *PostgreRepository) Close() {
	r.db.Close()
}

func (r *PostgreRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := r.db.Exec("INSERT INTO feeds (id, title, description) VALUES ($1, $2, $3)", feed.ID, feed.Title, feed.Description)
	return err
}

func (r *PostgreRepository) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := make([]*models.Feed, 0)

	for rows.Next() {
		feed := &models.Feed{}
		if err := rows.Scan(&feed.ID, &feed.Title, &feed.Description, &feed.CreatedAt); err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}
