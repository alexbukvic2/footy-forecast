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

type fakeTournamentGroupTableRepo struct {
	listFn func(context.Context, uuid.UUID) ([]*domain.TournamentGroupEntry, error)
}

func (f *fakeTournamentGroupTableRepo) ListByTournament(ctx context.Context, id uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
	return f.listFn(ctx, id)
}

type fakeGroupTableTournamentGetter struct {
	getByIDFn func(context.Context, uuid.UUID) (*domain.Tournament, error)
}

func (f *fakeGroupTableTournamentGetter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tournament, error) {
	return f.getByIDFn(ctx, id)
}

// ---------- tests ----------

func TestTournamentGroupTableService_ListByTournament(t *testing.T) {
	ctx := context.Background()
	tournamentID := uuid.Must(uuid.NewV7())
	validTournament := &domain.Tournament{ID: tournamentID}

	tests := []struct {
		name        string
		getByIDErr  error
		listEntries []*domain.TournamentGroupEntry
		listErr     error
		wantErr     error
		wantLen     int
	}{
		{
			name:       "tournament not found",
			getByIDErr: domain.ErrNotFound,
			wantErr:    domain.ErrNotFound,
		},
		{
			name:        "repo error",
			listEntries: nil,
			listErr:     errors.New("db failure"),
			wantErr:     nil, // generic error, not domain error
		},
		{
			name:        "empty table returns empty slice",
			listEntries: []*domain.TournamentGroupEntry{},
			wantLen:     0,
		},
		{
			name: "returns rows",
			listEntries: []*domain.TournamentGroupEntry{
				{TournamentID: tournamentID, TeamName: "Argentina", GroupLetter: "A", Position: 1, Points: 9},
			},
			wantLen: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tGetter := &fakeGroupTableTournamentGetter{
				getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Tournament, error) {
					if tc.getByIDErr != nil {
						return nil, tc.getByIDErr
					}
					return validTournament, nil
				},
			}
			repo := &fakeTournamentGroupTableRepo{
				listFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
					return tc.listEntries, tc.listErr
				},
			}

			svc := service.NewTournamentGroupTableService(repo, tGetter)
			entries, err := svc.ListByTournament(ctx, tournamentID)

			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
				return
			}
			if tc.listErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, entries, tc.wantLen)
		})
	}
}
