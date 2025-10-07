package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xm0-v1/sik6/internal/recipe"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const recipeColumns = `
    id,
    recipe_name,
    created_at,
    updated_at,
    deleted_at
`

type RecipeRepository struct {
	pool *pgxpool.Pool
}

func (r *RecipeRepository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

var _ recipe.Repository = (*RecipeRepository)(nil)

func NewRecipeRepository(pool *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{pool: pool}
}

func (r *RecipeRepository) Create(ctx context.Context, name string) (string, error) {
	const query = `
        INSERT INTO recipes (recipe_name)
        VALUES ($1)
        RETURNING id
    `

	var id uuid.UUID
	if err := r.pool.QueryRow(ctx, query, name).Scan(&id); err != nil {
		return "", fmt.Errorf("create recipe: %w", err)
	}

	return id.String(), nil
}

func (r *RecipeRepository) GetByID(ctx context.Context, id string) (*recipe.Recipe, error) {
	const query = `
        SELECT` + recipeColumns + `
        FROM recipes
        WHERE id = $1 AND deleted_at IS NULL
    `

	return r.fetchOne(ctx, query, id)
}

func (r *RecipeRepository) GetByName(ctx context.Context, name string) (*recipe.Recipe, error) {
	const query = `
        SELECT` + recipeColumns + `
        FROM recipes
        WHERE recipe_name = $1 AND deleted_at IS NULL
    `

	return r.fetchOne(ctx, query, name)
}

func (r *RecipeRepository) List(ctx context.Context, limit, offset int) ([]*recipe.Recipe, error) {
	const query = `
        SELECT` + recipeColumns + `
        FROM recipes
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	if limit <= 0 {
		limit = 20
	}

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list recipes: %w", err)
	}
	defer rows.Close()

	var recipes []*recipe.Recipe
	for rows.Next() {
		item, err := scanRecipe(rows)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list recipes: %w", err)
	}

	return recipes, nil
}

func (r *RecipeRepository) Rename(ctx context.Context, id string, newName string) error {
	const query = `
        UPDATE recipes
        SET recipe_name = $1
        WHERE id = $2 AND deleted_at IS NULL
    `

	tag, err := r.pool.Exec(ctx, query, newName, id)
	if err != nil {
		return fmt.Errorf("rename recipe: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *RecipeRepository) SoftDelete(ctx context.Context, id string) error {
	const query = `
        UPDATE recipes
        SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `

	tag, err := r.pool.Exec(ctx, query, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("soft delete recipe: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *RecipeRepository) fetchOne(ctx context.Context, query string, args ...any) (*recipe.Recipe, error) {
	row := r.pool.QueryRow(ctx, query, args...)
	rec, err := scanRecipe(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("fetch recipe: %w", err)
	}
	return rec, nil
}

func scanRecipe(row pgx.Row) (*recipe.Recipe, error) {
	var item recipe.Recipe
	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("scan recipe: %w", err)
	}
	return &item, nil
}
