package accountservice

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accounthistorymodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/storage"
)

type AccountService struct {
	accRepo     storage.AccountStorage
	ctxDuration time.Duration
}

func NewAccountService(
	stor storage.AccountStorage,
	ctxDur time.Duration,
) *AccountService {
	return &AccountService{
		accRepo:     stor,
		ctxDuration: ctxDur,
	}
}

func (s *AccountService) CreateAccount(
	account *accountmodel.Account,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	err := s.accRepo.CreateAccount(&ctx, account)
	if err != nil {
		return fmt.Errorf(
			"CreateAccount->s.repository.CreateAccount: %w",
			err)
	}

	return nil
}

func (s *AccountService) CreateAccountHistory(
	accHist *accounthistorymodel.AccountHistory,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	err := s.accRepo.CreateAccountHistory(&ctx, accHist)
	if err != nil {
		return fmt.Errorf(
			"CreateAccountHistory->accRepo.CreateAccountHistory: %w",
			err)
	}

	return nil
}

func (s *AccountService) GetAccountByClient(
	clientID int32,
) (*accountmodel.Account, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	acc, err := s.accRepo.GetAccountByClient(&ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf(
			"GetAccountByClient->repo.GetAccountByClient: %w",
			err)
	}

	return acc, nil
}

func (s *AccountService) GetAccountHistoryByClient(
	clientID int32,
) (*[]accounthistorymodel.AccountHistory,
	*[]error,
	error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	accHists,
		errors,
		err := s.accRepo.GetAccountHistoryByClient(
		&ctx,
		clientID)

	return accHists, errors, fmt.Errorf(
		"GetAccountHistoryByClient->GetAccHistByClient: %w",
		err)
}
