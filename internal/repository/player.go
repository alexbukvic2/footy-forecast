package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// PlayerRepository handles persistence of Player aggregates.
type PlayerRepository struct {
	q *dbgen.Queries
}

// NewPlayerRepository constructs a PlayerRepository backed by pool.
func NewPlayerRepository(pool *db.Pool) *PlayerRepository {
	return &PlayerRepository{q: dbgen.New(pool)}
}

// Search returns up to 5 players in the given tournament whose name fuzzy-matches query.
// escapedQuery has SQL wildcard characters escaped (for ILIKE); rawQuery is the original
// trimmed input used for similarity ranking. Both are provided by the caller (service layer).
// Returns an empty (non-nil) slice when nothing matches.
func (r *PlayerRepository) Search(
	ctx context.Context,
	tournamentID uuid.UUID,
	escapedQuery, rawQuery string,
) ([]*domain.PlayerSearchResult, error) {
	rows, err := r.q.SearchPlayers(ctx, dbgen.SearchPlayersParams{
		TournamentID: tournamentID,
		EscapedQuery: escapedQuery,
		RawQuery:     rawQuery,
	})
	if err != nil {
		return nil, fmt.Errorf("search players: %w", err)
	}

	out := make([]*domain.PlayerSearchResult, 0, len(rows))
	for _, row := range rows {
		p, err := playerSearchResultFromRow(row)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func playerSearchResultFromRow(row dbgen.SearchPlayersRow) (*domain.PlayerSearchResult, error) {
	var raw map[string]int
	if err := json.Unmarshal([]byte(row.Handicaps), &raw); err != nil {
		return nil, fmt.Errorf("unmarshal handicaps for player %s: %w", row.ID, err)
	}
	h := make(map[domain.PlayerHandicapCategory]int, len(raw))
	for k, v := range raw {
		h[domain.PlayerHandicapCategory(k)] = v
	}
	return &domain.PlayerSearchResult{
		ID:           row.ID,
		Name:         row.Name,
		TournamentID: row.TournamentID,
		TeamName:     row.TeamName,
		TeamLogo:     row.TeamLogo,
		Handicaps:    h,
	}, nil
}
