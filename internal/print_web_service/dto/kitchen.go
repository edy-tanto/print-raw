package dto

type KitchenDetail struct {
	Item string `json:"item"`
	Qty  uint   `json:"qty"`
}

type Kitchen struct {
	Id                  uint            `json:"id"`
	WaitressName        string          `json:"waitress_name"`
	Op                  string          `json:"op"`
	Code                string          `json:"code"`
	Outlet              string          `json:"outlet"`
	PrinterLocationName string          `json:"printer_location_name"`
	CustomerName        string          `json:"customer_name"`
	TableOrRoomNumber   string          `json:"table_or_room_number"`
	Date                string          `json:"date"`
	IsPrintAsCopy       bool            `json:"is_print_as_copy"`
	KitchenDetails      []KitchenDetail `json:"kitchen_details"`
}

type PrintKitchenRequestBody struct {
	Kitchen Kitchen `json:"kitchen"`
}
