package repository

import (
	"database/sql"

	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type AuthRepositoryInterface interface {
	GetByEmail(email string) (types.Credential, error)
}

type AuthRespotiory struct {
	db *sql.DB
}

func AuthRepositoryNew(db *sql.DB) *AuthRespotiory {
	return &AuthRespotiory{
		db: db,
	}
}

func (a *AuthRespotiory) GetByEmail(email string) (types.Credential, error) {
	var credentialReturned types.Credential
	err := a.db.QueryRow("SELECT email, password FROM auth where email = $1", email).Scan(
		&credentialReturned.Email, &credentialReturned.Password,
	)

	if err != nil {
		return types.Credential{}, err
	}

	return credentialReturned, err
}
