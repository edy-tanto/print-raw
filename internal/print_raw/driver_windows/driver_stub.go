//go:build !windows
// +build !windows

package driver_windows

func Print(data []byte) {
	println("stub Print called - not supported on this OS")
}

func ImageToBytes(filePath string, maxWidth int) ([]byte, int, int, error) {
	return []byte{}, 0, 0, nil
}
