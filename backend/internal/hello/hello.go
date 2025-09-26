// internal/hello/hello.go
package hello

import (
	"net/http"
	"time"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

func NewHelloHandler() http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data: map[string]any{
				"message":   "Hello 👋",
				"component": "api",
				"type":      "hello",
				"time":      time.Now().UTC().Format(time.RFC3339Nano),
			},
		}
	})

	return guard(h)
}
