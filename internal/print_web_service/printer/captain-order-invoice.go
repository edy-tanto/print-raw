package printer

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"edy-tanto/printer-pos/internal/print_web_service/utils"
	"fmt"
	"log"
	"strings"
	"time"
)

func ExecutePrintCaptainOrderInvoice(body dto.PrintCaptainOrderInvoiceRequestBody) {
	const MAX_WIDTH_IMAGE = 700 // sesuaikan dengan printer

	imageData, widthPixels, heightPixels, err := driver_windows.ImageToBytes("captain-order-receipt-header.bmp", MAX_WIDTH_IMAGE)
	if err != nil {
		log.Fatal(err)
	}

	x := (widthPixels + 7) / 8
	y := heightPixels
	xL := byte(x % 256)
	xH := byte(x / 256)
	yL := byte(y % 256)
	yH := byte(y / 256)

	data := []byte{
		0x1B, 0x40, // Initialize
		0x1B, 0x61, 0x00, // Left alignment
		0x1D, 0x76, 0x30, 0x00,
		xL, xH, yL, yH,
	}
	data = append(data, imageData...)
	data = append(data, 0x0A, 0x0A)       // Line feed tambahan
	data = append(data, 0x1D, 0x21, 0x00) // Normal size text
	data = append(data, 0x1B, 0x61, 0x00) // Kembali ke kiri
	data = append(data, []byte(strings.Repeat("=", 48))...)

	// Header
	codeWithOp := fmt.Sprintf("%-11s %-14s %9s %10s\n", "Table/Room:", body.CaptainOrderInvoice.TableOrRoomNumber, "Waitress:", body.CaptainOrderInvoice.WaitressName)
	data = append(data, []byte(codeWithOp)...)
	tableOrRoomNumberWithCustomerCount := fmt.Sprintf("%-26s %-14s %-2d/ %-2d\n", " ", "#Adult/#Child:", body.CaptainOrderInvoice.CustomerAdultCount, body.CaptainOrderInvoice.CustomerChildCount)
	data = append(data, []byte(tableOrRoomNumberWithCustomerCount)...)
	customerName := fmt.Sprintf("%-11s %-36s\n", "Guest Name:", body.CaptainOrderInvoice.CustomerName)
	data = append(data, []byte(customerName)...)

	data = append(data, []byte("\n")...)

	if body.CaptainOrderInvoice.IsPrintAsCopy == true {
		data = append(data, 0x1B, 0x61, 0x01) // Center alignment
		data = append(data, 0x1D, 0x21, 0x22) // height 3 width 3

		data = append(data, []byte("SALINAN")...)
		data = append(data, []byte("\n\n")...)

		data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size
		data = append(data, 0x1B, 0x61, 0x00) // Left alignment
	}

	// Content
	data = append(data, []byte(strings.Repeat("-", 48))...)
	columnName := fmt.Sprintf("%-20s %8s %18s\n", "Product", "Qty", "Subtotal")
	data = append(data, []byte(columnName)...)
	data = append(data, []byte(strings.Repeat("-", 48))...)

	for _, detail := range body.CaptainOrderInvoice.SalesDetails {
		// wrap product name to 20 chars; first line shows qty and subtotal, continuations omit them
		const itemColWidth = 20
		nameRunes := []rune(detail.Item)
		if len(nameRunes) <= itemColWidth {
			detailText := fmt.Sprintf("%-20s %8d %18s\n", detail.Item, detail.Qty, utils.FormatMoney(detail.SubtotalWithoutTax))
			data = append(data, []byte(detailText)...)
			continue
		}

		first := string(nameRunes[:itemColWidth])
		firstLine := fmt.Sprintf("%-20s %8d %18s\n", first, detail.Qty, utils.FormatMoney(detail.SubtotalWithoutTax))
		data = append(data, []byte(firstLine)...)

		for i := itemColWidth; i < len(nameRunes); i += itemColWidth {
			end := i + itemColWidth
			if end > len(nameRunes) {
				end = len(nameRunes)
			}
			part := string(nameRunes[i:end])
			contLine := fmt.Sprintf("%-20s %8s %18s\n", part, "", "")
			data = append(data, []byte(contLine)...)
		}
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	totalQty := fmt.Sprintf("%-9s %-3d %-6s", "Quantity:", body.CaptainOrderInvoice.TotalQty, " ")
	data = append(data, []byte(totalQty)...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off
	withSubtotal := fmt.Sprintf("%9s %18s\n", "Sub total", utils.FormatMoney(body.CaptainOrderInvoice.Subtotal))
	data = append(data, []byte(withSubtotal)...)

	discountAmount := fmt.Sprintf("%-20s %-8s %18s\n", " ", "Discount", utils.FormatMoney(body.CaptainOrderInvoice.DiscountAmount))
	data = append(data, []byte(discountAmount)...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	grandTotal := fmt.Sprintf(
		"%-17s %11s %18s\n\n",
		" ",
		"Grand Total",
		utils.FormatMoneyTwoDigitAfterComma(body.CaptainOrderInvoice.GrandTotal+body.CaptainOrderInvoice.CreditCardCharge),
	)
	data = append(data, []byte(grandTotal)...)

	data = append(data, []byte("\n")...)

	data = append(data, []byte(strings.Repeat("=", 48))...)
	data = append(data, 0x1B, 0x61, 0x01)
	paymentHeader := fmt.Sprintf("%s\n", "Payment")
	data = append(data, []byte(paymentHeader)...)
	data = append(data, 0x1B, 0x61, 0x00)
	data = append(data, []byte(strings.Repeat("=", 48))...)

	for _, payment := range body.CaptainOrderInvoice.Payments {
		amount := strings.TrimSpace(utils.FormatMoneyTwoDigitAfterComma(payment.Total))
		if len(amount) > 27 {
			amount = amount[len(amount)-27:]
		}
		paymentLine := fmt.Sprintf(
			"%-20s %27s\n",
			payment.Method,
			amount,
		)
		data = append(data, []byte(paymentLine)...)
	}

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	ccCharge := fmt.Sprintf("%-10s %-14s\n", "CC Charge:", utils.FormatMoneyTwoDigitAfterComma(body.CaptainOrderInvoice.CreditCardCharge))
	data = append(data, []byte(ccCharge)...)

	data = append(data, []byte("\n")...)

	remarkLabel := fmt.Sprintf("%-35s\n", "Remark:")
	data = append(data, []byte(remarkLabel)...)
	remark := fmt.Sprintf("%-48s\n", body.CaptainOrderInvoice.Note)
	data = append(data, []byte(remark)...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	data = append(data, []byte("\n\n\n\n\n")...)

	data = append(data, 0x1B, 0x61, 0x01) // Center alignment

	opCentered := utils.CenterInParentheses(body.CaptainOrderInvoice.Op, 20)
	guestSignature := utils.CenterInParentheses("____________________", 20)
	signatureLine := fmt.Sprintf("%s %s\n", opCentered, guestSignature)
	data = append(data, []byte(signatureLine)...)

	// Label "Cashier" dan "Guest", masing-masing diformat agar berada di tengah 22 karakter
	label1 := fmt.Sprintf("%-22s", utils.CenterText("Cashier", 22))
	label2 := fmt.Sprintf("%-22s", utils.CenterText("Guest", 22))
	signatureLabel := fmt.Sprintf("%s%s\n", label1, label2)
	data = append(data, []byte(signatureLabel)...)

	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	data = append(data, []byte("\n\n")...)

	parsedTime, _ := time.Parse(time.RFC3339Nano, body.CaptainOrderInvoice.Date)
	localTime := parsedTime.In(time.Local)
	postingDate := fmt.Sprintf("%-10s %-20s\n", "Posting Date:", localTime.Format("02/01/2006 15:04:05"))
	data = append(data, []byte(postingDate)...)
	
	// Printed: menggunakan waktu sekarang saat print
	printTime := time.Now()
	// Audit: menggunakan PostDate dari body
	var auditTime time.Time
	if body.CaptainOrderInvoice.PostDate != "" {
		parsedPostDate, err := time.Parse(time.RFC3339Nano, body.CaptainOrderInvoice.PostDate)
		if err == nil {
			auditTime = parsedPostDate.In(time.Local)
		} else {
			auditTime = printTime
		}
	} else {
		auditTime = printTime
	}
	printedWithAuditDate := fmt.Sprintf("%-8s %-20s %-5s %-17s\n", "Printed:", printTime.Format("02/01/2006 15:04:05"), "Audit:", auditTime.Format("02/01/2006"))
	data = append(data, []byte(printedWithAuditDate)...)

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data, body.PrinterName)
}
