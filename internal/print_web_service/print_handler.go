package print_web_service

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	Date             string        `json:"date"`
	DiscountAmount   float32       `json:"discount_amount"`
	Summary          float32       `json:"summary"`
	SalesDetails     []SalesDetail `json:"sales_details"`
}

type PrintRequestBody struct {
	Sales Sales `json:"sales"`
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

func ExecutePrint(body PrintRequestBody) {
	const MAX_WIDTH_IMAGE = 384

	imageData, widthPixels, heightPixels, _ := driver_windows.ImageToBytes("cat.bmp", MAX_WIDTH_IMAGE) // Adjust path and width
	x := (widthPixels + 7) / 8                                                                         // 1 byte per row
	y := heightPixels                                                                                  // 8 rows
	xL := byte(x % 256)                                                                                // 1
	xH := byte(x / 256)                                                                                // 0
	yL := byte(y % 256)                                                                                // 8
	yH := byte(y / 256)

	// ESC/POS commands
	data := []byte{
		0x1B, 0x40, // Initialize printer
		0x1B, 0x61, 0x01, // Center alignment
		0x1D, 0x76, 0x30, 0x00, // GS v 0 command
		xL, xH, yL, yH, // Width and height parameters
	}

	// Header
	data = append(data, imageData...) // Image data
	data = append(data, 0x0A)         // Line feed
	data = append(data, 0x1B, 0x61, 0x01)
	data = append(data, 0x1D, 0x21, 0x01) // Double height
	data = append(data, []byte(body.Sales.UnitBusinessName+"\n\n")...)
	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	date := fmt.Sprintf("%-6s : %-20s\n", "Date", body.Sales.Date)
	data = append(data, []byte(date)...)
	code := fmt.Sprintf("%-6s : %-20s\n", "Code", body.Sales.Code)
	data = append(data, []byte(code)...)
	op := fmt.Sprintf("%-6s : %-20s\n", "OP", body.Sales.Op)
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

	for _, detail := range body.Sales.SalesDetails {
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
	discountAmount := fmt.Sprintf("%-33s %14.0f\n", "Diskon", body.Sales.DiscountAmount)
	data = append(data, []byte(discountAmount)...)

	data = append(data, []byte("\n")...)

	data = append(data, 0x1D, 0x21, 0x11) // Double height

	grandTotal := fmt.Sprintf("%-9s %14.0f\n", "Total", body.Sales.Summary)
	data = append(data, []byte(grandTotal)...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data)
	// driver_linux.Print(data)
}
