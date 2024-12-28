package accountservice

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/storage"
)

type AccountService struct {
	repository  storage.AccountStorage
	ctxDuration time.Duration
}

func NewAccountService(
	stor storage.AccountStorage, ctxDur int,
) *AccountService {
	return &AccountService{
		repository:  stor,
		ctxDuration: time.Duration(ctxDur),
	}
}

func (s *AccountService) CreateAccount(
	account *accountmodel.Account,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	err := s.repository.CreateAccount(&ctx, account)
	if err != nil {
		return fmt.Errorf(
			"CreateAccount->s.repository.CreateAccount: %w",
			err)
	}

	return nil
}
