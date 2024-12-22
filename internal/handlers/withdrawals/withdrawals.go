package withdrawals

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdrawals struct {
	serv service.Service
}

func NewWithdrawalsHandler(
	s service.Service,
) *Withdrawals {
	return &Withdrawals{serv: s}
}

func (h *Withdrawals) GetWithdrawalsHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
