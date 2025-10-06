package health

import (
	"net/http"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

func NewLivenessHandler() http.Handler {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data:   response.NewMeta("api", "liveness"),
		}
	})

	return guard(h)
}
