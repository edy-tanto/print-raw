package handler

import (
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"edy-tanto/printer-pos/internal/print_web_service/printer"
	"encoding/json"
	"fmt"
	"net/http"
)

type PrintKitchenHandler struct{}

func (h *PrintKitchenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	var printReqBody dto.PrintKitchenRequestBody
	err := json.NewDecoder(r.Body).Decode(&printReqBody)

	if err != nil {
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	printer.ExecutePrintKitchen(printReqBody)

	json.NewEncoder(w).Encode(printReqBody)
}
