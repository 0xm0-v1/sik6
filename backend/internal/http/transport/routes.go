package httpx

import "net/http"

func BuildRoutes(
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
