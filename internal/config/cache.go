package config

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func GetCacheContext() context.Context {
	return ctx
}

func GetCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
