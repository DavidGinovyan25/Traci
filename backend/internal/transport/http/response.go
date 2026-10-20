package httptransport

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"traci/backend/internal/gen"
)

func writeError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	code, message string,
	logger *slog.Logger,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(gen.ErrorResponse{
		Timestamp: time.Now().UTC(),
		Status:    status,
		Error:     http.StatusText(status),
		Code:      code,
		Message:   message,
		Path:      r.URL.Path,
	}); err != nil {
		logger.ErrorContext(
			r.Context(),
			"JSON response encoding failed",
			"error", err,
		)
	}

}

func writeJSON(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	body any,
	logger *slog.Logger,
) {
	data, err := json.Marshal(body)
	if err != nil {
		logger.ErrorContext(
			r.Context(),
			"JSON serialization failed",
			"error", err,
		)
		writeError(w, r, http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"Internal server error", logger,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(append(data, '\n')); err != nil {
		logger.WarnContext(
			r.Context(),
			"HTTP response write failed",
			"error", err,
		)
	}
}
