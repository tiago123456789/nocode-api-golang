package service

import (
	"context"
	"errors"

	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"go.opentelemetry.io/otel/trace"
)

type EndpointService struct {
	tableService *TableService
	repository   repository.EndpointRepositoryInterface
	tracer       trace.Tracer
}

func EndpointServiceNew(
	tableService *TableService,
	repository repository.EndpointRepositoryInterface,
	tracer trace.Tracer,
) *EndpointService {
	return &EndpointService{
		tableService: tableService,
		repository:   repository,
		tracer:       tracer,
	}
}

func (e *EndpointService) Setup() error {
	return e.repository.Setup()
}

func (e *EndpointService) GetAllCreated(ctx context.Context) (map[string]types.Endpoint, error) {
	ctx, span := e.tracer.Start(ctx, "get-all-service")
	defer span.End()
	return e.repository.GetAllCreated(ctx)
}

func (e *EndpointService) Create(ctx context.Context, endpoint types.Endpoint) (types.Endpoint, error) {
	ctx, span := e.tracer.Start(ctx, "create-endpoint-service")
	defer span.End()

	table, _ := e.tableService.GetByName(ctx, endpoint.Table)
	if endpoint.Query == "" && len(table) == 0 {
		return types.Endpoint{}, errors.New("Table is not exists")
	}

	id, _ := e.repository.GetByPath(ctx, endpoint.Path)
	if id != 0 {
		return types.Endpoint{}, errors.New("Endpoint already exists")
	}

	id, err := e.repository.Create(ctx, endpoint)
	endpoint.ID = int64(id)
	return endpoint, err
}

func (e *EndpointService) Delete(id int64) (string, error) {
	return e.repository.Delete(id)
}
