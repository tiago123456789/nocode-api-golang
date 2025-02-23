package main

import (
	"testing"

	"github.com/tiago123456789/nocode-api-golang/internal/utils"
)

func TestIsValid(t *testing.T) {
	err := utils.IsValid(nil, "required", "field")
	if err.Error() == "The field field is required." {
		return
	}

	err = utils.IsValid(nil, "email", "field")
	if err.Error() == "The email format is invalid." {
		return
	}

	err = utils.IsValid(nil, "email", "field")
	if err.Error() == "The email format is invalid." {
		return
	}

	t.Fatalf("Validation broken")
}
