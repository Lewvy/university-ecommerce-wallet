package cache

import (
	// "context"
	"os"

	"github.com/valkey-io/valkey-go"
)

func New() (valkey.Client, error) {
	address, err := valkey.ParseURL(os.Getenv("CACHE_DSN"))
	if err != nil {
		return nil, err
	}
	client, err := valkey.NewClient(address)
	if err != nil {
		return nil, err
	}
	return client, nil
}
