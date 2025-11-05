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
	"OneNote for Windows 10":        {},
	"OneNote (Desktop)":             {},
	"Microsoft Print to PDF":        {},
	"Fax":                           {},
	"Microsoft XPS Document Writer": {},
}

// ListPrinters returns printer names that are currently Ready (idle) on the local machine.
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

	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-Command",
		script,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("enumerating printers: %w: %s", err, strings.TrimSpace(string(output)))
	}

	output = bytes.TrimSpace(output)
	output = sanitizeLegacyJSON(output)
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

const (
	modernPrinterQuery = `@(Get-CimInstance -ClassName Win32_Printer | Where-Object { $_.PrinterStatus -eq 3 -and -not $_.WorkOffline } | Select-Object -ExpandProperty Name) | ConvertTo-Json -Compress`
	legacyPrinterQuery = `'[' + ((Get-WmiObject -Class Win32_Printer | Where-Object { $_.PrinterStatus -eq 3 -and -not $_.WorkOffline } | ForEach-Object { $escaped = $_.Name.Replace('\', '\\').Replace('"', '\"'); '"' + $escaped + '"' }) -join ',') + ']'`
)

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
