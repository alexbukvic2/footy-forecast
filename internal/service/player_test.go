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

// ---------- fake ----------

type fakePlayerRepo struct {
	searchFn func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error)
}

func (f *fakePlayerRepo) Search(ctx context.Context, tournamentID uuid.UUID, escapedQuery, rawQuery string) ([]*domain.PlayerSearchResult, error) {
	return f.searchFn(ctx, tournamentID, escapedQuery, rawQuery)
}

// neverCallGetter returns a fakeTournamentGetter whose GetByID calls t.Fatal —
// used by validation tests to assert the service short-circuits before the DB.
func neverCallGetter(t *testing.T) *fakeTournamentGetter {
	t.Helper()
	return &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			t.Fatal("tournament getter must not be called when input is invalid")
			return nil, nil
		},
	}
}

// ---------- helpers ----------

func playerSvc(repo *fakePlayerRepo, getter *fakeTournamentGetter) *service.PlayerService {
	return service.NewPlayerService(repo, getter)
}

func tournamentExists(id uuid.UUID) *fakeTournamentGetter {
	return &fakeTournamentGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Tournament, error) {
			return upcomingTournament(id), nil
		},
	}
}

// ---------- validation ----------

func TestPlayerService_Search_EmptyQ(t *testing.T) {
	t.Parallel()

	called := false
	repo := &fakePlayerRepo{searchFn: func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error) {
		called = true
		return nil, nil
	}}
	svc := playerSvc(repo, neverCallGetter(t))

	_, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: uuid.New(), Query: ""})
	require.ErrorIs(t, err, domain.ErrInvalid)
	require.False(t, called)
}

func TestPlayerService_Search_WhitespaceOnlyQ(t *testing.T) {
	t.Parallel()

	called := false
	repo := &fakePlayerRepo{searchFn: func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error) {
		called = true
		return nil, nil
	}}
	svc := playerSvc(repo, neverCallGetter(t))

	_, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: uuid.New(), Query: "   "})
	require.ErrorIs(t, err, domain.ErrInvalid)
	require.False(t, called)
}

func TestPlayerService_Search_QTooLong(t *testing.T) {
	t.Parallel()

	called := false
	repo := &fakePlayerRepo{searchFn: func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error) {
		called = true
		return nil, nil
	}}
	svc := playerSvc(repo, neverCallGetter(t))

	_, err := svc.Search(context.Background(), domain.SearchPlayersInput{
		TournamentID: uuid.New(),
		Query:        string(make([]rune, 101)),
	})
	require.ErrorIs(t, err, domain.ErrInvalid)
	require.False(t, called)
}

// ---------- tournament checks ----------

func TestPlayerService_Search_TournamentNotFound(t *testing.T) {
	t.Parallel()

	repo := &fakePlayerRepo{}
	getter := &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			return nil, domain.ErrNotFound
		},
	}
	svc := playerSvc(repo, getter)

	_, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: uuid.New(), Query: "messi"})
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPlayerService_Search_TournamentGetterErrorPropagates(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("db down")
	repo := &fakePlayerRepo{}
	getter := &fakeTournamentGetter{
		getByIDFn: func(context.Context, uuid.UUID) (*domain.Tournament, error) {
			return nil, sentinel
		},
	}
	svc := playerSvc(repo, getter)

	_, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: uuid.New(), Query: "messi"})
	require.ErrorIs(t, err, sentinel)
}

// ---------- happy path ----------

func TestPlayerService_Search_ForwardsEscapedAndRawQuery(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	var gotTournamentID uuid.UUID
	var gotEscaped, gotRaw string

	repo := &fakePlayerRepo{searchFn: func(_ context.Context, tid uuid.UUID, escaped, raw string) ([]*domain.PlayerSearchResult, error) {
		gotTournamentID = tid
		gotEscaped = escaped
		gotRaw = raw
		return []*domain.PlayerSearchResult{}, nil
	}}
	svc := playerSvc(repo, tournamentExists(tournamentID))

	_, err := svc.Search(context.Background(), domain.SearchPlayersInput{
		TournamentID: tournamentID,
		Query:        "Mbap%pe_",
	})
	require.NoError(t, err)
	require.Equal(t, tournamentID, gotTournamentID)
	require.Equal(t, `Mbap\%pe\_`, gotEscaped, "escaped query should have wildcards escaped")
	require.Equal(t, "Mbap%pe_", gotRaw, "raw query must be unescaped for similarity ranking")
}

func TestPlayerService_Search_EmptyRepoResultPassedThrough(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	repo := &fakePlayerRepo{searchFn: func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error) {
		return []*domain.PlayerSearchResult{}, nil
	}}
	svc := playerSvc(repo, tournamentExists(tournamentID))

	players, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: tournamentID, Query: "x"})
	require.NoError(t, err)
	require.NotNil(t, players)
	require.Empty(t, players)
}

func TestPlayerService_Search_ReturnsPlayers(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	want := []*domain.PlayerSearchResult{
		{ID: uuid.New(), Name: "Kylian Mbappé", TeamName: "France", TeamLogo: "<svg/>"},
		{ID: uuid.New(), Name: "Lionel Messi", TeamName: "Argentina", TeamLogo: "<svg/>"},
	}
	repo := &fakePlayerRepo{searchFn: func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error) {
		return want, nil
	}}
	svc := playerSvc(repo, tournamentExists(tournamentID))

	got, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: tournamentID, Query: "messi"})
	require.NoError(t, err)
	require.Equal(t, want, got)
}
