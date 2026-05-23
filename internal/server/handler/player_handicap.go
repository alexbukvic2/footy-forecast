package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// PlayerHandicapService is the subset of service.PlayerHandicapService that the handler uses.
type PlayerHandicapService interface {
	Get(ctx context.Context, playerID uuid.UUID, category domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error)
}

// PlayerHandicap is the HTTP handler for the player-handicaps resource.
type PlayerHandicap struct {
	logger *slog.Logger
	svc    PlayerHandicapService
}

// NewPlayerHandicap constructs a PlayerHandicap handler.
func NewPlayerHandicap(logger *slog.Logger, svc PlayerHandicapService) *PlayerHandicap {
	return &PlayerHandicap{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type playerHandicapResponse struct {
	Points int `json:"points"`
}

// ---------- Handlers ----------

// Get handles GET /player-handicaps/{player_id}/{category}.
func (h *PlayerHandicap) Get(w http.ResponseWriter, r *http.Request) {
	playerID, err := parseUUIDPathValue(r, "player_id")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	catStr := r.PathValue("category")
	cat, err := domain.ParsePlayerHandicapCategory(catStr)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid category: %w", domain.ErrInvalid))
		return
	}

	handicap, err := h.svc.Get(r.Context(), playerID, cat)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	writeJSON(w, http.StatusOK, playerHandicapResponse{Points: handicap.Points})
}

// ensure service.PlayerHandicapService satisfies our interface at compile time
var _ PlayerHandicapService = (*service.PlayerHandicapService)(nil)
