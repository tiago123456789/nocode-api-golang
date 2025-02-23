package fakemock

import (
	"github.com/stretchr/testify/mock"
	"github.com/tiago123456789/nocode-api-golang/internal/types"
)

type MockAuthRepository struct {
	mock.Mock
}

func (a *MockAuthRepository) GetByEmail(email string) (types.Credential, error) {
	args := a.Called(email)
	return args.Get(0).(types.Credential), args.Error(1)
}
