package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// TournamentGroupTableSvc is the subset of service.TournamentGroupTableService the handler uses.
type TournamentGroupTableSvc interface {
	ListByTournament(
		ctx context.Context,
		tournamentID uuid.UUID,
	) ([]*domain.TournamentGroupEntry, error)
}

// TournamentGroupTable is the HTTP handler for group-stage standings routes.
type TournamentGroupTable struct {
	logger *slog.Logger
	svc    TournamentGroupTableSvc
}

// NewTournamentGroupTable constructs a TournamentGroupTable handler.
func NewTournamentGroupTable(
	logger *slog.Logger,
	svc TournamentGroupTableSvc,
) *TournamentGroupTable {
	return &TournamentGroupTable{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type tournamentGroupTableEntryResponse struct {
	ID           string `json:"id"`
	TournamentID string `json:"tournament_id"`
	TeamID       string `json:"team_id"`
	TeamName     string `json:"team_name"`
	GroupLetter  string `json:"group_letter"`
	Position     int    `json:"position"`
	Points       int    `json:"points"`
	Played       int    `json:"played"`
}

// ---------- Handlers ----------

// ListGroupTable handles GET /tournaments/{tournamentId}/group-table.
func (h *TournamentGroupTable) ListGroupTable(
	w http.ResponseWriter,
	r *http.Request,
) {
	tournamentID, err := parseUUIDPathValue(r, "tournamentId")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	entries, err := h.svc.ListByTournament(r.Context(), tournamentID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]tournamentGroupTableEntryResponse, 0, len(entries))
	for _, e := range entries {
		out = append(out, tournamentGroupTableEntryResponse{
			ID:           e.ID.String(),
			TournamentID: e.TournamentID.String(),
			TeamID:       e.TeamID.String(),
			TeamName:     e.TeamName,
			GroupLetter:  e.GroupLetter,
			Position:     e.Position,
			Points:       e.Points,
			Played:       e.Played,
		})
	}
	writeJSON(w, http.StatusOK, out)
}

// compile-time interface check
var _ TournamentGroupTableSvc = (*service.TournamentGroupTableService)(nil)
