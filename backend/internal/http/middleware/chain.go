package middleware

import "net/http"

// Middleware represents a function that wraps an http.Handler
// with additional behavior such as logging, CORS, or authentication.
type Middleware func(http.Handler) http.Handler

// Chain applies a sequence of middlewares around the given handler.
// Middlewares are applied in reverse order, so the first middleware
// in the list will be the outermost wrapper.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
