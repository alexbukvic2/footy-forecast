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

// repoWithHandicaps returns a fakePlayerRepo that returns one player whose Handicaps
// map is pre-populated with the provided values (entries omitted for absent categories).
func repoWithHandicaps(playerID uuid.UUID, h map[domain.PlayerHandicapCategory]int) *fakePlayerRepo {
	return &fakePlayerRepo{
		searchFn: func(context.Context, uuid.UUID, string, string) ([]*domain.PlayerSearchResult, error) {
			return []*domain.PlayerSearchResult{
				{ID: playerID, Name: "Player", TeamName: "Team", TeamLogo: "", Handicaps: h},
			}, nil
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

// ---------- handicap enrichment ----------

func TestPlayerService_Search_AllCategoriesAlwaysPresent(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	// repo returns a player with no handicap entries — service must fill all categories
	repo := repoWithHandicaps(uuid.New(), map[domain.PlayerHandicapCategory]int{})
	svc := playerSvc(repo, tournamentExists(tournamentID))

	got, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: tournamentID, Query: "x"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Len(t, got[0].Handicaps, len(domain.AllPlayerHandicapCategories))
	for _, cat := range domain.AllPlayerHandicapCategories {
		_, ok := got[0].Handicaps[cat]
		require.True(t, ok, "category %q must be present", cat)
	}
}

func TestPlayerService_Search_UsesDefaultWhenRepoReturnsNoHandicap(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	repo := repoWithHandicaps(uuid.New(), map[domain.PlayerHandicapCategory]int{})
	svc := playerSvc(repo, tournamentExists(tournamentID))

	got, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: tournamentID, Query: "x"})
	require.NoError(t, err)
	for _, cat := range domain.AllPlayerHandicapCategories {
		require.Equal(t, 20, got[0].Handicaps[cat], "category %q: expected default 20", cat)
	}
}

func TestPlayerService_Search_UsesRepoPointsWhenHandicapRowExists(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	repo := repoWithHandicaps(uuid.New(), map[domain.PlayerHandicapCategory]int{
		domain.PlayerHandicapCategoryGroupTopScorer: 7,
		domain.PlayerHandicapCategoryTotalTopScorer: 15,
	})
	svc := playerSvc(repo, tournamentExists(tournamentID))

	got, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: tournamentID, Query: "x"})
	require.NoError(t, err)
	require.Equal(t, 7, got[0].Handicaps[domain.PlayerHandicapCategoryGroupTopScorer])
	require.Equal(t, 15, got[0].Handicaps[domain.PlayerHandicapCategoryTotalTopScorer])
}

func TestPlayerService_Search_MixedHandicaps_DefaultAppliedOnlyForMissingCategories(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	// repo returns one category present, one absent
	repo := repoWithHandicaps(uuid.New(), map[domain.PlayerHandicapCategory]int{
		domain.PlayerHandicapCategoryGroupTopScorer: 3,
	})
	svc := playerSvc(repo, tournamentExists(tournamentID))

	got, err := svc.Search(context.Background(), domain.SearchPlayersInput{TournamentID: tournamentID, Query: "x"})
	require.NoError(t, err)
	require.Equal(t, 3, got[0].Handicaps[domain.PlayerHandicapCategoryGroupTopScorer])
	require.Equal(t, 20, got[0].Handicaps[domain.PlayerHandicapCategoryTotalTopScorer])
}
