package controller

import (
	"encoding/json"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tiago123456789/nocode-api-golang/internal/config"
	serviceModule "github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type EndpointController struct {
	service serviceModule.EndpointService
	cache   *redis.Client
	logger  *slog.Logger
}

func EndpointControllerNew(
	service serviceModule.EndpointService,
	cache *redis.Client,
	logger *slog.Logger,
) *EndpointController {
	return &EndpointController{
		service: service,
		cache:   cache,
		logger:  logger,
	}
}

func (e *EndpointController) GetAllCreated(c *fiber.Ctx) error {
	results, _ := e.service.GetAllCreated()
	return c.JSON(fiber.Map{
		"data": results,
	})
}

func (e *EndpointController) DeleteById(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	path, err := e.service.Delete(id)
	if err != nil {
		e.logger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	e.cache.Del(config.GetCacheContext(), path)
	return c.SendStatus(204)
}

func (e *EndpointController) Create(c *fiber.Ctx) error {
	var endpoint types.Endpoint
	c.BodyParser(&endpoint)

	_, err := e.service.Create(endpoint)
	if err != nil {
		e.logger.Error(err.Error())
		return c.Status(409).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	data, _ := json.Marshal(endpoint)
	e.cache.Set(config.GetCacheContext(), endpoint.Path, data, 0)
	return c.SendStatus(201)
}
