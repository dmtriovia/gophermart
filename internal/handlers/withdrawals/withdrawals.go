package withdrawals

import (
	"net/http"

	"github.com/dmitrovia/gophermart/internal/models/handlerattr/withdrawalsattr"
	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdrawals struct {
	orderService service.OrderService
	attr         *withdrawalsattr.WithdrawalsAttr
}

func NewWithdrawalsHandler(
	ords service.OrderService,
	inAttr *withdrawalsattr.WithdrawalsAttr,
) *Withdrawals {
	return &Withdrawals{
		orderService: ords,
		attr:         inAttr,
	}
}

func (h *Withdrawals) WithdrawalsHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	writer.WriteHeader(http.StatusOK)
}
