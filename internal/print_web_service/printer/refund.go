package printer

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"edy-tanto/printer-pos/internal/print_web_service/utils"
	"fmt"
	"strings"
)

func ExecutePrintCashRefund(body dto.PrintCashRefundRequestBody) {
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
	data = append(data, []byte("REFUND\n\n")...)
	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	date := fmt.Sprintf("%-14s : %-20s\n", "Date", utils.FormatDatetime(body.CashRefund.Date))
	data = append(data, []byte(date)...)
	op := fmt.Sprintf("%-14s : %-20s\n", "OP", body.CashRefund.Op)
	data = append(data, []byte(op)...)

	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on

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

	for _, detail := range body.CashRefund.CashRefundDetails {
		// wrap long item names to fit item column (20 chars)
		const itemColWidth = 20
		nameRunes := []rune(detail.Item)
		if len(nameRunes) <= itemColWidth {
			detailText := fmt.Sprintf(
				"%-20s %-1s %3d %-6s %14s\n",
				detail.Item,
				" ",
				detail.Qty,
				" ",
				utils.FormatMoney(detail.Subtotal),
			)
			data = append(data, []byte(detailText)...)
			continue
		}

		// First line with qty and price
		first := string(nameRunes[:itemColWidth])
		firstLine := fmt.Sprintf(
			"%-20s %-1s %3d %-6s %14s\n",
			first,
			" ",
			detail.Qty,
			" ",
			utils.FormatMoney(detail.Subtotal),
		)
		data = append(data, []byte(firstLine)...)

		// Continuation lines: only item column
		for i := itemColWidth; i < len(nameRunes); i += itemColWidth {
			end := i + itemColWidth
			if end > len(nameRunes) {
				end = len(nameRunes)
			}
			part := string(nameRunes[i:end])
			contLine := fmt.Sprintf(
				"%-20s %-1s %3s %-6s %14s\n",
				part,
				" ",
				"",
				" ",
				"",
			)
			data = append(data, []byte(contLine)...)
		}
	}

	data = append(data, []byte("\n")...)

	// Summary
	data = append(data, 0x1D, 0x21, 0x11) // Double height

	grandTotal := fmt.Sprintf("%-9s %14s\n", "Total", utils.FormatMoney(body.CashRefund.TotalRefund))
	data = append(data, []byte(grandTotal)...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data, body.PrinterName)
	// driver_linux.Print(data)
}
