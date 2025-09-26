package health

import (
	"net/http"
	"time"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

func NewReadinessHandler(check Checker) http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		// Bind checker to the request context with a short, caller-controlled deadline if needed.
		if err := check(r.Context()); err != nil {
			return http.StatusServiceUnavailable, response.Envelope{
				Status: "error",
				Error:  err.Error(),
				Data: map[string]any{
					"component": "api",
					"type":      "readiness",
					"time":      time.Now().UTC().Format(time.RFC3339Nano),
				},
			}
		}
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data: map[string]any{
				"component": "api",
				"type":      "readiness",
				"time":      time.Now().UTC().Format(time.RFC3339Nano),
			},
		}
	})

	return guard(h)
}
