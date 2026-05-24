package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// Clock allows injecting a deterministic time source for testing.
type Clock interface {
	Now() time.Time
}

// RealClock returns the actual current time.
type RealClock struct{}

// Now returns time.Now().
func (RealClock) Now() time.Time { return time.Now() }

// PlayerPredictionRepo is the data access interface for player predictions.
type PlayerPredictionRepo interface {
	UpsertPlayer(ctx context.Context, in domain.UpsertPlayerPredictionInput) (*domain.PlayerPrediction, error)
	ListPlayersByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.PlayerPrediction, error)
	ListPlayersByLeague(ctx context.Context, leagueID uuid.UUID) ([]*domain.PlayerLeaguePick, error)
}

// TeamPredictionRepo is the data access interface for team predictions.
type TeamPredictionRepo interface {
	UpsertTeam(ctx context.Context, in domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error)
	ListTeamsByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.TeamPrediction, error)
	ListTeamsByLeague(ctx context.Context, leagueID uuid.UUID) ([]*domain.TeamLeaguePick, error)
}

// PlayerGetter validates that a player exists.
type PlayerGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Player, error)
}

// TeamGetter validates that a team exists.
type TeamGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
}

// FixtureFirstKickoffGetter returns the first kickoff time for a tournament.
// Returns domain.ErrNotFound when no fixtures exist for the tournament.
type FixtureFirstKickoffGetter interface {
	GetFirstKickoffByTournament(ctx context.Context, tournamentID uuid.UUID) (time.Time, error)
}

// LeagueReader reads league and membership data needed for prediction views.
type LeagueReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.League, error)
	GetMember(ctx context.Context, leagueID, userID uuid.UUID) (*domain.LeagueMember, error)
	ListMembersForPredictions(ctx context.Context, leagueID uuid.UUID) ([]*domain.LeagueMemberDisplay, error)
}

// TournamentPredictionService orchestrates tournament-level prediction use cases.
type TournamentPredictionService struct {
	playerPredictions PlayerPredictionRepo
	teamPredictions   TeamPredictionRepo
	players           PlayerGetter
	teams             TeamGetter
	fixtures          FixtureFirstKickoffGetter
	leagues           LeagueReader
	clock             Clock
}

// NewTournamentPredictionService constructs a TournamentPredictionService.
func NewTournamentPredictionService(
	playerPredictions PlayerPredictionRepo,
	teamPredictions TeamPredictionRepo,
	players PlayerGetter,
	teams TeamGetter,
	fixtures FixtureFirstKickoffGetter,
	leagues LeagueReader,
	clock Clock,
) *TournamentPredictionService {
	return &TournamentPredictionService{
		playerPredictions: playerPredictions,
		teamPredictions:   teamPredictions,
		players:           players,
		teams:             teams,
		fixtures:          fixtures,
		leagues:           leagues,
		clock:             clock,
	}
}

// lockAt returns the lock time for a tournament and whether it has passed.
// When no fixtures exist, isLocked is false (predictions remain open).
func (s *TournamentPredictionService) lockAt(
	ctx context.Context,
	tournamentID uuid.UUID,
) (lockTime time.Time, isLocked bool, err error) {
	firstKickoff, err := s.fixtures.GetFirstKickoffByTournament(ctx, tournamentID)
	if errors.Is(err, domain.ErrNotFound) {
		return time.Time{}, false, nil
	}
	if err != nil {
		return time.Time{}, false, fmt.Errorf("get first kickoff: %w", err)
	}
	lockAt := firstKickoff.Add(-30 * time.Minute)
	return lockAt, !s.clock.Now().Before(lockAt), nil
}

// UpsertPlayerPrediction validates and saves a player tournament prediction.
// Returns domain.ErrNotFound if the player doesn't exist or belongs to a different tournament.
// Returns domain.ErrForbidden if predictions are locked.
func (s *TournamentPredictionService) UpsertPlayerPrediction(
	ctx context.Context,
	in domain.UpsertPlayerPredictionInput,
) (*domain.PlayerPrediction, error) {
	player, err := s.players.GetByID(ctx, in.Pick)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("player %s not found in this tournament: %w", in.Pick, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get player: %w", err)
	}
	if player.TournamentID != in.TournamentID {
		return nil, fmt.Errorf("player %s not found in this tournament: %w", in.Pick, domain.ErrNotFound)
	}

	_, locked, err := s.lockAt(ctx, in.TournamentID)
	if err != nil {
		return nil, err
	}
	if locked {
		return nil, fmt.Errorf("predictions are locked for this tournament: %w", domain.ErrForbidden)
	}

	return s.playerPredictions.UpsertPlayer(ctx, in)
}

// UpsertTeamPrediction validates and saves a team tournament prediction.
// Returns domain.ErrNotFound if the team doesn't exist or belongs to a different tournament.
// Returns domain.ErrForbidden if predictions are locked.
func (s *TournamentPredictionService) UpsertTeamPrediction(
	ctx context.Context,
	in domain.UpsertTeamPredictionInput,
) (*domain.TeamPrediction, error) {
	team, err := s.teams.GetByID(ctx, in.Pick)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("team %s not found in this tournament: %w", in.Pick, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get team: %w", err)
	}
	if team.TournamentID != in.TournamentID {
		return nil, fmt.Errorf("team %s not found in this tournament: %w", in.Pick, domain.ErrNotFound)
	}

	_, locked, err := s.lockAt(ctx, in.TournamentID)
	if err != nil {
		return nil, err
	}
	if locked {
		return nil, fmt.Errorf("predictions are locked for this tournament: %w", domain.ErrForbidden)
	}

	return s.teamPredictions.UpsertTeam(ctx, in)
}

// ListPlayerPredictionsForUser returns one entry per player category.
// Categories the user has not predicted yet have Prediction == nil.
func (s *TournamentPredictionService) ListPlayerPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.PlayerPredictionView, error) {
	rows, err := s.playerPredictions.ListPlayersByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return nil, fmt.Errorf("list player predictions: %w", err)
	}
	byCategory := make(map[domain.PlayerHandicapCategory]*domain.PlayerPrediction, len(rows))
	for _, r := range rows {
		byCategory[r.Category] = r
	}
	views := make([]*domain.PlayerPredictionView, 0, len(domain.AllPlayerHandicapCategories))
	for _, cat := range domain.AllPlayerHandicapCategories {
		views = append(views, &domain.PlayerPredictionView{
			Category:   cat,
			Prediction: byCategory[cat],
		})
	}
	return views, nil
}

// ListTeamPredictionsForUser returns one entry per team category.
// Categories the user has not predicted yet have Prediction == nil.
func (s *TournamentPredictionService) ListTeamPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.TeamPredictionView, error) {
	rows, err := s.teamPredictions.ListTeamsByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return nil, fmt.Errorf("list team predictions: %w", err)
	}
	byCategory := make(map[domain.TeamHandicapCategory]*domain.TeamPrediction, len(rows))
	for _, r := range rows {
		byCategory[r.Category] = r
	}
	views := make([]*domain.TeamPredictionView, 0, len(domain.AllTeamHandicapCategories))
	for _, cat := range domain.AllTeamHandicapCategories {
		views = append(views, &domain.TeamPredictionView{
			Category:   cat,
			Prediction: byCategory[cat],
		})
	}
	return views, nil
}

// ListLeaguePlayerPredictions returns all member picks grouped by player category.
// Requires the requesting user to be a league member and the lock time to have passed.
// The requesting user's pick appears first within each category; remaining members are
// sorted alphabetically by display name.
func (s *TournamentPredictionService) ListLeaguePlayerPredictions(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
) ([]*domain.LeaguePlayerCategoryView, error) {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get league: %w", err)
	}

	if _, err := s.leagues.GetMember(ctx, leagueID, requestingUserID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("not a member of league %s: %w", leagueID, domain.ErrForbidden)
		}
		return nil, fmt.Errorf("check membership: %w", err)
	}

	_, locked, err := s.lockAt(ctx, league.TournamentID)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, fmt.Errorf("league predictions not available until lock time has passed: %w", domain.ErrForbidden)
	}

	members, err := s.leagues.ListMembersForPredictions(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	picks, err := s.playerPredictions.ListPlayersByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list player predictions: %w", err)
	}

	return buildLeaguePlayerViews(members, picks, requestingUserID), nil
}

// ListLeagueTeamPredictions returns all member picks grouped by team category.
// Same membership and lock-time constraints as ListLeaguePlayerPredictions.
func (s *TournamentPredictionService) ListLeagueTeamPredictions(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
) ([]*domain.LeagueTeamCategoryView, error) {
	league, err := s.leagues.GetByID(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("get league: %w", err)
	}

	if _, err := s.leagues.GetMember(ctx, leagueID, requestingUserID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("not a member of league %s: %w", leagueID, domain.ErrForbidden)
		}
		return nil, fmt.Errorf("check membership: %w", err)
	}

	_, locked, err := s.lockAt(ctx, league.TournamentID)
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, fmt.Errorf("league predictions not available until lock time has passed: %w", domain.ErrForbidden)
	}

	members, err := s.leagues.ListMembersForPredictions(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	picks, err := s.teamPredictions.ListTeamsByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list team predictions: %w", err)
	}

	return buildLeagueTeamViews(members, picks, requestingUserID), nil
}

func buildLeaguePlayerViews(
	members []*domain.LeagueMemberDisplay,
	picks []*domain.PlayerLeaguePick,
	requestingUserID uuid.UUID,
) []*domain.LeaguePlayerCategoryView {
	// Index picks by (userID, category).
	type key struct {
		userID   uuid.UUID
		category domain.PlayerHandicapCategory
	}
	pickMap := make(map[key]*domain.PlayerLeaguePick, len(picks))
	for _, p := range picks {
		pickMap[key{p.UserID, p.Category}] = p
	}

	// Sort members: requesting user first, then alphabetically.
	sorted := sortedMembers(members, requestingUserID)

	views := make([]*domain.LeaguePlayerCategoryView, 0, len(domain.AllPlayerHandicapCategories))
	for _, cat := range domain.AllPlayerHandicapCategories {
		memberPicks := make([]domain.LeagueMemberPlayerPick, 0, len(sorted))
		for _, m := range sorted {
			pick := pickMap[key{m.UserID, cat}]
			mp := domain.LeagueMemberPlayerPick{
				UserID:      m.UserID,
				DisplayName: m.DisplayName,
			}
			if pick != nil {
				mp.PlayerID = &pick.PlayerID
				mp.PlayerName = &pick.PlayerName
				mp.Points = pick.Points
			}
			memberPicks = append(memberPicks, mp)
		}
		views = append(views, &domain.LeaguePlayerCategoryView{
			Category:    cat,
			Predictions: memberPicks,
		})
	}
	return views
}

func buildLeagueTeamViews(
	members []*domain.LeagueMemberDisplay,
	picks []*domain.TeamLeaguePick,
	requestingUserID uuid.UUID,
) []*domain.LeagueTeamCategoryView {
	type key struct {
		userID   uuid.UUID
		category domain.TeamHandicapCategory
	}
	pickMap := make(map[key]*domain.TeamLeaguePick, len(picks))
	for _, p := range picks {
		pickMap[key{p.UserID, p.Category}] = p
	}

	sorted := sortedMembers(members, requestingUserID)

	views := make([]*domain.LeagueTeamCategoryView, 0, len(domain.AllTeamHandicapCategories))
	for _, cat := range domain.AllTeamHandicapCategories {
		memberPicks := make([]domain.LeagueMemberTeamPick, 0, len(sorted))
		for _, m := range sorted {
			pick := pickMap[key{m.UserID, cat}]
			mp := domain.LeagueMemberTeamPick{
				UserID:      m.UserID,
				DisplayName: m.DisplayName,
			}
			if pick != nil {
				mp.TeamID = &pick.TeamID
				mp.TeamName = &pick.TeamName
				mp.Points = pick.Points
			}
			memberPicks = append(memberPicks, mp)
		}
		views = append(views, &domain.LeagueTeamCategoryView{
			Category:    cat,
			Predictions: memberPicks,
		})
	}
	return views
}

// compile-time interface checks
var (
	_ PlayerPredictionRepo      = (*repository.PlayerPredictionRepository)(nil)
	_ TeamPredictionRepo        = (*repository.TeamPredictionRepository)(nil)
	_ PlayerGetter              = (*repository.PlayerRepository)(nil)
	_ TeamGetter                = (*repository.TeamRepository)(nil)
	_ FixtureFirstKickoffGetter = (*repository.FixtureRepository)(nil)
	_ LeagueReader              = (*repository.LeagueRepository)(nil)
)

// sortedMembers returns members sorted with requestingUserID first, then alphabetically.
func sortedMembers(members []*domain.LeagueMemberDisplay, requestingUserID uuid.UUID) []*domain.LeagueMemberDisplay {
	sorted := make([]*domain.LeagueMemberDisplay, len(members))
	copy(sorted, members)
	sort.SliceStable(sorted, func(i, j int) bool {
		iIsMe := sorted[i].UserID == requestingUserID
		jIsMe := sorted[j].UserID == requestingUserID
		if iIsMe != jIsMe {
			return iIsMe
		}
		return sorted[i].DisplayName < sorted[j].DisplayName
	})
	return sorted
}
