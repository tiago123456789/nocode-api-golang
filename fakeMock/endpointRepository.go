package fakemock

import (
	"github.com/stretchr/testify/mock"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type MockEndpointRepository struct {
	mock.Mock
}

func (e *MockEndpointRepository) GetByPath(path string) (int, error) {
	args := e.Called(path)
	return args.Get(0).(int), args.Error(1)
}

func (e *MockEndpointRepository) Create(endpoint types.Endpoint) (int, error) {
	args := e.Called(endpoint)
	return args.Get(0).(int), args.Error(1)
}

func (e *MockEndpointRepository) GetAllCreated() (map[string]types.Endpoint, error) {
	args := e.Called()
	return args.Get(0).(map[string]types.Endpoint), args.Error(1)
}

func (e *MockEndpointRepository) Setup() error {
	args := e.Called()
	return args.Error(0)
}

func (e *MockEndpointRepository) Delete(id int64) (string, error) {
	args := e.Called(id)
	return args.Get(0).(string), args.Error(1)
}
