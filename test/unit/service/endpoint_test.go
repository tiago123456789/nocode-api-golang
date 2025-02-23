package main

import (
	"testing"

	fakemock "github.com/tiago123456789/nocode-api-golang/fakeMock"
	"github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

func TestEndpointServiceCreateTableNotFound(t *testing.T) {
	endpointRepository := new(fakemock.MockEndpointRepository)
	tableRepository := new(fakemock.MockTableRepository)
	tableService := service.TableServiceNew(tableRepository)
	endpointService := service.EndpointServiceNew(
		tableService,
		endpointRepository,
	)

	tableName := "fake_table"
	tableRepository.
		On("GetByName", tableName).Return(map[string]interface{}{}, nil).
		Once()
	_, err := endpointService.Create(types.Endpoint{
		Table: "fake_table",
		Query: "",
	})

	if err.Error() == "Table is not exists" {
		return
	}

	t.Fatal("Test failed")
}

func TestEndpointServiceEndpointAlreadyExists(t *testing.T) {
	endpointRepository := new(fakemock.MockEndpointRepository)
	tableRepository := new(fakemock.MockTableRepository)
	tableService := service.TableServiceNew(tableRepository)
	endpointService := service.EndpointServiceNew(
		tableService,
		endpointRepository,
	)

	tableName := "fake_table"
	path := "/fake-table"
	tableRepository.On("GetByName", tableName).
		Return(map[string]interface{}{
			"name": tableName,
		}, nil).
		Once()

	endpointRepository.On("GetByPath", path).Return(10, nil)
	_, err := endpointService.Create(types.Endpoint{
		Table: "fake_table",
		Query: "",
		Path:  path,
	})

	if err.Error() == "Endpoint already exists" {
		return
	}

	t.Fatal("Test failed")
}

func TestEndpointServiceEndpointSuccess(t *testing.T) {
	endpointRepository := new(fakemock.MockEndpointRepository)
	tableRepository := new(fakemock.MockTableRepository)
	tableService := service.TableServiceNew(tableRepository)
	endpointService := service.EndpointServiceNew(
		tableService,
		endpointRepository,
	)

	tableName := "fake_table"
	path := "/fake-table"
	tableRepository.On("GetByName", tableName).
		Return(map[string]interface{}{
			"name": tableName,
		}, nil).
		Once()

	fakeEnpoint := types.Endpoint{
		Table: "fake_table",
		Query: "",
		Path:  path,
	}
	endpointRepository.On("GetByPath", path).Return(0, nil).Once()
	endpointRepository.On("Create", fakeEnpoint).Return(0, nil).Once()

	_, err := endpointService.Create(fakeEnpoint)

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
		return
	}

	return

}
