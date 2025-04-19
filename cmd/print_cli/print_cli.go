package main

import (
	"edy-tanto/printer-pos/internal/print_web_service"
)

func main() {
	name := "Yandi"

	body := print_web_service.PrintRequestBody{
		Sales: print_web_service.Sales{
			Id:               1,
			UnitBusinessName: "ParadisQ",
			Code:             "#LC0551",
			Op:               "Kasir 1",
			CustomerName:     &name,
			PaymentMethod:    "Tunai",
			Date:             "2025-04-11T11:54:47.000",
			GrandTotal:       607000,
			SalesDetails: []print_web_service.SalesDetail{
				{Item: "Gelang Gelang1 Gelang2 Gelang3 Gelang4 Gelang5 Gelang6 Gelang7", Qty: 5, Subtotal: 125000},
				{Item: "CashQ", Qty: 1, Subtotal: 500000},
			},
		},
	}

	print_web_service.ExecutePrint(body)
}
