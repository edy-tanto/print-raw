package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
)

type PrinterListHandler struct{}

func (h *PrinterListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	printers, err := driver_windows.ListPrinters()
	if err != nil {
		fmt.Println(err.Error())

		status := http.StatusInternalServerError
		if errors.Is(err, driver_windows.ErrPrinterEnumerationUnsupported) {
			status = http.StatusNotImplemented
		}

		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "unable to list printers",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string][]string{
		"printers": printers,
	})
}
