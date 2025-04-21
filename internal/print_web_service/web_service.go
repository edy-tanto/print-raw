package print_web_service

import (
	"net/http"
)

func NewPrintMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", &PrintHandler{})
	mux.Handle("/cash-refund", &PrintCashRefundHandler{})
	mux.Handle("/kitchen", &PrintKitchenHandler{})
	mux.Handle("/table-check", &PrintTableCheckHandler{})

	return mux
}
