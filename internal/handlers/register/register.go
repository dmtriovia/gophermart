package register

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Register struct {
	serv service.Service
}

func NewRegisterHandler(
	s service.Service,
) *Register {
	return &Register{serv: s}
}

func (h *Register) GetRegisterHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
