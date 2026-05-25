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

type fakePredictionUpserter struct {
	upsertFn func(context.Context, domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
}

func (f *fakePredictionUpserter) Upsert(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
	return f.upsertFn(ctx, in)
}

type fakeFixtureGetter struct {
	getByIDFn func(context.Context, uuid.UUID) (*domain.Fixture, error)
}

func (f *fakeFixtureGetter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Fixture, error) {
	return f.getByIDFn(ctx, id)
}

type fakeFixtureForUserLister struct {
	listByTournamentFn func(context.Context, uuid.UUID, uuid.UUID) ([]*domain.UserFixtureView, error)
	listLockedFn       func(context.Context, uuid.UUID, uuid.UUID) ([]*domain.LeagueFixtureView, error)
}

func (f *fakeFixtureForUserLister) ListByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error) {
	return f.listByTournamentFn(ctx, tournamentID, userID)
}

func (f *fakeFixtureForUserLister) ListLockedByLeague(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeagueFixtureView, error) {
	return f.listLockedFn(ctx, leagueID, requestingUserID)
}

type fakeLeagueMemberChecker struct {
	getMemberFn func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error)
}

func (f *fakeLeagueMemberChecker) GetMember(ctx context.Context, leagueID, userID uuid.UUID) (*domain.LeagueMember, error) {
	return f.getMemberFn(ctx, leagueID, userID)
}

type frozenClock struct{ t time.Time }

func (c frozenClock) Now() time.Time { return c.t }

// ---------- PredictionService ----------

func TestPredictionService_UpsertScore_NegativeGoalsHome(t *testing.T) {
	t.Parallel()
	svc := service.NewPredictionService(&fakePredictionUpserter{}, &fakeFixtureGetter{}, frozenClock{})
	_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
		GoalsHome: -1, GoalsAway: 0,
	})
	require.ErrorIs(t, err, domain.ErrInvalid)
}

func TestPredictionService_UpsertScore_NegativeGoalsAway(t *testing.T) {
	t.Parallel()
	svc := service.NewPredictionService(&fakePredictionUpserter{}, &fakeFixtureGetter{}, frozenClock{})
	_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
		GoalsHome: 0, GoalsAway: -1,
	})
	require.ErrorIs(t, err, domain.ErrInvalid)
}

func TestPredictionService_UpsertScore_FixtureNotFound(t *testing.T) {
	t.Parallel()
	fixtureID := uuid.New()
	svc := service.NewPredictionService(
		&fakePredictionUpserter{},
		&fakeFixtureGetter{getByIDFn: func(_ context.Context, id uuid.UUID) (*domain.Fixture, error) {
			require.Equal(t, fixtureID, id)
			return nil, domain.ErrNotFound
		}},
		frozenClock{},
	)
	_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
		FixtureID: fixtureID, GoalsHome: 1, GoalsAway: 0,
	})
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPredictionService_UpsertScore_WithinLockWindow(t *testing.T) {
	t.Parallel()
	kickoff := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)
	// 29 minutes before kickoff = inside the 30-minute lock window
	now := kickoff.Add(-29 * time.Minute)

	svc := service.NewPredictionService(
		&fakePredictionUpserter{},
		&fakeFixtureGetter{getByIDFn: func(context.Context, uuid.UUID) (*domain.Fixture, error) {
			return &domain.Fixture{KickoffAt: kickoff}, nil
		}},
		frozenClock{t: now},
	)
	_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{GoalsHome: 1, GoalsAway: 0})
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestPredictionService_UpsertScore_AtExactLockTime(t *testing.T) {
	t.Parallel()
	kickoff := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)
	// exactly 30 minutes before = locked
	now := kickoff.Add(-30 * time.Minute)

	svc := service.NewPredictionService(
		&fakePredictionUpserter{},
		&fakeFixtureGetter{getByIDFn: func(context.Context, uuid.UUID) (*domain.Fixture, error) {
			return &domain.Fixture{KickoffAt: kickoff}, nil
		}},
		frozenClock{t: now},
	)
	_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{GoalsHome: 1, GoalsAway: 0})
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestPredictionService_UpsertScore_HappyPath(t *testing.T) {
	t.Parallel()
	kickoff := time.Date(2026, 6, 1, 15, 0, 0, 0, time.UTC)
	// 31 minutes before = open
	now := kickoff.Add(-31 * time.Minute)

	want := &domain.ScorePrediction{ID: uuid.New(), GoalsHome: 2, GoalsAway: 1}
	svc := service.NewPredictionService(
		&fakePredictionUpserter{upsertFn: func(_ context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
			require.Equal(t, 2, in.GoalsHome)
			require.Equal(t, 1, in.GoalsAway)
			return want, nil
		}},
		&fakeFixtureGetter{getByIDFn: func(context.Context, uuid.UUID) (*domain.Fixture, error) {
			return &domain.Fixture{KickoffAt: kickoff}, nil
		}},
		frozenClock{t: now},
	)
	got, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{GoalsHome: 2, GoalsAway: 1})
	require.NoError(t, err)
	require.Equal(t, want, got)
}

// ---------- FixtureService ----------

func TestFixtureService_ListForLeague_NotMember(t *testing.T) {
	t.Parallel()
	svc := service.NewFixtureService(
		&fakeFixtureForUserLister{},
		&fakeLeagueMemberChecker{getMemberFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error) {
			return nil, domain.ErrNotFound
		}},
	)
	_, err := svc.ListForLeague(context.Background(), uuid.New(), uuid.New())
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestFixtureService_ListForLeague_HappyPath(t *testing.T) {
	t.Parallel()
	leagueID, userID := uuid.New(), uuid.New()
	want := []*domain.LeagueFixtureView{{Fixture: domain.Fixture{ID: uuid.New()}}}

	svc := service.NewFixtureService(
		&fakeFixtureForUserLister{listLockedFn: func(_ context.Context, lid, uid uuid.UUID) ([]*domain.LeagueFixtureView, error) {
			require.Equal(t, leagueID, lid)
			require.Equal(t, userID, uid)
			return want, nil
		}},
		&fakeLeagueMemberChecker{getMemberFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error) {
			return &domain.LeagueMember{}, nil
		}},
	)
	got, err := svc.ListForLeague(context.Background(), leagueID, userID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestFixtureService_ListForUser_HappyPath(t *testing.T) {
	t.Parallel()
	tournamentID, userID := uuid.New(), uuid.New()
	want := []*domain.UserFixtureView{{Fixture: domain.Fixture{ID: uuid.New()}}}

	svc := service.NewFixtureService(
		&fakeFixtureForUserLister{listByTournamentFn: func(_ context.Context, tid, uid uuid.UUID) ([]*domain.UserFixtureView, error) {
			require.Equal(t, tournamentID, tid)
			require.Equal(t, userID, uid)
			return want, nil
		}},
		&fakeLeagueMemberChecker{},
	)
	got, err := svc.ListForUser(context.Background(), tournamentID, userID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
