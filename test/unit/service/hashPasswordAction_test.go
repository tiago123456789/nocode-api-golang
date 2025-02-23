package main

import (
	"testing"

	"github.com/tiago123456789/nocode-api-golang/internal/service"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHashPasswordActionServiceApplySuccess(t *testing.T) {
	received := service.HashPasswordActionServiceNew().Apply("teste teste")
	err := bcrypt.CompareHashAndPassword([]byte(received.(string)), []byte("teste teste"))
	if err != nil {
		t.Fatalf("Error: hash generate is invalid")
	}

	return
}
