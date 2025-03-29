package service

import (
	"context"
	"errors"

	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"go.opentelemetry.io/otel/trace"
)

type CustomEndpointService struct {
	repository repository.CustomEndpointInterface
	tracer     trace.Tracer
}

func CustomEndpointServiceNew(
	repository repository.CustomEndpointInterface,
	tracer trace.Tracer,

) *CustomEndpointService {
	return &CustomEndpointService{
		repository: repository,
		tracer:     tracer,
	}
}

func (c *CustomEndpointService) Put(
	ctx context.Context,
	dataToModify map[string]interface{},
	table string,
	idRegister string,
) error {
	ctx, span := c.tracer.Start(ctx, "update-register-service")
	defer span.End()
	items, _ := c.repository.GetById(context.TODO(), table, idRegister)

	if len(items) == 0 {
		return errors.New("Register not found")
	}

	return c.repository.Update(ctx, dataToModify, table, idRegister)
}

func (c *CustomEndpointService) Post(
	ctx context.Context,
	newRegister map[string]interface{},
	table string,
) (int64, error) {
	ctx, span := c.tracer.Start(ctx, "create-register-service")
	defer span.End()
	return c.repository.Create(ctx, newRegister, table)
}

func (c *CustomEndpointService) Delete(
	ctx context.Context,
	table string,
	id string,
) error {
	ctx, span := c.tracer.Start(ctx, "delete-register-service")
	defer span.End()
	items, _ := c.repository.GetById(ctx, table, id)

	if len(items) == 0 {
		return errors.New("Register not found")
	}

	return c.repository.Delete(ctx, table, id)
}

func (c *CustomEndpointService) GetById(
	ctx context.Context,
	table string,
	id string,
) ([]map[string]interface{}, error) {
	ctx, span := c.tracer.Start(ctx, "get-by-id-register-service")
	defer span.End()
	return c.repository.GetById(ctx, table, id)
}

func (c *CustomEndpointService) Get(
	ctx context.Context,
	endpoint types.Endpoint,
	params []interface{},
) ([]map[string]interface{}, error) {
	ctx, span := c.tracer.Start(ctx, "get-all-register-service")
	defer span.End()

	if endpoint.Query == "" {
		return c.repository.GetAll(ctx, endpoint)
	}

	return c.repository.GetAllByCustomQuery(ctx, endpoint, params)
}
