package fakemock

import (
	"github.com/stretchr/testify/mock"
)

type MockTableRepository struct {
	mock.Mock
}

func (t *MockTableRepository) GetByName(name string) (map[string]interface{}, error) {
	args := t.Called(name)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (t *MockTableRepository) GetAll() ([]string, error) {
	args := t.Called()
	return args.Get(0).([]string), args.Error(1)

}

func (t *MockTableRepository) GetColumnsFromTable(table string) ([]string, error) {
	args := t.Called(table)
	return args.Get(0).([]string), args.Error(1)
}
