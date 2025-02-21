package service

import "golang.org/x/crypto/bcrypt"

type HashPasswordActionService struct {
}

func HashPasswordActionServiceNew() *HashPasswordActionService {
	return &HashPasswordActionService{}
}

func (h *HashPasswordActionService) Apply(value interface{}) interface{} {
	if value == nil {
		return value
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(value.(string)), 14)
	return string(bytes)
}
