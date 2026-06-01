package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// OutcomesSvc is the subset of service.OutcomesService the handler uses.
type OutcomesSvc interface {
	ListByTournament(ctx context.Context, tournamentID uuid.UUID) (*domain.TournamentOutcomes, error)
}

// Outcomes is the HTTP handler for tournament-outcomes routes.
type Outcomes struct {
	logger *slog.Logger
	svc    OutcomesSvc
}

// NewOutcomes constructs an Outcomes handler.
func NewOutcomes(logger *slog.Logger, svc OutcomesSvc) *Outcomes {
	return &Outcomes{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type playerOutcomeResponse struct {
	ID         string    `json:"id"`
	Category   string    `json:"category"`
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	TeamID     string    `json:"team_id"`
	TeamName   string    `json:"team_name"`
	RecordedAt time.Time `json:"recorded_at"`
}

type teamOutcomeResponse struct {
	ID         string    `json:"id"`
	Category   string    `json:"category"`
	TeamID     string    `json:"team_id"`
	TeamName   string    `json:"team_name"`
	RecordedAt time.Time `json:"recorded_at"`
}

type tournamentOutcomesResponse struct {
	PlayerOutcomes []playerOutcomeResponse `json:"player_outcomes"`
	TeamOutcomes   []teamOutcomeResponse   `json:"team_outcomes"`
}

// ---------- Handlers ----------

// ListOutcomes handles GET /tournaments/{tournamentId}/outcomes.
func (h *Outcomes) ListOutcomes(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := parseUUIDPathValue(r, "tournamentId")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	outcomes, err := h.svc.ListByTournament(r.Context(), tournamentID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	players := make([]playerOutcomeResponse, 0, len(outcomes.PlayerOutcomes))
	for _, p := range outcomes.PlayerOutcomes {
		players = append(players, playerOutcomeResponse{
			ID:         p.ID.String(),
			Category:   string(p.Category),
			PlayerID:   p.PlayerID.String(),
			PlayerName: p.PlayerName,
			TeamID:     p.TeamID.String(),
			TeamName:   p.TeamName,
			RecordedAt: p.RecordedAt,
		})
	}

	teams := make([]teamOutcomeResponse, 0, len(outcomes.TeamOutcomes))
	for _, t := range outcomes.TeamOutcomes {
		teams = append(teams, teamOutcomeResponse{
			ID:         t.ID.String(),
			Category:   string(t.Category),
			TeamID:     t.TeamID.String(),
			TeamName:   t.TeamName,
			RecordedAt: t.RecordedAt,
		})
	}

	writeJSON(w, http.StatusOK, tournamentOutcomesResponse{
		PlayerOutcomes: players,
		TeamOutcomes:   teams,
	})
}

// compile-time interface check
var _ OutcomesSvc = (*service.OutcomesService)(nil)
