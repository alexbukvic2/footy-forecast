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

type fakeTournamentLockReader struct {
	isLockedFn func(context.Context, uuid.UUID) (bool, error)
}

func (f *fakeTournamentLockReader) IsPredictionsLocked(ctx context.Context, id uuid.UUID) (bool, error) {
	return f.isLockedFn(ctx, id)
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
	tournaments service.TournamentLockReader,
	leagues service.LeagueReader,
) *service.TournamentPredictionService {
	return service.NewTournamentPredictionService(ppRepo, tpRepo, players, teams, tGroups, tournaments, leagues)
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

func unlocked() *fakeTournamentLockReader {
	return &fakeTournamentLockReader{
		isLockedFn: func(_ context.Context, _ uuid.UUID) (bool, error) { return false, nil },
	}
}

func locked() *fakeTournamentLockReader {
	return &fakeTournamentLockReader{
		isLockedFn: func(_ context.Context, _ uuid.UUID) (bool, error) { return true, nil },
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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), locked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, teams, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, teams, noGroups(), locked(), nil)

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
	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), locked(), nil)

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
	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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
	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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
	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), players, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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
	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, teams, noGroups(), locked(), nil)

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
	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, groupsA(), unlocked(), nil)

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
	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), unlocked(), nil)

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

	svc := newSvc(ppRepo, defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, groupsA(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, &fakeTeamGetter{}, groupsA(), unlocked(), nil)

	_, views, err := svc.ListTeamPredictionsForUser(context.Background(), tournamentID, userID)
	require.NoError(t, err)
	require.Len(t, views, 9) // 1 group "A": 1 group_winner + 3 playoff + 4 semifinalist + 1 winner
	for _, v := range views {
		require.Nil(t, v.Prediction, "no predictions yet, all should be nil")
	}
}

// ---------- ListLeagueGroupPredictions tests ----------

func TestTournamentPredictionService_ListLeagueGroupPredictions_NotMember(t *testing.T) {
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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), unlocked(), leagues)

	_, err := svc.ListLeagueGroupPredictions(context.Background(), leagueID, userID, "A")
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTournamentPredictionService_ListLeagueGroupPredictions_BeforeLock(t *testing.T) {
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
		listMembersForPredictionsFn: func(_ context.Context, _ uuid.UUID) ([]*domain.LeagueMemberDisplay, error) {
			return nil, nil
		},
	}

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), unlocked(), leagues)

	_, err := svc.ListLeagueGroupPredictions(context.Background(), leagueID, userID, "A")
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTournamentPredictionService_ListLeagueGroupPredictions_AfterLock(t *testing.T) {
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

	teamID := uuid.New()
	playoffTeamID := uuid.New()
	tpRepo := &fakeTeamPredictionRepo{
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return []*domain.TeamLeaguePick{
				{UserID: me, Category: domain.TeamHandicapCategoryGroupWinner, GroupLetter: strPtr("A"), SlotIndex: 0, TeamID: teamID, TeamName: "Argentina"},
				{UserID: me, Category: domain.TeamHandicapCategoryPlayoff, GroupLetter: strPtr("A"), SlotIndex: 0, TeamID: playoffTeamID, TeamName: "Brazil"},
			}, nil
		},
	}

	svc := newSvc(ppRepo, tpRepo, &fakePlayerGetter{}, &fakeTeamGetter{}, groupsA(), locked(), leagues)

	result, err := svc.ListLeagueGroupPredictions(context.Background(), leagueID, me, "A")
	require.NoError(t, err)
	require.Equal(t, "A", result.Group)

	// group_winner + playoff (slots 0-2) views present
	require.Len(t, result.TeamPredictions, 4)
	require.Equal(t, domain.TeamHandicapCategoryGroupWinner, result.TeamPredictions[0].Category)
	require.Equal(t, 0, result.TeamPredictions[0].SlotIndex)
	require.Len(t, result.TeamPredictions[0].Predictions, 2)
	require.Equal(t, me, result.TeamPredictions[0].Predictions[0].UserID)
	require.NotNil(t, result.TeamPredictions[0].Predictions[0].TeamID)
	require.Nil(t, result.TeamPredictions[0].Predictions[1].TeamID)

	require.Equal(t, domain.TeamHandicapCategoryPlayoff, result.TeamPredictions[1].Category)
	require.Equal(t, 0, result.TeamPredictions[1].SlotIndex)
	require.NotNil(t, result.TeamPredictions[1].Predictions[0].TeamID)

	require.Equal(t, domain.TeamHandicapCategoryPlayoff, result.TeamPredictions[2].Category)
	require.Equal(t, 1, result.TeamPredictions[2].SlotIndex)
	require.Nil(t, result.TeamPredictions[2].Predictions[0].TeamID)

	require.Equal(t, domain.TeamHandicapCategoryPlayoff, result.TeamPredictions[3].Category)
	require.Equal(t, 2, result.TeamPredictions[3].SlotIndex)
	require.Nil(t, result.TeamPredictions[3].Predictions[0].TeamID)

	// group_top_scorer view present
	require.Len(t, result.PlayerPredictions, 1)
	require.Equal(t, domain.PlayerHandicapCategoryGroupTopScorer, result.PlayerPredictions[0].Category)
	require.Len(t, result.PlayerPredictions[0].Predictions, 2)
	require.Equal(t, me, result.PlayerPredictions[0].Predictions[0].UserID)
	require.NotNil(t, result.PlayerPredictions[0].Predictions[0].PlayerID)
	require.Nil(t, result.PlayerPredictions[0].Predictions[1].PlayerID)
}

// ---------- ListLeaguePlayoffPredictions tests ----------

func TestTournamentPredictionService_ListLeaguePlayoffPredictions_NotMember(t *testing.T) {
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

	svc := newSvc(defaultPlayerRepo(), defaultTeamRepo(), &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), unlocked(), leagues)

	_, err := svc.ListLeaguePlayoffPredictions(context.Background(), leagueID, userID)
	require.ErrorIs(t, err, domain.ErrForbidden)
}

func TestTournamentPredictionService_ListLeaguePlayoffPredictions_AfterLock(t *testing.T) {
	leagueID := uuid.New()
	tournamentID := uuid.New()
	me := uuid.New()

	leagues := &fakeLeagueReader{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*domain.League, error) {
			return &domain.League{ID: leagueID, TournamentID: tournamentID}, nil
		},
		getMemberFn: func(_ context.Context, _, _ uuid.UUID) (*domain.LeagueMember, error) {
			return &domain.LeagueMember{}, nil
		},
		listMembersForPredictionsFn: func(_ context.Context, _ uuid.UUID) ([]*domain.LeagueMemberDisplay, error) {
			return []*domain.LeagueMemberDisplay{{UserID: me, DisplayName: "me"}}, nil
		},
	}

	playerID := uuid.New()
	ppRepo := &fakePlayerPredictionRepo{
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.PlayerLeaguePick, error) {
			return []*domain.PlayerLeaguePick{
				{UserID: me, Category: domain.PlayerHandicapCategoryTotalTopScorer, PlayerID: playerID, PlayerName: "Neymar"},
			}, nil
		},
	}

	teamID := uuid.New()
	tpRepo := &fakeTeamPredictionRepo{
		listByLeagueFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TeamLeaguePick, error) {
			return []*domain.TeamLeaguePick{
				{UserID: me, Category: domain.TeamHandicapCategoryWinner, SlotIndex: 0, TeamID: teamID, TeamName: "Brazil"},
			}, nil
		},
	}

	svc := newSvc(ppRepo, tpRepo, &fakePlayerGetter{}, &fakeTeamGetter{}, noGroups(), locked(), leagues)

	result, err := svc.ListLeaguePlayoffPredictions(context.Background(), leagueID, me)
	require.NoError(t, err)

	// semifinalist (4) + winner (1) = 5 team category views
	require.Len(t, result.TeamPredictions, 5)

	// Find winner category
	var winnerView *domain.LeagueTeamCategoryView
	for _, v := range result.TeamPredictions {
		if v.Category == domain.TeamHandicapCategoryWinner {
			winnerView = v
			break
		}
	}
	require.NotNil(t, winnerView)
	require.Len(t, winnerView.Predictions, 1)
	require.NotNil(t, winnerView.Predictions[0].TeamID)

	// total_top_scorer player view present
	require.Len(t, result.PlayerPredictions, 1)
	require.Equal(t, domain.PlayerHandicapCategoryTotalTopScorer, result.PlayerPredictions[0].Category)
	require.NotNil(t, result.PlayerPredictions[0].Predictions[0].PlayerID)
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

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), unlocked(), nil)

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

	svc := newSvc(defaultPlayerRepo(), tpRepo, &fakePlayerGetter{}, teams, noGroups(), unlocked(), nil)

	// Batch tries to add the same team as a playoff pick.
	_, err := svc.BulkUpsertTeamPredictions(context.Background(), tournamentID, userID, []domain.BulkTeamPredictionItem{
		{Category: domain.TeamHandicapCategoryPlayoff, GroupLetter: &groupA, SlotIndex: 0, TeamID: &teamID},
	})
	require.True(t, errors.Is(err, domain.ErrInvalid), "expected ErrInvalid, got: %v", err)
}
