package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const MAX_WIDTH_IMAGE = 480

type Sales struct {
	CreatedAt string `json:"created_at"`
	SalesId   uint16 `json:"sales_id"`
}

func main() {
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &printHandler{})
	mux.Handle("/print-dynamic", &dynamicPrintHandler{})
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		salesJson := []byte(queryValues.Get("sales"))

		var sales Sales
		var err = json.Unmarshal(salesJson, &sales)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w.Write([]byte(sales.CreatedAt))
	})

	http.ListenAndServe(":8080", mux)

	// list_printers()
}

type printHandler struct{}

func (h *printHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageData, widthPixels, heightPixels, _ := imageToBytes("cat.bmp", MAX_WIDTH_IMAGE) // Adjust path and width
	x := (widthPixels + 7) / 8                                                          // 1 byte per row
	y := heightPixels                                                                   // 8 rows
	xL := byte(x % 256)                                                                 // 1
	xH := byte(x / 256)                                                                 // 0
	yL := byte(y % 256)                                                                 // 8
	yH := byte(y / 256)                                                                 // 0

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
	data = append(data, []byte("hello world\n")...)
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment
	data = append(data, []byte("hello world\n")...)
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	print(data)
	w.Write([]byte("print accepted"))
}

type dynamicPrintHandler struct{}

func (h *dynamicPrintHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sales := &Sales{}
	salesJsonData := []byte(r.URL.Query().Get("sales"))

	json.Unmarshal(salesJsonData, &sales)

	w.Write([]byte(sales.CreatedAt))
}
