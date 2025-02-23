package main

import (
	"testing"

	fakemock "github.com/tiago123456789/nocode-api-golang/fakeMock"
	"github.com/tiago123456789/nocode-api-golang/internal/service"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

func TestAuthServiceGetTokenEmailInvalid(t *testing.T) {
	fakeRepository := new(fakemock.MockAuthRepository)
	fakeRepository.On("GetByEmail", "test@gmail.com").Return(types.Credential{}, nil)
	authService := service.AuthServiceNew(fakeRepository)
	_, err := authService.GetToken(types.Credential{
		Email:    "test@gmail.com",
		Password: "1234456",
	})

	if err.Error() == "Credential is invalid!" {
		return
	}

	t.Fatalf("Test failed")
}

func TestAuthServiceGetTokenPasswordInvalid(t *testing.T) {
	fakeRepository := new(fakemock.MockAuthRepository)
	passwordHash := service.HashPasswordActionServiceNew().Apply("abcabc")
	fakeRepository.On("GetByEmail", "test@gmail.com").Return(types.Credential{
		Email:    "test@gmail.com",
		Password: passwordHash.(string),
	}, nil)
	authService := service.AuthServiceNew(fakeRepository)
	_, err := authService.GetToken(types.Credential{
		Email:    "test@gmail.com",
		Password: "1234456",
	})

	if err.Error() == "Credential is invalid!" {
		return
	}

	t.Fatalf("Test failed")
}

func TestAuthServiceGetTokenSuccess(t *testing.T) {
	fakeRepository := new(fakemock.MockAuthRepository)
	passwordHash := service.HashPasswordActionServiceNew().Apply("abcabc")
	fakeRepository.On("GetByEmail", "test@gmail.com").Return(types.Credential{
		Email:    "test@gmail.com",
		Password: passwordHash.(string),
	}, nil).Once()
	authService := service.AuthServiceNew(fakeRepository)
	_, err := authService.GetToken(types.Credential{
		Email:    "test@gmail.com",
		Password: "abcabc",
	})

	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	}

	return

}
