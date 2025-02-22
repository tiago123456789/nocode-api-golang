package types

import "github.com/golang-jwt/jwt/v5"

type TokePayload struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
