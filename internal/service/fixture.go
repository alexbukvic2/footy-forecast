package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// FixtureLister is the subset of FixtureRepository that FixtureService needs.
type FixtureLister interface {
	ListByTournamentForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.UserFixtureView, error)
	ListLockedByLeague(
		ctx context.Context,
		leagueID, requestingUserID uuid.UUID,
	) ([]*domain.LeagueFixtureView, error)
}

// LeagueMemberChecker is the subset of LeagueRepository that FixtureService needs.
type LeagueMemberChecker interface {
	IsMember(ctx context.Context, leagueID, userID uuid.UUID) (bool, error)
}

// FixtureService orchestrates fixture use cases.
type FixtureService struct {
	fixtures FixtureLister
	leagues  LeagueMemberChecker
}

// NewFixtureService constructs a FixtureService.
func NewFixtureService(fixtures FixtureLister, leagues LeagueMemberChecker) *FixtureService {
	return &FixtureService{fixtures: fixtures, leagues: leagues}
}

// ListForUser returns all fixtures for a tournament paired with the user's predictions.
func (s *FixtureService) ListForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.UserFixtureView, error) {
	return s.fixtures.ListByTournamentForUser(ctx, tournamentID, userID)
}

// ListForLeague returns locked fixtures for a league with all member predictions.
// Returns domain.ErrForbidden if userID is not a member of leagueID.
func (s *FixtureService) ListForLeague(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) ([]*domain.LeagueFixtureView, error) {
	ok, err := s.leagues.IsMember(ctx, leagueID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("user %s is not a member of league %s: %w", userID, leagueID, domain.ErrForbidden)
	}
	return s.fixtures.ListLockedByLeague(ctx, leagueID, userID)
}

// compile-time interface checks
var _ FixtureLister = (*repository.FixtureRepository)(nil)
var _ LeagueMemberChecker = (*repository.LeagueRepository)(nil)
