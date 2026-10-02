package httptransport

import (
	"encoding/json"
	"net/http"
	"time"

	"traci/backend/internal/gen"
)

func writeError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(gen.ErrorResponse{
		Timestamp: time.Now().UTC(), Status: status, Error: http.StatusText(status),
		Code: code, Message: message, Path: r.URL.Path,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	data, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(data, '\n'))
}
