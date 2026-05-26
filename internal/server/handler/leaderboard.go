package handler

import (
	"context"
	"fmt"
	"net/http"

	"log/slog"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
)

// LeaderboardService is the subset of service.LeaderboardService the handler uses.
type LeaderboardService interface {
	GetLeagueLeaderboard(
		ctx context.Context,
		leagueID, requesterID uuid.UUID,
	) ([]*domain.LeaderboardEntry, error)
	GetTournamentLeaderboard(
		ctx context.Context,
		tournamentID uuid.UUID,
	) ([]*domain.LeaderboardEntry, error)
}

// Leaderboard is the HTTP handler for leaderboard routes.
type Leaderboard struct {
	logger *slog.Logger
	svc    LeaderboardService
}

// NewLeaderboard constructs a Leaderboard handler.
func NewLeaderboard(
	logger *slog.Logger,
	svc LeaderboardService,
) *Leaderboard {
	return &Leaderboard{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type leaderboardEntryResponse struct {
	Position     int    `json:"position"`
	UserID       string `json:"user_id"`
	DisplayName  string `json:"display_name"`
	ScorePoints  int    `json:"score_points"`
	PlayerPoints int    `json:"player_points"`
	TeamPoints   int    `json:"team_points"`
	TotalPoints  int    `json:"total_points"`
}

func toLeaderboardEntryResponse(e *domain.LeaderboardEntry) leaderboardEntryResponse {
	return leaderboardEntryResponse{
		Position:     e.Position,
		UserID:       e.UserID.String(),
		DisplayName:  e.DisplayName,
		ScorePoints:  e.ScorePoints,
		PlayerPoints: e.PlayerPoints,
		TeamPoints:   e.TeamPoints,
		TotalPoints:  e.TotalPoints,
	}
}

// ---------- Handlers ----------

// GetForLeague handles GET /leagues/{id}/leaderboard.
func (h *Leaderboard) GetForLeague(
	w http.ResponseWriter,
	r *http.Request,
) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	leagueID, err := parseUUIDPathValue(r, "id")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	entries, err := h.svc.GetLeagueLeaderboard(r.Context(), leagueID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]leaderboardEntryResponse, 0, len(entries))
	for _, e := range entries {
		out = append(out, toLeaderboardEntryResponse(e))
	}
	writeJSON(w, http.StatusOK, map[string]any{"leaderboard": out})
}

// GetForTournament handles GET /tournaments/{id}/leaderboard.
func (h *Leaderboard) GetForTournament(
	w http.ResponseWriter,
	r *http.Request,
) {
	tournamentID, err := parseUUIDPathValue(r, "id")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	entries, err := h.svc.GetTournamentLeaderboard(r.Context(), tournamentID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]leaderboardEntryResponse, 0, len(entries))
	for _, e := range entries {
		out = append(out, toLeaderboardEntryResponse(e))
	}
	writeJSON(w, http.StatusOK, map[string]any{"leaderboard": out})
}
