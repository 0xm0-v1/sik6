package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
}

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

func defaultIfEmpty(v, def []string) []string {
	if len(v) == 0 {
		return def
	}
	return v
}
