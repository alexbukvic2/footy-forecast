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

// LeagueService is the subset of service.LeagueService that the handler uses.
type LeagueService interface {
	CreateLeague(ctx context.Context, userID uuid.UUID, in domain.CreateLeagueInput) (*domain.League, error)
	GetLeague(ctx context.Context, leagueID, requesterID uuid.UUID) (*domain.League, []*domain.LeagueMember, error)
	ListLeaguesForUser(ctx context.Context, userID uuid.UUID) ([]*domain.LeagueSummary, error)
	UpdateLeagueName(ctx context.Context, leagueID, requesterID uuid.UUID, name string) (*domain.League, error)
	DeleteLeague(ctx context.Context, leagueID, requesterID uuid.UUID) error
	RegenerateCode(ctx context.Context, leagueID, requesterID uuid.UUID) (*domain.League, error)
	JoinLeague(ctx context.Context, code string, userID uuid.UUID) (*domain.League, error)
	RemoveMember(ctx context.Context, leagueID, targetUserID, requesterID uuid.UUID) error
}

// League is the HTTP handler for the /leagues resource.
type League struct {
	logger *slog.Logger
	svc    LeagueService
}

// NewLeague constructs a League handler.
func NewLeague(logger *slog.Logger, svc LeagueService) *League {
	return &League{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type leagueResponse struct {
	ID           string    `json:"id"`
	TournamentID string    `json:"tournament_id"`
	OwnerID      string    `json:"owner_id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type leagueMemberResponse struct {
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type leagueListItemResponse struct {
	leagueResponse
	MemberCount int `json:"member_count"`
	MyPosition  int `json:"my_position"`
}

type leagueDetailResponse struct {
	leagueResponse
	Members []leagueMemberResponse `json:"members"`
}

func toLeagueResponse(l *domain.League) leagueResponse {
	return leagueResponse{
		ID:           l.ID.String(),
		TournamentID: l.TournamentID.String(),
		OwnerID:      l.OwnerID.String(),
		Name:         l.Name,
		Code:         l.Code,
		CreatedAt:    l.CreatedAt,
		UpdatedAt:    l.UpdatedAt,
	}
}

func toLeagueMemberResponse(m *domain.LeagueMember) leagueMemberResponse {
	return leagueMemberResponse{
		UserID:   m.UserID.String(),
		Role:     string(m.Role),
		JoinedAt: m.JoinedAt,
	}
}

// ---------- Handlers ----------

// Create handles POST /leagues.
func (h *League) Create(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	var req struct {
		TournamentID string `json:"tournament_id"`
		Name         string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}

	tournamentID, err := uuid.Parse(req.TournamentID)
	if err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid tournament_id: %w", domain.ErrInvalid))
		return
	}

	league, err := h.svc.CreateLeague(r.Context(), caller.ID, domain.CreateLeagueInput{
		TournamentID: tournamentID,
		Name:         req.Name,
	})
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusCreated, toLeagueResponse(league))
}

// List handles GET /leagues.
func (h *League) List(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	leagues, err := h.svc.ListLeaguesForUser(r.Context(), caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]leagueListItemResponse, 0, len(leagues))
	for _, l := range leagues {
		out = append(out, leagueListItemResponse{
			leagueResponse: toLeagueResponse(l.League),
			MemberCount:    l.MemberCount,
			MyPosition:     l.MyPosition,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"leagues": out})
}

// Get handles GET /leagues/{id}.
func (h *League) Get(w http.ResponseWriter, r *http.Request) {
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

	league, members, err := h.svc.GetLeague(r.Context(), leagueID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	memberResponses := make([]leagueMemberResponse, 0, len(members))
	for _, m := range members {
		memberResponses = append(memberResponses, toLeagueMemberResponse(m))
	}
	writeJSON(w, http.StatusOK, leagueDetailResponse{
		leagueResponse: toLeagueResponse(league),
		Members:        memberResponses,
	})
}

// UpdateName handles PATCH /leagues/{id}.
func (h *League) UpdateName(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}

	league, err := h.svc.UpdateLeagueName(r.Context(), leagueID, caller.ID, req.Name)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toLeagueResponse(league))
}

// Delete handles DELETE /leagues/{id}.
func (h *League) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.svc.DeleteLeague(r.Context(), leagueID, caller.ID); err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// RegenerateCode handles POST /leagues/{id}/code.
func (h *League) RegenerateCode(w http.ResponseWriter, r *http.Request) {
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

	league, err := h.svc.RegenerateCode(r.Context(), leagueID, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toLeagueResponse(league))
}

// Join handles POST /leagues/join.
func (h *League) Join(w http.ResponseWriter, r *http.Request) {
	caller, ok := ctxutil.UserFromCtx(r.Context())
	if !ok {
		writeError(w, r, h.logger, fmt.Errorf("auth user not in context: %w", domain.ErrUnauthorized))
		return
	}

	var req struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, r, h.logger, fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid))
		return
	}
	if req.Code == "" {
		writeError(w, r, h.logger, fmt.Errorf("code is required: %w", domain.ErrInvalid))
		return
	}
	if len(req.Code) > 20 {
		writeError(w, r, h.logger, fmt.Errorf("code is invalid: %w", domain.ErrInvalid))
		return
	}

	league, err := h.svc.JoinLeague(r.Context(), req.Code, caller.ID)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toLeagueResponse(league))
}

// RemoveMember handles DELETE /leagues/{id}/members/{userId}.
func (h *League) RemoveMember(w http.ResponseWriter, r *http.Request) {
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

	targetID, err := parseUUIDPathValue(r, "userId")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	if err := h.svc.RemoveMember(r.Context(), leagueID, targetID, caller.ID); err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ---------- helpers ----------

func parseUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	raw := r.PathValue(key)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s %q: %w", key, raw, domain.ErrInvalid)
	}
	return id, nil
}

// ensure service.LeagueService satisfies our interface at compile time
var _ LeagueService = (*service.LeagueService)(nil)
