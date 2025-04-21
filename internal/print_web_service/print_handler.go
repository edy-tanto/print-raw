package print_web_service

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
	GrandTotal       float32       `json:"grand_total"`
	SalesDetails     []SalesDetail `json:"sales_details"`
}

type PrintRequestBody struct {
	Sales Sales `json:"sales"`
}

type CashRefundDetail struct {
	Item     string  `json:"item"`
	Qty      uint    `json:"qty"`
	Subtotal float32 `json:"subtotal"`
}

type CashRefund struct {
	Id                uint               `json:"id"`
	Op                string             `json:"op"`
	Date              string             `json:"date"`
	TotalRefund       float32            `json:"total_refund"`
	CashRefundDetails []CashRefundDetail `json:"cash_refund_details"`
}

type PrintCashRefundRequestBody struct {
	CashRefund CashRefund `json:"cash_refund"`
}

type KitchenDetail struct {
	Item    string  `json:"item"`
	Qty     uint  	`json:"qty"`
	Note    string  `json:"note"`
}

type Kitchen struct {
	Id					uint				`json:"id"`
	Op					string				`json:"op"`
	Code				string				`json:"code"`
	Outlet				string				`json:"outlet"`
	CustomerName		string				`json:"customer_name"`
	TableNumber      	string        		`json:"table_number"`
	Date           string             		`json:"date"`
	KitchenDetails		[]KitchenDetail 	`json:"kitchen_details"`
}

type PrintKitchenRequestBody struct {
	Kitchen 	Kitchen		`json:"kitchen"`
}

type PrintHandler struct{}

func (h *PrintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

	if r.Method == http.MethodOptions {
		// handle preflight request
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	var printReqBody PrintRequestBody
	err := json.NewDecoder(r.Body).Decode(&printReqBody)

	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	ExecutePrint(printReqBody)

	json.NewEncoder(w).Encode(printReqBody)
}

type PrintCashRefundHandler struct{}

func (h *PrintCashRefundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

	if r.Method == http.MethodOptions {
		// handle preflight request
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	var printReqBody PrintCashRefundRequestBody
	err := json.NewDecoder(r.Body).Decode(&printReqBody)

	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	ExecutePrintCashRefund(printReqBody)

	json.NewEncoder(w).Encode(printReqBody)
}

type PrintKitchenHandler struct{}

func (h * PrintKitchenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

	if r.Method == http.MethodOptions {
		// handle preflight request
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	var printReqBody PrintKitchenRequestBody
	err := json.NewDecoder(r.Body).Decode(&printReqBody)

	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	ExecutePrintKitchen(printReqBody)

	json.NewEncoder(w).Encode(printReqBody)
}

func ExecutePrint(body PrintRequestBody) {
	const MAX_WIDTH_IMAGE = 384

	imageData, widthPixels, heightPixels, _ := driver_windows.ImageToBytes("paradis-q.bmp", MAX_WIDTH_IMAGE) // Adjust path and width
	x := (widthPixels + 7) / 8                                                                               // 1 byte per row
	y := heightPixels                                                                                        // 8 rows
	xL := byte(x % 256)                                                                                      // 1
	xH := byte(x / 256)                                                                                      // 0
	yL := byte(y % 256)                                                                                      // 8
	yH := byte(y / 256)

	// ESC/POS commands
	data := []byte{
		0x1B, 0x40, // Initialize printer
		0x1B, 0x61, 0x01, // Center alignment
		0x1D, 0x76, 0x30, 0x00, // GS v 0 command
		xL, xH, yL, yH, // Width and height parameters
	}

	// Header
	data = append(data, imageData...)     // Image data
	data = append(data, 0x0A)             // Line feed
	data = append(data, 0x1B, 0x61, 0x01) // Center alignment
	data = append(data, 0x1D, 0x21, 0x01) // Double height
	data = append(data, []byte(body.Sales.UnitBusinessName+"\n\n")...)
	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	date := fmt.Sprintf("%-14s : %-20s\n", "Date", formatDatetime(body.Sales.Date))
	data = append(data, []byte(date)...)
	code := fmt.Sprintf("%-14s : %-20s\n", "Code", body.Sales.Code)
	data = append(data, []byte(code)...)
	op := fmt.Sprintf("%-14s : %-20s\n", "OP", body.Sales.Op)
	data = append(data, []byte(op)...)

	if body.Sales.CustomerName != "" {
		customerName := fmt.Sprintf("%-14s : %-20s\n", "Nama Pelanggan", body.Sales.CustomerName)
		data = append(data, []byte(customerName)...)
	}

	if body.Sales.TableNumber != "" {
		tableNumber := fmt.Sprintf("%-14s : %-20s\n", "Nomor Meja", body.Sales.TableNumber)
		data = append(data, []byte(tableNumber)...)
	}

	paymentMethod := fmt.Sprintf("%-14s : %-20s\n", "Metode Bayar", body.Sales.PaymentMethod)
	data = append(data, []byte(paymentMethod)...)

	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on

	if body.Sales.IsPrintAsCopy == true {
		data = append(data, 0x1B, 0x61, 0x01) // Center alignment
		data = append(data, 0x1D, 0x21, 0x22) // height 3 width 3

		data = append(data, []byte("SALINAN")...)

		data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
		data = append(data, 0x1B, 0x61, 0x00) // Left alignment
		data = append(data, []byte("\n\n")...)
	}

	// Content
	columnName := fmt.Sprintf(
		"%-23s %-9s %14s\n",
		"Item",
		"Qty",
		"Price",
	)
	data = append(data, []byte(columnName)...)

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	for _, detail := range body.Sales.SalesDetails {
		// TODO: handle if item name to long, solution new line or truncate at the last
		detailText := fmt.Sprintf(
			"%-20s %-1s %3d %-6s %14.0f\n",
			detail.Item,
			" ",
			detail.Qty,
			" ",
			detail.Subtotal,
		)
		data = append(data, []byte(detailText)...)
	}

	// Summary
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment``
	data = append(data, []byte("\n\n")...)

	data = append(data, 0x1D, 0x21, 0x11) // Double height

	grandTotal := fmt.Sprintf("%-9s %14.0f\n", "Total", body.Sales.GrandTotal)
	data = append(data, []byte(grandTotal)...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	// Footer
	if body.Sales.Footnote != "" {
		data = append(data, []byte("\n")...)
		data = append(data, []byte(strings.Repeat("-", 48))...)
		data = append(data, 0x1B, 0x4D, 0x01) // Change font

		data = append(data, []byte("\n\n")...)
		data = append(data, []byte(body.Sales.Footnote)...)
	}

	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data)
	// driver_linux.Print(data)
}

func ExecutePrintCashRefund(body PrintCashRefundRequestBody) {
	const MAX_WIDTH_IMAGE = 384

	imageData, widthPixels, heightPixels, _ := driver_windows.ImageToBytes("paradis-q.bmp", MAX_WIDTH_IMAGE) // Adjust path and width
	x := (widthPixels + 7) / 8                                                                               // 1 byte per row
	y := heightPixels                                                                                        // 8 rows
	xL := byte(x % 256)                                                                                      // 1
	xH := byte(x / 256)                                                                                      // 0
	yL := byte(y % 256)                                                                                      // 8
	yH := byte(y / 256)

	// ESC/POS commands
	data := []byte{
		0x1B, 0x40, // Initialize printer
		0x1B, 0x61, 0x01, // Center alignment
		0x1D, 0x76, 0x30, 0x00, // GS v 0 command
		xL, xH, yL, yH, // Width and height parameters
	}

	// Header
	data = append(data, imageData...)     // Image data
	data = append(data, 0x0A)             // Line feed
	data = append(data, 0x1B, 0x61, 0x01) // Center alignment
	data = append(data, 0x1D, 0x21, 0x01) // Double height
	data = append(data, []byte("REFUND\n\n")...)
	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	date := fmt.Sprintf("%-14s : %-20s\n", "Date", formatDatetime(body.CashRefund.Date))
	data = append(data, []byte(date)...)
	op := fmt.Sprintf("%-14s : %-20s\n", "OP", body.CashRefund.Op)
	data = append(data, []byte(op)...)

	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on

	// Content
	columnName := fmt.Sprintf(
		"%-23s %-9s %14s\n",
		"Item",
		"Qty",
		"Price",
	)
	data = append(data, []byte(columnName)...)

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	for _, detail := range body.CashRefund.CashRefundDetails {
		// TODO: handle if item name to long, solution new line or truncate at the last
		detailText := fmt.Sprintf(
			"%-20s %-1s %3d %-6s %14.0f\n",
			detail.Item,
			" ",
			detail.Qty,
			" ",
			detail.Subtotal,
		)
		data = append(data, []byte(detailText)...)
	}

	data = append(data, []byte("\n")...)

	// Summary
	data = append(data, 0x1D, 0x21, 0x11) // Double height

	grandTotal := fmt.Sprintf("%-9s %14.0f\n", "Total", body.CashRefund.TotalRefund)
	data = append(data, []byte(grandTotal)...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data)
	// driver_linux.Print(data)
}

func ExecutePrintKitchen(body PrintKitchenRequestBody) {
	// ESC/POS commands
	data := []byte{
		0x1B, 0x40,       // Initialize printer
		0x1B, 0x61, 0x00, // Left Allignment
	}

	data = append(data, []byte(strings.Repeat("-", 48))...)

	// Header
	codeWithOutlet := fmt.Sprintf("%-14s : %-15s %10s\n", "CO ID", body.Kitchen.Code, body.Kitchen.Outlet)
	data = append(data, []byte(codeWithOutlet)...)
	op := fmt.Sprintf("%-14s : %-20s\n", "Waitress", body.Kitchen.Op)
	data = append(data, []byte(op)...)
	tableNumber := fmt.Sprintf("%-14s : %-20s\n", "Table Number", body.Kitchen.TableNumber)
	data = append(data, []byte(tableNumber)...)
	customerName := fmt.Sprintf("%-14s : %-20s\n", "Table Number", body.Kitchen.CustomerName)
	data = append(data, []byte(customerName)...)
	location := fmt.Sprintf("%-14s : %-20s\n", "Location", "Kitchen")
	data = append(data, []byte(location)...)

	// Content
	data = append(data, []byte(strings.Repeat("-", 48))...)
	columnName := fmt.Sprintf("%9s %-25s\n", "Qty", "Product")
	data = append(data, []byte(columnName)...)
	data = append(data, []byte(strings.Repeat("-", 48))...)

	for _, detail := range body.Kitchen.KitchenDetails {
		detailText := fmt.Sprintf(
			"%9d %-25s\n",
			detail.Qty,
			detail.Item,
		)
		data = append(data, []byte(detailText)...)
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)


	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	data = append(data, []byte("\n\n")...)

	parsedTime, _ := time.Parse(time.RFC3339Nano, body.Kitchen.Date)
	localTime := parsedTime.In(time.Local)
	date := fmt.Sprintf("%-10s %-20s\n", "Printed:", localTime.Format("02/01/2006 15:04:05"))
	data = append(data, []byte(date)...)

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data)
}

func formatDatetime(dateString string) string {
	layout := "2006-01-02T15:04:05.000"

	parsedTime, _ := time.ParseInLocation(layout, dateString, time.Local)

	localTime := parsedTime.In(time.Local)
	formattedDate := localTime.Format("02/01/2006 15:04:05")

	return formattedDate
}
