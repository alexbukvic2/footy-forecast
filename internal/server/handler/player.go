package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// PlayerService is the subset of service.PlayerService that the handler uses.
type PlayerService interface {
	Search(ctx context.Context, in domain.SearchPlayersInput) ([]*domain.PlayerSearchResult, error)
}

// Player is the HTTP handler for the players resource.
type Player struct {
	logger *slog.Logger
	svc    PlayerService
}

// NewPlayer constructs a Player handler.
func NewPlayer(logger *slog.Logger, svc PlayerService) *Player {
	return &Player{logger: logger, svc: svc}
}

// ---------- DTOs ----------

type playerResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TeamName string `json:"team_name"`
	TeamLogo string `json:"team_logo"`
}

func toPlayerResponse(p *domain.PlayerSearchResult) playerResponse {
	return playerResponse{
		ID:       p.ID.String(),
		Name:     p.Name,
		TeamName: p.TeamName,
		TeamLogo: p.TeamLogo,
	}
}

// ---------- Handlers ----------

// Search handles GET /tournaments/{tournament_id}/players/search?q=<string>.
func (h *Player) Search(w http.ResponseWriter, r *http.Request) {
	tournamentID, err := parseUUIDPathValue(r, "tournament_id")
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	q := r.URL.Query().Get("q")

	players, err := h.svc.Search(r.Context(), domain.SearchPlayersInput{
		TournamentID: tournamentID,
		Query:        q,
	})
	if err != nil {
		writeError(w, r, h.logger, err)
		return
	}

	out := make([]playerResponse, 0, len(players))
	for _, p := range players {
		out = append(out, toPlayerResponse(p))
	}
	writeJSON(w, http.StatusOK, map[string]any{"players": out})
}

// ensure service.PlayerService satisfies our interface at compile time
var _ PlayerService = (*service.PlayerService)(nil)
