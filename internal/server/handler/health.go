// Package handler contains HTTP handlers grouped by resource.
package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// Health serves liveness and readiness probes.
type Health struct {
	logger *slog.Logger
}

// NewHealth constructs a Health handler.
func NewHealth(logger *slog.Logger) *Health {
	return &Health{logger: logger}
}

// Live is for "is the process running" — used by load balancers / systemd.
func (h *Health) Live(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready is for "is the process able to serve traffic" — checks downstream deps later (DB, etc).
func (h *Health) Ready(
	w http.ResponseWriter,
	_ *http.Request,
) {
	// TODO: ping DB, Redis, etc. once wired up.
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
