package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func Connect(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{Addr: addr})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return client, nil
}
