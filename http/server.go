package http

import (
	"net/http"
	"time"

	"io"

	"context"

	"github.com/gorilla/handlers"
)

type Server struct {
	ImageHandler ImageHandler
	Addr         string
	server       *http.Server
}

// Start creates an http.Server and calls ListenAndServes (blocking).
func (s *Server) Start(logWriter io.Writer) error {
	s.server = &http.Server{
		Addr:         s.Addr,
		Handler:      handlers.LoggingHandler(logWriter, s.ImageHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop stops the server gracefully (hopefully)image_handler_test.go.
func (s *Server) Stop(ctx context.Context) error {
	s.server.SetKeepAlivesEnabled(false)
	return s.server.Shutdown(ctx)
}
