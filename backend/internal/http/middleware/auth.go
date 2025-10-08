package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// BearerConfig captures the configuration for bearer token enforcement.
type BearerConfig struct {
	Token        string
	Methods      []string
	Unauthorized http.Handler
}

// RequireBearerToken returns a middleware enforcing bearer token authentication
// for the configured HTTP methods. If Token is empty the middleware is a no-op.
func RequireBearerToken(cfg BearerConfig) Middleware {
	token := strings.TrimSpace(cfg.Token)
	if token == "" {
		return func(next http.Handler) http.Handler { return next }
	}

	methodSet := make(map[string]struct{}, len(cfg.Methods))
	for _, m := range cfg.Methods {
		if m == "" {
			continue
		}
		methodSet[strings.ToUpper(m)] = struct{}{}
	}
	requireAll := len(methodSet) == 0

	unauthorized := cfg.Unauthorized
	if unauthorized == nil {
		unauthorized = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"status":"error","error":"unauthorized"}`))
		})
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !requireAll {
				if _, ok := methodSet[r.Method]; !ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			header := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(header, "Bearer ") {
				unauthorized.ServeHTTP(w, r)
				return
			}

			candidate := strings.TrimSpace(header[len("Bearer "):])
			if subtle.ConstantTimeCompare([]byte(candidate), []byte(token)) != 1 {
				unauthorized.ServeHTTP(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
