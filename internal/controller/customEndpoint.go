package controller

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	serviceModule "github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

type CustomEndpointController struct {
	service              serviceModule.CustomEndpointService
	actionsBeforePersist map[string]types.ActionInterface
	logger               *slog.Logger
}

func CustomEndpointControllerNew(
	service serviceModule.CustomEndpointService,
	actionsBeforePersist map[string]types.ActionInterface,
	logger *slog.Logger,
) *CustomEndpointController {
	return &CustomEndpointController{
		service:              service,
		actionsBeforePersist: actionsBeforePersist,
		logger:               logger,
	}
}

func (cE *CustomEndpointController) Put(c *fiber.Ctx) error {
	newRegister := map[string]interface{}{}
	if err := c.BodyParser(&newRegister); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Data provided is invalid",
		})
	}

	endpoint := c.Locals("endpoint").(types.Endpoint)
	err := cE.service.Put(
		newRegister, endpoint.Table, c.Params("id"),
	)
	if err != nil {
		cE.logger.Error(err.Error())
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "ok",
		"id":      c.Params("id"),
	})
}

func (cE *CustomEndpointController) Post(c *fiber.Ctx) error {
	endpoint := c.Locals("endpoint").(types.Endpoint)

	newRegister := map[string]interface{}{}
	c.BodyParser(&newRegister)

	if endpoint.Validations != nil || len(endpoint.Validations) > 0 {
		validationErrors := []string{}
		for _, value := range endpoint.Validations {
			err := utils.IsValid(
				newRegister[value.Field],
				value.Rules,
				value.Field,
			)

			if err != nil {
				validationErrors = append(validationErrors, err.Error())
			}
		}

		if len(validationErrors) > 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": validationErrors,
			})
		}
	}

	if len(endpoint.ActionsBeforePersist) > 0 {
		for _, item := range endpoint.ActionsBeforePersist {
			newRegister[item.Field] = cE.actionsBeforePersist[item.Action].Apply(
				newRegister[item.Field],
			)
		}
	}

	id, err := cE.service.Post(newRegister, endpoint.Table)
	if err != nil {
		cE.logger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"message": "Interval server error",
		})
	}

	return c.JSON(fiber.Map{
		"message": "ok",
		"id":      id,
	})
}

func (cE *CustomEndpointController) Delete(c *fiber.Ctx) error {
	endpoint := c.Locals("endpoint").(types.Endpoint)
	err := cE.service.Delete(endpoint.Table, c.Params("id"))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(204)
}

func (cE *CustomEndpointController) GetById(c *fiber.Ctx) error {
	endpoint := c.Locals("endpoint").(types.Endpoint)
	results, err := cE.service.GetById(endpoint.Table, c.Params("id"))
	if err != nil {
		cE.logger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if len(results) == 0 {
		return c.SendStatus(404)
	}

	return c.JSON(fiber.Map{
		"data": results[0],
	})
}

func (cE *CustomEndpointController) GetAll(c *fiber.Ctx) error {
	endpoint := c.Locals("endpoint").(types.Endpoint)
	var params []interface{}

	if endpoint.Query != "" {
		for _, value := range endpoint.QueryParams {
			queryStringValue := c.Query(value)
			if queryStringValue != "" {
				params = append(params, queryStringValue)
			} else if c.Params(value) != "" {
				params = append(params, c.Params(value))
			}
		}

		if len(params) != len(endpoint.QueryParams) {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf(
					"You need to provide the following params via querystring: %s",
					strings.Join(endpoint.QueryParams, ","),
				),
			})
		}

		results, err := cE.service.Get(endpoint, params)
		if err != nil {
			cE.logger.Error(err.Error())
			return c.Status(500).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		return c.JSON(fiber.Map{
			"data": results,
		})
	}

	results, err := cE.service.Get(endpoint, params)
	if err != nil {
		cE.logger.Error(err.Error())
		return c.Status(500).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	return c.JSON(fiber.Map{
		"data": results,
	})
}
