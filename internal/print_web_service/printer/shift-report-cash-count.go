package printer

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"edy-tanto/printer-pos/internal/print_web_service/utils"
	"fmt"
	"strings"
)

func ExecutePrintShiftReportCashCount(body dto.PrintShiftReportCashCountRequestBody) {
	data := []byte{
		0x1B, 0x40, // Initialize
		0x1B, 0x61, 0x00, // Left alignment
	}
	data = append(data, 0x1D, 0x21, 0x00) // Normal size text

	// Shift
	data = append(data, []byte(fmt.Sprintf("%-20s %27d\n", "Shift ID", body.ShiftReportCashCount.Shift.Id))...)
	data = append(data, []byte(fmt.Sprintf("%-20s %27s\n", "User", body.ShiftReportCashCount.Shift.UserName))...)
	data = append(data, []byte(fmt.Sprintf("%-20s %27s\n", "Time-In", utils.FormatDatetime(body.ShiftReportCashCount.Shift.ShiftStartAt)))...)
	data = append(data, []byte(fmt.Sprintf("%-20s %27s\n", "Time-Out", utils.FormatDatetime(body.ShiftReportCashCount.Shift.ShiftEndAt)))...)
	// NOTE: Primus "dihilangkan tidak apa"
	// data = append(data, []byte(fmt.Sprintf("%-20s %27s\n", "Opening Balance", utils.FormatMoney(body.ShiftReportCashCount.Shift.OpeningBalance)))...)

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x61, 0x01) // Center alignment
	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	data = append(data, []byte("Payment\n")...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	// Payment Summary
	for _, payment := range body.ShiftReportCashCount.PaymentSummary {
		detailText := fmt.Sprintf(
			"%-20s %27s\n",
			payment.PaymentMethod,
			utils.FormatMoney(payment.GrandTotal),
		)
		data = append(data, []byte(detailText)...)
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	data = append(data, []byte(fmt.Sprintf("%-20s %27s\n", "Total Payment", utils.FormatMoney(body.ShiftReportCashCount.TotalPayment)))...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n\n")...)

	data = append(data, 0x1B, 0x61, 0x01) // Center alignment
	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	data = append(data, []byte("Cash Count\n")...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off
	data = append(data, 0x1B, 0x61, 0x00) // Left alignment

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	for _, cashCount := range body.ShiftReportCashCount.CashCounts {
		detailText := fmt.Sprintf(
			"Rp %-17s %8d %18s\n",
			utils.FormatMoney(cashCount.Nominal),
			cashCount.Qty,
			utils.FormatMoney(cashCount.Total),
		)
		data = append(data, []byte(detailText)...)
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	data = append(data, 0x1B, 0x45, 0x01) // Turn bold on
	data = append(data, []byte(fmt.Sprintf("%-20s %8d %18s\n", "Total Cash Count", body.ShiftReportCashCount.CashCountsQtyTotal, utils.FormatMoney(body.ShiftReportCashCount.CashCountsGrandTotal)))...)
	data = append(data, 0x1B, 0x45, 0x00) // Turn bold off

	data = append(data, []byte(strings.Repeat("-", 48))...)
	data = append(data, []byte("\n")...)

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.Print(data)
}
