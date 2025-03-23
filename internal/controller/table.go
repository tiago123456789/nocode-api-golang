package controller

import (
	"github.com/gofiber/fiber/v2"
	serviceModule "github.com/tiago123456789/nocode-api-golang/internal/service"
	"go.opentelemetry.io/otel/trace"
)

type TableController struct {
	service serviceModule.TableService
	tracer  trace.Tracer
}

func TableControllerNew(
	service serviceModule.TableService,
	tracer trace.Tracer,
) *TableController {
	return &TableController{
		service: service,
		tracer:  tracer,
	}
}

func (a *TableController) GetAll(c *fiber.Ctx) error {
	ctx, span := a.tracer.Start(c.UserContext(), "get-tables")
	defer span.End()

	results, _ := a.service.GetAll(ctx)
	return c.JSON(fiber.Map{
		"data": results,
	})
}

func (a *TableController) GetColumnsFromTable(c *fiber.Ctx) error {
	ctx, span := a.tracer.Start(c.UserContext(), "get-columns-of-table")
	defer span.End()

	results, _ := a.service.GetColumnsFromTable(ctx, c.Params("table"))
	return c.JSON(fiber.Map{
		"data": results,
	})
}
