package printer

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"fmt"
	"strings"
	"time"
)

func ExecutePrintTableCheck(body dto.PrintTableCheckRequestBody) {
	// ESC/POS commands
	data := []byte{
		0x1B, 0x40, // Initialize printer
		0x1B, 0x61, 0x00, // Left Allignment
	}

	data = append(data, []byte(strings.Repeat("=", 48))...)

	// Header
	codeWithOp := fmt.Sprintf("%-6s %-15s %-8s %-14s\n", "CO:", body.TableCheck.Code, "Waitress:", body.TableCheck.Op)
	data = append(data, []byte(codeWithOp)...)
	tableOrRoomNumberWithCustomerCount := fmt.Sprintf("%-6s %-10s %-13s %-2d/ %-2d\n", "Table/Room:", body.TableCheck.TableOrRoomNumber, "#Adult/#Child:", body.TableCheck.CustomerAdultCount, body.TableCheck.CustomerChildCount)
	data = append(data, []byte(tableOrRoomNumberWithCustomerCount)...)
	customerName := fmt.Sprintf("%-6s %-15s\n", "Guest Name:", body.TableCheck.CustomerName)
	data = append(data, []byte(customerName)...)

	data = append(data, []byte("\n\n")...)

	// Content
	data = append(data, []byte(strings.Repeat("-", 48))...)
	columnName := fmt.Sprintf("%-30s %-14s\n", "Product", "Qty")
	data = append(data, []byte(columnName)...)
	data = append(data, []byte(strings.Repeat("-", 48))...)

	for _, detail := range body.TableCheck.TableCheckDetails {
		detailText := fmt.Sprintf(
			"%-30s %-14d\n",
			detail.Item,
			detail.Qty,
		)
		data = append(data, []byte(detailText)...)
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)
	totalQty := fmt.Sprintf("%-10s %-15d\n", "Quantity:", body.TableCheck.TotalQty)
	data = append(data, []byte(totalQty)...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	data = append(data, []byte("\n\n")...)

	parsedTime, _ := time.Parse(time.RFC3339Nano, body.TableCheck.Date)
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