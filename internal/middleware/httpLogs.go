package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
)

var logger *slog.Logger = config.GetLogger()

func HttpLogs(c *fiber.Ctx) error {
	logger.Info(
		"Starting the request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
	)

	err := c.Next()

	logger.Info(
		"Finished the request",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
	)

	return err
}
