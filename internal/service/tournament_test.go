package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// fakeRepo is a hand-written fake of service.TournamentRepo.
// Each field lets a test customize behavior for one method.
type fakeRepo struct {
	createFn func(
		context.Context,
		repository.CreateTournamentParams,
	) (*domain.Tournament, error)
	getByIDFn func(
		context.Context,
		uuid.UUID,
	) (*domain.Tournament, error)
	listFn func(context.Context) ([]*domain.Tournament, error)
}

func (f *fakeRepo) Create(
	ctx context.Context,
	p repository.CreateTournamentParams,
) (*domain.Tournament, error) {
	return f.createFn(ctx, p)
}
func (f *fakeRepo) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Tournament, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeRepo) List(ctx context.Context) ([]*domain.Tournament, error) {
	return f.listFn(ctx)
}

// validInput returns a valid CreateTournamentInput for use as a baseline.
// Tests can mutate one field at a time to test validation.
func validInput() domain.CreateTournamentInput {
	return domain.CreateTournamentInput{
		Slug:     "world-cup-2026",
		Name:     "FIFA World Cup 2026",
		StartsAt: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC),
	}
}

func TestTournamentService_Create_Validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		mutate func(*domain.CreateTournamentInput)
	}{
		{"empty slug", func(in *domain.CreateTournamentInput) { in.Slug = "" }},
		{"too short slug", func(in *domain.CreateTournamentInput) { in.Slug = "ab" }},
		{"slug with spaces", func(in *domain.CreateTournamentInput) { in.Slug = "world cup" }},
		{"slug with double hyphens", func(in *domain.CreateTournamentInput) { in.Slug = "world--cup" }},
		{"slug with leading hyphen", func(in *domain.CreateTournamentInput) { in.Slug = "-cup" }},
		{"empty name", func(in *domain.CreateTournamentInput) { in.Name = "" }},
		{"name too long", func(in *domain.CreateTournamentInput) { in.Name = string(make([]byte, 201)) }},
		{"zero starts_at", func(in *domain.CreateTournamentInput) { in.StartsAt = time.Time{} }},
		{"zero ends_at", func(in *domain.CreateTournamentInput) { in.EndsAt = time.Time{} }},
		{"ends_at equals starts_at", func(in *domain.CreateTournamentInput) { in.EndsAt = in.StartsAt }},
		{"ends_at before starts_at", func(in *domain.CreateTournamentInput) { in.EndsAt = in.StartsAt.Add(-time.Hour) }},
	}

	for _, tc := range cases {
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()

				repoCalled := false
				repo := &fakeRepo{
					createFn: func(
						context.Context,
						repository.CreateTournamentParams,
					) (*domain.Tournament, error) {
						repoCalled = true
						return nil, errors.New("repo should not be called for invalid input")
					},
				}
				svc := service.NewTournamentService(repo)

				in := validInput()
				tc.mutate(&in)

				_, err := svc.Create(context.Background(), in)

				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrInvalid)
				require.False(t, repoCalled, "repo should not be reached on invalid input")
			},
		)
	}
}

func TestTournamentService_Create_Success(t *testing.T) {
	t.Parallel()

	var receivedParams repository.CreateTournamentParams
	expected := &domain.Tournament{
		ID:     uuid.New(),
		Slug:   "world-cup-2026",
		Name:   "FIFA World Cup 2026",
		Status: domain.TournamentStatusUpcoming,
	}

	repo := &fakeRepo{
		createFn: func(
			_ context.Context,
			p repository.CreateTournamentParams,
		) (*domain.Tournament, error) {
			receivedParams = p
			return expected, nil
		},
	}
	svc := service.NewTournamentService(repo)

	in := validInput()
	got, err := svc.Create(context.Background(), in)

	require.NoError(t, err)
	require.Equal(t, expected, got)

	// Verify the service passed normalized params through.
	require.Equal(t, "world-cup-2026", receivedParams.Slug)
	require.Equal(t, "FIFA World Cup 2026", receivedParams.Name)
	// Dates should be UTC.
	require.Equal(t, time.UTC, receivedParams.StartsAt.Location())
	require.Equal(t, time.UTC, receivedParams.EndsAt.Location())
}

func TestTournamentService_Create_NormalizesInput(t *testing.T) {
	t.Parallel()

	var receivedParams repository.CreateTournamentParams
	repo := &fakeRepo{
		createFn: func(
			_ context.Context,
			p repository.CreateTournamentParams,
		) (*domain.Tournament, error) {
			receivedParams = p
			return &domain.Tournament{}, nil
		},
	}
	svc := service.NewTournamentService(repo)

	_, err := svc.Create(
		context.Background(), domain.CreateTournamentInput{
			Slug:     "  WORLD-CUP-2026  ",
			Name:     "  FIFA World Cup 2026  ",
			StartsAt: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
			EndsAt:   time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC),
		},
	)

	require.NoError(t, err)
	require.Equal(t, "world-cup-2026", receivedParams.Slug, "slug should be trimmed and lowercased")
	require.Equal(t, "FIFA World Cup 2026", receivedParams.Name, "name should be trimmed but case preserved")
}

func TestTournamentService_Create_PropagatesConflict(t *testing.T) {
	t.Parallel()

	repo := &fakeRepo{
		createFn: func(
			context.Context,
			repository.CreateTournamentParams,
		) (*domain.Tournament, error) {
			return nil, domain.ErrConflict
		},
	}
	svc := service.NewTournamentService(repo)

	_, err := svc.Create(context.Background(), validInput())
	require.ErrorIs(t, err, domain.ErrConflict)
}
