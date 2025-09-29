package app

import (
	"context"
	stdhttp "net/http"
	"os"
	"strings"

	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/health"
	"github.com/0xm0-v1/sik6/internal/hello"
	"github.com/0xm0-v1/sik6/internal/http/middleware"
	"github.com/0xm0-v1/sik6/internal/httpserver"
	"github.com/0xm0-v1/sik6/internal/root"
)

func NewHTTPHandler(cfg *config.Config) stdhttp.Handler {
	checker := func(ctx context.Context) error { return nil }

	mux := httpserver.NewRouter(
		health.NewLivenessHandler(),
		health.NewReadinessHandler(checker),
		map[string]stdhttp.Handler{
			"/":          root.NewRootHandler(),
			"GET /hello": hello.NewHelloHandler(),
		},
	)

	allowed := strings.Split(getenv("CORS_ALLOWED_ORIGINS", "http://localhost:4200"), ",")
	for i := range allowed {
		allowed[i] = strings.TrimSpace(allowed[i])
	}

	return middleware.Chain(mux,
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins:   allowed,
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
			AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "X-Requested-With"},
			ExposedHeaders:   []string{"Content-Length", "Content-Type"},
			AllowCredentials: getenv("CORS_ALLOW_CREDENTIALS", "false") == "true",
		}),
	)
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
