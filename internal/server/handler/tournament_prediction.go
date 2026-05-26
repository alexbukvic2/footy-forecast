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
	BulkUpsertPlayerPredictions(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
		items []domain.BulkPlayerPredictionItem,
	) ([]*domain.PlayerPredictionView, error)
	BulkUpsertTeamPredictions(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
		items []domain.BulkTeamPredictionItem,
	) ([]*domain.TeamPredictionView, error)
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

type playerPredictionViewResponse struct {
	Category   string  `json:"category"`
	Group      *string `json:"group,omitempty"`
	PlayerID   *string `json:"player_id"`
	PlayerName *string `json:"player_name"`
	Points     *int    `json:"points"`
}

type teamPredictionViewResponse struct {
	Category string  `json:"category"`
	Group    *string `json:"group,omitempty"`
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
	Group       *string                          `json:"group,omitempty"`
	Predictions []leagueMemberPlayerPickResponse `json:"predictions"`
}

type leagueTeamCategoryViewResponse struct {
	Category    string                         `json:"category"`
	Group       *string                        `json:"group,omitempty"`
	Predictions []leagueMemberTeamPickResponse `json:"predictions"`
}

// ---------- request DTOs ----------

type bulkPlayerPredictionItemRequest struct {
	Category    string  `json:"category"`
	GroupLetter *string `json:"group_letter"`
	PlayerID    *string `json:"player_id"`
}

type bulkTeamPredictionItemRequest struct {
	Category    string  `json:"category"`
	GroupLetter *string `json:"group_letter"`
	SlotIndex   int     `json:"slot_index"`
	TeamID      *string `json:"team_id"`
}

func toPlayerPredictionViewResponse(v *domain.PlayerPredictionView) playerPredictionViewResponse {
	resp := playerPredictionViewResponse{
		Category: string(v.Category),
		Group:    v.GroupLetter,
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
		Group:    v.GroupLetter,
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

// BulkUpsertPlayerPredictions handles PUT /tournaments/{tournamentId}/predictions/players.
func (h *TournamentPrediction) BulkUpsertPlayerPredictions(
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

	var rawItems []bulkPlayerPredictionItemRequest
	if err := json.NewDecoder(r.Body).Decode(&rawItems); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}

	items := make([]domain.BulkPlayerPredictionItem, 0, len(rawItems))
	for i, raw := range rawItems {
		cat, err := domain.ParsePlayerHandicapCategory(raw.Category)
		if err != nil {
			writeError(w, r, h.logger, fmt.Errorf("item %d: invalid category %q: %w", i, raw.Category, domain.ErrInvalid))
			return
		}
		item := domain.BulkPlayerPredictionItem{
			Category:    cat,
			GroupLetter: raw.GroupLetter,
		}
		if raw.PlayerID != nil {
			id, err := uuid.Parse(*raw.PlayerID)
			if err != nil {
				writeError(w, r, h.logger, fmt.Errorf("item %d: invalid player_id: %w", i, domain.ErrInvalid))
				return
			}
			item.PlayerID = &id
		}
		items = append(items, item)
	}

	views, err := h.svc.BulkUpsertPlayerPredictions(r.Context(), tournamentID, caller.ID, items)
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

// BulkUpsertTeamPredictions handles PUT /tournaments/{tournamentId}/predictions/teams.
func (h *TournamentPrediction) BulkUpsertTeamPredictions(
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

	var rawItems []bulkTeamPredictionItemRequest
	if err := json.NewDecoder(r.Body).Decode(&rawItems); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}

	items := make([]domain.BulkTeamPredictionItem, 0, len(rawItems))
	for i, raw := range rawItems {
		cat, err := domain.ParseTeamHandicapCategory(raw.Category)
		if err != nil {
			writeError(w, r, h.logger, fmt.Errorf("item %d: invalid category %q: %w", i, raw.Category, domain.ErrInvalid))
			return
		}
		item := domain.BulkTeamPredictionItem{
			Category:    cat,
			GroupLetter: raw.GroupLetter,
			SlotIndex:   raw.SlotIndex,
		}
		if raw.TeamID != nil {
			id, err := uuid.Parse(*raw.TeamID)
			if err != nil {
				writeError(w, r, h.logger, fmt.Errorf("item %d: invalid team_id: %w", i, domain.ErrInvalid))
				return
			}
			item.TeamID = &id
		}
		items = append(items, item)
	}

	views, err := h.svc.BulkUpsertTeamPredictions(r.Context(), tournamentID, caller.ID, items)
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

// ---------- Handlers ----------

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
				Group:       v.GroupLetter,
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
				Group:       v.GroupLetter,
				Predictions: preds,
			},
		)
	}
	writeJSON(w, http.StatusOK, out)
}

// compile-time interface checks
var _ TournamentPredictionSvc = (*service.TournamentPredictionService)(nil)
