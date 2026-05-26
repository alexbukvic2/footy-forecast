package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/service"
)

// ---------- fakes ----------

type fakePlayerPredictionRepo struct {
	upsertFn                  func(context.Context, domain.UpsertPlayerPredictionInput) (*domain.PlayerPrediction, error)
	deleteFn                  func(context.Context, uuid.UUID, uuid.UUID, domain.PlayerHandicapCategory, *string) error
	listByTournamentForUserFn func(context.Context, uuid.UUID, uuid.UUID) ([]*domain.PlayerPrediction, error)
	listByLeagueFn            func(context.Context, uuid.UUID) ([]*domain.PlayerLeaguePick, error)
}

func (f *fakePlayerPredictionRepo) UpsertPlayer(ctx context.Context, in domain.UpsertPlayerPredictionInput) (*domain.PlayerPrediction, error) {
	return f.upsertFn(ctx, in)
}
func (f *fakePlayerPredictionRepo) DeletePlayer(ctx context.Context, userID, tournamentID uuid.UUID, category domain.PlayerHandicapCategory, groupLetter *string) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, userID, tournamentID, category, groupLetter)
	}
	return nil
}
func (f *fakePlayerPredictionRepo) ListPlayersByTournamentForUser(ctx context.Context, tID, uID uuid.UUID) ([]*domain.PlayerPrediction, error) {
	return f.listByTournamentForUserFn(ctx, tID, uID)
}
func (f *fakePlayerPredictionRepo) ListPlayersByLeague(ctx context.Context, lID uuid.UUID) ([]*domain.PlayerLeaguePick, error) {
	return f.listByLeagueFn(ctx, lID)
}

type fakeTeamPredictionRepo struct {
	upsertFn                  func(context.Context, domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error)
	deleteFn                  func(context.Context, uuid.UUID, uuid.UUID, domain.TeamHandicapCategory, *string, int) error
	listByTournamentForUserFn func(context.Context, uuid.UUID, uuid.UUID) ([]*domain.TeamPrediction, error)
	listByLeagueFn            func(context.Context, uuid.UUID) ([]*domain.TeamLeaguePick, error)
	countWildcardsFn          func(context.Context, uuid.UUID, uuid.UUID) (int, error)
}

func (f *fakeTeamPredictionRepo) UpsertTeam(ctx context.Context, in domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error) {
	return f.upsertFn(ctx, in)
}
func (f *fakeTeamPredictionRepo) DeleteTeam(ctx context.Context, userID, tournamentID uuid.UUID, category domain.TeamHandicapCategory, groupLetter *string, slotIndex int) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, userID, tournamentID, category, groupLetter, slotIndex)
	}
	return nil
}
func (f *fakeTeamPredictionRepo) ListTeamsByTournamentForUser(ctx context.Context, tID, uID uuid.UUID) ([]*domain.TeamPrediction, error) {
	return f.listByTournamentForUserFn(ctx, tID, uID)
}
func (f *fakeTeamPredictionRepo) ListTeamsByLeague(ctx context.Context, lID uuid.UUID) ([]*domain.TeamLeaguePick, error) {
	return f.listByLeagueFn(ctx, lID)
}
func (f *fakeTeamPredictionRepo) CountPlayoffWildcards(ctx context.Context, tournamentID, userID uuid.UUID) (int, error) {
	if f.countWildcardsFn != nil {
		return f.countWildcardsFn(ctx, tournamentID, userID)
	}
	return 0, nil
}

type fakePlayerGetter struct {
	getByIDFn func(context.Context, uuid.UUID) (*domain.Player, error)
}

func (f *fakePlayerGetter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Player, error) {
	return f.getByIDFn(ctx, id)
}

type fakeTeamGetter struct {
	getByIDFn func(context.Context, uuid.UUID) (*domain.Team, error)
}

func (f *fakeTeamGetter) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	return f.getByIDFn(ctx, id)
}

type fakeFixtureRepo struct {
	getFirstKickoffFn func(context.Context, uuid.UUID) (time.Time, error)
}

func (f *fakeFixtureRepo) GetFirstKickoffByTournament(ctx context.Context, tID uuid.UUID) (time.Time, error) {
	return f.getFirstKickoffFn(ctx, tID)
}

type fakeLeagueReader struct {
	getByIDFn                   func(context.Context, uuid.UUID) (*domain.League, error)
	getMemberFn                 func(context.Context, uuid.UUID, uuid.UUID) (*domain.LeagueMember, error)
	listMembersForPredictionsFn func(context.Context, uuid.UUID) ([]*domain.LeagueMemberDisplay, error)
}

func (f *fakeLeagueReader) GetByID(ctx context.Context, id uuid.UUID) (*domain.League, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeLeagueReader) GetMember(ctx context.Context, leagueID, userID uuid.UUID) (*domain.LeagueMember, error) {
	return f.getMemberFn(ctx, leagueID, userID)
}
func (f *fakeLeagueReader) ListMembersForPredictions(ctx context.Context, leagueID uuid.UUID) ([]*domain.LeagueMemberDisplay, error) {
	return f.listMembersForPredictionsFn(ctx, leagueID)
}

type fakeClock struct{ t time.Time }

func (c fakeClock) Now() time.Time { return c.t }

type fakeTeamGroupLister struct {
	listFn func(context.Context, uuid.UUID) ([]string, error)
}

func (f *fakeTeamGroupLister) ListGroupLettersByTournament(ctx context.Context, tID uuid.UUID) ([]string, error) {
	return f.listFn(ctx, tID)
}

func noGroups() *fakeTeamGroupLister {
	return &fakeTeamGroupLister{listFn: func(_ context.Context, _ uuid.UUID) ([]string, error) {
		return []string{}, nil
	}}
}

func groupsA() *fakeTeamGroupLister {
	return &fakeTeamGroupLister{listFn: func(_ context.Context, _ uuid.UUID) ([]string, error) {
		return []string{"A"}, nil
	}}
}

func strPtr(s string) *string { return &s }

// ---------- helpers ----------

func newSvc(
	ppRepo service.PlayerPredictionRepo,
	tpRepo service.TeamPredictionRepo,
	players service.PlayerGetter,
	teams service.TeamGetter,
	tGroups service.TeamGroupLister,
	fixtures service.FixtureFirstKickoffGetter,
	leagues service.LeagueReader,
	clock service.Clock,
) *service.TournamentPredictionService {
	return service.NewTournamentPredictionService(ppRepo, tpRepo, players, teams, tGroups, fixtures, leagues, clock)
}

func defaultPlayerRepo() *fakePlayerPredictionRepo {
	return &fakePlayerPredictionRepo{
		upsertFn: func(_ context.Context, in domain.UpsertPlayerPredictionInput) (*domain.PlayerPrediction, error) {
			return &domain.PlayerPrediction{ID: uuid.New(), UserID: in.UserID, TournamentID: in.TournamentID, Category: in.Category, Pick: in.Pick, PickName: "Messi"}, nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.PlayerPrediction, error) {
			return nil, nil
		},
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.PlayerLeaguePick, error) {
			return nil, nil
		},
	}
}

func defaultTeamRepo() *fakeTeamPredictionRepo {
	return &fakeTeamPredictionRepo{
		upsertFn: func(_ context.Context, in domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error) {
			return &domain.TeamPrediction{ID: uuid.New(), UserID: in.UserID, TournamentID: in.TournamentID, Category: in.Category, Pick: in.Pick, PickName: "Argentina"}, nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return nil, nil
		},
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return nil, nil
		},
	}
}

func noFixtures() *fakeFixtureRepo {
	return &fakeFixtureRepo{
		getFirstKickoffFn: func(_ context.Context, _ uuid.UUID) (time.Time, error) {
			return time.Time{}, domain.ErrNotFound
		},
	}
}

func kickoffAt(t time.Time) *fakeFixtureRepo {
	return &fakeFixtureRepo{
		getFirstKickoffFn: func(_ context.Context, _ uuid.UUID) (time.Time, error) {
			return t, nil
		},
	}
}

// ---------- UpsertPlayerPrediction tests ----------

func TestTournamentPredictionService_UpsertPlayer_WrongTournament(t *testing.T) {
	tournamentID := uuid.New()
	playerID := uuid.New()
	otherTournamentID := uuid.New()

	players := &fakePlayerGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Player, error) {
			return &domain.Player{ID: playerID, TournamentID: otherTournamentID}, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.UpsertPlayerPrediction(context.Background(), domain.UpsertPlayerPredictionInput{
		UserID:       uuid.New(),
		TournamentID: tournamentID,
		Category:     domain.PlayerHandicapCategoryGroupTopScorer,
		Pick:         playerID,
	})
	require.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestTournamentPredictionService_UpsertPlayer_Locked(t *testing.T) {
	tournamentID := uuid.New()
	playerID := uuid.New()

	players := &fakePlayerGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Player, error) {
			return &domain.Player{ID: playerID, TournamentID: tournamentID, GroupLetter: strPtr("A")}, nil
		},
	}
	// Lock time: 1 hour ago, so we are now 30+ min past it.
	firstKickoff := time.Now().Add(-90 * time.Minute)

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), kickoffAt(firstKickoff), nil, fakeClock{time.Now()})

	_, err := svc.UpsertPlayerPrediction(context.Background(), domain.UpsertPlayerPredictionInput{
		UserID:       uuid.New(),
		TournamentID: tournamentID,
		Category:     domain.PlayerHandicapCategoryGroupTopScorer,
		GroupLetter:  strPtr("A"),
		Pick:         playerID,
	})
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestTournamentPredictionService_UpsertPlayer_Happy(t *testing.T) {
	tournamentID := uuid.New()
	playerID := uuid.New()
	userID := uuid.New()

	players := &fakePlayerGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Player, error) {
			return &domain.Player{ID: playerID, TournamentID: tournamentID, GroupLetter: strPtr("A")}, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	pred, err := svc.UpsertPlayerPrediction(context.Background(), domain.UpsertPlayerPredictionInput{
		UserID:       userID,
		TournamentID: tournamentID,
		Category:     domain.PlayerHandicapCategoryGroupTopScorer,
		GroupLetter:  strPtr("A"),
		Pick:         playerID,
	})
	require.NoError(t, err)
	require.Equal(t, playerID, pred.Pick)
}

func TestTournamentPredictionService_UpsertPlayer_Happy_FixturesExistNotLocked(t *testing.T) {
	tournamentID := uuid.New()
	playerID := uuid.New()
	userID := uuid.New()

	players := &fakePlayerGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Player, error) {
			return &domain.Player{ID: playerID, TournamentID: tournamentID, GroupLetter: strPtr("A")}, nil
		},
	}
	// Kickoff is 2h in the future; lock is 1.5h in the future — not yet reached.
	fixtures := kickoffAt(time.Now().Add(2 * time.Hour))

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), fixtures, nil, fakeClock{time.Now()})

	pred, err := svc.UpsertPlayerPrediction(context.Background(), domain.UpsertPlayerPredictionInput{
		UserID:       userID,
		TournamentID: tournamentID,
		Category:     domain.PlayerHandicapCategoryGroupTopScorer,
		GroupLetter:  strPtr("A"),
		Pick:         playerID,
	})
	require.NoError(t, err)
	require.Equal(t, playerID, pred.Pick)
}

// ---------- UpsertTeamPrediction tests ----------

func TestTournamentPredictionService_UpsertTeam_WrongTournament(t *testing.T) {
	tournamentID := uuid.New()
	teamID := uuid.New()

	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: uuid.New()}, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, teams, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.UpsertTeamPrediction(context.Background(), domain.UpsertTeamPredictionInput{
		UserID:       uuid.New(),
		TournamentID: tournamentID,
		Category:     domain.TeamHandicapCategoryWinner,
		Pick:         teamID,
	})
	require.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestTournamentPredictionService_UpsertTeam_Locked(t *testing.T) {
	tournamentID := uuid.New()
	teamID := uuid.New()

	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: tournamentID}, nil
		},
	}
	firstKickoff := time.Now().Add(-90 * time.Minute)

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, teams, noGroups(), kickoffAt(firstKickoff), nil, fakeClock{time.Now()})

	_, err := svc.UpsertTeamPrediction(context.Background(), domain.UpsertTeamPredictionInput{
		UserID:       uuid.New(),
		TournamentID: tournamentID,
		Category:     domain.TeamHandicapCategoryWinner,
		Pick:         teamID,
	})
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

// ---------- BulkUpsertPlayerPredictions tests ----------

func TestTournamentPredictionService_BulkUpsertPlayers_Locked(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	playerID := uuid.New()

	players := &fakePlayerGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Player, error) {
			return &domain.Player{ID: playerID, TournamentID: tournamentID}, nil
		},
	}
	firstKickoff := time.Now().Add(-90 * time.Minute)
	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), kickoffAt(firstKickoff), nil, fakeClock{time.Now()})

	_, err := svc.BulkUpsertPlayerPredictions(context.Background(), tournamentID, userID, []domain.BulkPlayerPredictionItem{
		{Category: domain.PlayerHandicapCategoryTotalTopScorer, PlayerID: &playerID},
	})
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestTournamentPredictionService_BulkUpsertPlayers_EmptyBatch(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()

	ppRepo := &fakePlayerPredictionRepo{
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.PlayerPrediction, error) {
			return nil, nil
		},
	}
	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	views, err := svc.BulkUpsertPlayerPredictions(context.Background(), tournamentID, userID, nil)
	require.NoError(t, err)
	require.NotNil(t, views)
}

func TestTournamentPredictionService_BulkUpsertPlayers_ClearSlot(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()

	deleted := false
	ppRepo := &fakePlayerPredictionRepo{
		deleteFn: func(_ context.Context, _, _ uuid.UUID, _ domain.PlayerHandicapCategory, _ *string) error {
			deleted = true
			return nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.PlayerPrediction, error) {
			return nil, nil
		},
	}
	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.BulkUpsertPlayerPredictions(context.Background(), tournamentID, userID, []domain.BulkPlayerPredictionItem{
		{Category: domain.PlayerHandicapCategoryTotalTopScorer, PlayerID: nil},
	})
	require.NoError(t, err)
	require.True(t, deleted)
}

func TestTournamentPredictionService_BulkUpsertPlayers_PlayerWrongTournament(t *testing.T) {
	tournamentID := uuid.New()
	playerID := uuid.New()

	players := &fakePlayerGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Player, error) {
			return &domain.Player{ID: playerID, TournamentID: uuid.New()}, nil
		},
	}
	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.BulkUpsertPlayerPredictions(context.Background(), tournamentID, uuid.New(), []domain.BulkPlayerPredictionItem{
		{Category: domain.PlayerHandicapCategoryTotalTopScorer, PlayerID: &playerID},
	})
	require.True(t, errors.Is(err, domain.ErrNotFound))
}

// ---------- BulkUpsertTeamPredictions tests ----------

func TestTournamentPredictionService_BulkUpsertTeams_Locked(t *testing.T) {
	tournamentID := uuid.New()
	teamID := uuid.New()

	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: tournamentID}, nil
		},
	}
	firstKickoff := time.Now().Add(-90 * time.Minute)
	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, teams, noGroups(), kickoffAt(firstKickoff), nil, fakeClock{time.Now()})

	_, err := svc.BulkUpsertTeamPredictions(context.Background(), tournamentID, uuid.New(), []domain.BulkTeamPredictionItem{
		{Category: domain.TeamHandicapCategoryWinner, SlotIndex: 0, TeamID: &teamID},
	})
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestTournamentPredictionService_BulkUpsertTeams_WildcardCapExceeded(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()

	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			groupA := "A"
			return &domain.Team{ID: teamID, TournamentID: tournamentID, GroupLetter: &groupA}, nil
		},
	}
	tpRepo := &fakeTeamPredictionRepo{
		countWildcardsFn: func(_ context.Context, _, _ uuid.UUID) (int, error) { return 8, nil },
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return nil, nil
		},
	}
	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, groupsA(), noFixtures(), nil, fakeClock{time.Now()})

	groupA := "A"
	_, err := svc.BulkUpsertTeamPredictions(context.Background(), tournamentID, userID, []domain.BulkTeamPredictionItem{
		{Category: domain.TeamHandicapCategoryPlayoff, GroupLetter: &groupA, SlotIndex: 2, TeamID: &teamID},
	})
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestTournamentPredictionService_BulkUpsertTeams_ClearSlot(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()

	deleted := false
	tpRepo := &fakeTeamPredictionRepo{
		deleteFn: func(_ context.Context, _, _ uuid.UUID, _ domain.TeamHandicapCategory, _ *string, _ int) error {
			deleted = true
			return nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return nil, nil
		},
	}
	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.BulkUpsertTeamPredictions(context.Background(), tournamentID, userID, []domain.BulkTeamPredictionItem{
		{Category: domain.TeamHandicapCategoryWinner, SlotIndex: 0, TeamID: nil},
	})
	require.NoError(t, err)
	require.True(t, deleted)
}

// ---------- ListPlayerPredictionsForUser tests ----------

func TestTournamentPredictionService_ListPlayerPredictionsForUser_AllCategories(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	playerID := uuid.New()

	ppRepo := &fakePlayerPredictionRepo{
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.PlayerPrediction, error) {
			return []*domain.PlayerPrediction{
				{Category: domain.PlayerHandicapCategoryGroupTopScorer, GroupLetter: strPtr("A"), Pick: playerID, PickName: "Messi"},
			}, nil
		},
	}

	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, groupsA(), noFixtures(), nil, fakeClock{time.Now()})

	_, views, err := svc.ListPlayerPredictionsForUser(context.Background(), tournamentID, userID)
	require.NoError(t, err)
	require.Len(t, views, len(domain.AllPlayerHandicapCategories))

	found := false
	for _, v := range views {
		if v.Category == domain.PlayerHandicapCategoryGroupTopScorer {
			require.NotNil(t, v.Prediction)
			found = true
		} else {
			require.Nil(t, v.Prediction, "unpredicted category %s should be nil", v.Category)
		}
	}
	require.True(t, found)
}

// ---------- ListTeamPredictionsForUser tests ----------

func TestTournamentPredictionService_ListTeamPredictionsForUser_AllCategories(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()

	tpRepo := &fakeTeamPredictionRepo{
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return nil, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, &fakeTeamGetter{}, groupsA(), noFixtures(), nil, fakeClock{time.Now()})

	_, views, err := svc.ListTeamPredictionsForUser(context.Background(), tournamentID, userID)
	require.NoError(t, err)
	require.Len(t, views, 9) // 1 group "A": 1 group_winner + 3 playoff + 4 semifinalist + 1 winner
	for _, v := range views {
		require.Nil(t, v.Prediction, "no predictions yet, all should be nil")
	}
}

// ---------- ListLeaguePlayerPredictions tests ----------

func TestTournamentPredictionService_ListLeaguePlayerPredictions_NotMember(t *testing.T) {
	leagueID := uuid.New()
	tournamentID := uuid.New()
	userID := uuid.New()

	leagues := &fakeLeagueReader{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, TournamentID: tournamentID}, nil
		},
		getMemberFn: func(_ context.Context, _, _ uuid.UUID) (*domain.LeagueMember, error) {
			return nil, domain.ErrNotFound
		},
	}

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), noFixtures(), leagues, fakeClock{time.Now()})

	_, err := svc.ListLeaguePlayerPredictions(context.Background(), leagueID, userID)
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestTournamentPredictionService_ListLeaguePlayerPredictions_BeforeLock(t *testing.T) {
	leagueID := uuid.New()
	tournamentID := uuid.New()
	userID := uuid.New()

	leagues := &fakeLeagueReader{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, TournamentID: tournamentID}, nil
		},
		getMemberFn: func(_ context.Context, _, _ uuid.UUID) (*domain.LeagueMember, error) {
			return &domain.LeagueMember{}, nil
		},
	}

	// Kickoff is 2 hours in the future, lock is 1.5h in the future — not yet locked.
	fixtures := kickoffAt(time.Now().Add(2 * time.Hour))

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), fixtures, leagues, fakeClock{time.Now()})

	_, err := svc.ListLeaguePlayerPredictions(context.Background(), leagueID, userID)
	require.True(t, errors.Is(err, domain.ErrForbidden))
}

func TestTournamentPredictionService_ListLeaguePlayerPredictions_AfterLock(t *testing.T) {
	leagueID := uuid.New()
	tournamentID := uuid.New()
	me := uuid.New()
	alice := uuid.New()

	leagues := &fakeLeagueReader{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, TournamentID: tournamentID}, nil
		},
		getMemberFn: func(_ context.Context, _, _ uuid.UUID) (*domain.LeagueMember, error) {
			return &domain.LeagueMember{}, nil
		},
		listMembersForPredictionsFn: func(_ context.Context, _ uuid.UUID) ([]*domain.LeagueMemberDisplay, error) {
			return []*domain.LeagueMemberDisplay{
				{UserID: me, DisplayName: "me"},
				{UserID: alice, DisplayName: "alice"},
			}, nil
		},
	}

	playerID := uuid.New()
	ppRepo := &fakePlayerPredictionRepo{
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.PlayerLeaguePick, error) {
			return []*domain.PlayerLeaguePick{
				{UserID: me, Category: domain.PlayerHandicapCategoryGroupTopScorer, GroupLetter: strPtr("A"), PlayerID: playerID, PlayerName: "Messi"},
			}, nil
		},
	}

	// Kickoff was 1h ago → lock was 90 min ago → locked.
	fixtures := kickoffAt(time.Now().Add(-time.Hour))

	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, groupsA(), fixtures, leagues, fakeClock{time.Now()})

	views, err := svc.ListLeaguePlayerPredictions(context.Background(), leagueID, me)
	require.NoError(t, err)
	require.Len(t, views, len(domain.AllPlayerHandicapCategories))

	// Find the group_top_scorer category.
	var found *domain.LeaguePlayerCategoryView
	for _, v := range views {
		if v.Category == domain.PlayerHandicapCategoryGroupTopScorer {
			found = v
			break
		}
	}
	require.NotNil(t, found)
	require.Len(t, found.Predictions, 2)
	// Requesting user first.
	require.Equal(t, me, found.Predictions[0].UserID)
	require.NotNil(t, found.Predictions[0].PlayerID)
	// Alice has no prediction.
	require.Equal(t, alice, found.Predictions[1].UserID)
	require.Nil(t, found.Predictions[1].PlayerID)
}

// ---------- Group cross-conflict tests ----------

func TestTournamentPredictionService_UpsertTeam_CrossConflict_PlayoffBlocksGroupWinner(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	groupA := "A"

	// Existing: the same team is already a playoff pick in group A.
	tpRepo := &fakeTeamPredictionRepo{
		upsertFn: func(_ context.Context, _ domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error) {
			return &domain.TeamPrediction{}, nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return []*domain.TeamPrediction{
				{
					Category:    domain.TeamHandicapCategoryPlayoff,
					GroupLetter: &groupA,
					SlotIndex:   0,
					Pick:        teamID,
				},
			}, nil
		},
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return nil, nil
		},
	}
	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: tournamentID, GroupLetter: &groupA}, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.UpsertTeamPrediction(context.Background(), domain.UpsertTeamPredictionInput{
		UserID:       userID,
		TournamentID: tournamentID,
		Category:     domain.TeamHandicapCategoryGroupWinner,
		Pick:         teamID,
		GroupLetter:  &groupA,
		SlotIndex:    0,
	})
	require.True(t, errors.Is(err, domain.ErrInvalid), "expected ErrInvalid, got: %v", err)
}

func TestTournamentPredictionService_UpsertTeam_CrossConflict_GroupWinnerBlocksPlayoff(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	groupA := "A"

	// Existing: the same team is already picked as group_winner in group A.
	tpRepo := &fakeTeamPredictionRepo{
		upsertFn: func(_ context.Context, _ domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error) {
			return &domain.TeamPrediction{}, nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return []*domain.TeamPrediction{
				{
					Category:    domain.TeamHandicapCategoryGroupWinner,
					GroupLetter: &groupA,
					SlotIndex:   0,
					Pick:        teamID,
				},
			}, nil
		},
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return nil, nil
		},
	}
	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: tournamentID, GroupLetter: &groupA}, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	_, err := svc.UpsertTeamPrediction(context.Background(), domain.UpsertTeamPredictionInput{
		UserID:       userID,
		TournamentID: tournamentID,
		Category:     domain.TeamHandicapCategoryPlayoff,
		Pick:         teamID,
		GroupLetter:  &groupA,
		SlotIndex:    0,
	})
	require.True(t, errors.Is(err, domain.ErrInvalid), "expected ErrInvalid, got: %v", err)
}

func TestTournamentPredictionService_BulkUpsertTeams_CrossConflict_InBatch(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	groupA := "A"

	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: tournamentID, GroupLetter: &groupA}, nil
		},
	}
	tpRepo := &fakeTeamPredictionRepo{
		upsertFn: func(_ context.Context, _ domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error) {
			return &domain.TeamPrediction{}, nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return nil, nil // no prior picks
		},
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return nil, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	// Same team picked as group_winner AND playoff in the same batch.
	_, err := svc.BulkUpsertTeamPredictions(context.Background(), tournamentID, userID, []domain.BulkTeamPredictionItem{
		{Category: domain.TeamHandicapCategoryGroupWinner, GroupLetter: &groupA, SlotIndex: 0, TeamID: &teamID},
		{Category: domain.TeamHandicapCategoryPlayoff, GroupLetter: &groupA, SlotIndex: 1, TeamID: &teamID},
	})
	require.True(t, errors.Is(err, domain.ErrInvalid), "expected ErrInvalid, got: %v", err)
}

func TestTournamentPredictionService_BulkUpsertTeams_CrossConflict_AgainstExisting(t *testing.T) {
	tournamentID := uuid.New()
	userID := uuid.New()
	teamID := uuid.New()
	groupA := "A"

	teams := &fakeTeamGetter{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.Team, error) {
			return &domain.Team{ID: teamID, TournamentID: tournamentID, GroupLetter: &groupA}, nil
		},
	}
	// DB already has a group_winner pick for teamID in group A.
	tpRepo := &fakeTeamPredictionRepo{
		upsertFn: func(_ context.Context, _ domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error) {
			return &domain.TeamPrediction{}, nil
		},
		listByTournamentForUserFn: func(_ context.Context, _, _ uuid.UUID) ([]*domain.TeamPrediction, error) {
			return []*domain.TeamPrediction{
				{Category: domain.TeamHandicapCategoryGroupWinner, GroupLetter: &groupA, SlotIndex: 0, Pick: teamID},
			}, nil
		},
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return nil, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), noFixtures(), nil, fakeClock{time.Now()})

	// Batch tries to add the same team as a playoff pick.
	_, err := svc.BulkUpsertTeamPredictions(context.Background(), tournamentID, userID, []domain.BulkTeamPredictionItem{
		{Category: domain.TeamHandicapCategoryPlayoff, GroupLetter: &groupA, SlotIndex: 0, TeamID: &teamID},
	})
	require.True(t, errors.Is(err, domain.ErrInvalid), "expected ErrInvalid, got: %v", err)
}
