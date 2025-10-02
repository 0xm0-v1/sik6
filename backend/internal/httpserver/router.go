// internal/httpserver/router.go
package httpserver

import (
	"net/http"

	httpx "github.com/0xm0-v1/sik6/internal/http/transport"
)

// NewRouter builds and returns the application's router.
// Accept handlers as http.Handler (not http.HandlerFunc).
func NewRouter(
	livez http.Handler,
	readyz http.Handler,
	extra map[string]http.Handler,
) *http.ServeMux {
	return httpx.BuildRoutes(livez, readyz, extra)
}
