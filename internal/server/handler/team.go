package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// TeamService is the subset of service.TeamService that the handler uses.
type TeamService interface {
	ListWithHandicaps(ctx context.Context, tournamentID uuid.UUID) ([]domain.TeamWithHandicaps, error)
}

// Team is the HTTP handler for the teams resource.
type Team struct {
	logger *slog.Logger
	svc    TeamService
}

// NewTeam constructs a Team handler.
func NewTeam(logger *slog.Logger, svc TeamService) *Team {
	return &Team{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type teamHandicapItemResponse struct {
	Category string `json:"category"`
	Points   int    `json:"points"`
}

type teamWithHandicapsResponse struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Logo        string                     `json:"logo"`
	GroupLetter *string                    `json:"group_letter,omitempty"`
	Handicaps   []teamHandicapItemResponse `json:"handicaps"`
}

func toTeamWithHandicapsResponse(t domain.TeamWithHandicaps) teamWithHandicapsResponse {
	handicaps := make([]teamHandicapItemResponse, 0, len(t.Handicaps))
	for _, h := range t.Handicaps {
		handicaps = append(handicaps, teamHandicapItemResponse{
			Category: string(h.Category),
			Points:   h.Points,
		})
	}
	return teamWithHandicapsResponse{
		ID:          t.ID.String(),
		Name:        t.Name,
		Logo:        t.Logo,
		GroupLetter: t.GroupLetter,
		Handicaps:   handicaps,
	}
}

// ---------- Handlers ----------

// List handles GET /tournaments/{tournament_id}/teams.
func (h *Team) List(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := parseUUIDPathValue(r, "tournament_id")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	teams, err := h.svc.ListWithHandicaps(r.Context(), tournamentID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]teamWithHandicapsResponse, 0, len(teams))
	for _, t := range teams {
		out = append(out, toTeamWithHandicapsResponse(t))
	}
	writeJSON(w, http.StatusOK, map[string]any{"teams": out})
}

// ensure service.TeamService satisfies our interface at compile time
var _ TeamService = (*service.TeamService)(nil)
