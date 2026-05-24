package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// FixtureService is the subset of service.FixtureService that the handler uses.
type FixtureService interface {
	ListForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.UserFixtureView, error)
	ListForLeague(
		ctx context.Context,
		leagueID, userID uuid.UUID,
	) ([]*domain.LeagueFixtureView, error)
}

// Fixture is the HTTP handler for fixture endpoints.
type Fixture struct {
	logger *slog.Logger
	svc    FixtureService
}

// NewFixture constructs a Fixture handler.
func NewFixture(logger *slog.Logger, svc FixtureService) *Fixture {
	return &Fixture{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type fixtureResponse struct {
	ID           string    `json:"id"`
	ExternalID   int64     `json:"external_id"`
	TournamentID string    `json:"tournament_id"`
	HomeTeamID   string    `json:"home_team_id"`
	AwayTeamID   string    `json:"away_team_id"`
	KickoffAt    time.Time `json:"kickoff_at"`
	Status       string    `json:"status"`
	GoalsHome    *int      `json:"goals_home"`
	GoalsAway    *int      `json:"goals_away"`
}

func toFixtureResponse(f domain.Fixture) fixtureResponse {
	return fixtureResponse{
		ID:           f.ID.String(),
		ExternalID:   f.ExternalID,
		TournamentID: f.TournamentID.String(),
		HomeTeamID:   f.HomeTeamID.String(),
		AwayTeamID:   f.AwayTeamID.String(),
		KickoffAt:    f.KickoffAt,
		Status:       string(f.Status),
		GoalsHome:    f.GoalsHome,
		GoalsAway:    f.GoalsAway,
	}
}

type memberPredictionResponse struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	GoalsHome   *int   `json:"goals_home"`
	GoalsAway   *int   `json:"goals_away"`
	Points      *int   `json:"points"`
}

type fixtureWithPrediction struct {
	fixtureResponse
	Prediction *scorePredictionResponse `json:"prediction"`
}

type fixtureWithMemberPredictions struct {
	fixtureResponse
	Predictions []memberPredictionResponse `json:"predictions"`
}

// ---------- Handlers ----------

// ListForUser handles GET /tournaments/{tournamentId}/fixtures.
func (h *Fixture) ListForUser(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context"))
		return
	}

	idStr := r.PathValue("tournamentId")
	tournamentID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid tournamentId %q: %w", idStr, domain.ErrInvalid))
		return
	}

	views, err := h.svc.ListForUser(r.Context(), tournamentID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]fixtureWithPrediction, 0, len(views))
	for _, v := range views {
		item := fixtureWithPrediction{fixtureResponse: toFixtureResponse(v.Fixture)}
		if v.Prediction != nil {
			resp := toScorePredictionResponse(*v.Prediction)
			item.Prediction = &resp
		}
		out = append(out, item)
	}
	writeJSON(w, http.StatusOK, map[string]any{"fixtures": out})
}

// ListForLeague handles GET /leagues/{leagueId}/predictions.
func (h *Fixture) ListForLeague(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context"))
		return
	}

	idStr := r.PathValue("leagueId")
	leagueID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid leagueId %q: %w", idStr, domain.ErrInvalid))
		return
	}

	views, err := h.svc.ListForLeague(r.Context(), leagueID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]fixtureWithMemberPredictions, 0, len(views))
	for _, v := range views {
		preds := make([]memberPredictionResponse, 0, len(v.Predictions))
		for _, p := range v.Predictions {
			preds = append(preds, memberPredictionResponse{
				UserID:      p.UserID.String(),
				DisplayName: p.DisplayName,
				GoalsHome:   p.GoalsHome,
				GoalsAway:   p.GoalsAway,
				Points:      p.Points,
			})
		}
		out = append(out, fixtureWithMemberPredictions{
			fixtureResponse: toFixtureResponse(v.Fixture),
			Predictions:     preds,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"fixtures": out})
}

// compile-time check
var _ FixtureService = (*service.FixtureService)(nil)
