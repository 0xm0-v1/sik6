package httpserver

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/0xm0-v1/sik6/internal/config"
)

// ErrServerClosed mirrors http.ErrServerClosed for external checks.
var ErrServerClosed = http.ErrServerClosed

// Server wraps http.Server with configuration-derived settings.
type Server struct {
	srv *http.Server
}

// NewServer builds a configured http.Server with timeouts.
func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              cfg.Addr(),
			Handler:           handler,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			BaseContext: func(net.Listener) context.Context {
				return context.Background()
			},
		},
	}
}

// ListenAndServe starts the server and returns when it stops.
func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the server within the given context deadline.
func (s *Server) Shutdown(ctx context.Context) error {
	// http.Server.Shutdown stops accepting new connections and waits for in-flight
	// requests up to the context deadline. See Go docs.
	return s.srv.Shutdown(ctx)
}

// Close forces the server to close immediately (fallback after a timed-out Shutdown).
func (s *Server) Close() error {
	// After a Shutdown timeout, force-close to avoid dangling resources.
	// Callers should prefer Shutdown first.
	err := s.srv.Close()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
