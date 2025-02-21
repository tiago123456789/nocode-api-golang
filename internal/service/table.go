package service

import (
	"github.com/tiago123456789/nocode-api-golang/internal/repository"
)

type TableService struct {
	repository repository.TableRepositoryInterface
}

func TableServiceNew(
	repository repository.TableRepositoryInterface,
) *TableService {
	return &TableService{
		repository: repository,
	}
}

func (t *TableService) GetColumnsFromTable(table string) ([]string, error) {
	return t.repository.GetColumnsFromTable(table)
}

func (t *TableService) GetAll() ([]string, error) {
	return t.repository.GetAll()

}

func (t *TableService) GetByName(name string) (map[string]interface{}, error) {
	return t.repository.GetByName(name)
}
