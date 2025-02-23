package service

import (
	"errors"

	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type EndpointService struct {
	tableService *TableService
	repository   repository.EndpointRepositoryInterface
}

func EndpointServiceNew(
	tableService *TableService,
	repository repository.EndpointRepositoryInterface,
) *EndpointService {
	return &EndpointService{
		tableService: tableService,
		repository:   repository,
	}
}

func (e *EndpointService) Setup() error {
	return e.repository.Setup()
}

func (e *EndpointService) GetAllCreated() (map[string]types.Endpoint, error) {
	return e.repository.GetAllCreated()
}

func (e *EndpointService) Create(endpoint types.Endpoint) (types.Endpoint, error) {
	table, _ := e.tableService.GetByName(endpoint.Table)
	if endpoint.Query == "" && len(table) == 0 {
		return types.Endpoint{}, errors.New("Table is not exists")
	}

	id, _ := e.repository.GetByPath(endpoint.Path)
	if id != 0 {
		return types.Endpoint{}, errors.New("Endpoint already exists")
	}

	id, err := e.repository.Create(endpoint)
	endpoint.ID = int64(id)
	return endpoint, err
}

func (e *EndpointService) Delete(id int64) (string, error) {
	return e.repository.Delete(id)
}
