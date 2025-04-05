package main

import (
	"fmt"
	"syscall"
	"unsafe"
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

func print() {
	// Convert printer name to UTF-16 pointer
	printerName, err := syscall.UTF16PtrFromString("EPSON TM-T82 Receipt")
	if err != nil {
		fmt.Println("Error converting printer name:", err)
		return
	}

	// Open the printer
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

	// Prepare DOC_INFO_1 structure
	docName, err := syscall.UTF16PtrFromString("My Document")
	if err != nil {
		fmt.Println("Error converting document name:", err)
		return
	}
	dataType, err := syscall.UTF16PtrFromString("RAW")
	if err != nil {
		fmt.Println("Error converting data type:", err)
		return
	}
	di := DOC_INFO_1{
		pDocName:    docName,
		pOutputFile: nil,
		pDatatype:   dataType,
	}

	// Start the print document
	ret, _, err = startDocPrinter.Call(
		uintptr(handle),
		1,
		uintptr(unsafe.Pointer(&di)),
	)
	if ret == 0 {
		fmt.Println("Error starting document:", err)
		return
	}

	// Data to print: initialize, print text, feed paper, and cut
	data := append([]byte{0x1B, 0x40}, "hello world\n"...) // ESC @ to initialize, print text with newline
	data = append(data, 0x1B, 0x64, 0x05)                  // ESC d 5: Feed 5 lines
	data = append(data, 0x1D, 0x56, 0x00)                  // GS V 0: Full cut

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

	// End the print document
	ret, _, err = endDocPrinter.Call(uintptr(handle))
	if ret == 0 {
		fmt.Println("Error ending document:", err)
		return
	}

	fmt.Println("Successfully printed 'hello world', fed the paper, and cut.")
}
