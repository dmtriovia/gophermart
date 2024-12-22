package getorders

import (
	"fmt"
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
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
