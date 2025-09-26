// internal/root/root.go
package root

import (
	"net/http"
	"time"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

func NewRootHandler() http.Handler {
	guard := response.MethodGuard(http.MethodGet)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data: map[string]any{
				"message":   "Welcome to the Root URL 👋",
				"component": "api",
				"type":      "root",
				"time":      time.Now().UTC().Format(time.RFC3339Nano),
			},
		}
	})

	return guard(h)
}
