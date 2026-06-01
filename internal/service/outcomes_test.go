package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// ---------- fakes ----------

type fakeOutcomesRepo struct {
	listPlayersFn func(context.Context, uuid.UUID) ([]*domain.PlayerOutcome, error)
	listTeamsFn   func(context.Context, uuid.UUID) ([]*domain.TeamOutcome, error)
}

func (f *fakeOutcomesRepo) ListPlayerOutcomes(ctx context.Context, id uuid.UUID) ([]*domain.PlayerOutcome, error) {
	if f.listPlayersFn != nil {
		return f.listPlayersFn(ctx, id)
	}
	return nil, nil
}

func (f *fakeOutcomesRepo) ListTeamOutcomes(ctx context.Context, id uuid.UUID) ([]*domain.TeamOutcome, error) {
	if f.listTeamsFn != nil {
		return f.listTeamsFn(ctx, id)
	}
	return nil, nil
}

type fakeOutcomesTournamentGetter struct {
	getByIDFn func(context.Context, uuid.UUID) (*domain.Tournament, error)
}

func (f *fakeOutcomesTournamentGetter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tournament, error) {
	return f.getByIDFn(ctx, id)
}

// ---------- tests ----------

func TestOutcomesService_ListByTournament(t *testing.T) {
	ctx := context.Background()
	tournamentID := uuid.Must(uuid.NewV7())
	validTournament := &domain.Tournament{ID: tournamentID}

	playerID := uuid.Must(uuid.NewV7())
	teamID := uuid.Must(uuid.NewV7())

	tests := []struct {
		name           string
		getByIDErr     error
		playerOutcome  *domain.PlayerOutcome
		teamOutcome    *domain.TeamOutcome
		listPlayersErr error
		listTeamsErr   error
		wantErr        error
		wantPlayers    int
		wantTeams      int
	}{
		{
			name:       "tournament not found",
			getByIDErr: domain.ErrNotFound,
			wantErr:    domain.ErrNotFound,
		},
		{
			name:           "player repo error",
			listPlayersErr: errors.New("db failure"),
		},
		{
			name:         "team repo error",
			listTeamsErr: errors.New("db failure"),
		},
		{
			name:        "empty outcomes",
			wantPlayers: 0,
			wantTeams:   0,
		},
		{
			name: "returns outcomes",
			playerOutcome: &domain.PlayerOutcome{
				PlayerID:   playerID,
				PlayerName: "Messi",
				TeamID:     teamID,
				TeamName:   "Argentina",
				Category:   domain.PlayerHandicapCategoryTotalTopScorer,
			},
			teamOutcome: &domain.TeamOutcome{
				TeamID:   teamID,
				TeamName: "Argentina",
				Category: domain.TeamHandicapCategoryWinner,
			},
			wantPlayers: 1,
			wantTeams:   1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tGetter := &fakeOutcomesTournamentGetter{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Tournament, error) {
					if tc.getByIDErr != nil {
						return nil, tc.getByIDErr
					}
					return validTournament, nil
				},
			}

			repo := &fakeOutcomesRepo{
				listPlayersFn: func(_ context.Context, _ uuid.UUID) ([]*domain.PlayerOutcome, error) {
					if tc.listPlayersErr != nil {
						return nil, tc.listPlayersErr
					}
					if tc.playerOutcome != nil {
						return []*domain.PlayerOutcome{tc.playerOutcome}, nil
					}
					return []*domain.PlayerOutcome{}, nil
				},
				listTeamsFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamOutcome, error) {
					if tc.listTeamsErr != nil {
						return nil, tc.listTeamsErr
					}
					if tc.teamOutcome != nil {
						return []*domain.TeamOutcome{tc.teamOutcome}, nil
					}
					return []*domain.TeamOutcome{}, nil
				},
			}

			svc := service.NewOutcomesService(repo, tGetter)
			out, err := svc.ListByTournament(ctx, tournamentID)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			if tc.listPlayersErr != nil || tc.listTeamsErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, out.PlayerOutcomes, tc.wantPlayers)
			require.Len(t, out.TeamOutcomes, tc.wantTeams)
		})
	}
}
