package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

func IsInternalAuthorized(c *fiber.Ctx) error {
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
