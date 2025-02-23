package main

import (
	"testing"

	fakemock "github.com/tiago123456789/nocode-api-golang/fakeMock"
	"github.com/tiago123456789/nocode-api-golang/internal/service"
)

func TestCustomEndpointUpdateRegisterNotFound(t *testing.T) {
	fakeRepository := new(fakemock.MockCustomEndpointRepository)
	fakeRepository.On("GetById", "fake_table", "1").
		Return([]map[string]interface{}{}, nil).
		Once()
	customEndpointService := service.CustomEndpointServiceNew(fakeRepository)

	err := customEndpointService.Put(map[string]interface{}{}, "fake_table", "1")
	if err.Error() == "Register not found" {
		return
	}

	t.Fatalf("Test failed")
}

func TestCustomEndpointDeleteRegisterNotFound(t *testing.T) {
	fakeRepository := new(fakemock.MockCustomEndpointRepository)
	fakeRepository.On("GetById", "fake_table", "1").
		Return([]map[string]interface{}{}, nil).
		Once()
	customEndpointService := service.CustomEndpointServiceNew(fakeRepository)

	err := customEndpointService.Delete("fake_table", "1")
	if err.Error() == "Register not found" {
		return
	}

	t.Fatalf("Test failed")
}

func TestCustomEndpointDeleteSuccess(t *testing.T) {
	var items []map[string]interface{}
	items = append(items, map[string]interface{}{
		"id": 1,
	})
	fakeRepository := new(fakemock.MockCustomEndpointRepository)

	fakeRepository.On("GetById", "fake_table", "1").
		Return(items, nil).
		Once()
	fakeRepository.On("Delete", "fake_table", "1").
		Return(nil).
		Once()

	customEndpointService := service.CustomEndpointServiceNew(fakeRepository)

	err := customEndpointService.Delete("fake_table", "1")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	return
}

func TestCustomEndpointUpdateSuccess(t *testing.T) {
	var items []map[string]interface{}
	items = append(items, map[string]interface{}{
		"id": 1,
	})
	fakeRepository := new(fakemock.MockCustomEndpointRepository)

	fakeRepository.On("GetById", "fake_table", "1").
		Return(items, nil).
		Once()
	fakeRepository.On("Update", map[string]interface{}{}, "fake_table", "1").
		Return(nil).
		Once()

	customEndpointService := service.CustomEndpointServiceNew(fakeRepository)

	err := customEndpointService.Put(map[string]interface{}{}, "fake_table", "1")
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	return
}
