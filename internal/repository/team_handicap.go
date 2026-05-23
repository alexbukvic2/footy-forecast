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

// TeamHandicapRepository handles persistence of TeamHandicap aggregates.
type TeamHandicapRepository struct {
	q *dbgen.Queries
}

// NewTeamHandicapRepository constructs a TeamHandicapRepository backed by pool.
func NewTeamHandicapRepository(pool *db.Pool) *TeamHandicapRepository {
	return &TeamHandicapRepository{q: dbgen.New(pool)}
}

// Get returns the handicap for the given team and category.
// Returns domain.ErrNotFound when no row exists.
func (r *TeamHandicapRepository) Get(
	ctx context.Context,
	teamID uuid.UUID,
	category domain.TeamHandicapCategory,
) (*domain.TeamHandicap, error) {
	row, err := r.q.GetTeamHandicap(ctx, dbgen.GetTeamHandicapParams{
		TeamID:   teamID,
		Category: dbgen.TeamHandicapCategory(category),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get team handicap: %w", err)
	}
	return &domain.TeamHandicap{
		ID:       row.ID,
		TeamID:   row.TeamID,
		Category: domain.TeamHandicapCategory(row.Category),
		Points:   int(row.Points),
	}, nil
}
