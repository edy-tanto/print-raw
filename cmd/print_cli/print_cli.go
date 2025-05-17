package main

import (
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"edy-tanto/printer-pos/internal/print_web_service/printer"
)

func main() {
	footnote := "I acknowledge and agree that\nINSURANCE COVERAGE IS UP TO 60 YEARS OLD. Certain accident risk\nnot guaranteed Children under 12 y/o is not permitted to enter\nwithout supervision. Swimwear is compulsory in water facility.\nNo food and drink from outside. No drug or\nweapon/dangerous subtance"

	body := dto.PrintRequestBody{
		Sales: dto.Sales{
			Id:               1,
			UnitBusinessName: "ParadisQ",
			Code:             "#LC0551",
			Op:               "Kasir 1",
			CustomerName:     "Yandi",
			PaymentMethod:    "Tunai",
			Date:             "2025-04-11T11:54:47.000",
			GrandTotal:       607000,
			IsPrintAsCopy:    false,
			Footnote:         footnote,
			FootnoteAlign:    "CENTER",
			SalesDetails: []dto.SalesDetail{
				{Item: "Gelang", Qty: 5, TotalFinal: 125000},
				{Item: "CashQ", Qty: 1, TotalFinal: 500000},
			},
		},
	}

	printer.ExecutePrint(body)
}
