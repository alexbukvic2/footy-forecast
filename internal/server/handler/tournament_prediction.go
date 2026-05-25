package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// TournamentPredictionSvc is the subset of service.TournamentPredictionService the handler uses.
type TournamentPredictionSvc interface {
	UpsertPlayerPrediction(
		ctx context.Context,
		in domain.UpsertPlayerPredictionInput,
	) (*domain.PlayerPrediction, error)
	UpsertTeamPrediction(
		ctx context.Context,
		in domain.UpsertTeamPredictionInput,
	) (*domain.TeamPrediction, error)
	ListPlayerPredictionsForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.PlayerPredictionView, error)
	ListTeamPredictionsForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.TeamPredictionView, error)
	ListLeaguePlayerPredictions(
		ctx context.Context,
		leagueID, requestingUserID uuid.UUID,
	) ([]*domain.LeaguePlayerCategoryView, error)
	ListLeagueTeamPredictions(
		ctx context.Context,
		leagueID, requestingUserID uuid.UUID,
	) ([]*domain.LeagueTeamCategoryView, error)
}

// TournamentPrediction is the HTTP handler for tournament-level prediction routes.
type TournamentPrediction struct {
	logger *slog.Logger
	svc    TournamentPredictionSvc
}

// NewTournamentPrediction constructs a TournamentPrediction handler.
func NewTournamentPrediction(
	logger *slog.Logger,
	svc TournamentPredictionSvc,
) *TournamentPrediction {
	return &TournamentPrediction{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type playerPredictionResponse struct {
	ID           string `json:"id"`
	TournamentID string `json:"tournament_id"`
	Category     string `json:"category"`
	PlayerID     string `json:"player_id"`
	PlayerName   string `json:"player_name"`
	Points       *int   `json:"points"`
}

type teamPredictionResponse struct {
	ID           string `json:"id"`
	TournamentID string `json:"tournament_id"`
	Category     string `json:"category"`
	TeamID       string `json:"team_id"`
	TeamName     string `json:"team_name"`
	Points       *int   `json:"points"`
}

type playerPredictionViewResponse struct {
	Category   string  `json:"category"`
	PlayerID   *string `json:"player_id"`
	PlayerName *string `json:"player_name"`
	Points     *int    `json:"points"`
}

type teamPredictionViewResponse struct {
	Category string  `json:"category"`
	TeamID   *string `json:"team_id"`
	TeamName *string `json:"team_name"`
	Points   *int    `json:"points"`
}

type leagueMemberPlayerPickResponse struct {
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	PlayerID    *string `json:"player_id"`
	PlayerName  *string `json:"player_name"`
	Points      *int    `json:"points"`
}

type leagueMemberTeamPickResponse struct {
	UserID      string  `json:"user_id"`
	DisplayName string  `json:"display_name"`
	TeamID      *string `json:"team_id"`
	TeamName    *string `json:"team_name"`
	Points      *int    `json:"points"`
}

type leaguePlayerCategoryViewResponse struct {
	Category    string                           `json:"category"`
	Predictions []leagueMemberPlayerPickResponse `json:"predictions"`
}

type leagueTeamCategoryViewResponse struct {
	Category    string                         `json:"category"`
	Predictions []leagueMemberTeamPickResponse `json:"predictions"`
}

func toPlayerPredictionResponse(p *domain.PlayerPrediction) playerPredictionResponse {
	return playerPredictionResponse{
		ID:           p.ID.String(),
		TournamentID: p.TournamentID.String(),
		Category:     string(p.Category),
		PlayerID:     p.Pick.String(),
		PlayerName:   p.PickName,
		Points:       p.Points,
	}
}

func toTeamPredictionResponse(p *domain.TeamPrediction) teamPredictionResponse {
	return teamPredictionResponse{
		ID:           p.ID.String(),
		TournamentID: p.TournamentID.String(),
		Category:     string(p.Category),
		TeamID:       p.Pick.String(),
		TeamName:     p.PickName,
		Points:       p.Points,
	}
}

func toPlayerPredictionViewResponse(v *domain.PlayerPredictionView) playerPredictionViewResponse {
	resp := playerPredictionViewResponse{
		Category: string(v.Category),
	}
	if v.Prediction != nil {
		id := v.Prediction.Pick.String()
		resp.PlayerID = &id
		resp.PlayerName = &v.Prediction.PickName
		resp.Points = v.Prediction.Points
	}
	return resp
}

func toTeamPredictionViewResponse(v *domain.TeamPredictionView) teamPredictionViewResponse {
	resp := teamPredictionViewResponse{
		Category: string(v.Category),
	}
	if v.Prediction != nil {
		id := v.Prediction.Pick.String()
		resp.TeamID = &id
		resp.TeamName = &v.Prediction.PickName
		resp.Points = v.Prediction.Points
	}
	return resp
}

func toLeagueMemberPlayerPickResponse(p domain.LeagueMemberPlayerPick) leagueMemberPlayerPickResponse {
	resp := leagueMemberPlayerPickResponse{
		UserID:      p.UserID.String(),
		DisplayName: p.DisplayName,
		Points:      p.Points,
	}
	if p.PlayerID != nil {
		id := p.PlayerID.String()
		resp.PlayerID = &id
	}
	if p.PlayerName != nil {
		name := *p.PlayerName
		resp.PlayerName = &name
	}
	return resp
}

func toLeagueMemberTeamPickResponse(p domain.LeagueMemberTeamPick) leagueMemberTeamPickResponse {
	resp := leagueMemberTeamPickResponse{
		UserID:      p.UserID.String(),
		DisplayName: p.DisplayName,
		Points:      p.Points,
	}
	if p.TeamID != nil {
		id := p.TeamID.String()
		resp.TeamID = &id
	}
	if p.TeamName != nil {
		name := *p.TeamName
		resp.TeamName = &name
	}
	return resp
}

// ---------- Handlers ----------

// UpsertPlayerPredictions handles PUT /tournaments/{tournamentId}/predictions/players/{category}.
func (h *TournamentPrediction) UpsertPlayerPredictions(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	category, err := domain.ParsePlayerHandicapCategory(r.PathValue("category"))
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid category %q: %w", r.PathValue("category"), domain.ErrInvalid))
		return
	}

	var req struct {
		PlayerID string `json:"player_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}
	playerID, err := uuid.Parse(req.PlayerID)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid player_id: %w", domain.ErrInvalid))
		return
	}

	pred, err := h.svc.UpsertPlayerPrediction(
		r.Context(), domain.UpsertPlayerPredictionInput{
			UserID:       caller.ID,
			TournamentID: tournamentID,
			Category:     category,
			Pick:         playerID,
		},
	)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toPlayerPredictionResponse(pred))
}

// UpsertTeamPredictions handles PUT /tournaments/{tournamentId}/predictions/teams/{category}.
func (h *TournamentPrediction) UpsertTeamPredictions(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	category, err := domain.ParseTeamHandicapCategory(r.PathValue("category"))
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid category %q: %w", r.PathValue("category"), domain.ErrInvalid))
		return
	}

	var req struct {
		TeamID string `json:"team_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}
	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid team_id: %w", domain.ErrInvalid))
		return
	}

	pred, err := h.svc.UpsertTeamPrediction(
		r.Context(), domain.UpsertTeamPredictionInput{
			UserID:       caller.ID,
			TournamentID: tournamentID,
			Category:     category,
			Pick:         teamID,
		},
	)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toTeamPredictionResponse(pred))
}

// ListMyPlayerPredictions handles GET /tournaments/{tournamentId}/predictions/players.
func (h *TournamentPrediction) ListMyPlayerPredictions(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	views, err := h.svc.ListPlayerPredictionsForUser(r.Context(), tournamentID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]playerPredictionViewResponse, 0, len(views))
	for _, v := range views {
		out = append(out, toPlayerPredictionViewResponse(v))
	}
	writeJSON(w, http.StatusOK, out)
}

// ListMyTeamPredictions handles GET /tournaments/{tournamentId}/predictions/teams.
func (h *TournamentPrediction) ListMyTeamPredictions(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	views, err := h.svc.ListTeamPredictionsForUser(r.Context(), tournamentID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]teamPredictionViewResponse, 0, len(views))
	for _, v := range views {
		out = append(out, toTeamPredictionViewResponse(v))
	}
	writeJSON(w, http.StatusOK, out)
}

// ListLeaguePlayerPredictions handles GET /leagues/{leagueId}/predictions/players.
func (h *TournamentPrediction) ListLeaguePlayerPredictions(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	views, err := h.svc.ListLeaguePlayerPredictions(r.Context(), leagueID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]leaguePlayerCategoryViewResponse, 0, len(views))
	for _, v := range views {
		preds := make([]leagueMemberPlayerPickResponse, 0, len(v.Predictions))
		for _, p := range v.Predictions {
			preds = append(preds, toLeagueMemberPlayerPickResponse(p))
		}
		out = append(
			out, leaguePlayerCategoryViewResponse{
				Category:    string(v.Category),
				Predictions: preds,
			},
		)
	}
	writeJSON(w, http.StatusOK, out)
}

// ListLeagueTeamPredictions handles GET /leagues/{leagueId}/predictions/teams.
func (h *TournamentPrediction) ListLeagueTeamPredictions(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	views, err := h.svc.ListLeagueTeamPredictions(r.Context(), leagueID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]leagueTeamCategoryViewResponse, 0, len(views))
	for _, v := range views {
		preds := make([]leagueMemberTeamPickResponse, 0, len(v.Predictions))
		for _, p := range v.Predictions {
			preds = append(preds, toLeagueMemberTeamPickResponse(p))
		}
		out = append(
			out, leagueTeamCategoryViewResponse{
				Category:    string(v.Category),
				Predictions: preds,
			},
		)
	}
	writeJSON(w, http.StatusOK, out)
}

// compile-time interface checks
var _ TournamentPredictionSvc = (*service.TournamentPredictionService)(nil)
