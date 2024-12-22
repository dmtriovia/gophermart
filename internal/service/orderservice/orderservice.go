package orderservice

import "github.com/dmitrovia/gophermart/internal/storage"

type OrderService struct {
	repository storage.Storage
}

func NewOrderService(
	stor storage.Storage,
) *OrderService {
	return &OrderService{repository: stor}
}
