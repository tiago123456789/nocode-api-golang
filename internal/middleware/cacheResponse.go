package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

const defaultCacheTime = 5

func CacheResponse(cache *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		endpoint := c.Locals("endpoint").(types.Endpoint)
		cacheKey := utils.GetCacheKeyByEndpoint(endpoint, c)
		dataCached, _ := cache.Get(config.GetCacheContext(), cacheKey).Result()
		if dataCached != "" && len(dataCached) > 0 {
			c.Set("content-type", "application/json; charset=utf-8")
			return c.SendString(dataCached)
		}

		err := c.Next()

		if c.Response().StatusCode() == 200 {
			body := c.Response().Body()
			c.Response().StatusCode()
			howManyTimeCache := defaultCacheTime
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
