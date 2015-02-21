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
	tcpL      *gracefulListener
	waitGroup *sync.WaitGroup
}

// New creates a HTTP service listening on the specified port
func New(port int) (*Service, error) {
	s := &Service{
		ServeMux:  http.NewServeMux(),
		waitGroup: &sync.WaitGroup{},
	}

	tcpL, err := newGracefulListener(port)
	if err != nil {
		return nil, err
	}
	s.tcpL = tcpL

	return s, nil
}

// Run starts the HTTP service responding to requests
func (s *Service) Run() {
	server := http.Server{Handler: s.ServeMux}

	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		server.Serve(s.tcpL)
	}()
}

// Stop shuts down the HTTP service
func (s *Service) Stop() {
	s.tcpL.Stop()
	s.waitGroup.Wait()
}
