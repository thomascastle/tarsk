package messaging

import (
	"context"

	"github.com/go-redis/redis/v8"
)

func NewClient() (*redis.Client, error) {
	r_client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	if _, e := r_client.Ping(context.Background()).Result(); e != nil {
		return nil, e
	}

	return r_client, nil
}
