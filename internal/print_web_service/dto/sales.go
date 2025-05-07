package dto

type SalesDetail struct {
	Item     string  `json:"item"`
	Qty      uint    `json:"qty"`
	Subtotal float32 `json:"subtotal"`
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
	Footnote         string        `json:"footnote"`
	FootnoteAlign    string        `json:"footnote_align"`
	GrandTotal       float32       `json:"grand_total"`
	SalesDetails     []SalesDetail `json:"sales_details"`
}

type PrintRequestBody struct {
	Sales Sales `json:"sales"`
}
