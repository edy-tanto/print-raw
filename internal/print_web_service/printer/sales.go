package printer

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"edy-tanto/printer-pos/internal/print_web_service/utils"
	"fmt"
	"strings"
)

type FootnoteAlignEnum string

const (
	FootnoteAlignLeft   FootnoteAlignEnum = "LEFT"
	FootnoteAlignCenter FootnoteAlignEnum = "CENTER"
	FootnoteAlignRight  FootnoteAlignEnum = "RIGHT"
)

func ExecutePrint(body dto.PrintRequestBody) {
	const MAX_WIDTH_IMAGE = 384

	imageData, widthPixels, heightPixels, _ := driver_windows.ImageToBytes("paradis-q.bmp", MAX_WIDTH_IMAGE) // Adjust path and width
	x := (widthPixels + 7) / 8                                                                               // 1 byte per row
	y := heightPixels                                                                                        // 8 rows
	xL := byte(x % 256)                                                                                      // 1
	xH := byte(x / 256)                                                                                      // 0
	yL := byte(y % 256)                                                                                      // 8
	yH := byte(y / 256)

	// ESC/POS commands
	data := []byte{
		0x1B, 0x40, // Initialize printer
		0x1B, 0x61, 0x01, // Center alignment
		0x1D, 0x76, 0x30, 0x00, // GS v 0 command
		xL, xH, yL, yH, // Width and height parameters
	}

	// Header
	data = append(data, imageData...)     // Image data
	data = append(data, 0x0A)             // Line feed
	data = append(data, 0x1B, 0x61, 0x01) // Center alignment
	data = append(data, 0x1D, 0x21, 0x01) // Double height
	data = append(data, []byte(body.Sales.UnitBusinessName+"\n\n")...)
	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	date := fmt.Sprintf("%-14s : %-20s\n", "Date", utils.FormatDatetime(body.Sales.Date))
	data = append(data, []byte(date)...)
	code := fmt.Sprintf("%-14s : %-20s\n", "Code", body.Sales.Code)
	data = append(data, []byte(code)...)
	op := fmt.Sprintf("%-14s : %-20s\n", "OP", body.Sales.Op)
	data = append(data, []byte(op)...)

	if body.Sales.CustomerName != "" {
		customerName := fmt.Sprintf("%-14s : %-20s\n", "Nama Pelanggan", body.Sales.CustomerName)
		data = append(data, []byte(customerName)...)
	}

	if body.Sales.TableNumber != "" {
		tableNumber := fmt.Sprintf("%-14s : %-20s\n", "Nomor Meja", body.Sales.TableNumber)
		data = append(data, []byte(tableNumber)...)
	}

	paymentMethod := fmt.Sprintf("%-14s : %-20s\n", "Metode Bayar", body.Sales.PaymentMethod)
	data = append(data, []byte(paymentMethod)...)

	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on

	if body.Sales.IsPrintAsCopy == true {
		data = append(data, 0x1B, 0x61, 0x01) // Center alignment
		data = append(data, 0x1D, 0x21, 0x22) // height 3 width 3

		data = append(data, []byte("SALINAN")...)

		data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
		data = append(data, 0x1B, 0x61, 0x00) // Left alignment
		data = append(data, []byte("\n\n")...)
	}

	// Content
	columnName := fmt.Sprintf(
		"%-23s %-9s %14s\n",
		"Item",
		"Qty",
		"Price",
	)
	data = append(data, []byte(columnName)...)

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	for _, detail := range body.Sales.SalesDetails {
		// TODO: handle if item name to long, solution new line or truncate at the last
		detailText := fmt.Sprintf(
			"%-20s %-1s %3d %-6s %14s\n",
			detail.Item,
			" ",
			detail.Qty,
			" ",

			utils.FormatMoney(detail.TotalFinal),
		)
		data = append(data, []byte(detailText)...)
	}

	// Summary
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment``
	data = append(data, []byte(strings.Repeat("-", 48))...)
	ccCharge := fmt.Sprintf("%-10s %15s\n\n", "CC Charge:", utils.FormatMoney(body.Sales.CreditCardCharge))
	data = append(data, []byte(ccCharge)...)

	data = append(data, 0x1D, 0x21, 0x11) // Double height

	grandTotal := fmt.Sprintf("%-9s %14s\n", "Total", utils.FormatMoney(body.Sales.GrandTotal+body.Sales.CreditCardCharge))
	data = append(data, []byte(grandTotal)...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	// Footer
	if body.Sales.Footnote != "" {
		data = append(data, []byte("\n")...)
		data = append(data, []byte(strings.Repeat("-", 48))...)
		data = append(data, 0x1B, 0x4D, 0x01) // Change font
		data = append(data, []byte("\n\n")...)

		switch body.Sales.FootnoteAlign {
		case string(FootnoteAlignCenter):
			data = append(data, 0x1B, 0x61, 0x01)
		case string(FootnoteAlignRight):
			data = append(data, 0x1B, 0x61, 0x02)
		default: // FootnoteAlignLeft or unknown
			data = append(data, 0x1B, 0x61, 0x00)
		}

		data = append(data, []byte(body.Sales.Footnote)...)
	}

	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data, body.PrinterName)
	// driver_linux.Print(data)
}
