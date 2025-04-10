package print_web_service

import (
	"net/http"
)

func NewPrintMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", &PrintHandler{})

	return mux
}
