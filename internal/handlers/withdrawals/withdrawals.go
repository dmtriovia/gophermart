package withdrawals

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdrawals struct {
	serv service.AccountService
}

func NewWithdrawalsHandler(
	s service.AccountService,
) *Withdrawals {
	return &Withdrawals{serv: s}
}

func (h *Withdrawals) WithdrawalsHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
