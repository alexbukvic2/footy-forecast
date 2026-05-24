package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// FixtureRepository handles persistence of Fixture data.
type FixtureRepository struct {
	q *dbgen.Queries
}

// NewFixtureRepository constructs a FixtureRepository backed by pool.
func NewFixtureRepository(pool *db.Pool) *FixtureRepository {
	return &FixtureRepository{q: dbgen.New(pool)}
}

// GetFirstKickoffByTournament returns the earliest kickoff_at across all fixtures
// for the given tournament. Returns domain.ErrNotFound if no fixtures exist.
func (r *FixtureRepository) GetFirstKickoffByTournament(
	ctx context.Context,
	tournamentID uuid.UUID,
) (time.Time, error) {
	t, err := r.q.GetFirstKickoffByTournament(ctx, tournamentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, fmt.Errorf("no fixtures for tournament %s: %w", tournamentID, domain.ErrNotFound)
		}
		return time.Time{}, fmt.Errorf("get first kickoff: %w", err)
	}
	return t, nil
}
