package main

import (
	"testing"

	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

func TestIsValidToken(t *testing.T) {
	err := utils.IsValidToken("adsfjaklsdlfasd")

	if err != nil {
		return
	}

	t.Fatalf("Token can't be valid")
}
