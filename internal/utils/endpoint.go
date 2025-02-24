package utils

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

var cache *redis.Client

func SetEndpointsInCache(endpoints map[string]types.Endpoint) {
	if cache == nil {
		cache = config.GetCache()
	}

	for key := range endpoints {
		cache.Del(config.GetCacheContext(), key)
		item := endpoints[key]
		data, _ := json.Marshal(item)
		cache.Set(config.GetCacheContext(), key, data, 0)
	}
}
