package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

func IsValid(value interface{}, rules string, field string) error {
	validate := validator.New()
	err := validate.Var(value, rules)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var msg string
			switch e.Tag() {
			case "required":
				msg = fmt.Sprintf("The field %s is required.", strings.ToLower(field))
			case "email":
				msg = "The email format is invalid."
			default:
				msg = fmt.Sprintf("The field %s is invalid input.", strings.ToLower(field))
			}

			return errors.New(msg)
		}
	}

	return nil
}
