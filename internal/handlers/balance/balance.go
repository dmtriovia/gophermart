package balance

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Balance struct {
	serv service.AccountService
}

func NewGetOrderHandler(
	s service.AccountService,
) *Balance {
	return &Balance{serv: s}
}

func (h *Balance) BalanceHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
