package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

func IsAuthorized() fiber.Handler {
	return func(c *fiber.Ctx) error {
		table := c.Params("table")
		var path string
		if table != "" {
			path = fmt.Sprintf("/%s", table)
		} else {
			path = c.Path()
		}

		cacheKey, _ := config.GetCache().Exists(config.GetCacheContext(), path).Result()
		if cacheKey == 0 {
			return c.Status(404).JSON(fiber.Map{
				"message": "Endpoint not found",
			})
		}

		cacheValue, _ := config.GetCache().Get(config.GetCacheContext(), path).Result()
		var endpoint types.Endpoint
		json.Unmarshal([]byte(cacheValue), &endpoint)
		c.Locals("endpoint", endpoint)
		if endpoint.IsPublic == true {
			return c.Next()
		}

		apiKey := os.Getenv("API_KEY")
		if apiKey == c.Get("api-key") {
			return c.Next()
		}

		accessToken := c.Get("Authorization")
		accessToken = strings.ReplaceAll(accessToken, "Bearer ", "")

		if err := utils.IsValidToken(accessToken); err == nil {
			return c.Next()
		}

		return c.Status(403).JSON(fiber.Map{
			"message": "You need to provide a valida api key or accessToken.",
		})
	}
}
