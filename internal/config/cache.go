package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func GetCacheContext() context.Context {
	return ctx
}

func InitCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func GetCache() *redis.Client {
	return rdb
}
