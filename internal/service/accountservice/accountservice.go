package accountservice

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
)

type AccountService struct {
	accRepo     storage.AccountStorage
	accOrder    storage.OrderStorage
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

func (s *AccountService) CalculatePoints(
	acc *accountmodel.Account,
	order *ordermodel.Order,
	points float32,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	tranz, err := s.pgxConn.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->s.pgxConn.Begin %w", err)
	}

	newValuePoints := acc.GetPoints() - points
	newValueWithdrawn := acc.GetWithdrawn() + points
	newValueOrderPWO := order.GetpointsWriteOff() + points

	isUpt, err := s.accRepo.UpdateAccountPointsByID(
		&ctx, acc.GetID(), newValuePoints)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->UpdateAccountPointsByID: %w",
			err)
	}

	if isUpt {
		acc.SetPoints(newValuePoints)
	}

	isUpt, err = s.accRepo.UpdateAccountWithdrawnByID(
		&ctx, acc.GetID(), newValueWithdrawn)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->UpdateAccountPointsByID: %w", err)
	}

	if isUpt {
		acc.SetWithdrawn(newValueWithdrawn)
	}

	isUpt, err = s.accOrder.UpdateOrderPointsWriteOffByID(
		&ctx, acc.GetID(), newValueWithdrawn)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->UpdateOrderPointsWriteOffByID: %w",
			err)
	}

	if isUpt {
		order.SetpointsWriteOff(newValueOrderPWO)
	}

	err = tranz.Commit(ctx)
	if err != nil {
		return fmt.Errorf("CalculatePoints->tranz.Commit %w", err)
	}

	return nil
}
