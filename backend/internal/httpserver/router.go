// internal/httpserver/router.go
package httpserver

import "net/http"

// NewRouter builds and returns the application's router.
// Accept handlers as http.Handler (not http.HandlerFunc).
func NewRouter(
	livez http.Handler,
	readyz http.Handler,
	extra map[string]http.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/livez", livez)
	mux.Handle("/readyz", readyz)

	for pattern, h := range extra {
		mux.Handle(pattern, h)
	}

	return mux
}
