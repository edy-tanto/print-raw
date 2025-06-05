package printer

import (
	"edy-tanto/printer-pos/internal/print_raw/driver_windows"
	"edy-tanto/printer-pos/internal/print_web_service/dto"
	"fmt"
	"strings"
	"time"
)

func ExecutePrintKitchenEth(body dto.PrintKitchenEthRequestBody) {
	// ESC/POS commands
	data := []byte{
		0x1B, 0x40, // Initialize printer
		0x1B, 0x61, 0x00, // Left Allignment
	}

	data = append(data, []byte(strings.Repeat("-", 48))...)

	// Header
	codeWithOutlet := fmt.Sprintf("%-14s : %-15s %10s\n", "CO ID", body.Kitchen.Code, body.Kitchen.Outlet)
	data = append(data, []byte(codeWithOutlet)...)
	op := fmt.Sprintf("%-14s : %-20s\n", "Waitress", body.Kitchen.Op)
	data = append(data, []byte(op)...)
	tableNumber := fmt.Sprintf("%-14s : %-20s\n", "Table/Room", body.Kitchen.TableOrRoomNumber)
	data = append(data, []byte(tableNumber)...)
	customerName := fmt.Sprintf("%-14s : %-20s\n", "Table Number", body.Kitchen.CustomerName)
	data = append(data, []byte(customerName)...)
	location := fmt.Sprintf("%-14s : %-20s\n", "Location", "Kitchen")
	data = append(data, []byte(location)...)

	data = append(data, []byte("\n\n")...)

	// Content
	data = append(data, []byte(strings.Repeat("-", 48))...)
	columnName := fmt.Sprintf("%9s %-25s\n", "Qty", "Product")
	data = append(data, []byte(columnName)...)
	data = append(data, []byte(strings.Repeat("-", 48))...)

	for _, detail := range body.Kitchen.KitchenDetails {
		detailText := fmt.Sprintf(
			"%9d %-25s\n",
			detail.Qty,
			detail.Item,
		)
		data = append(data, []byte(detailText)...)
	}
	data = append(data, []byte(strings.Repeat("-", 48))...)

	data = append(data, 0x1D, 0x21, 0x00) // Reset to normal size

	data = append(data, []byte("\n\n")...)

	parsedTime, _ := time.Parse(time.RFC3339Nano, body.Kitchen.Date)
	localTime := parsedTime.In(time.Local)
	date := fmt.Sprintf("%-10s %-20s\n", "Printed:", localTime.Format("02/01/2006 15:04:05"))
	data = append(data, []byte(date)...)

	// Footer
	data = append(data, 0x1B, 0x64, 0x04) // Feed 4 lines
	data = append(data, 0x1D, 0x56, 0x00) // Full cut

	driver_windows.PrintEth(data, body.Kitchen.TargetPrinter)
}
