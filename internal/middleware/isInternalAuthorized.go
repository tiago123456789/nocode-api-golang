package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func IsInternalAuthorized(c *fiber.Ctx) error {
	apiKey := os.Getenv("API_KEY")
	if apiKey == c.Get("api-key") {
		return c.Next()
	}

	return c.Status(403).JSON(fiber.Map{
		"message": "You need to provide a valida api key.",
	})
}
