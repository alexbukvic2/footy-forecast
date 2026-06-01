package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// TournamentGroupTableRepository handles persistence for the group-stage standings table.
type TournamentGroupTableRepository struct {
	q *dbgen.Queries
}

// NewTournamentGroupTableRepository constructs a TournamentGroupTableRepository backed by pool.
func NewTournamentGroupTableRepository(pool *db.Pool) *TournamentGroupTableRepository {
	return &TournamentGroupTableRepository{q: dbgen.New(pool)}
}

// ListByTournament returns all group-table entries for a tournament, sorted by
// group_letter ASC, position ASC. Returns an empty slice (not an error) when
// the tournament exists but has no entries yet.
func (r *TournamentGroupTableRepository) ListByTournament(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]*domain.TournamentGroupEntry, error) {
	rows, err := r.q.ListGroupTableByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list group table: %w", err)
	}
	out := make([]*domain.TournamentGroupEntry, 0, len(rows))
	for _, row := range rows {
		out = append(out, &domain.TournamentGroupEntry{
			ID:           row.ID,
			TournamentID: row.TournamentID,
			TeamID:       row.TeamID,
			TeamName:     row.TeamName,
			GroupLetter:  row.GroupLetter,
			Position:     int(row.Position),
			Points:       int(row.Points),
			Played:       int(row.Played),
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		})
	}
	return out, nil
}
