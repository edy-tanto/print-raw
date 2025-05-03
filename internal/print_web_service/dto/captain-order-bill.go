package dto

type CaptainOrderBillDetail struct {
	Item     string  `json:"item"`
	Qty      uint    `json:"qty"`
	Subtotal float32 `json:"subtotal"`
}

type CaptainOrderBill struct {
	Id                      uint                     `json:"id"`
	Op                      string                   `json:"op"`
	Code                    string                   `json:"code"`
	CustomerName            string                   `json:"customer_name"`
	TableOrRoomNumber       string                   `json:"table_or_room_number"`
	CustomerAdultCount      uint                     `json:"customer_adult_count"`
	CustomerChildCount      uint                     `json:"customer_child_count"`
	TotalQty                uint                     `json:"total_qty"`
	DiscountAmount          float32                  `json:"discount_amount"`
	TotalGross              float32                  `json:"total_gross"`
	TotalNet                float32                  `json:"total_net"`
	GrandTotal              float32                  `json:"grand_total"`
	Date                    string                   `json:"date"`
	IsPrintAsCopy           bool                     `json:"is_print_as_copy"`
	CaptainOrderBillDetails []CaptainOrderBillDetail `json:"captain_order_bill_details"`
}

type PrintCaptainOrderBillRequestBody struct {
	CaptainOrderBill CaptainOrderBill `json:"captain_order_bill"`
}