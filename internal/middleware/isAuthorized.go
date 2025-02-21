package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

func IsAuthorized(endpoints map[string]types.Endpoint) fiber.Handler {
	return func(c *fiber.Ctx) error {
		table := c.Params("table")
		var path string
		if table != "" {
			path = fmt.Sprintf("/%s", table)
		} else {
			path = c.Path()
		}

		endpoint := endpoints[path]

		if endpoint.Table == "" || len(endpoint.Table) == 0 {
			return c.Status(404).JSON(fiber.Map{
				"message": "Endpoint not found",
			})

		}

		if endpoint.IsPublic == true {
			return c.Next()
		}

		apiKey := os.Getenv("API_KEY")
		if apiKey == c.Get("api-key") {
			return c.Next()
		}

		return c.Status(403).JSON(fiber.Map{
			"message": "You need to provide a valida api key.",
		})
	}
}
