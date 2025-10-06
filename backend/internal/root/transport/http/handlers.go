package roothttp

import (
	"net/http"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

type Handlers struct {
	Root http.Handler
}

func NewHandlers() Handlers {
	guard := response.MethodGuard(http.MethodGet, http.MethodHead)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		payload := struct {
			Message string        `json:"message"`
			Meta    response.Meta `json:"meta"`
		}{
			Message: "Welcome to the Root URL",
			Meta:    response.NewMeta("api", "root"),
		}

		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data:   payload,
		}
	})

	return Handlers{Root: guard(h)}
}
