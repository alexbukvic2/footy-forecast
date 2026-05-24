package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// ErrorResponse is the JSON shape returned to clients for any error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// writeError translates a domain error into an HTTP response.
//
// The mapping:
//   - domain.ErrInvalid   → 400 Bad Request   (message exposed)
//   - domain.ErrNotFound  → 404 Not Found     (generic message)
//   - domain.ErrConflict  → 409 Conflict      (message exposed)
//   - anything else       → 500 Internal      (generic message, full error logged)
//
// We expose validation messages to clients but never internal-error details —
// those go in logs, attached to request_id.
func writeError(
	w http.ResponseWriter,
	r *http.Request,
	logger *slog.Logger,
	err error,
) {
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
	case errors.Is(err, domain.ErrInvalid):
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
	case errors.Is(err, domain.ErrConflict):
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrForbidden):
		writeJSON(w, http.StatusForbidden, ErrorResponse{Error: "forbidden"})
	default:
		logger.ErrorContext(r.Context(), "internal error", "err", err)
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
	}
}
