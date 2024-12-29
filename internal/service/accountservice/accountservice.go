package accountservice

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
)

type AccountService struct {
	accRepo     storage.AccountStorage
	ctxDuration time.Duration
	pgxConn     *pgx.Conn
}

func NewAccountService(
	stor storage.AccountStorage,
	ctxDur int,
	pgxC *pgx.Conn,
) *AccountService {
	return &AccountService{
		pgxConn:     pgxC,
		accRepo:     stor,
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

	err := s.accRepo.CreateAccount(&ctx, account)
	if err != nil {
		return fmt.Errorf(
			"CreateAccount->s.repository.CreateAccount: %w",
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
