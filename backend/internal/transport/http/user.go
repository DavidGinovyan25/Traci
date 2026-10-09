package httptransport

import (
	"net/http"

	"github.com/google/uuid"

	"traci/backend/internal/gen"
)

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request, params gen.GetUsersParams) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	limit, offset := pagination(params.Limit, params.Offset)
	page, err := h.service.ListUsers(r.Context(), actor, limit, offset)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	items := make([]gen.User, len(page.Items))
	for i, user := range page.Items {
		items[i] = userResponse(user)
	}
	writeJSON(w, http.StatusOK, gen.UserPage{Items: items, Limit: page.Limit, Offset: page.Offset, Total: page.Total})
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	user, err := h.service.GetUser(r.Context(), actor, userID)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, userResponse(*user))
}

func (h *Handler) UpdateUserBlockStatus(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	body, ok := decodeBody[gen.UpdateUserRequest](w, r)
	if !ok {
		return
	}
	user, err := h.service.SetUserBlocked(r.Context(), actor, userID, body.IsBlocked)
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, userResponse(*user))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	actor := h.actor(w, r, true)
	if actor == nil {
		return
	}
	if err := h.service.DeleteUser(r.Context(), actor, userID); err != nil {
		h.fail(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
