package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// FixtureForUserLister retrieves fixtures paired with a user's score predictions.
type FixtureForUserLister interface {
	ListByTournamentForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.UserFixtureView, error)
	ListLockedByLeague(
		ctx context.Context,
		leagueID, requestingUserID uuid.UUID,
	) ([]*domain.LeagueFixtureView, error)
	GetLockedFixtureDates(
		ctx context.Context,
		leagueID uuid.UUID,
	) ([]time.Time, error)
	ListLockedByLeagueAndDates(
		ctx context.Context,
		leagueID, requestingUserID uuid.UUID,
		dates []time.Time,
	) ([]*domain.LeagueFixtureView, error)
}

// LeagueMemberChecker verifies league membership.
type LeagueMemberChecker interface {
	GetMember(
		ctx context.Context,
		leagueID, userID uuid.UUID,
	) (*domain.LeagueMember, error)
}

// FixtureService handles fixture listing use cases.
type FixtureService struct {
	fixtures FixtureForUserLister
	leagues  LeagueMemberChecker
}

// NewFixtureService constructs a FixtureService.
func NewFixtureService(
	fixtures FixtureForUserLister,
	leagues LeagueMemberChecker,
) *FixtureService {
	return &FixtureService{fixtures: fixtures, leagues: leagues}
}

// ListForUser returns all fixtures for a tournament paired with the user's score prediction.
func (s *FixtureService) ListForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.UserFixtureView, error) {
	views, err := s.fixtures.ListByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return nil, fmt.Errorf("list fixtures for user: %w", err)
	}
	return views, nil
}

// ListForLeague returns all locked fixtures for a league with every member's predictions.
// Returns domain.ErrForbidden if the requesting user is not a member of the league.
func (s *FixtureService) ListForLeague(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) ([]*domain.LeagueFixtureView, error) {
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

// ListForLeaguePaged returns locked fixtures for a league paginated by day.
// n is the number of days-with-fixtures to fetch; skip is the number of most-recent
// days to skip before fetching.
// Returns domain.ErrForbidden if the requesting user is not a member of the league.
func (s *FixtureService) ListForLeaguePaged(
	ctx context.Context,
	leagueID, userID uuid.UUID,
	n, skip int,
) ([]*domain.LeagueFixtureView, error) {
	if _, err := s.leagues.GetMember(ctx, leagueID, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("not a member of league %s: %w", leagueID, domain.ErrForbidden)
		}
		return nil, fmt.Errorf("get league member: %w", err)
	}

	allDates, err := s.fixtures.GetLockedFixtureDates(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get locked fixture dates: %w", err)
	}

	if skip >= len(allDates) {
		return []*domain.LeagueFixtureView{}, nil
	}
	end := skip + n
	if end > len(allDates) {
		end = len(allDates)
	}
	pageDates := allDates[skip:end]

	views, err := s.fixtures.ListLockedByLeagueAndDates(ctx, leagueID, userID, pageDates)
	if err != nil {
		return nil, fmt.Errorf("list locked fixtures by dates: %w", err)
	}
	return views, nil
}

// ---------- PredictionService ----------

// PredictionUpserter persists a score prediction.
type PredictionUpserter interface {
	Upsert(
		ctx context.Context,
		in domain.UpsertScorePredictionInput,
	) (*domain.ScorePrediction, error)
}

// FixtureGetter fetches a single fixture by ID.
type FixtureGetter interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Fixture, error)
}

// PredictionService handles score prediction use cases.
type PredictionService struct {
	predictions PredictionUpserter
	fixtures    FixtureGetter
}

// NewPredictionService constructs a PredictionService.
func NewPredictionService(
	predictions PredictionUpserter,
	fixtures FixtureGetter,
) *PredictionService {
	return &PredictionService{predictions: predictions, fixtures: fixtures}
}

// UpsertScore validates and stores a user's score prediction.
// Returns domain.ErrInvalid for negative goals, missing winner on knockout fixtures,
// or a winner that isn't one of the fixture's two teams.
// Returns domain.ErrNotFound if the fixture does not exist.
// Returns domain.ErrForbidden if predictions are locked (at or past kickoff).
func (s *PredictionService) UpsertScore(
	ctx context.Context,
	in domain.UpsertScorePredictionInput,
) (*domain.ScorePrediction, error) {
	if in.GoalsHome < 0 || in.GoalsAway < 0 {
		return nil, fmt.Errorf("goals must be non-negative: %w", domain.ErrInvalid)
	}

	fixture, err := s.fixtures.GetByID(ctx, in.FixtureID)
	if err != nil {
		return nil, fmt.Errorf("get fixture: %w", err)
	}

	if fixture.PredictionLocked {
		return nil, fmt.Errorf("predictions locked close to kickoff: %w", domain.ErrForbidden)
	}

	// Knockout fixtures (group == nil) require the user to specify which team advances.
	isKnockout := fixture.Group == nil
	if isKnockout {
		if in.Winner == nil {
			return nil, fmt.Errorf("winner is required for knockout fixtures: %w", domain.ErrInvalid)
		}
		if *in.Winner != fixture.HomeTeamID && *in.Winner != fixture.AwayTeamID {
			return nil, fmt.Errorf("winner must be the home or away team of the fixture: %w", domain.ErrInvalid)
		}
	}

	pred, err := s.predictions.Upsert(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("upsert score prediction: %w", err)
	}
	return pred, nil
}
