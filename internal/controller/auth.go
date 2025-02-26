package controller

import (
	"log/slog"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	serviceModule "github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type AuthController struct {
	service serviceModule.AuthService
	logger  *slog.Logger
}

func AuthControllerNew(
	service serviceModule.AuthService,
	logger *slog.Logger,
) *AuthController {
	return &AuthController{
		service: service,
		logger:  logger,
	}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var credential types.Credential
	if err := c.BodyParser(&credential); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Data provided is invalid",
		})
	}

	validate := validator.New()
	err := validate.Struct(credential)
	if err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = "failed validation: " + e.Tag()
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	token, err := a.service.GetToken(credential)
	if err != nil {
		a.logger.Error(err.Error())
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"accessToken": token,
	})
}
