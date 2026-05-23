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

// TeamHandicapService is the subset of service.TeamHandicapService that the handler uses.
type TeamHandicapService interface {
	Get(ctx context.Context, teamID uuid.UUID, category domain.TeamHandicapCategory) (*domain.TeamHandicap, error)
}

// TeamHandicap is the HTTP handler for the team-handicaps resource.
type TeamHandicap struct {
	logger *slog.Logger
	svc    TeamHandicapService
}

// NewTeamHandicap constructs a TeamHandicap handler.
func NewTeamHandicap(logger *slog.Logger, svc TeamHandicapService) *TeamHandicap {
	return &TeamHandicap{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type teamHandicapResponse struct {
	Points int `json:"points"`
}

// ---------- Handlers ----------

// Get handles GET /team-handicaps/{team_id}/{category}.
func (h *TeamHandicap) Get(w http.ResponseWriter, r *http.Request) {
	teamID, err := parseUUIDPathValue(r, "team_id")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	catStr := r.PathValue("category")
	cat, err := domain.ParseTeamHandicapCategory(catStr)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid category: %w", domain.ErrInvalid))
		return
	}

	handicap, err := h.svc.Get(r.Context(), teamID, cat)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	writeJSON(w, http.StatusOK, teamHandicapResponse{Points: handicap.Points})
}

// ensure service.TeamHandicapService satisfies our interface at compile time
var _ TeamHandicapService = (*service.TeamHandicapService)(nil)
