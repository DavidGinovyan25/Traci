package httptransport

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"traci/backend/internal/application"
	"traci/backend/internal/domain"
	"traci/backend/internal/gen"
)

type Handler struct {
	service       *application.Service
	logger        *slog.Logger
	tokenLifetime time.Duration
}

var _ gen.ServerInterface = (*Handler)(nil)

func NewHandler(service *application.Service, logger *slog.Logger, tokenLifetime time.Duration) *Handler {
	return &Handler{service: service, logger: logger, tokenLifetime: tokenLifetime}
}

func (h *Handler) actor(w http.ResponseWriter, r *http.Request, admin bool) *domain.User {
	fields := strings.Fields(r.Header.Get("Authorization"))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		h.fail(w, r, application.ErrUnauthorized)
		return nil
	}
	actor, err := h.service.Authenticate(r.Context(), fields[1])
	if err != nil {
		h.fail(w, r, err)
		return nil
	}
	if admin && actor.Role != domain.RoleAdmin {
		h.fail(w, r, application.ErrForbidden)
		return nil
	}
	return actor
}

func (h *Handler) fail(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, application.ErrUnauthorized):
		w.Header().Set("WWW-Authenticate", "Bearer")
		writeError(w, r, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid credentials or access token")
	case errors.Is(err, application.ErrBlocked):
		writeError(w, r, http.StatusForbidden, "USER_BLOCKED", "User is blocked")
	case errors.Is(err, application.ErrForbidden):
		writeError(w, r, http.StatusForbidden, "FORBIDDEN", "Access denied")
	case errors.Is(err, application.ErrSelfModification):
		writeError(w, r, http.StatusConflict, "SELF_MODIFICATION", err.Error())
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, r, http.StatusNotFound, "NOT_FOUND", "Resource not found")
	case errors.Is(err, domain.ErrAlreadyExists):
		writeError(w, r, http.StatusConflict, "ALREADY_EXISTS", "Resource already exists")
	default:
		h.logger.Error("request failed", "method", r.Method, "path", r.URL.Path, "error", err)
		writeError(w, r, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}

func decodeBody[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var body T
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid JSON body")
		return body, false
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Expected a single JSON object")
		return body, false
	}
	return body, true
}

func pagination(limit, offset *int) (int, int) {
	l, o := 20, 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}
	return l, o
}
