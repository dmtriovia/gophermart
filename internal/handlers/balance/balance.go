package balance

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Balance struct {
	serv service.Service
}

func NewGetOrderHandler(
	s service.Service,
) *Balance {
	return &Balance{serv: s}
}

func (h *Balance) GetBalanceHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
