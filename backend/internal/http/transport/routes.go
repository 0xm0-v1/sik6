package httpx

import "net/http"

// BuildRoutes constructs an HTTP mux with standard health endpoints
// (/livez and /readyz) and any additional routes provided in the map.
func BuildRoutes(
	livez http.Handler,
	readyz http.Handler,
	extra map[string]http.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/livez", livez)
	mux.Handle("/readyz", readyz)
	for pattern, h := range extra {
		if pattern == "/" {
			mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/" {
					http.NotFound(w, r)
					return
				}
				h.ServeHTTP(w, r)
			}))
		} else {
			mux.Handle(pattern, h)
		}
	}
	return mux
}
