package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// PredictionService is the subset of service.PredictionService that the handler uses.
type PredictionService interface {
	UpsertScore(
		ctx context.Context,
		in domain.UpsertScorePredictionInput,
	) (*domain.ScorePrediction, error)
}

// Prediction is the HTTP handler for prediction endpoints.
type Prediction struct {
	logger *slog.Logger
	svc    PredictionService
}

// NewPrediction constructs a Prediction handler.
func NewPrediction(logger *slog.Logger, svc PredictionService) *Prediction {
	return &Prediction{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type scorePredictionResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	FixtureID string    `json:"fixture_id"`
	GoalsHome int       `json:"goals_home"`
	GoalsAway int       `json:"goals_away"`
	Points    *int      `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toScorePredictionResponse(p domain.ScorePrediction) scorePredictionResponse {
	return scorePredictionResponse{
		ID:        p.ID.String(),
		UserID:    p.UserID.String(),
		FixtureID: p.FixtureID.String(),
		GoalsHome: p.GoalsHome,
		GoalsAway: p.GoalsAway,
		Points:    p.Points,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ---------- Handlers ----------

// UpsertScore handles PUT /predictions/{fixtureId}.
func (h *Prediction) UpsertScore(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context"))
		return
	}

	idStr := r.PathValue("fixtureId")
	fixtureID, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid fixtureId %q: %w", idStr, domain.ErrInvalid))
		return
	}

	var req struct {
		GoalsHome *int `json:"goals_home"`
		GoalsAway *int `json:"goals_away"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}
	if req.GoalsHome == nil {
		writeError(w, r, h.logger, fmt.Errorf("goals_home is required: %w", domain.ErrInvalid))
		return
	}
	if req.GoalsAway == nil {
		writeError(w, r, h.logger, fmt.Errorf("goals_away is required: %w", domain.ErrInvalid))
		return
	}

	pred, err := h.svc.UpsertScore(r.Context(), domain.UpsertScorePredictionInput{
		UserID:    caller.ID,
		FixtureID: fixtureID,
		GoalsHome: *req.GoalsHome,
		GoalsAway: *req.GoalsAway,
	})
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toScorePredictionResponse(*pred))
}

// compile-time check
var _ PredictionService = (*service.PredictionService)(nil)
