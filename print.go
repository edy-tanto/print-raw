package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/nfnt/resize" // For resizing the image
	"golang.org/x/image/bmp" // Correct import for BMP decoding
)

// Windows API functions from winspool.drv
var (
	winspool        = syscall.NewLazyDLL("winspool.drv")
	openPrinter     = winspool.NewProc("OpenPrinterW")
	closePrinter    = winspool.NewProc("ClosePrinter")
	startDocPrinter = winspool.NewProc("StartDocPrinterW")
	endDocPrinter   = winspool.NewProc("EndDocPrinter")
	writePrinter    = winspool.NewProc("WritePrinter")
)

// DOC_INFO_1 structure
type DOC_INFO_1 struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

func imageToBytes(filePath string, maxWidth int) ([]byte, int, int, error) {
	// Open the image file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	// Decode the image (assuming BMP for simplicity; adjust for PNG/JPEG if needed)
	img, err := bmp.Decode(file)
	if err != nil {
		return nil, 0, 0, err
	}

	// Resize to fit printer width (e.g., 384 pixels)
	resized := resize.Resize(uint(maxWidth), 0, img, resize.Lanczos3)
	bounds := resized.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert to monochrome
	bytesPerRow := (width + 7) / 8
	imageData := make([]byte, bytesPerRow*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			// Simple threshold: if brightness < 50%, make it black
			if (r+g+b)/3 < 0x8000 {
				byteIndex := y*bytesPerRow + x/8
				bitIndex := 7 - (x % 8)
				imageData[byteIndex] |= 1 << bitIndex
			}
		}
	}

	return imageData, width, height, nil
}

func print() {
	// Printer name
	printerName, err := syscall.UTF16PtrFromString("EPSON TM-T82 Receipt")
	if err != nil {
		fmt.Println("Error converting printer name:", err)
		return
	}

	// Open printer
	var handle syscall.Handle
	ret, _, err := openPrinter.Call(
		uintptr(unsafe.Pointer(printerName)),
		uintptr(unsafe.Pointer(&handle)),
		0,
	)
	if ret == 0 {
		fmt.Println("Error opening printer:", err)
		return
	}
	defer closePrinter.Call(uintptr(handle))

	// Start document
	docName, _ := syscall.UTF16PtrFromString("My Document")
	dataType, _ := syscall.UTF16PtrFromString("RAW")
	di := DOC_INFO_1{
		pDocName:  docName,
		pDatatype: dataType,
	}
	ret, _, err = startDocPrinter.Call(
		uintptr(handle),
		1,
		uintptr(unsafe.Pointer(&di)),
	)
	if ret == 0 {
		fmt.Println("Error starting document:", err)
		return
	}

	imageData, widthPixels, heightPixels, _ := imageToBytes("cat.bmp", 480) // Adjust path and width
	x := (widthPixels + 7) / 8                                              // 1 byte per row
	y := heightPixels                                                       // 8 rows
	xL := byte(x % 256)                                                     // 1
	xH := byte(x / 256)                                                     // 0
	yL := byte(y % 256)                                                     // 8
	yH := byte(y / 256)                                                     // 0

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

	// Write to printer
	var written uint32
	ret, _, err = writePrinter.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if ret == 0 {
		fmt.Println("Error writing to printer:", err)
		return
	}

	// End document
	ret, _, err = endDocPrinter.Call(uintptr(handle))
	if ret == 0 {
		fmt.Println("Error ending document:", err)
		return
	}

	fmt.Println("Print successful!")
}
