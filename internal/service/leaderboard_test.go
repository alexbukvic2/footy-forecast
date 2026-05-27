package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// --- fakes ---

type fakeLBRepo struct {
	forLeague        []*domain.LeaderboardEntry
	forTournament    []*domain.LeaderboardEntry
	forLeagueErr     error
	forTournamentErr error
}

func (f *fakeLBRepo) GetForLeague(_ context.Context, _ uuid.UUID) ([]*domain.LeaderboardEntry, error) {
	return f.forLeague, f.forLeagueErr
}
func (f *fakeLBRepo) GetForTournament(_ context.Context, _ uuid.UUID) ([]*domain.LeaderboardEntry, error) {
	return f.forTournament, f.forTournamentErr
}
func (f *fakeLBRepo) GetUserPositionsInLeagues(_ context.Context, _ uuid.UUID, _ []uuid.UUID) (map[uuid.UUID]int, error) {
	return nil, nil
}

type fakeLBLeagueGetter struct {
	league    *domain.League
	leagueErr error
	member    *domain.LeagueMember
	memberErr error
}

func (f *fakeLBLeagueGetter) GetByID(_ context.Context, _ uuid.UUID) (*domain.League, error) {
	return f.league, f.leagueErr
}
func (f *fakeLBLeagueGetter) GetMember(_ context.Context, _, _ uuid.UUID) (*domain.LeagueMember, error) {
	return f.member, f.memberErr
}

type fakeLBTournamentGetter struct {
	tournament *domain.Tournament
	err        error
}

func (f *fakeLBTournamentGetter) GetByID(_ context.Context, _ uuid.UUID) (*domain.Tournament, error) {
	return f.tournament, f.err
}

// --- tests ---

func TestLeaderboardService_GetLeagueLeaderboard(t *testing.T) {
	ctx := context.Background()
	leagueID := uuid.Must(uuid.NewV7())
	requesterID := uuid.Must(uuid.NewV7())
	league := &domain.League{ID: leagueID}
	member := &domain.LeagueMember{UserID: requesterID}

	t.Run("league not found returns ErrNotFound", func(t *testing.T) {
		svc := service.NewLeaderboardService(
			&fakeLBRepo{},
			&fakeLBLeagueGetter{leagueErr: domain.ErrNotFound},
			&fakeLBTournamentGetter{},
		)
		_, err := svc.GetLeagueLeaderboard(ctx, leagueID, requesterID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("non-member returns ErrForbidden", func(t *testing.T) {
		svc := service.NewLeaderboardService(
			&fakeLBRepo{},
			&fakeLBLeagueGetter{league: league, memberErr: domain.ErrNotFound},
			&fakeLBTournamentGetter{},
		)
		_, err := svc.GetLeagueLeaderboard(ctx, leagueID, requesterID)
		if !errors.Is(err, domain.ErrForbidden) {
			t.Fatalf("expected ErrForbidden, got %v", err)
		}
	})

	t.Run("member gets leaderboard result", func(t *testing.T) {
		want := []*domain.LeaderboardEntry{{Position: 1, UserID: requesterID, TotalPoints: 10, WinnerPts: 10}}
		svc := service.NewLeaderboardService(
			&fakeLBRepo{forLeague: want},
			&fakeLBLeagueGetter{league: league, member: member},
			&fakeLBTournamentGetter{},
		)
		got, err := svc.GetLeagueLeaderboard(ctx, leagueID, requesterID)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].Position != 1 {
			t.Fatalf("unexpected result: %v", got)
		}
	})
}

func TestLeaderboardService_GetTournamentLeaderboard(t *testing.T) {
	ctx := context.Background()
	tournamentID := uuid.Must(uuid.NewV7())

	t.Run("tournament not found returns ErrNotFound", func(t *testing.T) {
		svc := service.NewLeaderboardService(
			&fakeLBRepo{},
			&fakeLBLeagueGetter{},
			&fakeLBTournamentGetter{err: domain.ErrNotFound},
		)
		_, err := svc.GetTournamentLeaderboard(ctx, tournamentID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("valid tournament returns leaderboard", func(t *testing.T) {
		want := []*domain.LeaderboardEntry{{Position: 1, TotalPoints: 5, WinnerPts: 5}}
		svc := service.NewLeaderboardService(
			&fakeLBRepo{forTournament: want},
			&fakeLBLeagueGetter{},
			&fakeLBTournamentGetter{tournament: &domain.Tournament{ID: tournamentID}},
		)
		got, err := svc.GetTournamentLeaderboard(ctx, tournamentID)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].Position != 1 {
			t.Fatalf("unexpected result: %v", got)
		}
	})
}
