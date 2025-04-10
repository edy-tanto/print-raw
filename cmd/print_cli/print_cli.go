package main

import (
	"edy-tanto/printer-pos/internal/print_web_service"
)

func main() {
	body := print_web_service.PrintRequestBody{
		Sales: print_web_service.Sales{
			Id:             1,
			Code:           "10000",
			Op:             "Kasir 1",
			Date:           "2025-05-01 18:59:59",
			DiscountAmount: 20000,
			Summary:        605000,
			SalesDetails: []print_web_service.SalesDetail{
				{Item: "Gelang", Qty: 5, Subtotal: 125000},
				{Item: "CashQ", Qty: 1, Subtotal: 500000},
			},
		},
	}

	print_web_service.ExecutePrint(body)
}
