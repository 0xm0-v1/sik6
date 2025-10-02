package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CORSConfig defines the allowed origins, methods, headers, and credentials
// for configuring Cross-Origin Resource Sharing (CORS).
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
}

// CORS creates a Middleware that applies CORS settings to HTTP requests
// according to the provided configuration.
func CORS(cfg CORSConfig) Middleware {
	opts := cors.Options{
		AllowedOrigins:   defaultIfEmpty(cfg.AllowedOrigins, []string{"http://localhost:4200"}),
		AllowedMethods:   defaultIfEmpty(cfg.AllowedMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}),
		AllowedHeaders:   defaultIfEmpty(cfg.AllowedHeaders, []string{"Authorization", "Content-Type", "Accept", "X-Requested-With"}),
		ExposedHeaders:   cfg.ExposedHeaders,
		AllowCredentials: cfg.AllowCredentials,
	}
	c := cors.New(opts)

	return func(next http.Handler) http.Handler {
		return c.Handler(next)
	}
}

// defaultIfEmpty returns the default slice if the provided one is empty.
func defaultIfEmpty(v, def []string) []string {
	if len(v) == 0 {
		return def
	}
	return v
}
