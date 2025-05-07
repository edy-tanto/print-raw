package driver_windows

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// Printer represents a simplified structure for printer data from PowerShell.
type Printer struct {
	Name string `json:"Name"`
}

func list_printers() {
	// Execute PowerShell command to get printer list in JSON format.
	cmd := exec.Command("powershell", "-Command", "Get-Printer | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}

	// Parse the JSON output into a slice of Printer structs.
	var printers []Printer
	err = json.Unmarshal(output, &printers)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check if any printers were found and list them.
	if len(printers) == 0 {
		fmt.Println("No printers found.")
	} else {
		fmt.Println("Available printers:")
		for _, printer := range printers {
			fmt.Println(printer.Name)
		}
	}
}
