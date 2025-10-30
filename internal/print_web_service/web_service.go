package print_web_service

import (
	"edy-tanto/printer-pos/internal/print_web_service/handler"
	"net/http"
)

func NewPrintMuxServer() *http.ServeMux {
	mux := http.NewServeMux()

	// default service to print receipt for POS
	mux.Handle("/", &handler.PrintHandler{})

	// printer list
	mux.Handle("/printers", &handler.PrinterListHandler{})

	// cash refund at waterpark
	mux.Handle("/cash-refund", &handler.PrintCashRefundHandler{})

	// print for Patio and Dimsum
	mux.Handle("/kitchen", &handler.PrintKitchenHandler{})
	mux.Handle("/kitchen-eth", &handler.PrintKitchenEthHandler{})
	mux.Handle("/table-check", &handler.PrintTableCheckHandler{})
	mux.Handle("/captain-order-bill", &handler.PrintCaptainOrderBillHandler{})
	mux.Handle("/captain-order-invoice", &handler.PrintCaptainOrderInvoiceHandler{})

	mux.Handle("/shift/report/cash-count", &handler.PrintShiftReportCashCountHandler{})

	return mux
}
