// Package handler contains HTTP handlers grouped by resource.
package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// Pinger is anything that can verify it's reachable. The DB pool implements it.
type Pinger interface {
	HealthCheck(ctx context.Context) error
}

// Health serves liveness and readiness probes.
type Health struct {
	logger *slog.Logger
	db     Pinger
}

// NewHealth constructs a Health handler.
func NewHealth(
	logger *slog.Logger,
	db Pinger,
) *Health {
	return &Health{logger: logger, db: db}
}

// Live reports whether the process is running.
// It does not check downstream dependencies.
func (h *Health) Live(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready reports whether the process can serve traffic.
// Returns 503 if any downstream dependency is unreachable.
func (h *Health) Ready(
	w http.ResponseWriter,
	r *http.Request,
) {
	if err := h.db.HealthCheck(r.Context()); err != nil {
		h.logger.Warn("readiness check failed", "err", err)
		writeJSON(
			w, http.StatusServiceUnavailable, map[string]string{
				"status": "not_ready",
				"reason": "database unreachable",
			},
		)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	body any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
