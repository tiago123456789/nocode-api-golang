package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client
var logger = GetLogger()

func GetCacheContext() context.Context {
	return ctx
}

func InitCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	err := rdb.ConfigSet(
		context.Background(),
		"maxmemory",
		os.Getenv("REDIS_LIMIT_MEMORY"),
	).Err()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to set maxmemory: %v", err))
		panic(err)
	}

	err = rdb.ConfigSet(context.Background(), "maxmemory-policy", "allkeys-lru").Err()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to set LRU policy: %v", err))
		panic(err)
	}
}

func GetCache() *redis.Client {
	return rdb
}
