package main

import (
	"log"
	"net/http"

	"edy-tanto/printer-pos/internal/print_web_service"
)

const (
	serviceName = "PrintRawWebService"
	serviceAddr = ":8080"
)

func main() {
	server := newHTTPServer()

	runningAsService, err := runWindowsService(server)
	if err != nil {
		log.Fatalf("failed to initialize Windows service: %v", err)
	}

	if runningAsService {
		return
	}

	log.Printf("starting HTTP server on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %v", err)
	}
}

func newHTTPServer() *http.Server {
	mux := print_web_service.NewPrintMuxServer()

	return &http.Server{
		Addr:    serviceAddr,
		Handler: mux,
	}
}
