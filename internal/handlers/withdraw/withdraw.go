package withdraw

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Withdraw struct {
	serv service.Service
}

func NewWithdrawHandler(
	s service.Service,
) *Withdraw {
	return &Withdraw{serv: s}
}

func (h *Withdraw) GetWithdrawHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
