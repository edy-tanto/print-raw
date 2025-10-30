//go:build windows
// +build windows

package driver_windows

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
)

var (
	winspool        = syscall.NewLazyDLL("winspool.drv")
	openPrinter     = winspool.NewProc("OpenPrinterW")
	closePrinter    = winspool.NewProc("ClosePrinter")
	startDocPrinter = winspool.NewProc("StartDocPrinterW")
	endDocPrinter   = winspool.NewProc("EndDocPrinter")
	writePrinter    = winspool.NewProc("WritePrinter")
)

type DOC_INFO_1 struct {
	pDocName    *uint16
	pOutputFile *uint16
	pDatatype   *uint16
}

func ImageToBytes(filePath string, maxWidth int) ([]byte, int, int, error) {
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

func Print(data []byte, printerNameRequest string) {
	// Printer name
	printerName, err := syscall.UTF16PtrFromString(printerNameRequest)
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
	docName, _ := syscall.UTF16PtrFromString("Receipt Print")
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

// Target printer value example "192.168.10.252:9100"
func PrintEth(data []byte, targetPrinter string) {
	// Dial a TCP connection to the printer
	conn, err := net.DialTimeout("tcp", targetPrinter, 3*time.Second)
	if err != nil {
		fmt.Printf("Failed to connect to printer at %s: %v\n", targetPrinter, err)
		return
	}
	defer conn.Close()

	// Write to printer
	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Error writing to printer: %v\n", err)
		return
	}

	fmt.Println("Print successful!")
}
