package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// ---------- service interfaces ----------

// ScorePredictionSvc is the service contract for score predictions.
type ScorePredictionSvc interface {
	UpsertScore(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
}

// FixtureSvc is the service contract for fixture listings.
type FixtureSvc interface {
	ListForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error)
	ListForLeague(ctx context.Context, leagueID, userID uuid.UUID) ([]*domain.LeagueFixtureView, error)
	ListForLeaguePaged(ctx context.Context, leagueID, userID uuid.UUID, n, skip int) ([]*domain.LeagueFixtureView, error)
}

// ---------- DTOs ----------

type scorePredictionResponse struct {
	ID        string `json:"id"`
	FixtureID string `json:"fixture_id"`
	GoalsHome int    `json:"goals_home"`
	GoalsAway int    `json:"goals_away"`
	Points    *int   `json:"points"`
}

type fixtureResponse struct {
	ID               string    `json:"id"`
	ExternalID       int64     `json:"external_id"`
	TournamentID     string    `json:"tournament_id"`
	HomeTeamName     string    `json:"home_team_name"`
	AwayTeamName     string    `json:"away_team_name"`
	Round            string    `json:"round"`
	KickoffAt        time.Time `json:"kickoff_at"`
	Status           string    `json:"status"`
	PredictionLocked bool      `json:"prediction_locked"`
	GoalsHome        *int      `json:"goals_home"`
	GoalsAway        *int      `json:"goals_away"`
}

type userFixtureViewResponse struct {
	fixtureResponse
	Prediction *scorePredictionResponse `json:"prediction"`
}

type leagueMemberScorePrediction struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	GoalsHome   *int   `json:"goals_home"`
	GoalsAway   *int   `json:"goals_away"`
	Points      *int   `json:"points"`
}

type leagueFixtureViewResponse struct {
	fixtureResponse
	Predictions []leagueMemberScorePrediction `json:"predictions"`
}

// ---------- ScorePrediction handler ----------

// ScorePrediction handles score prediction routes.
type ScorePrediction struct {
	logger *slog.Logger
	svc    ScorePredictionSvc
}

// NewScorePrediction constructs a ScorePrediction handler.
func NewScorePrediction(logger *slog.Logger, svc ScorePredictionSvc) *ScorePrediction {
	return &ScorePrediction{logger: logger, svc: svc}
}

// UpsertScore handles PUT /predictions/{fixtureId}.
func (h *ScorePrediction) UpsertScore(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	fixtureID, err := parseUUIDPathValue(r, "fixtureId")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	var req struct {
		GoalsHome int `json:"goals_home"`
		GoalsAway int `json:"goals_away"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid request body: %w", domain.ErrInvalid))
		return
	}

	pred, err := h.svc.UpsertScore(r.Context(), domain.UpsertScorePredictionInput{
		UserID:    caller.ID,
		FixtureID: fixtureID,
		GoalsHome: req.GoalsHome,
		GoalsAway: req.GoalsAway,
	})
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toScorePredictionResponse(pred))
}

// ---------- Fixture handler ----------

// Fixture handles fixture listing routes.
type Fixture struct {
	logger *slog.Logger
	svc    FixtureSvc
}

// NewFixture constructs a Fixture handler.
func NewFixture(logger *slog.Logger, svc FixtureSvc) *Fixture {
	return &Fixture{logger: logger, svc: svc}
}

// ListForUser handles GET /tournaments/{tournamentId}/fixtures.
func (h *Fixture) ListForUser(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	tournamentID, err := parseUUIDPathValue(r, "tournamentId")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	views, err := h.svc.ListForUser(r.Context(), tournamentID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]userFixtureViewResponse, 0, len(views))
	for _, v := range views {
		out = append(out, toUserFixtureViewResponse(v))
	}
	writeJSON(w, http.StatusOK, out)
}

// ListForLeague handles GET /leagues/{leagueId}/predictions.
func (h *Fixture) ListForLeague(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	leagueID, err := parseUUIDPathValue(r, "leagueId")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	nStr := r.URL.Query().Get("n")
	skipStr := r.URL.Query().Get("skip")

	// skip without n is invalid
	if skipStr != "" && nStr == "" {
		writeError(w, r, h.logger, fmt.Errorf("n is required when skip is provided: %w", domain.ErrInvalid))
		return
	}

	var views []*domain.LeagueFixtureView
	if nStr != "" {
		n, parseErr := strconv.Atoi(nStr)
		if parseErr != nil || n <= 0 {
			writeError(w, r, h.logger, fmt.Errorf("n must be a positive integer: %w", domain.ErrInvalid))
			return
		}
		skip := 0
		if skipStr != "" {
			skip, parseErr = strconv.Atoi(skipStr)
			if parseErr != nil || skip < 0 {
				writeError(w, r, h.logger, fmt.Errorf("skip must be a non-negative integer: %w", domain.ErrInvalid))
				return
			}
		}
		views, err = h.svc.ListForLeaguePaged(r.Context(), leagueID, caller.ID, n, skip)
	} else {
		views, err = h.svc.ListForLeague(r.Context(), leagueID, caller.ID)
	}
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]leagueFixtureViewResponse, 0, len(views))
	for _, v := range views {
		out = append(out, toLeagueFixtureViewResponse(v))
	}
	writeJSON(w, http.StatusOK, out)
}

// ---------- compile-time interface checks ----------

var _ ScorePredictionSvc = (*service.PredictionService)(nil)
var _ FixtureSvc = (*service.FixtureService)(nil)

// ---------- response converters ----------

func toScorePredictionResponse(p *domain.ScorePrediction) scorePredictionResponse {
	return scorePredictionResponse{
		ID:        p.ID.String(),
		FixtureID: p.FixtureID.String(),
		GoalsHome: p.GoalsHome,
		GoalsAway: p.GoalsAway,
		Points:    p.Points,
	}
}

func toFixtureResponse(f domain.Fixture) fixtureResponse {
	return fixtureResponse{
		ID:               f.ID.String(),
		ExternalID:       f.ExternalID,
		TournamentID:     f.TournamentID.String(),
		HomeTeamName:     f.HomeTeamName,
		AwayTeamName:     f.AwayTeamName,
		Round:            f.Round,
		KickoffAt:        f.KickoffAt,
		Status:           string(f.Status),
		PredictionLocked: f.PredictionLocked,
		GoalsHome:        f.GoalsHome,
		GoalsAway:        f.GoalsAway,
	}
}

func toUserFixtureViewResponse(v *domain.UserFixtureView) userFixtureViewResponse {
	r := userFixtureViewResponse{fixtureResponse: toFixtureResponse(v.Fixture)}
	if v.Prediction != nil {
		p := toScorePredictionResponse(v.Prediction)
		r.Prediction = &p
	}
	return r
}

func toLeagueFixtureViewResponse(v *domain.LeagueFixtureView) leagueFixtureViewResponse {
	preds := make([]leagueMemberScorePrediction, 0, len(v.Predictions))
	for _, p := range v.Predictions {
		preds = append(preds, leagueMemberScorePrediction{
			UserID:      p.UserID.String(),
			DisplayName: p.DisplayName,
			GoalsHome:   p.GoalsHome,
			GoalsAway:   p.GoalsAway,
			Points:      p.Points,
		})
	}
	return leagueFixtureViewResponse{
		fixtureResponse: toFixtureResponse(v.Fixture),
		Predictions:     preds,
	}
}
