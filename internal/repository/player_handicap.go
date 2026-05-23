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

// PlayerHandicapRepository handles persistence of PlayerHandicap aggregates.
type PlayerHandicapRepository struct {
	q *dbgen.Queries
}

// NewPlayerHandicapRepository constructs a PlayerHandicapRepository backed by pool.
func NewPlayerHandicapRepository(pool *db.Pool) *PlayerHandicapRepository {
	return &PlayerHandicapRepository{q: dbgen.New(pool)}
}

// Get returns the handicap for the given player and category.
// Returns domain.ErrNotFound when no row exists.
func (r *PlayerHandicapRepository) Get(
	ctx context.Context,
	playerID uuid.UUID,
	category domain.PlayerHandicapCategory,
) (*domain.PlayerHandicap, error) {
	row, err := r.q.GetPlayerHandicap(ctx, dbgen.GetPlayerHandicapParams{
		PlayerID: playerID,
		Category: dbgen.PlayerHandicapCategory(category),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get player handicap: %w", err)
	}
	return &domain.PlayerHandicap{
		ID:       row.ID,
		PlayerID: row.PlayerID,
		Category: domain.PlayerHandicapCategory(row.Category),
		Points:   int(row.Points),
	}, nil
}
