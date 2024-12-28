package withdraw

import (
	"net/http"

	"github.com/dmitrovia/gophermart/internal/models/handlerattr/withdrawattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdraw struct {
	serv service.AccountService
	attr *withdrawattr.WithdrawAttr
}

func NewWithdrawHandler(
	s service.AccountService,
	inAttr *withdrawattr.WithdrawAttr,
) *Withdraw {
	return &Withdraw{serv: s, attr: inAttr}
}

func (h *Withdraw) WithdrawHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	writer.WriteHeader(http.StatusOK)
}
