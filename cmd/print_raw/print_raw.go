package main

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
)

func main() {
	data := []byte{}

	driver_windows.Print(data)
}
