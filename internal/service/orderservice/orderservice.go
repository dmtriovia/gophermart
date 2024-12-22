package orderservice

import (
	"time"

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
