package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var dateString string

	// Cek apakah ada argument dari command line
	if len(os.Args) >= 2 {
		dateString = os.Args[1]
	} else {
		// Jika tidak ada argument, minta input dari user
		fmt.Println("=== Test Timezone Conversion ===")
		fmt.Println("\nMasukkan string tanggal (contoh: 2026-01-03T18:08:47.000Z atau 2026-01-03T18:08:47.000+07:00)")
		fmt.Print("Input: ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ Error membaca input: %v\n", err)
			os.Exit(1)
		}

		// Hapus newline dan whitespace
		dateString = strings.TrimSpace(input)

		if dateString == "" {
			fmt.Println("\n⚠️  Input kosong, menggunakan contoh default...")
			dateString = "2026-01-03T18:08:47.000Z"
		}
		fmt.Println()
	}

	fmt.Println("=== Test Timezone Conversion ===\n")
	fmt.Printf("Input string: %s\n\n", dateString)

	// Parse waktu sesuai dengan format RFC3339Nano (seperti di kode asli)
	parsedTime, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		fmt.Printf("❌ Error parsing: %v\n", err)
		fmt.Println("\nPastikan format sesuai RFC3339Nano, contoh:")
		fmt.Println("  - 2026-01-03T18:08:47.000Z (UTC)")
		fmt.Println("  - 2026-01-03T18:08:47.000+07:00 (Local time)")
		os.Exit(1)
	}

	// Informasi waktu setelah parsing
	fmt.Println("=== Hasil Parsing ===")
	fmt.Printf("Waktu yang di-parse: %s\n", parsedTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Timezone: %s\n", parsedTime.Location().String())
	fmt.Printf("Offset: %s\n", parsedTime.Format("-07:00"))
	fmt.Println()

	// Konversi ke local time (seperti di kode asli baris 168)
	localTime := parsedTime.In(time.Local)

	fmt.Println("=== Setelah .In(time.Local) ===")
	fmt.Printf("Waktu lokal: %s\n", localTime.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("Timezone: %s\n", localTime.Location().String())
	fmt.Printf("Selisih waktu: %s\n", localTime.Sub(parsedTime))
	fmt.Println()

	// Format output seperti di captain-order-invoice.go
	fmt.Println("=== Format Output (seperti di captain-order-invoice.go) ===")
	postingDate := fmt.Sprintf("%-10s %-20s", "Posting Date:", localTime.Format("02/01/2006 15:04:05"))
	fmt.Println(postingDate)

	printedWithAuditDate := fmt.Sprintf("%-8s %-20s %-5s %-17s",
		"Printed:",
		localTime.Format("02/01/2006 15:04:05"),
		"Audit:",
		localTime.Format("02/01/2006"))
	fmt.Println(printedWithAuditDate)
	fmt.Println()

	// Kesimpulan
	fmt.Println("=== Kesimpulan ===")
	trimmedDateString := strings.TrimSpace(dateString)
	if len(trimmedDateString) > 0 && trimmedDateString[len(trimmedDateString)-1] == 'Z' {
		fmt.Println("⚠️  Format dengan 'Z' (UTC) akan dikonversi ke timezone lokal")
		fmt.Printf("   Waktu berubah dari %s menjadi %s\n",
			parsedTime.Format("15:04:05"),
			localTime.Format("15:04:05"))
	} else {
		fmt.Println("✓ Format dengan timezone offset (seperti +07:00) sudah dalam timezone lokal")
		fmt.Printf("   Waktu tetap %s (tidak dikonversi lagi)\n", localTime.Format("15:04:05"))
	}
}
