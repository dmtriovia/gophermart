package setorders

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type SetOrders struct {
	serv service.Service
}

func NewSetOrderHandler(
	s service.Service,
) *SetOrders {
	return &SetOrders{serv: s}
}

func (h *SetOrders) SetOrderHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
