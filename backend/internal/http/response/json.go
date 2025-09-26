package response

import (
	"encoding/json"
	"net/http"
)

// Envelope is a small, typed JSON envelope for consistent responses.
type Envelope struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// WriteJSON writes a JSON response with the given status code.
// It sets Content-Type *before* calling WriteHeader.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	// It's OK for HEAD to reach here if a handler accidentally calls WriteJSON,
	// because most clients ignore the body of HEAD responses. We still prefer
	// to skip encoding in HEAD paths via HeadAware (see methods.go).
	_ = json.NewEncoder(w).Encode(v)
}

// WriteNoBody writes only headers and status (used for HEAD or empty bodies).
func WriteNoBody(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}
