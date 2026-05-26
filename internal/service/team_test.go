package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// ---------- fake ----------

type fakeTeamRepo struct {
	listFn func(context.Context, uuid.UUID) ([]domain.TeamWithHandicaps, error)
}

func (f *fakeTeamRepo) ListWithHandicapsByTournament(ctx context.Context, tid uuid.UUID) ([]domain.TeamWithHandicaps, error) {
	return f.listFn(ctx, tid)
}

// ---------- tests ----------

func TestTeamService_ListWithHandicaps_HappyPath(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	want := []domain.TeamWithHandicaps{
		{
			Team: domain.Team{ID: uuid.New(), Name: "France", TournamentID: tournamentID},
			Handicaps: []domain.TeamHandicapItem{
				{Category: domain.TeamHandicapCategoryWinner, Points: 20},
			},
		},
	}

	svc := service.NewTeamService(&fakeTeamRepo{
		listFn: func(_ context.Context, tid uuid.UUID) ([]domain.TeamWithHandicaps, error) {
			require.Equal(t, tournamentID, tid)
			return want, nil
		},
	})

	got, err := svc.ListWithHandicaps(context.Background(), tournamentID)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestTeamService_ListWithHandicaps_RepoErrPropagates(t *testing.T) {
	t.Parallel()

	svc := service.NewTeamService(&fakeTeamRepo{
		listFn: func(context.Context, uuid.UUID) ([]domain.TeamWithHandicaps, error) {
			return nil, domain.ErrNotFound
		},
	})

	_, err := svc.ListWithHandicaps(context.Background(), uuid.New())
	require.ErrorIs(t, err, domain.ErrNotFound)
}
