package dto

type TableCheckDetail struct {
	Item string `json:"item"`
	Qty  uint   `json:"qty"`
}

type TableCheck struct {
	Id                 uint               `json:"id"`
	WaitressName       string             `json:"waitress_name"`
	Op                 string             `json:"op"`
	Code               string             `json:"code"`
	CustomerName       string             `json:"customer_name"`
	TableOrRoomNumber  string             `json:"table_or_room_number"`
	CustomerAdultCount uint               `json:"customer_adult_count"`
	CustomerChildCount uint               `json:"customer_child_count"`
	TotalQty           uint               `json:"total_qty"`
	Date               string             `json:"date"`
	IsPrintAsCopy      bool               `json:"is_print_as_copy"`
	TableCheckDetails  []TableCheckDetail `json:"table_check_details"`
}

type PrintTableCheckRequestBody struct {
	TableCheck  TableCheck `json:"table_check"`
	PrinterName string     `json:"printer_name"`
}
