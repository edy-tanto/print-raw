package print_web_service

import (
	"edy-tanto/printer-pos/internal/print_web_service/handler"
	"net/http"
)

func NewPrintMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", &handler.PrintHandler{})
	mux.Handle("/cash-refund", &handler.PrintCashRefundHandler{})
	mux.Handle("/kitchen", &handler.PrintKitchenHandler{})
	mux.Handle("/table-check", &handler.PrintTableCheckHandler{})
	mux.Handle("/captain-order-bill", &handler.PrintCaptainOrderBillHandler{})
	mux.Handle("/captain-order-invoice", &handler.PrintCaptainOrderInvoiceHandler{})

	return mux
}
