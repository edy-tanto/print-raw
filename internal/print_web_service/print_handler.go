package print_web_service

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_linux"
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"encoding/json"
	"fmt"
	"net/http"
)

type SalesDetail struct {
}

type Sales struct {
	Id             uint          `json:"id"`
	Code           string        `json:"code"`
	DiscountAmount float32       `json:"discount_amount"`
	Summary        float32       `json:"summary"`
	SalesDetails   []SalesDetail `json:"sales_details"`
}

type PrintRequestBody struct {
	Sales Sales `json:"sales"`
}

type PrintHandler struct{}

func (h *PrintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	executePrint(printReqBody)

	json.NewEncoder(w).Encode(printReqBody)
}

func executePrint(body PrintRequestBody) {
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

	data = append(data, imageData...) // Image data
	data = append(data, 0x0A)         // Line feed
	data = append(data, 0x1B, 0x61, 0x01)
	data = append(data, []byte("Qubu Resort Waterpark\n")...)
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment
	data = append(data, []byte("2025-05-01 18:59:59")...)

	summary := fmt.Sprintf("Summary %14.0f\n", body.Sales.Summary)
	data = append(data, []byte(summary)...)
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	// driver_windows.Print(data)
	driver_linux.Print(data)
}
