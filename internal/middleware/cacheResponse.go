package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

func CacheResponse(cache *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		extraCacheKey := ""
		endpoint := c.Locals(c.Path()).(types.Endpoint)
		if len(endpoint.QueryParams) > 0 {
			for _, key := range endpoint.QueryParams {
				value := c.Query(key)
				extraCacheKey += fmt.Sprintf("%s_%s", key, value)
			}
		}

		cacheKey := fmt.Sprintf("response_cached_%s_%s", c.Path(), extraCacheKey)
		dataCached, _ := cache.Get(config.GetCacheContext(), cacheKey).Result()
		if dataCached != "" && len(dataCached) > 0 {
			c.Set("content-type", "application/json; charset=utf-8")
			return c.SendString(dataCached)
		}

		err := c.Next()

		if c.Response().StatusCode() == 200 {
			body := c.Response().Body()
			c.Response().StatusCode()
			howManyTimeCache := 5
			if endpoint.IsCacheable {
				howManyTimeCache = endpoint.CacheTtl
			}
			cache.Set(
				config.GetCacheContext(),
				cacheKey,
				body,
				time.Second*time.Duration(howManyTimeCache),
			)
		}

		return err
	}
}
