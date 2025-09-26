package httpserver

import (
	"net/http"
)

// NewRouter builds and returns the application's router.
// Accept handlers as http.Handler (not http.HandlerFunc).
func NewRouter(livez http.Handler, readyz http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/livez", livez)
	mux.Handle("/readyz", readyz)
	return mux
}
