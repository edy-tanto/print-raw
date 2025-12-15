package driver_windows

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

var ignoredPrinters = map[string]struct{}{
	"OneNote for Windows 10": {},
	"OneNote (Desktop)":      {},
	"Microsoft Print to PDF":        {},
	"Fax":                           {},
	"Microsoft XPS Document Writer": {},
}

// ListPrinters returns printer names that are currently Ready (idle) on the local machine.
// This includes both local printers and shared/network printers. For shared printers,
// the function uses more flexible filtering to account for Windows status reporting quirks.
func ListPrinters() ([]string, error) {
	if runtime.GOOS != "windows" {
		return nil, ErrPrinterEnumerationUnsupported
	}

	script := legacyPrinterQuery
	if major, err := windowsMajorVersion(); err != nil {
		fmt.Fprintf(os.Stderr, "windows version detection failed: %v\n", err)
	} else {
		if major == 10 {
			script = modernPrinterQuery
		}
	}

	useModernQuery := script == modernPrinterQuery
	names, err := fetchPrinterNames(script, !useModernQuery)
	if err != nil && useModernQuery {
		fmt.Fprintf(os.Stderr, "modern printer query failed (%v); falling back to legacy WMI\n", err)
		names, err = fetchPrinterNames(legacyPrinterQuery, true)
	}
	if err != nil {
		return nil, err
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

const (
	// modernPrinterQuery uses CIM for Windows 10+ with improved shared printer support.
	// Filter logic:
	// - Local printers: Status 3 (Idle) and not WorkOffline (strict)
	// - Shared/Network printers: Status 1,2,3,4,5,6 (Other, Unknown, Idle, Printing, Warmup, Stopped printing) and not Status 7 (Offline), WorkOffline ignored
	//   Note: Status 4,5,6 are included to handle printer during/after print job completion
	modernPrinterQuery = `@(Get-CimInstance -ClassName Win32_Printer | Where-Object { ($_.PrinterStatus -eq 3 -and -not $_.WorkOffline) -or (($_.Shared -eq $true -or $_.Network -eq $true) -and $_.PrinterStatus -in @(1,2,3,4,5,6) -and $_.PrinterStatus -ne 7) } | Select-Object -ExpandProperty Name) | ConvertTo-Json -Compress`
	// legacyPrinterQuery uses WMI for Windows 7 with improved shared printer support.
	// Same filter logic as modernPrinterQuery but using WMI syntax.
	legacyPrinterQuery = `'[' + ((Get-WmiObject -Class Win32_Printer | Where-Object { ($_.PrinterStatus -eq 3 -and -not $_.WorkOffline) -or (($_.Shared -eq $true -or $_.Network -eq $true) -and ($_.PrinterStatus -eq 1 -or $_.PrinterStatus -eq 2 -or $_.PrinterStatus -eq 3 -or $_.PrinterStatus -eq 4 -or $_.PrinterStatus -eq 5 -or $_.PrinterStatus -eq 6) -and $_.PrinterStatus -ne 7) } | ForEach-Object { $escaped = $_.Name.Replace('\', '\\').Replace('"', '\"'); '"' + $escaped + '"' }) -join ',') + ']'`
)

func fetchPrinterNames(script string, sanitize bool) ([]string, error) {
	output, err := runPowerShell(script)
	if err != nil {
		return nil, err
	}

	output = bytes.TrimSpace(output)
	if sanitize {
		output = sanitizeLegacyJSON(output)
	}
	if len(output) == 0 {
		return []string{}, nil
	}

	var names []string
	if err := json.Unmarshal(output, &names); err != nil {
		return nil, fmt.Errorf("parsing printer list: %w", err)
	}

	return names, nil
}

func windowsMajorVersion() (int, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return 0, fmt.Errorf("open version registry key: %w", err)
	}
	defer key.Close()

	if major, _, err := key.GetIntegerValue("CurrentMajorVersionNumber"); err == nil {
		return int(major), nil
	}

	versionString, _, err := key.GetStringValue("CurrentVersion")
	if err != nil {
		return 0, fmt.Errorf("read windows version: %w", err)
	}

	dot := strings.Index(versionString, ".")
	if dot == -1 {
		dot = len(versionString)
	}

	major, err := strconv.Atoi(versionString[:dot])
	if err != nil {
		return 0, fmt.Errorf("parse windows major version: %w", err)
	}

	return major, nil
}

func runPowerShell(script string) ([]byte, error) {
	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-Command",
		script,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("enumerating printers: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	if stderr.Len() > 0 {
		return nil, fmt.Errorf("enumerating printers: %s", strings.TrimSpace(stderr.String()))
	}

	return stdout.Bytes(), nil
}

func sanitizeLegacyJSON(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	out := make([]byte, 0, len(data))
	for i := 0; i < len(data); i++ {
		b := data[i]
		if b != '\\' {
			if b == '\r' {
				out = append(out, '\\', 'r')
				continue
			}
			if b == '\n' {
				out = append(out, '\\', 'n')
				continue
			}
			out = append(out, b)
			continue
		}

		if i+1 >= len(data) {
			out = append(out, '\\', '\\')
			break
		}

		next := data[i+1]
		switch next {
		case '\\', '"', '/', 'b', 'f', 'n', 'r', 't':
			out = append(out, '\\', next)
			i++
		case 'u':
			if i+5 < len(data) {
				out = append(out, '\\', 'u', data[i+2], data[i+3], data[i+4], data[i+5])
				i += 5
			} else {
				out = append(out, '\\', '\\')
			}
		case '\r':
			out = append(out, '\\', 'r')
			i++
		case '\n':
			out = append(out, '\\', 'n')
			i++
		default:
			out = append(out, '\\', '\\', next)
			i++
		}
	}

	return out
}
