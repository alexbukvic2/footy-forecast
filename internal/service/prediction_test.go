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

type fakePredictionRepo struct {
	upsertFn func(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
}

func (f *fakePredictionRepo) Upsert(
	ctx context.Context,
	in domain.UpsertScorePredictionInput,
) (*domain.ScorePrediction, error) {
	return f.upsertFn(ctx, in)
}

type fakeFixtureGetter struct {
	getByIDFn func(ctx context.Context, id uuid.UUID) (*domain.Fixture, error)
}

func (f *fakeFixtureGetter) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Fixture, error) {
	return f.getByIDFn(ctx, id)
}

type frozenClock struct{ t time.Time }

func (c frozenClock) Now() time.Time { return c.t }

// ---------- tests ----------

func TestPredictionService_UpsertScore(t *testing.T) {
	t.Parallel()

	t.Run("returns ErrInvalid for negative goals_home", func(t *testing.T) {
		t.Parallel()
		svc := service.NewPredictionService(&fakePredictionRepo{}, &fakeFixtureGetter{}, frozenClock{t: time.Now()})
		_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
			UserID:    uuid.New(),
			FixtureID: uuid.New(),
			GoalsHome: -1,
			GoalsAway: 0,
		})
		require.ErrorIs(t, err, domain.ErrInvalid)
	})

	t.Run("returns ErrInvalid for negative goals_away", func(t *testing.T) {
		t.Parallel()
		svc := service.NewPredictionService(&fakePredictionRepo{}, &fakeFixtureGetter{}, frozenClock{t: time.Now()})
		_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
			UserID:    uuid.New(),
			FixtureID: uuid.New(),
			GoalsHome: 0,
			GoalsAway: -1,
		})
		require.ErrorIs(t, err, domain.ErrInvalid)
	})

	t.Run("returns ErrNotFound when fixture does not exist", func(t *testing.T) {
		t.Parallel()
		svc := service.NewPredictionService(
			&fakePredictionRepo{},
			&fakeFixtureGetter{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Fixture, error) {
					return nil, domain.ErrNotFound
				},
			},
			frozenClock{t: time.Now()},
		)
		_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
			UserID:    uuid.New(),
			FixtureID: uuid.New(),
			GoalsHome: 1,
			GoalsAway: 0,
		})
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("returns ErrForbidden when fixture has already kicked off", func(t *testing.T) {
		t.Parallel()
		kickoff := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
		now := kickoff.Add(time.Minute) // one minute after kickoff
		fixture := &domain.Fixture{ID: uuid.New(), KickoffAt: kickoff}
		svc := service.NewPredictionService(
			&fakePredictionRepo{},
			&fakeFixtureGetter{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Fixture, error) {
					return fixture, nil
				},
			},
			frozenClock{t: now},
		)
		_, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
			UserID:    uuid.New(),
			FixtureID: fixture.ID,
			GoalsHome: 2,
			GoalsAway: 1,
		})
		require.ErrorIs(t, err, domain.ErrForbidden)
	})

	t.Run("delegates to repo for valid future fixture", func(t *testing.T) {
		t.Parallel()
		kickoff := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)
		now := kickoff.Add(-time.Hour) // one hour before kickoff
		fixture := &domain.Fixture{ID: uuid.New(), KickoffAt: kickoff}
		expected := &domain.ScorePrediction{ID: uuid.New(), GoalsHome: 2, GoalsAway: 1}

		svc := service.NewPredictionService(
			&fakePredictionRepo{
				upsertFn: func(_ context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error) {
					require.Equal(t, 2, in.GoalsHome)
					require.Equal(t, 1, in.GoalsAway)
					return expected, nil
				},
			},
			&fakeFixtureGetter{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Fixture, error) {
					return fixture, nil
				},
			},
			frozenClock{t: now},
		)

		got, err := svc.UpsertScore(context.Background(), domain.UpsertScorePredictionInput{
			UserID:    uuid.New(),
			FixtureID: fixture.ID,
			GoalsHome: 2,
			GoalsAway: 1,
		})
		require.NoError(t, err)
		require.Equal(t, expected, got)
	})
}
