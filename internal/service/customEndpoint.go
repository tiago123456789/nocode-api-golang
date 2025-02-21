package service

import (
	"database/sql"
	"errors"

	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type CustomEndpointService struct {
	db         *sql.DB
	repository repository.CustomEndpointInterface
}

func CustomEndpointServiceNew(
	db *sql.DB,
	repository repository.CustomEndpointInterface,
) *CustomEndpointService {
	return &CustomEndpointService{
		db:         db,
		repository: repository,
	}
}

func (c *CustomEndpointService) Put(
	dataToModify map[string]interface{},
	table string,
	idRegister string,
) error {
	items, _ := c.repository.GetById(table, idRegister)

	if len(items) == 0 {
		return errors.New("Register not found")
	}

	return c.repository.Update(dataToModify, table, idRegister)
}

func (c *CustomEndpointService) Post(
	newRegister map[string]interface{},
	table string,
) (int64, error) {
	return c.repository.Create(newRegister, table)
}

func (c *CustomEndpointService) Delete(
	table string,
	id string,
) error {
	items, _ := c.repository.GetById(table, id)

	if len(items) == 0 {
		return errors.New("Register not found")
	}

	return c.repository.Delete(table, id)
}

func (c *CustomEndpointService) GetById(
	table string,
	id string,
) ([]map[string]interface{}, error) {
	return c.repository.GetById(table, id)
}

func (c *CustomEndpointService) Get(
	endpoint types.Endpoint,
	params []interface{},
) ([]map[string]interface{}, error) {
	if endpoint.Query == "" {
		return c.repository.GetAll(endpoint)
	}

	return c.repository.GetAllByCustomQuery(endpoint, params)
}
