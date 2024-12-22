package accountservice

import (
	"github.com/dmitrovia/gophermart/internal/storage"
)

type AccountService struct {
	repository storage.Storage
}

func NewAccountService(
	stor storage.Storage,
) *AccountService {
	return &AccountService{repository: stor}
}
