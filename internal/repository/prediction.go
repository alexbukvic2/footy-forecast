package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// PredictionRepository handles persistence of ScorePrediction aggregates.
type PredictionRepository struct {
	q *dbgen.Queries
}

// NewPredictionRepository constructs a PredictionRepository backed by pool.
func NewPredictionRepository(pool *db.Pool) *PredictionRepository {
	return &PredictionRepository{q: dbgen.New(pool)}
}

// Upsert inserts or updates a score prediction for a fixture.
func (r *PredictionRepository) Upsert(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
	winner := uuid.Nil
	if in.Winner != nil {
		winner = *in.Winner
	}
	row, err := r.q.UpsertScorePrediction(ctx, dbgen.UpsertScorePredictionParams{
		UserID:    in.UserID,
		FixtureID: in.FixtureID,
		GoalsHome: int32(in.GoalsHome), //nolint:gosec
		GoalsAway: int32(in.GoalsAway), //nolint:gosec
		Winner:    winner,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert score prediction: %w", err)
	}
	pred := &domain.ScorePrediction{
		ID:        row.ID,
		UserID:    row.UserID,
		FixtureID: row.FixtureID,
		GoalsHome: int(row.GoalsHome),
		GoalsAway: int(row.GoalsAway),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.Winner != (uuid.UUID{}) {
		w := row.Winner
		pred.Winner = &w
	}
	if row.Points != nil {
		v := int(*row.Points)
		pred.Points = &v
	}
	return pred, nil
}
