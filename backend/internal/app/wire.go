package app

import (
	"context"
	stdhttp "net/http"
	"strings"

	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/health"
	"github.com/0xm0-v1/sik6/internal/http/middleware"
	"github.com/0xm0-v1/sik6/internal/httpserver"
	"github.com/0xm0-v1/sik6/internal/root"
)

// NewHTTPHandler constructs the main HTTP handler for the application.
// It configures liveness and readiness endpoints, the root handler,
// Middleware such as CORS is applied to the resulting handler.
func NewHTTPHandler(cfg *config.Config) stdhttp.Handler {
	checker := func(ctx context.Context) error { return nil }

	mux := httpserver.NewRouter(
		health.NewLivenessHandler(),
		health.NewReadinessHandler(checker),
		map[string]stdhttp.Handler{
			"/": root.NewRootHandler(),
		},
	)

	allowed := strings.Split(config.GetEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200"), ",")
	for i := range allowed {
		allowed[i] = strings.TrimSpace(allowed[i])
	}

	return middleware.Chain(mux,
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins:   allowed,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "X-Requested-With"},
			ExposedHeaders:   []string{"Content-Length", "Content-Type"},
			AllowCredentials: config.GetEnv("CORS_ALLOW_CREDENTIALS", "false") == "true",
		}),
	)
}
