package notallowed

import (
	"fmt"
	"net/http"
)

type NotAllowed struct{}

func (h NotAllowed) ServeHTTP(
	rw http.ResponseWriter, r *http.Request,
) {
	MethodNotAllowedHandler(rw, r)
}

func MethodNotAllowedHandler(
	rw http.ResponseWriter, _ *http.Request,
) {
	rw.WriteHeader(http.StatusNotFound)

	Body := "Method not allowed!\n"
	fmt.Fprintf(rw, "%s", Body)
}
