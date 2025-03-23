package service

import (
	"context"

	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"go.opentelemetry.io/otel/trace"
)

type TableService struct {
	repository repository.TableRepositoryInterface
	tracer     trace.Tracer
}

func TableServiceNew(
	repository repository.TableRepositoryInterface,
	tracer trace.Tracer,
) *TableService {
	return &TableService{
		repository: repository,
		tracer:     tracer,
	}
}

func (t *TableService) GetColumnsFromTable(ctx context.Context, table string) ([]string, error) {
	ctx, span := t.tracer.Start(ctx, "get-columns-of-table-service")
	defer span.End()
	return t.repository.GetColumnsFromTable(ctx, table)
}

func (t *TableService) GetAll(ctx context.Context) ([]string, error) {
	ctx, span := t.tracer.Start(ctx, "get-tables-service")
	defer span.End()
	return t.repository.GetAll(ctx)

}

func (t *TableService) GetByName(ctx context.Context, name string) (map[string]interface{}, error) {
	ctx, span := t.tracer.Start(ctx, "get-by-name-service")
	defer span.End()

	return t.repository.GetByName(ctx, name)
}
