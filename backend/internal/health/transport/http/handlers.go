package healthhttp

import (
	"context"
	"net/http"

	"github.com/0xm0-v1/sik6/internal/health"
	"github.com/0xm0-v1/sik6/internal/http/response"
)

type Handlers struct {
	Liveness  http.Handler
	Readiness http.Handler
}

func NewHandlers(checker health.Checker) Handlers {
	liveness := buildLiveness()
	readiness := buildReadiness(checker)

	return Handlers{
		Liveness:  liveness,
		Readiness: readiness,
	}
}

func buildLiveness() http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data:   response.NewMeta("api", "liveness"),
		}
	})

	return guard(h)
}

func buildReadiness(checker health.Checker) http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		if checker == nil {
			checker = func(context.Context) error { return nil }
		}
		if err := checker(r.Context()); err != nil {
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
