package login

import (
	"fmt"
	"net/http"

	"github.com/dmitrovia/gophermart/internal/service"
)

type Login struct {
	serv service.Service
}

func NewLoginHandler(
	s service.Service,
) *Login {
	return &Login{serv: s}
}

func (h *Login) GetLoginandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	status := http.StatusOK

	fmt.Println(writer)
	fmt.Println(req)

	writer.WriteHeader(status)
}
