package search

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	ES *elasticsearch.Client
}

func NewClient(dsn string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{dsn},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{ES: es}, nil
}
