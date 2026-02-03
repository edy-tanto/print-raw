package dto

import "encoding/json"

// FootnoteItem is one footnote line with optional per-item alignment.
type FootnoteItem struct {
	Footnote  string `json:"footnote"`
	Alignment string `json:"alignment"`
}

// FootnoteList supports unmarshaling from either:
// - new format: [{"footnote": "text", "alignment": "LEFT"}, ...]
// - old format: ["text1", "text2"] (alignment from footnote_align or CENTER)
type FootnoteList []FootnoteItem

// UnmarshalJSON decodes footnote as []FootnoteItem or, on failure, []string (legacy).
func (l *FootnoteList) UnmarshalJSON(data []byte) error {
	var items []FootnoteItem
	if err := json.Unmarshal(data, &items); err == nil {
		*l = items
		return nil
	}
	var strs []string
	if err := json.Unmarshal(data, &strs); err != nil {
		return err
	}
	*l = make(FootnoteList, len(strs))
	for i, s := range strs {
		(*l)[i] = FootnoteItem{Footnote: s, Alignment: ""}
	}
	return nil
}

type SalesDetail struct {
	Item               string  `json:"item"`
	Qty                uint    `json:"qty"`
	TotalFinal         float32 `json:"total_final"`
	SubtotalWithoutTax float32 `json:"subtotal_with_tax"`
}

type Sales struct {
	Id               uint          `json:"id"`
	UnitBusinessName string        `json:"unit_business_name"`
	Code             string        `json:"code"`
	Op               string        `json:"op"`
	CustomerName     string        `json:"customer_name"`
	TableNumber      string        `json:"table_number"`
	PaymentMethod    string        `json:"payment_method"`
	Date             string        `json:"date"`
	IsPrintAsCopy    bool          `json:"is_print_as_copy"`
	Footnote         FootnoteList   `json:"footnote"`
	FootnoteAlign    string        `json:"footnote_align"`
	GrandTotal       float32       `json:"grand_total"`
	CreditCardCharge float32       `json:"credit_card_charge"`
	CashQBalance     *float32      `json:"cash_q_balance"`
	SalesDetails     []SalesDetail `json:"sales_details"`
}

type PrintRequestBody struct {
	Sales       Sales  `json:"sales"`
	PrinterName string `json:"printer_name"`
}
