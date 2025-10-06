package recipehttp

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/0xm0-v1/sik6/internal/http/response"
	"github.com/0xm0-v1/sik6/internal/recipe"
	"github.com/jackc/pgx/v5"
)

type Handlers struct {
	Collection http.Handler
	Resource   http.Handler
}

func NewHandlers(repo recipe.Repository, token string) Handlers {
	h := handler{repo: repo, token: strings.TrimSpace(token)}

	collection := response.MethodGuard(http.MethodGet, http.MethodHead, http.MethodPost)(http.HandlerFunc(h.collection))
	resource := response.MethodGuard(http.MethodGet, http.MethodHead, http.MethodPatch, http.MethodDelete)(http.HandlerFunc(h.resource))

	return Handlers{
		Collection: collection,
		Resource:   resource,
	}
}

type handler struct {
	repo  recipe.Repository
	token string
}

const (
	defaultListLimit = 20
	maxListLimit     = 100
	bearerPrefix     = "Bearer "
)

func (h handler) collection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		response.HeadAware(h.listResponder).ServeHTTP(w, r)
	case http.MethodPost:
		if !h.authorize(w, r) {
			return
		}
		h.create(w, r)
	default:
		response.WriteNoBody(w, http.StatusMethodNotAllowed)
	}
}

func (h handler) resource(w http.ResponseWriter, r *http.Request) {
	id, ok := extractID(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found")
		return
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		response.HeadAware(h.getResponder(id)).ServeHTTP(w, r)
	case http.MethodPatch:
		if !h.authorize(w, r) {
			return
		}
		h.rename(w, r, id)
	case http.MethodDelete:
		if !h.authorize(w, r) {
			return
		}
		h.softDelete(w, r, id)
	default:
		response.WriteNoBody(w, http.StatusMethodNotAllowed)
	}
}

func (h handler) listResponder(r *http.Request) (int, any) {
	limit, offset, err := parsePagination(r)
	if err != nil {
		return http.StatusBadRequest, response.Envelope{
			Status: "error",
			Error:  err.Error(),
			Data:   response.NewMeta("api", "recipes:list"),
		}
	}

	items, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		return http.StatusInternalServerError, response.Envelope{
			Status: "error",
			Error:  "could not list recipes",
			Data:   response.NewMeta("api", "recipes:list"),
		}
	}

	payload := struct {
		Recipes []*recipe.Recipe `json:"recipes"`
		Meta    response.Meta    `json:"meta"`
	}{
		Recipes: items,
		Meta:    response.NewMeta("api", "recipes:list"),
	}

	return http.StatusOK, response.Envelope{
		Status: "ok",
		Data:   payload,
	}
}

func (h handler) getResponder(id string) response.JSONResponder {
	return func(r *http.Request) (int, any) {
		item, err := h.repo.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return http.StatusNotFound, response.Envelope{
					Status: "error",
					Error:  "recipe not found",
					Data:   response.NewMeta("api", "recipes:get"),
				}
			}
			return http.StatusInternalServerError, response.Envelope{
				Status: "error",
				Error:  "could not fetch recipe",
				Data:   response.NewMeta("api", "recipes:get"),
			}
		}

		payload := struct {
			Recipe *recipe.Recipe `json:"recipe"`
			Meta   response.Meta  `json:"meta"`
		}{
			Recipe: item,
			Meta:   response.NewMeta("api", "recipes:get"),
		}

		return http.StatusOK, response.Envelope{
			Status: "ok",
			Data:   payload,
		}
	}
}

func (h handler) create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"recipe_name"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		writeError(w, http.StatusBadRequest, "recipe name is required")
		return
	}
	if len(req.Name) > 200 {
		writeError(w, http.StatusBadRequest, "recipe name must be 200 characters or fewer")
		return
	}

	id, err := h.repo.Create(r.Context(), req.Name)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create recipe")
		return
	}

	item, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch created recipe")
		return
	}

	payload := struct {
		Recipe *recipe.Recipe `json:"recipe"`
		Meta   response.Meta  `json:"meta"`
	}{
		Recipe: item,
		Meta:   response.NewMeta("api", "recipes:create"),
	}

	response.WriteJSON(w, http.StatusCreated, response.Envelope{
		Status: "ok",
		Data:   payload,
	})
}

func (h handler) rename(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name string `json:"recipe_name"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		writeError(w, http.StatusBadRequest, "recipe name is required")
		return
	}
	if len(req.Name) > 200 {
		writeError(w, http.StatusBadRequest, "recipe name must be 200 characters or fewer")
		return
	}

	if err := h.repo.Rename(r.Context(), id, req.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not rename recipe")
		return
	}

	item, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not fetch recipe")
		return
	}

	payload := struct {
		Recipe *recipe.Recipe `json:"recipe"`
		Meta   response.Meta  `json:"meta"`
	}{
		Recipe: item,
		Meta:   response.NewMeta("api", "recipes:rename"),
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Status: "ok",
		Data:   payload,
	})
}

func (h handler) softDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.repo.SoftDelete(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeError(w, http.StatusNotFound, "recipe not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "could not delete recipe")
		return
	}

	response.WriteNoBody(w, http.StatusNoContent)
}

func (h handler) authorize(w http.ResponseWriter, r *http.Request) bool {
	if h.token == "" {
		return true
	}

	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(header, bearerPrefix) {
		writeError(w, http.StatusUnauthorized, "missing or invalid token")
		return false
	}

	candidate := strings.TrimSpace(header[len(bearerPrefix):])
	if subtle.ConstantTimeCompare([]byte(candidate), []byte(h.token)) != 1 {
		writeError(w, http.StatusUnauthorized, "missing or invalid token")
		return false
	}

	return true
}

func extractID(path string) (string, bool) {
	if !strings.HasPrefix(path, "/recipes/") {
		return "", false
	}

	remainder := strings.TrimPrefix(path, "/recipes/")
	if remainder == "" || strings.Contains(remainder, "/") {
		return "", false
	}

	return remainder, true
}

func parsePagination(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	limit := defaultListLimit
	if raw := strings.TrimSpace(query.Get("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			return 0, 0, errors.New("limit must be a positive integer")
		}
		if n > maxListLimit {
			n = maxListLimit
		}
		limit = n
	}

	offset := 0
	if raw := strings.TrimSpace(query.Get("offset")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			return 0, 0, errors.New("offset must be a non-negative integer")
		}
		offset = n
	}

	return limit, offset, nil
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return errors.New("invalid JSON payload")
	}
	return nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	response.WriteJSON(w, status, response.Envelope{
		Status: "error",
		Error:  message,
		Data:   response.NewMeta("api", "recipes"),
	})
}
