//go:build windows

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

func runWindowsService(server *http.Server) (bool, error) {
	isService, err := svc.IsWindowsService()
	if err != nil {
		return false, fmt.Errorf("detecting Windows service context failed: %w", err)
	}

	if !isService {
		return false, nil
	}

	logger := newServiceLogger(serviceName)
	defer logger.Close()

	handler := &serviceHandler{
		server: server,
		logger: logger,
	}

	if err := svc.Run(serviceName, handler); err != nil {
		return false, fmt.Errorf("run windows service: %w", err)
	}

	return true, nil
}

type serviceHandler struct {
	server *http.Server
	logger *serviceLogger
}

func (h *serviceHandler) Execute(_ []string, requests <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	status <- svc.Status{State: svc.StartPending}

	errCh := make(chan error, 1)

	go func() {
		err := h.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		errCh <- err
	}()

	status <- svc.Status{
		State:   svc.Running,
		Accepts: svc.AcceptStop | svc.AcceptShutdown,
	}
	_ = h.logger.Info("service started")

	for {
		select {
		case c := <-requests:
			switch c.Cmd {
			case svc.Interrogate:
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				status <- svc.Status{State: svc.StopPending}

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err := h.server.Shutdown(ctx)
				cancel()
				if err != nil {
					_ = h.logger.Error(fmt.Sprintf("graceful shutdown failed: %v", err))
				}

				if err = <-errCh; err != nil {
					_ = h.logger.Error(fmt.Sprintf("http server stopped with error: %v", err))
				} else {
					_ = h.logger.Info("service stopped")
				}

				status <- svc.Status{State: svc.Stopped}
				return false, 0
			default:
			}
		case err := <-errCh:
			if err != nil {
				_ = h.logger.Error(fmt.Sprintf("http server error: %v", err))
				status <- svc.Status{State: svc.Stopped}
				return false, 1
			}

			_ = h.logger.Info("service exited")
			status <- svc.Status{State: svc.Stopped}
			return false, 0
		}
	}
}

type serviceLogger struct {
	eventLog *eventlog.Log
}

func newServiceLogger(name string) *serviceLogger {
	l, err := eventlog.Open(name)
	if err != nil {
		log.Printf("warn: unable to open Windows event log for %q: %v", name, err)
		return &serviceLogger{}
	}

	return &serviceLogger{eventLog: l}
}

func (l *serviceLogger) Close() {
	if l.eventLog != nil {
		_ = l.eventLog.Close()
	}
}

func (l *serviceLogger) Info(message string) error {
	if l.eventLog != nil {
		return l.eventLog.Info(1, message)
	}

	log.Printf("INFO: %s", message)
	return nil
}

func (l *serviceLogger) Error(message string) error {
	if l.eventLog != nil {
		return l.eventLog.Error(1, message)
	}

	log.Printf("ERROR: %s", message)
	return nil
}
