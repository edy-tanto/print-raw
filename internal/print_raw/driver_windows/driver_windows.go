//go:build windows
// +build windows

package driver_windows

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
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

func resolveAssetPath(filePath string) string {
	// Services start with C:\Windows\System32 as the working directory, so ensure assets resolve relative to the executable.
	if filepath.IsAbs(filePath) {
		return filePath
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		candidate := filepath.Join(exeDir, filePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return filePath
}

func ImageToBytes(filePath string, maxWidth int) ([]byte, int, int, error) {
	filePath = resolveAssetPath(filePath)

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
	// Normalize printer name (trim whitespace, handle UNC paths for shared printers)
	printerNameRequest = strings.TrimSpace(printerNameRequest)
	if printerNameRequest == "" {
		fmt.Println("Error: printer name is empty")
		return
	}

	// Printer name
	printerName, err := syscall.UTF16PtrFromString(printerNameRequest)
	if err != nil {
		fmt.Printf("Error converting printer name %q: %v\n", printerNameRequest, err)
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
		// Get Windows error code for more detailed error information
		errCode := syscall.GetLastError()
		fmt.Printf("Error opening printer %q: %v (Windows error code: %d)\n", printerNameRequest, err, errCode)
		fmt.Println("Note: For shared printers, ensure:")
		fmt.Println("  - Printer is accessible via network")
		fmt.Println("  - Service account has permission to access the printer")
		fmt.Println("  - Printer name format is correct (e.g., \\\\server\\printer for shared printers)")
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
		errCode := syscall.GetLastError()
		fmt.Printf("Error starting document on printer %q: %v (Windows error code: %d)\n", printerNameRequest, err, errCode)
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
		errCode := syscall.GetLastError()
		fmt.Printf("Error writing to printer %q: %v (Windows error code: %d)\n", printerNameRequest, err, errCode)
		return
	}

	// End document
	ret, _, err = endDocPrinter.Call(uintptr(handle))
	if ret == 0 {
		errCode := syscall.GetLastError()
		fmt.Printf("Error ending document on printer %q: %v (Windows error code: %d)\n", printerNameRequest, err, errCode)
		return
	}

	fmt.Printf("Print successful to printer: %q\n", printerNameRequest)
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
