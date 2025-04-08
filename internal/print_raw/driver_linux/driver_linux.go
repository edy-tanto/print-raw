//go:build linux
// +build linux

package driver_linux

import (
	"log"
	"os"
	"os/exec"
)

func Print(data []byte) {
	// Write the receipt data to a temporary file
	tmpFile, err := os.CreateTemp("", "receipt-*.escpos")
	if err != nil {
		log.Fatalf("Unable to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		log.Fatalf("Unable to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatalf("Unable to close temporary file: %v", err)
	}

	// Execute the lp command to print using raw mode
	cmd := exec.Command("lp", "-d", "TM-T82-S-A", "-o", "raw", tmpFile.Name())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to execute lp command: %v\nOutput: %s", err, string(out))
	}

	log.Println("Receipt sent to printer successfully.")
}
