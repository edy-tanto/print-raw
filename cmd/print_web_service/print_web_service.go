package main

import (
	"edy-tanto/printer-pos/internal/print_web_service"
	"net/http"
)

func main() {
	mux := print_web_service.NewPrintMuxServer()
	http.ListenAndServe(":8080", mux)
}
