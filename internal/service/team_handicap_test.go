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

type fakeTeamHandicapRepo struct {
	getFn func(context.Context, uuid.UUID, domain.TeamHandicapCategory) (*domain.TeamHandicap, error)
}

func (f *fakeTeamHandicapRepo) Get(ctx context.Context, teamID uuid.UUID, category domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
	return f.getFn(ctx, teamID, category)
}

// ---------- tests ----------

func TestTeamHandicapService_Get_HappyPath(t *testing.T) {
	t.Parallel()

	teamID := uuid.New()
	want := &domain.TeamHandicap{
		ID:       uuid.New(),
		TeamID:   teamID,
		Category: domain.TeamHandicapCategoryWinner,
		Points:   20,
	}
	repo := &fakeTeamHandicapRepo{
		getFn: func(_ context.Context, id uuid.UUID, cat domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
			require.Equal(t, teamID, id)
			require.Equal(t, domain.TeamHandicapCategoryWinner, cat)
			return want, nil
		},
	}
	svc := service.NewTeamHandicapService(repo)

	got, err := svc.Get(context.Background(), teamID, domain.TeamHandicapCategoryWinner)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestTeamHandicapService_Get_RepoErrNotFoundPropagates(t *testing.T) {
	t.Parallel()

	repo := &fakeTeamHandicapRepo{
		getFn: func(context.Context, uuid.UUID, domain.TeamHandicapCategory) (*domain.TeamHandicap, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := service.NewTeamHandicapService(repo)

	_, err := svc.Get(context.Background(), uuid.New(), domain.TeamHandicapCategoryPlayoff)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
