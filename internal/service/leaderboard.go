package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// LeaderboardRepo is the subset of the leaderboard repository the service needs.
type LeaderboardRepo interface {
	GetForLeague(
		ctx context.Context,
		leagueID uuid.UUID,
	) ([]*domain.LeaderboardEntry, error)
	GetForTournament(
		ctx context.Context,
		tournamentID uuid.UUID,
	) ([]*domain.LeaderboardEntry, error)
	GetUserPositionsInLeagues(
		ctx context.Context,
		userID uuid.UUID,
		leagueIDs []uuid.UUID,
	) (map[uuid.UUID]int, error)
}

// LeaderboardService orchestrates leaderboard use cases.
type LeaderboardService struct {
	lb          LeaderboardRepo
	leagues     leagueMemberGetter
	tournaments TournamentGetter
}

// leagueMemberGetter is the narrow interface used for membership checks.
type leagueMemberGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.League, error)
	GetMember(ctx context.Context, leagueID, userID uuid.UUID) (*domain.LeagueMember, error)
}

// NewLeaderboardService constructs a LeaderboardService.
func NewLeaderboardService(
	lb LeaderboardRepo,
	leagues leagueMemberGetter,
	tournaments TournamentGetter,
) *LeaderboardService {
	return &LeaderboardService{lb: lb, leagues: leagues, tournaments: tournaments}
}

// GetLeagueLeaderboard returns the ranked leaderboard for a league.
// The requesting user must be a member of the league.
// Returns ErrNotFound when the league does not exist.
// Returns ErrForbidden when the requester is not a member.
func (s *LeaderboardService) GetLeagueLeaderboard(
	ctx context.Context,
	leagueID, requesterID uuid.UUID,
) ([]*domain.LeaderboardEntry, error) {
	// Verify league exists.
	if _, err := s.leagues.GetByID(ctx, leagueID); err != nil {
		return nil, fmt.Errorf("get league: %w", err)
	}

	// Verify membership.
	_, err := s.leagues.GetMember(ctx, leagueID, requesterID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("requester is not a member of league %s: %w", leagueID, domain.ErrForbidden)
		}
		return nil, fmt.Errorf("check membership: %w", err)
	}

	entries, err := s.lb.GetForLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get league leaderboard: %w", err)
	}
	return entries, nil
}

// GetTournamentLeaderboard returns the global ranked leaderboard for a tournament.
// Returns ErrNotFound when the tournament does not exist.
func (s *LeaderboardService) GetTournamentLeaderboard(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]*domain.LeaderboardEntry, error) {
	if _, err := s.tournaments.GetByID(ctx, tournamentID); err != nil {
		return nil, fmt.Errorf("get tournament: %w", err)
	}

	entries, err := s.lb.GetForTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("get tournament leaderboard: %w", err)
	}
	return entries, nil
}

// compile-time interface check
var _ LeaderboardRepo = (*repository.LeaderboardRepository)(nil)
