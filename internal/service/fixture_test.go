package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// ---------- fakes ----------

type fakeFixtureLister struct {
	listByTournamentFn func(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error)
	listLockedFn       func(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeagueFixtureView, error)
}

func (f *fakeFixtureLister) ListByTournamentForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.UserFixtureView, error) {
	return f.listByTournamentFn(ctx, tournamentID, userID)
}

func (f *fakeFixtureLister) ListLockedByLeague(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
) ([]*domain.LeagueFixtureView, error) {
	return f.listLockedFn(ctx, leagueID, requestingUserID)
}

type fakeLeagueMemberChecker struct {
	isMemberFn func(ctx context.Context, leagueID, userID uuid.UUID) (bool, error)
}

func (f *fakeLeagueMemberChecker) IsMember(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) (bool, error) {
	return f.isMemberFn(ctx, leagueID, userID)
}

// ---------- tests ----------

func TestFixtureService_ListForLeague(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	userID := uuid.New()

	t.Run("returns ErrForbidden when user is not a member", func(t *testing.T) {
		t.Parallel()
		svc := service.NewFixtureService(
			&fakeFixtureLister{},
			&fakeLeagueMemberChecker{
				isMemberFn: func(_ context.Context, _, _ uuid.UUID) (bool, error) {
					return false, nil
				},
			},
		)

		_, err := svc.ListForLeague(context.Background(), leagueID, userID)
		require.ErrorIs(t, err, domain.ErrForbidden)
	})

	t.Run("delegates to fixture lister when user is a member", func(t *testing.T) {
		t.Parallel()
		fixture := domain.Fixture{
			ID:        uuid.New(),
			KickoffAt: time.Now().Add(time.Hour),
			Status:    domain.FixtureStatusInProgress,
		}
		views := []*domain.LeagueFixtureView{{Fixture: fixture}}

		svc := service.NewFixtureService(
			&fakeFixtureLister{
				listLockedFn: func(_ context.Context, gLeagueID, gUserID uuid.UUID) ([]*domain.LeagueFixtureView, error) {
					require.Equal(t, leagueID, gLeagueID)
					require.Equal(t, userID, gUserID)
					return views, nil
				},
			},
			&fakeLeagueMemberChecker{
				isMemberFn: func(_ context.Context, _, _ uuid.UUID) (bool, error) {
					return true, nil
				},
			},
		)

		got, err := svc.ListForLeague(context.Background(), leagueID, userID)
		require.NoError(t, err)
		require.Equal(t, views, got)
	})
}
