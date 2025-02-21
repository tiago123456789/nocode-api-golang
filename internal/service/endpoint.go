package service

import (
	"database/sql"
	"errors"

	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type EndpointService struct {
	db           *sql.DB
	tableService *TableService
	repository   repository.EndpointRepositoryInterface
}

func EndpointServiceNew(
	db *sql.DB,
	tableService *TableService,
	repository repository.EndpointRepositoryInterface,
) *EndpointService {
	return &EndpointService{
		db:           db,
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
	if endpoint.Query == "" && table == nil {
		return types.Endpoint{}, errors.New("Table is not exists")
	}

	id, _ := e.repository.GetByPath(endpoint.Path)
	if id != 0 {
		return types.Endpoint{}, errors.New("Endpoint already exists")
	}

	_, err := e.repository.Create(endpoint)
	return endpoint, err
}
