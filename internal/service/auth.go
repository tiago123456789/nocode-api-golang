package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tiago123456789/nocode-api-golang/internal/repository"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repository repository.AuthRepositoryInterface
}

func AuthServiceNew(repository repository.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (a *AuthService) GetToken(credential types.Credential) (string, error) {
	credentialReturned, err := a.repository.GetByEmail(credential.Email)
	if credentialReturned.Email == "" || len(credentialReturned.Email) == 0 {
		return "", errors.New("Credential is invalid!")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(credentialReturned.Password), []byte(credential.Password),
	)

	if err != nil {
		return "", errors.New("Credential is invalid!")
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.TokePayload{
		Email: credential.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return tokenString, nil
}
