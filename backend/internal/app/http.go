package app

import (
	"context"
	stdhttp "net/http"

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
	RecipeService    recipe.Service
}

// NewHTTPHandler constructs the main HTTP handler for the application.
// It configures liveness and readiness endpoints, the root handler,
// and applies middleware such as CORS to the resulting handler.

func NewHTTPHandler(cfg *config.Config, deps Dependencies) stdhttp.Handler {
	recipeSvc := deps.RecipeService
	if recipeSvc == nil && deps.RecipeRepository != nil {
		recipeSvc = recipe.NewService(deps.RecipeRepository)
	}

	var checker health.Checker = func(ctx context.Context) error {
		if recipeSvc == nil {
			return recipe.ErrRepositoryUnavailable
		}
		return recipeSvc.Ping(ctx)
	}

	healthHandlers := healthhttp.NewHandlers(checker)
	recipesHandlers := recipehttp.NewHandlers(recipeSvc, cfg.API.Token)
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

	return middleware.Chain(mux,
		middleware.CORS(middleware.CORSConfig{
			AllowedOrigins:   cfg.CORS.AllowedOrigins,
			AllowedMethods:   cfg.CORS.AllowedMethods,
			AllowedHeaders:   cfg.CORS.AllowedHeaders,
			ExposedHeaders:   cfg.CORS.ExposedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
		}),
	)
}
