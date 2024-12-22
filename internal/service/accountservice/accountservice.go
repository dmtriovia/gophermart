package accountservice

import (
	"time"

	"github.com/dmitrovia/gophermart/internal/storage"
)

type AccountService struct {
	repository  storage.Storage
	ctxDuration time.Duration
}

func NewAccountService(
	stor storage.Storage, ctxDur int,
) *AccountService {
	return &AccountService{
		repository:  stor,
		ctxDuration: time.Duration(ctxDur),
	}
}
