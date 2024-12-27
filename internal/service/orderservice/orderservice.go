package orderservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrovia/gophermart/internal/models/bizmodels/ordermodel"
	"github.com/dmitrovia/gophermart/internal/storage"
)

type OrderService struct {
	repository  storage.Storage
	ctxDuration time.Duration
}

func NewOrderService(
	stor storage.Storage, ctxDur int,
) *OrderService {
	return &OrderService{
		repository:  stor,
		ctxDuration: time.Duration(ctxDur),
	}
}

func (s *OrderService) OrderIsExist(ident string) (
	bool, *ordermodel.Order, error,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	order, err := s.repository.GetOrder(&ctx, ident)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil, nil
	}

	if err != nil {
		return false, nil, fmt.Errorf(
			"OrderIsExist->GetOrder: %w",
			err)
	}

	return true, order, nil
}

func (s *OrderService) CreateOrder(
	order *ordermodel.Order,
) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		s.ctxDuration)
	defer cancel()

	err := s.repository.CreateOrder(&ctx, order)
	if err != nil {
		return fmt.Errorf(
			"CreateOrder->s.repository.CreateOrder: %w",
			err)
	}

	return nil
}
