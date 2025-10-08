package recipehttp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/0xm0-v1/sik6/internal/http/middleware"
	"github.com/0xm0-v1/sik6/internal/http/response"
	"github.com/0xm0-v1/sik6/internal/recipe"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handlers struct {
	Collection http.Handler
	Resource   http.Handler
}

func NewHandlers(svc recipe.Service, token string) Handlers {
	if svc == nil {
		svc = recipe.NewService(nil)
	}

	h := handler{svc: svc}
	token = strings.TrimSpace(token)

	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusUnauthorized, "missing or invalid token", "recipes:auth")
	})

	collection := response.MethodGuard(http.MethodGet, http.MethodHead, http.MethodPost)(http.HandlerFunc(h.collection))
	collection = middleware.Chain(collection, middleware.RequireBearerToken(middleware.BearerConfig{
		Token:        token,
		Methods:      []string{http.MethodPost},
		Unauthorized: unauthorized,
	}))

	resource := response.MethodGuard(http.MethodGet, http.MethodHead, http.MethodPatch, http.MethodDelete)(http.HandlerFunc(h.resource))
	resource = middleware.Chain(resource, middleware.RequireBearerToken(middleware.BearerConfig{
		Token:        token,
		Methods:      []string{http.MethodPatch, http.MethodDelete},
		Unauthorized: unauthorized,
	}))

	return Handlers{
		Collection: collection,
		Resource:   resource,
	}
}

type handler struct {
	svc recipe.Service
}

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

const metaComponent = "api"

func meta(kind string) response.Meta {
	return response.NewMeta(metaComponent, kind)
}

func (h handler) collection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		response.HeadAware(h.listResponder).ServeHTTP(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		response.WriteNoBody(w, http.StatusMethodNotAllowed)
	}
}

func (h handler) resource(w http.ResponseWriter, r *http.Request) {
	id, ok := extractID(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "resource not found", "recipes")
		return
	}

	switch r.Method {
	case http.MethodGet, http.MethodHead:
		response.HeadAware(h.getResponder(id)).ServeHTTP(w, r)
	case http.MethodPatch:
		h.rename(w, r, id)
	case http.MethodDelete:
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
			Data:   meta("recipes:list"),
		}
	}

	items, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		if errors.Is(err, recipe.ErrRepositoryUnavailable) {
			log.Printf("recipes:list unavailable: %v", err)
			return http.StatusInternalServerError, response.Envelope{
				Status: "error",
				Error:  "recipe service unavailable",
				Data:   meta("recipes:list"),
			}
		}
		log.Printf("recipes:list error: %v", err)
		return http.StatusInternalServerError, response.Envelope{
			Status: "error",
			Error:  "could not list recipes",
			Data:   meta("recipes:list"),
		}
	}

	payload := struct {
		Recipes []*recipe.Recipe `json:"recipes"`
		Meta    response.Meta    `json:"meta"`
	}{
		Recipes: items,
		Meta:    meta("recipes:list"),
	}

	return http.StatusOK, response.Envelope{
		Status: "ok",
		Data:   payload,
	}
}

func (h handler) getResponder(id string) response.JSONResponder {
	return func(r *http.Request) (int, any) {
		item, err := h.svc.Get(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return http.StatusNotFound, response.Envelope{
					Status: "error",
					Error:  "recipe not found",
					Data:   meta("recipes:get"),
				}
			}
			if errors.Is(err, recipe.ErrRepositoryUnavailable) {
				log.Printf("recipes:get unavailable id=%s: %v", id, err)
				return http.StatusInternalServerError, response.Envelope{
					Status: "error",
					Error:  "recipe service unavailable",
					Data:   meta("recipes:get"),
				}
			}
			log.Printf("recipes:get error id=%s: %v", id, err)
			return http.StatusInternalServerError, response.Envelope{
				Status: "error",
				Error:  "could not fetch recipe",
				Data:   meta("recipes:get"),
			}
		}

		payload := struct {
			Recipe *recipe.Recipe `json:"recipe"`
			Meta   response.Meta  `json:"meta"`
		}{
			Recipe: item,
			Meta:   meta("recipes:get"),
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
		writeError(w, http.StatusBadRequest, err.Error(), "recipes:create")
		return
	}

	name, errMsg := normalizeRecipeName(req.Name)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg, "recipes:create")
		return
	}

	item, err := h.svc.Create(r.Context(), name)
	if h.handleServiceError(w, err, serviceErrorConfig{
		Kind:            "recipes:create",
		ConflictMessage: "recipe name already exists",
		FallbackMessage: "could not create recipe",
	}) {
		return
	}

	respondWithRecipe(w, http.StatusCreated, item, "recipes:create")
}

func (h handler) rename(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		Name string `json:"recipe_name"`
	}

	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error(), "recipes:rename")
		return
	}

	name, errMsg := normalizeRecipeName(req.Name)
	if errMsg != "" {
		writeError(w, http.StatusBadRequest, errMsg, "recipes:rename")
		return
	}

	item, err := h.svc.Rename(r.Context(), id, name)
	if h.handleServiceError(w, err, serviceErrorConfig{
		Kind:            "recipes:rename",
		NotFoundMessage: "recipe not found",
		ConflictMessage: "recipe name already exists",
		FallbackMessage: "could not rename recipe",
	}) {
		return
	}

	respondWithRecipe(w, http.StatusOK, item, "recipes:rename")
}

func (h handler) softDelete(w http.ResponseWriter, r *http.Request, id string) {
	if h.handleServiceError(w, h.svc.SoftDelete(r.Context(), id), serviceErrorConfig{
		Kind:            "recipes:delete",
		NotFoundMessage: "recipe not found",
		FallbackMessage: "could not delete recipe",
	}) {
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Envelope{
		Status: "ok",
		Data: struct {
			Meta response.Meta `json:"meta"`
		}{
			Meta: meta("recipes:delete"),
		},
	})
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

func writeError(w http.ResponseWriter, status int, message, kind string) {
	response.WriteJSON(w, status, response.Envelope{
		Status: "error",
		Error:  message,
		Data:   meta(kind),
	})
}

func respondWithRecipe(w http.ResponseWriter, status int, item *recipe.Recipe, kind string) {
	payload := struct {
		Recipe *recipe.Recipe `json:"recipe"`
		Meta   response.Meta  `json:"meta"`
	}{
		Recipe: item,
		Meta:   meta(kind),
	}

	response.WriteJSON(w, status, response.Envelope{
		Status: "ok",
		Data:   payload,
	})
}

type serviceErrorConfig struct {
	Kind            string
	NotFoundMessage string
	ConflictMessage string
	FallbackMessage string
}

func (h handler) handleServiceError(w http.ResponseWriter, err error, cfg serviceErrorConfig) bool {
	if err == nil {
		return false
	}

	log.Printf("recipes handler error kind=%s: %v", cfg.Kind, err)

	switch {
	case cfg.NotFoundMessage != "" && errors.Is(err, pgx.ErrNoRows):
		writeError(w, http.StatusNotFound, cfg.NotFoundMessage, cfg.Kind)
	case cfg.ConflictMessage != "" && isUniqueViolation(err):
		writeError(w, http.StatusConflict, cfg.ConflictMessage, cfg.Kind)
	case errors.Is(err, recipe.ErrRepositoryUnavailable):
		writeError(w, http.StatusInternalServerError, "recipe service unavailable", cfg.Kind)
	default:
		msg := cfg.FallbackMessage
		if msg == "" {
			msg = "operation failed"
		}
		writeError(w, http.StatusInternalServerError, msg, cfg.Kind)
	}
	return true
}

func normalizeRecipeName(raw string) (string, string) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", "recipe name is required"
	}
	if len(name) > 200 {
		return "", "recipe name must be 200 characters or fewer"
	}
	return name, ""
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
