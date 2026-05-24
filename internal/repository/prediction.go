package repository

import (
	"context"
	"fmt"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// PredictionRepository handles persistence of ScorePrediction aggregates.
type PredictionRepository struct {
	pool *db.Pool
	q    *dbgen.Queries
}

// NewPredictionRepository constructs a PredictionRepository backed by pool.
func NewPredictionRepository(pool *db.Pool) *PredictionRepository {
	return &PredictionRepository{pool: pool, q: dbgen.New(pool)}
}

// Upsert creates or updates a score prediction, returning the persisted record.
func (r *PredictionRepository) Upsert(
	ctx context.Context,
	in domain.UpsertScorePredictionInput,
) (*domain.ScorePrediction, error) {
	row, err := r.q.UpsertScorePrediction(ctx, dbgen.UpsertScorePredictionParams{
		UserID:    in.UserID,
		FixtureID: in.FixtureID,
		GoalsHome: int32(in.GoalsHome), //nolint:gosec // validated >= 0 by service
		GoalsAway: int32(in.GoalsAway), //nolint:gosec // validated >= 0 by service
	})
	if err != nil {
		return nil, fmt.Errorf("upsert score prediction: %w", err)
	}
	return scorePredictionFromRow(row), nil
}

func scorePredictionFromRow(row dbgen.ScorePrediction) *domain.ScorePrediction {
	p := &domain.ScorePrediction{
		ID:        row.ID,
		UserID:    row.UserID,
		FixtureID: row.FixtureID,
		GoalsHome: int(row.GoalsHome),
		GoalsAway: int(row.GoalsAway),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.Points != nil {
		v := int(*row.Points)
		p.Points = &v
	}
	return p
}
