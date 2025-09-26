package response

import (
	"net/http"
	"strings"
)

// MethodGuard returns a middleware that allows only the provided methods.
// If the request method is not allowed, it responds with 405 and sets "Allow".
func MethodGuard(allowed ...string) func(http.Handler) http.Handler {
	// Build a fast lookup set.
	set := make(map[string]struct{}, len(allowed))
	for _, m := range allowed {
		set[strings.ToUpper(m)] = struct{}{}
	}

	allowHeader := strings.Join(allowed, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := set[r.Method]; !ok {
				w.Header().Set("Allow", allowHeader)
				WriteNoBody(w, http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// JSONResponder produces the HTTP status and payload for a response.
type JSONResponder func(r *http.Request) (int, any)

// HeadAware wraps a JSONResponder so that HEAD mirrors GET headers/status,
// but omits the body, per HTTP semantics.
func HeadAware(responder JSONResponder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, payload := responder(r)
		if r.Method == http.MethodHead {
			WriteNoBody(w, status)
			return
		}
		WriteJSON(w, status, payload)
	}
}
