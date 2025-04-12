package main

import (
	"edy-tanto/printer-pos/internal/print_web_service"
)

func main() {
	body := print_web_service.PrintRequestBody{
		Sales: print_web_service.Sales{
			Id:               1,
			UnitBusinessName: "ParadisQ",
			Code:             "#LC0551",
			Op:               "Kasir 1",
			ConsumerName:     "Yandi",
			PaymentMethod:    "Tunai",
			Date:             "2025-04-11T11:54:47.000Z",
			DiscountAmount:   20000,
			Summary:          605000,
			SalesDetails: []print_web_service.SalesDetail{
				{Item: "Gelang", Qty: 5, Subtotal: 125000},
				{Item: "CashQ", Qty: 1, Subtotal: 500000},
			},
		},
	}

	print_web_service.ExecutePrint(body)
}
