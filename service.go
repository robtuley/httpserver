// Package httpserver adds graceful shutdown & health-check
// utilities to the standard http.Server
package httpserver

import (
	"net/http"
	"sync"
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

// Err from service running if any
func (s *Service) Err() error {
	return s.runErr
}
