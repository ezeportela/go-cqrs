package search

import (
	"bytes"
	"context"
	"encoding/json"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/ezeportela/go-cqrs/models"
	"github.com/pkg/errors"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{client}, nil
}

func (r *ElasticSearchRepository) Close() {
	//
}

func (r *ElasticSearchRepository) IndexFeed(ctx context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := r.client.Index(
		"feeds",
		bytes.NewReader(body),
		r.client.Index.WithDocumentID(feed.ID),
		r.client.Index.WithContext(ctx),
		r.client.Index.WithRefresh("wait_for"),
	)
	return err
}

type MapAny map[string]any
type MapInterface map[string]interface{}

func (r *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) (results []*models.Feed, err error) {
	searchQuery := MapAny{
		"query": MapAny{
			"multi_match": MapAny{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzziness":        3,
				"cutoff_frequency": 0.001,
			},
		},
	}
	var buff bytes.Buffer
	if err = json.NewEncoder(&buff).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex("feeds"),
		r.client.Search.WithBody(&buff),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eRes MapAny
	if err = json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}
	feeds := make([]*models.Feed, 0)
	for _, hit := range eRes["hits"].(MapAny)["hits"].([]any) {
		feed := models.Feed{}
		source := hit.(MapAny)["_source"]
		marshal, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(marshal, &feed); err != nil {
			continue
		}
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}
