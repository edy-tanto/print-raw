//go:build !linux
// +build !linux

package driver_linux

import "fmt"

func Print(data []byte) {
	fmt.Println("driver linux")
}
