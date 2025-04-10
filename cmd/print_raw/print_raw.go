package main

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
)

const MAX_WIDTH_IMAGE = 480

func main() {
	imageData, widthPixels, heightPixels, _ := driver_windows.ImageToBytes("cat.bmp", MAX_WIDTH_IMAGE)
	x := (widthPixels + 7) / 8 // 1 byte per row
	y := heightPixels          // 8 rows
	xL := byte(x % 256)        // 1
	xH := byte(x / 256)        // 0
	yL := byte(y % 256)        // 8
	yH := byte(y / 256)        // 0

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

	driver_windows.Print(data)
}
