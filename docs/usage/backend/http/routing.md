# HTTP Routing Guide

This document explains how HTTP routes are organized in the BE ends and how to add new ones.

## Handler

```go
func NewHelloHandler() http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data: map[string]any{
				"message":   "Hello 👋",
				"component": "api",
				"type":      "hello",
				"time":      time.Now().UTC().Format(time.RFC3339Nano),
			},
		}
	})

	return guard(h)
}
```

## Wiring

Path: `backend/internal/app/wire.go` —— inside `NewHTTPHandler()` 
_(Make sure the feature package is imported at the top, e.g. `github.com/0xm0-v1/sik6/internal/hello`.)_

```go
	mux := httpserver.NewRouter(
		health.NewLivenessHandler(),
		health.NewReadinessHandler(checker),
		map[string]stdhttp.Handler{
			"/":          root.NewRootHandler(),
			"GET /hello": hello.NewHelloHandler(), // <- insert new route
		},
	)
```