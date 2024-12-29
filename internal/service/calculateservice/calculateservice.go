package calculateservice

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/accountmodel"
	"github.com/dmitrovia/gophermart/internal/storage"
	"github.com/jackc/pgx/v5"
)

type CalculateService struct {
	accRepo     storage.AccountStorage
	accOrder    storage.OrderStorage
	ctxDuration time.Duration
	pgxConn     *pgx.Conn
}

func NewCalculateService(
	stor storage.AccountStorage,
	ctxDur int,
	pgxC *pgx.Conn,
) *CalculateService {
	return &CalculateService{
		pgxConn:     pgxC,
		accRepo:     stor,
		ctxDuration: time.Duration(ctxDur),
	}
}

func (s *CalculateService) CalculatePoints(
	acc *accountmodel.Account,
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

	_, err = s.accRepo.MinusPointsByID(
		&ctx, acc.GetID(), points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->MinusPointsByID: %w",
			err)
	}

	_, err = s.accRepo.PlusWithdrawnByID(
		&ctx, acc.GetID(), points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->PlusWithdrawnByID: %w", err)
	}

	_, err = s.accOrder.PlusPointsWriteOffByID(
		&ctx, acc.GetID(), points)
	if err != nil {
		return fmt.Errorf(
			"CalculatePoints->PlusPointsWriteOffByID: %w",
			err)
	}

	err = tranz.Commit(ctx)
	if err != nil {
		return fmt.Errorf("CalculatePoints->tranz.Commit %w", err)
	}

	return nil
}
