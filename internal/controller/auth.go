package controller

import (
	"github.com/gofiber/fiber/v2"
	serviceModule "github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type AuthController struct {
	service serviceModule.AuthService
}

func AuthControllerNew(
	service serviceModule.AuthService,
) *AuthController {
	return &AuthController{
		service: service,
	}
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	var credential types.Credential
	c.BodyParser(&credential)

	token, err := a.service.GetToken(credential)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"accessToken": token,
	})
}
