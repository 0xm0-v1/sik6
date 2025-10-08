package main

import (
	"context"
	"log"

	"github.com/0xm0-v1/sik6/internal/app"
	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/httpserver"
	"github.com/0xm0-v1/sik6/internal/recipe"
	"github.com/0xm0-v1/sik6/internal/storage/postgres"
	recipepg "github.com/0xm0-v1/sik6/internal/storage/postgres/repository/recipe"
)

func main() {
	if err := config.LoadDevDotEnv(); err != nil {
		log.Printf("warning: could not load .env.development: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	ctx := context.Background()

	pool, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("database connect error: %v", err)
	}
	defer pool.Close()

	recipeRepo := recipepg.NewRecipeRepository(pool)
	recipeSvc := recipe.NewService(recipeRepo)
	handler := app.NewHTTPHandler(cfg, app.Dependencies{
		RecipeRepository: recipeRepo,
		RecipeService:    recipeSvc,
	})

	if err := httpserver.Run(ctx, cfg.Server, handler); err != nil {
		log.Fatalf("application run error: %v", err)
	}
}
