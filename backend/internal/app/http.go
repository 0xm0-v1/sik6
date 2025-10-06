package app

import (
	"context"
	"errors"
	stdhttp "net/http"
	"strings"

	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/health"
	healthhttp "github.com/0xm0-v1/sik6/internal/health/transport/http"
	"github.com/0xm0-v1/sik6/internal/http/middleware"
	"github.com/0xm0-v1/sik6/internal/httpserver"
	"github.com/0xm0-v1/sik6/internal/recipe"
	recipehttp "github.com/0xm0-v1/sik6/internal/recipe/transport/http"
	roothttp "github.com/0xm0-v1/sik6/internal/root/transport/http"
)

type Dependencies struct {
	RecipeRepository recipe.Repository
}

// NewHTTPHandler constructs the main HTTP handler for the application.
// It configures liveness and readiness endpoints, the root handler,
// and applies middleware such as CORS to the resulting handler.
func NewHTTPHandler(cfg *config.Config, deps Dependencies) stdhttp.Handler {
	var checker health.Checker = func(ctx context.Context) error {
		if deps.RecipeRepository == nil {
			return errors.New("recipe repository not configured")
		}
		return deps.RecipeRepository.Ping(ctx)
	}

	healthHandlers := healthhttp.NewHandlers(checker)
	apiToken := config.GetEnv("API_TOKEN", "")
	recipesHandlers := recipehttp.NewHandlers(deps.RecipeRepository, apiToken)
	rootHandlers := roothttp.NewHandlers()

	mux := httpserver.NewRouter(
		healthHandlers.Liveness,
		healthHandlers.Readiness,
		map[string]stdhttp.Handler{
			"/":         rootHandlers.Root,
			"/recipes":  recipesHandlers.Collection,
			"/recipes/": recipesHandlers.Resource,
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
