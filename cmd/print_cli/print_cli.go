package main

import (
	"edy-tanto/printer-pos/internal/print_web_service"
)

func main() {
	name := "Yandi"
	foodnote := "I acknowledge and agree that\nINSURANCE COVERAGE IS UP TO 60 YEARS OLD\nCertain accident risk not guaranteed\nChildren under 12 y/o is not permitted to enter without\nsupervision\nSwimwear is compulsory in water facility\nNo food and drink from outside\nNo drug or weapon/dangerous subtance"

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
			IsPrintAsCopy:    true,
			Foodnote:         &foodnote,
			SalesDetails: []print_web_service.SalesDetail{
				{Item: "Gelang", Qty: 5, Subtotal: 125000},
				{Item: "CashQ", Qty: 1, Subtotal: 500000},
			},
		},
	}

	print_web_service.ExecutePrint(body)
}
