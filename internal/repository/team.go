package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// TeamRepository handles persistence of Team aggregates.
type TeamRepository struct {
	q *dbgen.Queries
}

// NewTeamRepository constructs a TeamRepository backed by pool.
func NewTeamRepository(pool *db.Pool) *TeamRepository {
	return &TeamRepository{q: dbgen.New(pool)}
}

// GetByID fetches a team by its UUID.
// Returns domain.ErrNotFound if no row exists.
func (r *TeamRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Team, error) {
	row, err := r.q.GetTeamByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("team %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get team: %w", err)
	}
	return &domain.Team{
		ID:           row.ID,
		Name:         row.Name,
		Logo:         row.Logo,
		TournamentID: row.TournamentID,
	}, nil
}
