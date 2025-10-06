package health

import (
	"net/http"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

func NewReadinessHandler(check Checker) http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		if err := check(r.Context()); err != nil {
			return http.StatusServiceUnavailable, response.Envelope{
				Status: "error",
				Error:  err.Error(),
				Data:   response.NewMeta("api", "readiness"),
			}
		}
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data:   response.NewMeta("api", "readiness"),
		}
	})

	return guard(h)
}
