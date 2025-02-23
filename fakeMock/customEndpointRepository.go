package fakemock

import (
	"github.com/stretchr/testify/mock"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type MockCustomEndpointRepository struct {
	mock.Mock
}

func (c *MockCustomEndpointRepository) GetById(
	table string,
	id string,
) ([]map[string]interface{}, error) {
	args := c.Called(table, id)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (c *MockCustomEndpointRepository) GetAll(
	endpoint types.Endpoint,
) ([]map[string]interface{}, error) {
	args := c.Called(endpoint)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (c *MockCustomEndpointRepository) GetAllByCustomQuery(
	endpoint types.Endpoint,
	params []interface{},
) ([]map[string]interface{}, error) {
	args := c.Called(endpoint, params)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

func (c *MockCustomEndpointRepository) Delete(
	table string,
	id string,
) error {
	args := c.Called(table, id)
	return args.Error(0)
}

func (c *MockCustomEndpointRepository) Create(
	newRegister map[string]interface{},
	table string,
) (int64, error) {
	args := c.Called(newRegister, table)
	return args.Get(0).(int64), args.Error(1)
}

func (c *MockCustomEndpointRepository) Update(
	newRegister map[string]interface{},
	table string,
	idRegister string,
) error {
	args := c.Called(newRegister, table, idRegister)
	return args.Error(0)
}
