package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// Clock allows injecting a deterministic time source for testing.
type Clock interface {
	Now() time.Time
}

// RealClock returns the actual current time.
type RealClock struct{}

// Now returns time.Now().
func (RealClock) Now() time.Time { return time.Now() }

// PredictionRepo is the subset of PredictionRepository that PredictionService needs.
type PredictionRepo interface {
	Upsert(
		ctx context.Context,
		in domain.UpsertScorePredictionInput,
	) (*domain.ScorePrediction, error)
}

// FixtureGetter is the subset of FixtureRepository that PredictionService needs.
type FixtureGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Fixture, error)
}

// PredictionService orchestrates prediction use cases.
type PredictionService struct {
	predictions PredictionRepo
	fixtures    FixtureGetter
	clock       Clock
}

// NewPredictionService constructs a PredictionService.
func NewPredictionService(
	predictions PredictionRepo,
	fixtures FixtureGetter,
	clock Clock,
) *PredictionService {
	return &PredictionService{predictions: predictions, fixtures: fixtures, clock: clock}
}

// UpsertScore creates or updates a score prediction.
//
// Returns:
//   - domain.ErrInvalid if goals_home or goals_away are negative
//   - domain.ErrNotFound if the fixture does not exist
//   - domain.ErrForbidden if the fixture status is not upcoming, or within 30 minutes of kickoff
func (s *PredictionService) UpsertScore(
	ctx context.Context,
	in domain.UpsertScorePredictionInput,
) (*domain.ScorePrediction, error) {
	if in.GoalsHome < 0 {
		return nil, fmt.Errorf("goals_home must be >= 0: %w", domain.ErrInvalid)
	}
	if in.GoalsAway < 0 {
		return nil, fmt.Errorf("goals_away must be >= 0: %w", domain.ErrInvalid)
	}

	fixture, err := s.fixtures.GetByID(ctx, in.FixtureID)
	if err != nil {
		return nil, fmt.Errorf("get fixture: %w", err)
	}

	if fixture.Status != domain.FixtureStatusUpcoming {
		return nil, fmt.Errorf("fixture %s is locked for predictions: %w", in.FixtureID, domain.ErrForbidden)
	}

	lockAt := fixture.KickoffAt.Add(-30 * time.Minute)
	if !s.clock.Now().Before(lockAt) {
		return nil, fmt.Errorf("fixture %s is locked for predictions: %w", in.FixtureID, domain.ErrForbidden)
	}

	return s.predictions.Upsert(ctx, in)
}

// compile-time interface checks
var _ PredictionRepo = (*repository.PredictionRepository)(nil)
var _ FixtureGetter = (*repository.FixtureRepository)(nil)
