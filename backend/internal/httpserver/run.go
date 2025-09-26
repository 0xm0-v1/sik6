package httpserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/0xm0-v1/sik6/internal/config"
)

// Run starts the HTTP server, waits for an interrupt, and performs a graceful shutdown.
func Run(ctx context.Context, cfg *config.Config, handler http.Handler) error {
	srv := NewServer(cfg, handler)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	log.Printf("server starting... On port: %d", cfg.Port)

	// Derive a context canceled by SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	select {
	case <-ctx.Done():
		// proceed to shutdown
	case err := <-errCh:
		return err
	}

	// Graceful shutdown with timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
		_ = srv.Close() // best-effort force close
	}

	log.Println("server stopped cleanly")
	return nil
}
