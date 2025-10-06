// internal/root/root.go
package root

import (
	"net/http"

	"github.com/0xm0-v1/sik6/internal/http/response"
)

type rootPayload struct {
	Message string `json:"message"`
	response.Meta
}

func NewRootHandler() http.Handler {
	guard := response.MethodGuard(http.MethodGet)

	h := response.HeadAware(func(r *http.Request) (int, any) {
		payload := rootPayload{
			Message: "Welcome to the Root URL",
			Meta:    response.NewMeta("api", "root"),
		}

		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data:   payload,
		}
	})

	return guard(h)
}
