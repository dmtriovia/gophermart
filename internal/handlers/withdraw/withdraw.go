package withdraw

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdraw struct {
	serv service.AccountService
}

func NewWithdrawHandler(
	s service.AccountService,
) *Withdraw {
	return &Withdraw{serv: s}
}

func (h *Withdraw) WithdrawHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
