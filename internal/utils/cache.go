package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

func GetCacheKeyByEndpoint(endpoint types.Endpoint, c *fiber.Ctx) string {
	extraCacheKey := ""
	if len(endpoint.QueryParams) > 0 {
		for _, key := range endpoint.QueryParams {
			value := c.Query(key)
			extraCacheKey += fmt.Sprintf("%s_%s", key, value)
		}
	}

	cacheKey := fmt.Sprintf("response_cached_%s_%s", c.Path(), extraCacheKey)
	return cacheKey
}
