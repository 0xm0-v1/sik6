package recipehttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xm0-v1/sik6/internal/recipe"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestSoftDeleteRespondsWithEnvelope(t *testing.T) {
	h := handler{
		repo:  stubRecipeRepo{},
		token: "",
	}

	req := httptest.NewRequest(http.MethodDelete, "/recipes/a1b2c3", nil)
	rec := httptest.NewRecorder()

	h.resource(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var payload responseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Status != "ok" {
		t.Fatalf("expected status ok, got %q", payload.Status)
	}
	if payload.Data.Meta.Type != "recipes:delete" {
		t.Errorf("expected meta.type recipes:delete, got %q", payload.Data.Meta.Type)
	}
	if payload.Data.Meta.Component != metaComponent {
		t.Errorf("expected meta.component %q, got %q", metaComponent, payload.Data.Meta.Component)
	}
}

func TestSoftDeleteNotFound(t *testing.T) {
	h := handler{
		repo:  stubRecipeRepo{softDeleteErr: pgx.ErrNoRows},
		token: "",
	}

	req := httptest.NewRequest(http.MethodDelete, "/recipes/missing", nil)
	rec := httptest.NewRecorder()

	h.resource(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	var payload responseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Status != "error" {
		t.Fatalf("expected response status error, got %q", payload.Status)
	}
	if payload.Error != "recipe not found" {
		t.Fatalf("expected error message 'recipe not found', got %q", payload.Error)
	}
}

func TestIsUniqueViolation(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", &pgconn.PgError{Code: "23505"})
	if !isUniqueViolation(err) {
		t.Fatal("expected unique violation to be detected")
	}

	if isUniqueViolation(errors.New("nope")) {
		t.Fatal("unexpected unique violation detected")
	}
}

type stubRecipeRepo struct {
	softDeleteErr error
}

func (s stubRecipeRepo) Create(context.Context, string) (string, error) { return "", nil }
func (s stubRecipeRepo) GetByID(context.Context, string) (*recipe.Recipe, error) {
	return nil, errors.New("not implemented")
}

func (s stubRecipeRepo) GetByName(context.Context, string) (*recipe.Recipe, error) {
	return nil, errors.New("not implemented")
}

func (s stubRecipeRepo) List(context.Context, int, int) ([]*recipe.Recipe, error) {
	return nil, errors.New("not implemented")
}

func (s stubRecipeRepo) Rename(context.Context, string, string) error {
	return errors.New("not implemented")
}
func (s stubRecipeRepo) SoftDelete(context.Context, string) error { return s.softDeleteErr }
func (s stubRecipeRepo) Ping(context.Context) error               { return nil }

type responseEnvelope struct {
	Status string `json:"status"`
	Data   struct {
		Meta struct {
			Component string `json:"component"`
			Type      string `json:"type"`
		} `json:"meta"`
	} `json:"data"`
	Error string `json:"error"`
}
