package controller

import (
	"github.com/gofiber/fiber/v2"
	serviceModule "github.com/tiago123456789/nocode-api-golang/internal/service"
)

type TableController struct {
	service serviceModule.TableService
}

func TableControllerNew(
	service serviceModule.TableService,
) *TableController {
	return &TableController{
		service: service,
	}
}

func (a *TableController) GetAll(c *fiber.Ctx) error {
	results, _ := a.service.GetAll()
	return c.JSON(fiber.Map{
		"data": results,
	})
}

func (a *TableController) GetColumnsFromTable(c *fiber.Ctx) error {
	results, _ := a.service.GetColumnsFromTable(c.Params("table"))
	return c.JSON(fiber.Map{
		"data": results,
	})
}
