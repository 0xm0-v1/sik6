package cismoke

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/0xm0-v1/sik6/internal/cli/shared"
	"github.com/0xm0-v1/sik6/internal/config"
	"github.com/0xm0-v1/sik6/internal/recipe"
)

// Run executes the CI smoke checks for the recipes endpoint.
func Run(args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	fs := flag.NewFlagSet("cismoke", flag.ContinueOnError)

	defaultExpected := strings.Join(cfg.Smoke.Expected, ",")
	url := fs.String("url", cfg.Smoke.RecipesURL, "recipes endpoint URL to validate")
	expectedRaw := fs.String("expected", defaultExpected, "comma-separated list of recipe names that must be present")
	timeout := fs.Duration("timeout", cfg.Smoke.Timeout, "HTTP request timeout")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	expected := shared.SanitizeList(strings.Split(*expectedRaw, ","))
	if len(expected) == 0 {
		expected = append([]string(nil), cfg.Smoke.Expected...)
	}

	client := &http.Client{Timeout: *timeout}
	return verifyRecipes(context.Background(), client, *url, expected)
}

func verifyRecipes(ctx context.Context, client *http.Client, url string, expected []string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("perform request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("warning: close response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	var envelope struct {
		Status string `json:"status"`
		Error  string `json:"error"`
		Data   struct {
			Recipes []*recipe.Recipe `json:"recipes"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if envelope.Status != "ok" {
		return fmt.Errorf("api status %q (error: %s)", envelope.Status, envelope.Error)
	}

	if len(envelope.Data.Recipes) == 0 {
		return fmt.Errorf("no recipes returned")
	}

	available := make(map[string]struct{}, len(envelope.Data.Recipes))
	for _, item := range envelope.Data.Recipes {
		if item != nil && item.Name != "" {
			available[item.Name] = struct{}{}
		}
	}

	var missing []string
	for _, want := range expected {
		if _, ok := available[want]; !ok {
			missing = append(missing, want)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing expected recipes: %s", strings.Join(missing, ", "))
	}

	return nil
}
