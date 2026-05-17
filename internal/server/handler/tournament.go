package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexbukvic2/footy-forecast/internal/service"
	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// TournamentService is the subset of service.TournamentService that the handler uses.
type TournamentService interface {
	Create(
		ctx context.Context,
		in domain.CreateTournamentInput,
	) (*domain.Tournament, error)
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Tournament, error)
	GetBySlug(
		ctx context.Context,
		slug string,
	) (*domain.Tournament, error)
	List(ctx context.Context) ([]*domain.Tournament, error)
}

// Tournament is the HTTP handler for the /tournaments resource.
type Tournament struct {
	logger *slog.Logger
	svc    TournamentService
}

// NewTournament constructs a Tournament handler.
func NewTournament(
	logger *slog.Logger,
	svc TournamentService,
) *Tournament {
	return &Tournament{logger: logger, svc: svc}
}

// ---------- DTOs ----------

// tournamentResponse is the JSON shape returned to clients for a single tournament.
type tournamentResponse struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toTournamentResponse(t *domain.Tournament) tournamentResponse {
	return tournamentResponse{
		ID:        t.ID.String(),
		Slug:      t.Slug,
		Name:      t.Name,
		Status:    string(t.Status),
		StartsAt:  t.StartsAt,
		EndsAt:    t.EndsAt,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// createTournamentRequest is the JSON shape clients post to create a tournament.
type createTournamentRequest struct {
	Slug     string    `json:"slug"`
	Name     string    `json:"name"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

// ---------- Handlers ----------

// Create handles POST /tournaments.
func (h *Tournament) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req createTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w, r, h.logger,
			fmt.Errorf("invalid JSON body: %w", domain.ErrInvalid),
		)
		return
	}

	t, err := h.svc.Create(
		r.Context(), domain.CreateTournamentInput{
			Slug:     req.Slug,
			Name:     req.Name,
			StartsAt: req.StartsAt,
			EndsAt:   req.EndsAt,
		},
	)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusCreated, toTournamentResponse(t))
}

// List handles GET /tournaments.
func (h *Tournament) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	ts, err := h.svc.List(r.Context())
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	out := make([]tournamentResponse, 0, len(ts))
	for _, t := range ts {
		out = append(out, toTournamentResponse(t))
	}
	writeJSON(w, http.StatusOK, map[string]any{"tournaments": out})
}

// GetByID handles GET /tournaments/{id}.
func (h *Tournament) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(
			w, r, h.logger,
			fmt.Errorf("invalid id %q: %w", idStr, domain.ErrInvalid),
		)
		return
	}

	t, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toTournamentResponse(t))
}

// GetBySlug handles GET /tournaments/slug/{slug}.
func (h *Tournament) GetBySlug(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")
	t, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}
	writeJSON(w, http.StatusOK, toTournamentResponse(t))
}

// ensure service.TournamentService satisfies our interface at compile time
var _ TournamentService = (*service.TournamentService)(nil)
