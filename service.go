// Package httpserver adds graceful shutdown & health-check
// utilities to the standard http.Server
package httpserver

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Service presents a http.ServeMux interface
type Service struct {
	*http.ServeMux
	HasStoppedC chan bool

	runErr    error
	tcpL      *gracefulListener
	startOnce sync.Once
}

// New creates a HTTP service listening on the specified port
func New(port int) (*Service, error) {
	s := &Service{
		ServeMux:    http.NewServeMux(),
		HasStoppedC: make(chan bool),
	}

	tcpL, err := newGracefulListener(port)
	if err != nil {
		return nil, err
	}
	s.tcpL = tcpL

	return s, nil
}

// Start the HTTP service responding to requests
func (s *Service) Start() {
	s.startOnce.Do(func() {
		go func() {
			err := http.Serve(s.tcpL, s.ServeMux)
			close(s.HasStoppedC)
			if err != nil && err != s.tcpL.StoppedErr {
				s.runErr = err
			}
		}()
	})
}

// Stop shuts down the HTTP service
func (s *Service) Stop() {
	s.tcpL.Stop()
	<-s.HasStoppedC
}

// WaitStop waits for a OS term signal, then stops
func (s *Service) WaitStop() os.Signal {
	var sig os.Signal
	sigC := make(chan os.Signal)
	signal.Notify(sigC,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

WaitForStop:
	for {
		select {
		case sig = <-sigC:
			s.Stop()
		case <-s.HasStoppedC:
			break WaitForStop
		}
	}

	return sig
}

// Err from service running if any
func (s *Service) Err() error {
	return s.runErr
}
