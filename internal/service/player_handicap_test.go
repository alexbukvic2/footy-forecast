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

type fakePlayerHandicapRepo struct {
	getFn func(context.Context, uuid.UUID, domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error)
}

func (f *fakePlayerHandicapRepo) Get(ctx context.Context, playerID uuid.UUID, category domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
	return f.getFn(ctx, playerID, category)
}

// ---------- tests ----------

func TestPlayerHandicapService_Get_HappyPath(t *testing.T) {
	t.Parallel()

	playerID := uuid.New()
	want := &domain.PlayerHandicap{
		ID:       uuid.New(),
		PlayerID: playerID,
		Category: domain.PlayerHandicapCategoryGroupTopScorer,
		Points:   5,
	}
	repo := &fakePlayerHandicapRepo{
		getFn: func(_ context.Context, id uuid.UUID, cat domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
			require.Equal(t, playerID, id)
			require.Equal(t, domain.PlayerHandicapCategoryGroupTopScorer, cat)
			return want, nil
		},
	}
	svc := service.NewPlayerHandicapService(repo)

	got, err := svc.Get(context.Background(), playerID, domain.PlayerHandicapCategoryGroupTopScorer)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestPlayerHandicapService_Get_RepoErrNotFoundPropagates(t *testing.T) {
	t.Parallel()

	repo := &fakePlayerHandicapRepo{
		getFn: func(context.Context, uuid.UUID, domain.PlayerHandicapCategory) (*domain.PlayerHandicap, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := service.NewPlayerHandicapService(repo)

	_, err := svc.Get(context.Background(), uuid.New(), domain.PlayerHandicapCategoryTotalTopScorer)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
