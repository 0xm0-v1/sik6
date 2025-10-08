package fixtures

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	recipesFixture    = "recipes.sql"
	devRecipesFixture = "recipes_dev.sql"
)

// SeedRecipes applies the default recipes seed data used by integration tests.
func SeedRecipes(ctx context.Context, pool *pgxpool.Pool) error {
	return Seed(ctx, pool, recipesFixture)
}

// SeedDevRecipes applies non-destructive seed data suited for local development.
func SeedDevRecipes(ctx context.Context, pool *pgxpool.Pool) error {
	return Seed(ctx, pool, devRecipesFixture)
}
