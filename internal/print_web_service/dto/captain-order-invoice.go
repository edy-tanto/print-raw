package dto

type CaptainOrderInvoice struct {
	SalesId            uint          `json:"sales_id"`
	Op                 string        `json:"op"`
	CustomerName       string        `json:"customer_name"`
	TableOrRoomNumber  string        `json:"table_or_room_number"`
	CustomerAdultCount uint          `json:"customer_adult_count"`
	CustomerChildCount uint          `json:"customer_child_count"`
	TotalQty           uint          `json:"total_qty"`
	DiscountAmount     float32       `json:"discount_amount"`
	TotalGross         float32       `json:"total_gross"`
	TotalNet           float32       `json:"total_net"`
	GrandTotal         float32       `json:"grand_total"`
	PaymentMethod      string        `json:"payment_method"`
	Date               string        `json:"date"`
	Note               string        `json:"note"`
	IsPrintAsCopy      bool          `json:"is_print_as_copy"`
	SalesDetails       []SalesDetail `json:"sales_details"`
}

type PrintCaptainOrderInvoiceRequestBody struct {
	CaptainOrderInvoice CaptainOrderInvoice `json:"captain_order_invoice"`
}
