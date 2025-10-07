# HTTP Routing Guide

This document explains how API routes are organised in the backend and how to add a new endpoint.

## Architecture Overview

```
internal/
  app/http.go                // wires handlers, middleware, and router
  http/transport/routes.go   // shared helpers for stdlib ServeMux
  <feature>/transport/http   // feature-specific handlers
  http/middleware            // reusable middlewares (CORS, auth, etc.)
```

Each feature exposes a `transport/http` package that returns an `http.Handler`. The application aggregates all feature handlers in `app.NewHTTPHandler`.

## Building a Handler

Create a new package under `internal/<feature>/transport/http` and return strongly-typed handlers:

```go
package examplehttp

import (
	"net/http"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

type Handlers struct {
	Collection http.Handler
}

func NewHandlers(repo Repository) Handlers {
	return Handlers{
		Collection: response.MethodGuard(http.MethodGet, http.MethodHead)(
			response.HeadAware(func(r *http.Request) (int, any) {
				items, err := repo.List(r.Context())
				if err != nil {
					return http.StatusInternalServerError, response.Envelope{
						Status: "error",
						Error:  "could not list examples",
						Data:   response.NewMeta("api", "examples:list"),
					}
				}

				payload := struct {
					Items any           `json:"items"`
					Meta  response.Meta `json:"meta"`
				}{
					Items: items,
					Meta:  response.NewMeta("api", "examples:list"),
				}

				return http.StatusOK, response.Envelope{
					Status: "ok",
					Data:   payload,
				}
			}),
		),
	}
}
```

Reusable helpers:

- `response.MethodGuard` keeps HTTP method checks consistent (405 + `Allow` header).
- `response.HeadAware` mirrors GET responses for HEAD requests.
- `response.Envelope` standardises the JSON structure (`status`, `data`, `error`).

When business logic needs authentication, wrap handlers with middleware from `internal/http/middleware`.

## Wiring the Route

Register your handler inside `app.NewHTTPHandler` (`backend/internal/app/http.go`):

```go
exampleHandlers := examplehttp.NewHandlers(deps.ExampleRepository)

mux := httpserver.NewRouter(
	healthHandlers.Liveness,
	healthHandlers.Readiness,
	map[string]http.Handler{
		"/":            rootHandlers.Root,
		"/recipes":     recipesHandlers.Collection,
		"/recipes/":    recipesHandlers.Resource,
		"/examples":    exampleHandlers.Collection, // new route
	},
)
```

Guidelines:

- Use `/path` for collection handlers and `/path/` (note the trailing slash) when you need to serve resources such as `/path/{id}`. The helper in `internal/http/transport/routes.go` keeps the root route (`/`) strict to avoid accidental shadowing.
- Keep dependencies in `app.Dependencies`; wire repositories or services there so new features follow the same pattern.
- If you introduce new middleware (auth, logging, rate limiting), extend `middleware.Chain` in `app.NewHTTPHandler` or wrap specific handlers directly inside your feature package.

## Database & Repository Layer

For a new table:

1. Create a migration under `internal/storage/postgres/migrations` (see `0001_init.sql` as a reference).
2. Add a repository implementation under `internal/storage/postgres/repository/<feature>`.
3. Expose a domain interface in `internal/<feature>/repository.go` so the rest of the code interacts through well-defined contracts.
4. Extend `app.Dependencies` and the constructor in `cmd/sik6/main.go` to initialise the repository.

Following this structure allows any new table to be surfaced through the API and later consumed by the frontend with minimal duplication.
