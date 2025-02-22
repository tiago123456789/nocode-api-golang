package utils

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

func IsValidToken(tokenString string) error {
	claims := &types.TokePayload{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
