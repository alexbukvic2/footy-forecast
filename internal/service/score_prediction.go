package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

const lockWindow = 30 * time.Minute

// FixtureForUserLister retrieves fixtures paired with a user's score predictions.
type FixtureForUserLister interface {
	ListByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error)
	ListLockedByLeague(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeagueFixtureView, error)
}

// LeagueMemberChecker verifies league membership.
type LeagueMemberChecker interface {
	GetMember(ctx context.Context, leagueID, userID uuid.UUID) (*domain.LeagueMember, error)
}

// FixtureService handles fixture listing use cases.
type FixtureService struct {
	fixtures FixtureForUserLister
	leagues  LeagueMemberChecker
}

// NewFixtureService constructs a FixtureService.
func NewFixtureService(fixtures FixtureForUserLister, leagues LeagueMemberChecker) *FixtureService {
	return &FixtureService{fixtures: fixtures, leagues: leagues}
}

// ListForUser returns all fixtures for a tournament paired with the user's score prediction.
func (s *FixtureService) ListForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error) {
	views, err := s.fixtures.ListByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return nil, fmt.Errorf("list fixtures for user: %w", err)
	}
	return views, nil
}

// ListForLeague returns all locked fixtures for a league with every member's predictions.
// Returns domain.ErrForbidden if the requesting user is not a member of the league.
func (s *FixtureService) ListForLeague(ctx context.Context, leagueID, userID uuid.UUID) ([]*domain.LeagueFixtureView, error) {
	if _, err := s.leagues.GetMember(ctx, leagueID, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("not a member of league %s: %w", leagueID, domain.ErrForbidden)
		}
		return nil, fmt.Errorf("get league member: %w", err)
	}
	views, err := s.fixtures.ListLockedByLeague(ctx, leagueID, userID)
	if err != nil {
		return nil, fmt.Errorf("list locked fixtures for league: %w", err)
	}
	return views, nil
}

// ---------- PredictionService ----------

// PredictionUpserter persists a score prediction.
type PredictionUpserter interface {
	Upsert(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
}

// FixtureGetter fetches a single fixture by ID.
type FixtureGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Fixture, error)
}

// PredictionService handles score prediction use cases.
type PredictionService struct {
	predictions PredictionUpserter
	fixtures    FixtureGetter
	clock       Clock
}

// NewPredictionService constructs a PredictionService.
func NewPredictionService(predictions PredictionUpserter, fixtures FixtureGetter, clock Clock) *PredictionService {
	return &PredictionService{predictions: predictions, fixtures: fixtures, clock: clock}
}

// UpsertScore validates and stores a user's score prediction.
// Returns domain.ErrInvalid for negative goals.
// Returns domain.ErrNotFound if the fixture does not exist.
// Returns domain.ErrForbidden if within 30 minutes of kickoff (or past it).
func (s *PredictionService) UpsertScore(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
	if in.GoalsHome < 0 || in.GoalsAway < 0 {
		return nil, fmt.Errorf("goals must be non-negative: %w", domain.ErrInvalid)
	}

	fixture, err := s.fixtures.GetByID(ctx, in.FixtureID)
	if err != nil {
		return nil, fmt.Errorf("get fixture: %w", err)
	}

	lockAt := fixture.KickoffAt.Add(-lockWindow)
	if !s.clock.Now().Before(lockAt) {
		return nil, fmt.Errorf("predictions lock 30 minutes before kickoff: %w", domain.ErrForbidden)
	}

	pred, err := s.predictions.Upsert(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("upsert score prediction: %w", err)
	}
	return pred, nil
}
