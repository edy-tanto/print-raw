package driver_windows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var ignoredPrinters = map[string]struct{}{
	"OneNote for Windows 10": {},
	"OneNote (Desktop)":      {},
	"Microsoft Print to PDF": {},
	"Fax":                    {},
}

// ListPrinters returns printer names that are currently Ready (idle) on the local machine.
func ListPrinters() ([]string, error) {
	if runtime.GOOS != "windows" {
		return nil, ErrPrinterEnumerationUnsupported
	}

	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-Command",
		"@(Get-CimInstance -ClassName Win32_Printer | Where-Object { $_.PrinterStatus -eq 3 -and -not $_.WorkOffline } | Select-Object -ExpandProperty Name) | ConvertTo-Json -Compress",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("enumerating printers: %w: %s", err, strings.TrimSpace(string(output)))
	}

	output = bytes.TrimSpace(output)
	if len(output) == 0 {
		return []string{}, nil
	}

	var names []string
	if err := json.Unmarshal(output, &names); err != nil {
		return nil, fmt.Errorf("parsing printer list: %w", err)
	}

	filtered := names[:0]
	for _, name := range names {
		if _, skip := ignoredPrinters[name]; skip {
			continue
		}
		filtered = append(filtered, name)
	}

	return filtered, nil
}
