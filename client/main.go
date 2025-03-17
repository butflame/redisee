package client

import "github.com/redis/go-redis/v9"

func New(options redis.Options) *redis.Client {
	return redis.NewClient(&options)
}
