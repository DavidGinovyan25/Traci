package httptransport

import (
	"net/http"
	"net/mail"
	"strings"

	"traci/backend/internal/application"
	"traci/backend/internal/gen"
)

func validEmail(value string) bool {
	parsed, err := mail.ParseAddress(value)
	return err == nil && parsed.Address == value
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, ok := decodeBody[gen.LoginRequest](w, r)
	if !ok {
		return
	}
	if !validEmail(string(body.Email)) {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid email")
		return
	}
	token, err := h.service.Login(r.Context(), application.Login{Email: string(body.Email), Password: body.Password})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, gen.LoginResponse{AccessToken: token, TokenType: "Bearer", ExpiresIn: int(h.tokenLifetime.Seconds())})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	body, ok := decodeBody[gen.RegisterRequest](w, r)
	if !ok {
		return
	}
	if !validEmail(string(body.Email)) || strings.TrimSpace(body.Username) == "" || strings.TrimSpace(body.FirstName) == "" || strings.TrimSpace(body.SecondName) == "" {
		writeError(w, r, http.StatusBadRequest, "BAD_REQUEST", "Invalid registration fields")
		return
	}
	user, err := h.service.Register(r.Context(), application.Register{
		Username: body.Username, FirstName: body.FirstName, SecondName: body.SecondName,
		Email: string(body.Email), Password: body.Password,
	})
	if err != nil {
		h.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, userResponse(*user))
}
