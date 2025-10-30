//go:build !windows

package main

import "net/http"

func runWindowsService(_ *http.Server) (bool, error) {
	return false, nil
}
