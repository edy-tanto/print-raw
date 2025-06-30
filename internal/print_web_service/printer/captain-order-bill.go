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

func ExecutePrintCaptainOrderBill(body dto.PrintCaptainOrderBillRequestBody) {
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
		0x1B, 0x40,         // Initialize
		0x1B, 0x61, 0x00,   // Left alignment
		0x1D, 0x76, 0x30, 0x00,
		xL, xH, yL, yH,
	}
	data = append(data, imageData...)
	data = append(data, 0x0A, 0x0A) // Line feed tambahan
	data = append(data, 0x1D, 0x21, 0x00) // Normal size text
	data = append(data, 0x1B, 0x61, 0x00) // Kembali ke kiri
	data = append(data, []byte(strings.Repeat("=", 48))...)

	// Header
	codeWithOp := fmt.Sprintf("%-3s %-22s %9s %10s\n", "CO:", body.CaptainOrderBill.Code, "Waitress:", body.CaptainOrderBill.WaitressName)
	data = append(data, []byte(codeWithOp)...)
	tableOrRoomNumberWithCustomerCount := fmt.Sprintf("%-11s %-14s %-14s %-2d/ %-2d\n", "Table/Room:", body.CaptainOrderBill.TableOrRoomNumber, "#Adult/#Child:", body.CaptainOrderBill.CustomerAdultCount, body.CaptainOrderBill.CustomerChildCount)
	data = append(data, []byte(tableOrRoomNumberWithCustomerCount)...)
	customerName := fmt.Sprintf("%-6s %-15s\n", "Guest Name:", body.CaptainOrderBill.CustomerName)
	data = append(data, []byte(customerName)...)

	data = append(data, []byte("\n\n")...)

	if body.CaptainOrderBill.IsPrintAsCopy == true {
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

	for _, detail := range body.CaptainOrderBill.CaptainOrderBillDetails {
		detailText := fmt.Sprintf(
			"%-20s %8d %18s\n",
			detail.Item,
			detail.Qty,
			utils.FormatMoney(detail.Subtotal),
		)
		data = append(data, []byte(detailText)...)
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	totalQty := fmt.Sprintf("%-9s %-3d %-6s", "Quantity:", body.CaptainOrderBill.TotalQty, " ")
	data = append(data, []byte(totalQty)...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off
	withSubtotal := fmt.Sprintf("%9s %18s\n", "Sub total", utils.FormatMoney(body.CaptainOrderBill.TotalGross))
	data = append(data, []byte(withSubtotal)...)

	discountAmount := fmt.Sprintf("%-20s %-8s %18s\n", " ", "Discount", utils.FormatMoney(body.CaptainOrderBill.DiscountAmount))
	data = append(data, []byte(discountAmount)...)

	separator := fmt.Sprintf("%48s\n", "--------------------------------")
	data = append(data, []byte(separator)...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	grandTotal := fmt.Sprintf("%-17s %11s %18s\n", " ", "Grand Total", utils.FormatMoney(body.CaptainOrderBill.TotalNet))
	data = append(data, []byte(grandTotal)...)

	data = append(data, []byte("\n\n")...)

	ccChargeLabel := fmt.Sprintf("%-10s %-5s\n", "CC Charge:", "0,00")
	data = append(data, []byte(ccChargeLabel)...)

	data = append(data, []byte("\n")...)

	remarkLabel := fmt.Sprintf("%-35s\n", "Remark:")
	data = append(data, []byte(remarkLabel)...)
	remark := fmt.Sprintf("%-48s\n", body.CaptainOrderBill.Note)
	data = append(data, []byte(remark)...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	data = append(data, []byte("\n\n\n\n\n")...)

	data = append(data, 0x1B, 0x61, 0x01) // Center alignment

	opCentered := utils.CenterInParentheses(body.CaptainOrderBill.Op, 20)
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

	parsedTime, _ := time.Parse(time.RFC3339Nano, body.CaptainOrderBill.Date)
	localTime := parsedTime.In(time.Local)
	postingDate := fmt.Sprintf("%-10s %-20s\n", "Posting Date:", localTime.Format("02/01/2006 15:04:05"))
	data = append(data, []byte(postingDate)...)
	printedWithAuditDate := fmt.Sprintf("%-8s %-20s %-5s %-17s\n", "Printed:", localTime.Format("02/01/2006 15:04:05"), "Audit:", localTime.Format("02/01/2006"))
	data = append(data, []byte(printedWithAuditDate)...)

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data)
}
