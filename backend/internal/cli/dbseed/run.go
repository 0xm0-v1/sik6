package dbseed

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/0xm0-v1/sik6/internal/cli/shared"
	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/storage/postgres"
	"github.com/0xm0-v1/sik6/internal/storage/postgres/fixtures"
)

// Run executes the database seeding workflow using the provided CLI arguments.
func Run(args []string) error {
	if err := config.LoadDevDotEnv(); err != nil {
		log.Printf("warning: could not load .env.development: %v", err)
	}

	fs := flag.NewFlagSet("dbseed", flag.ContinueOnError)
	var fixtureList string
	fs.StringVar(&fixtureList, "fixtures", "", "comma-separated list of fixture filenames (defaults to recipes.sql)")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx := context.Background()

	pool, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("database connect error: %w", err)
	}
	defer pool.Close()

	names := shared.SanitizeList(fs.Args())
	if len(names) == 0 {
		names = parseFixtureList(fixtureList)
	}

	if len(names) == 0 {
		if err := fixtures.SeedRecipes(ctx, pool); err != nil {
			return err
		}
		log.Println("applied default fixtures: recipes.sql")
		return nil
	}

	if err := fixtures.Seed(ctx, pool, names...); err != nil {
		return err
	}

	log.Printf("applied fixtures: %v", names)
	return nil
}

func parseFixtureList(list string) []string {
	list = strings.TrimSpace(list)
	if list == "" {
		return nil
	}
	parts := strings.Split(list, ",")
	return shared.SanitizeList(parts)
}
