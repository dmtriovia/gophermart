package getorder

import (
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type GetOrders struct {
	serv service.OrderService
}

func NewGetOrderHandler(
	s service.OrderService,
) *GetOrders {
	return &GetOrders{serv: s}
}

func (h *GetOrders) GetOrderHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
