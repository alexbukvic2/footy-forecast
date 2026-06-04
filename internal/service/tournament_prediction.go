package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

// PlayerPredictionRepo is the data access interface for player predictions.
type PlayerPredictionRepo interface {
	UpsertPlayer(
		ctx context.Context,
		in domain.UpsertPlayerPredictionInput,
	) (*domain.PlayerPrediction, error)
	DeletePlayer(
		ctx context.Context,
		userID, tournamentID uuid.UUID,
		category domain.PlayerHandicapCategory,
		groupLetter *string,
	) error
	ListPlayersByTournamentForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.PlayerPrediction, error)
	ListPlayersByLeague(
		ctx context.Context,
		leagueID uuid.UUID,
	) ([]*domain.PlayerLeaguePick, error)
}

// TeamPredictionRepo is the data access interface for team predictions.
type TeamPredictionRepo interface {
	UpsertTeam(
		ctx context.Context,
		in domain.UpsertTeamPredictionInput,
	) (*domain.TeamPrediction, error)
	DeleteTeam(
		ctx context.Context,
		userID, tournamentID uuid.UUID,
		category domain.TeamHandicapCategory,
		groupLetter *string,
		slotIndex int,
	) error
	ListTeamsByTournamentForUser(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) ([]*domain.TeamPrediction, error)
	ListTeamsByLeague(
		ctx context.Context,
		leagueID uuid.UUID,
	) ([]*domain.TeamLeaguePick, error)
	CountPlayoffWildcards(
		ctx context.Context,
		tournamentID, userID uuid.UUID,
	) (int, error)
}

// PlayerGetter validates that a player exists.
type PlayerGetter interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Player, error)
}

// TeamGetter validates that a team exists.
type TeamGetter interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.Team, error)
}

// TeamGroupLister returns the sorted group letters for a tournament.
type TeamGroupLister interface {
	ListGroupLettersByTournament(
		ctx context.Context,
		tournamentID uuid.UUID,
	) ([]string, error)
}

// TournamentLockReader checks whether predictions are locked for a tournament.
type TournamentLockReader interface {
	IsPredictionsLocked(ctx context.Context, tournamentID uuid.UUID) (bool, error)
}

// LeagueReader reads league and membership data needed for prediction views.
type LeagueReader interface {
	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*domain.League, error)
	GetMember(
		ctx context.Context,
		leagueID, userID uuid.UUID,
	) (*domain.LeagueMember, error)
	ListMembersForPredictions(
		ctx context.Context,
		leagueID uuid.UUID,
	) ([]*domain.LeagueMemberDisplay, error)
}

// TournamentPredictionService orchestrates tournament-level prediction use cases.
type TournamentPredictionService struct {
	playerPredictions PlayerPredictionRepo
	teamPredictions   TeamPredictionRepo
	players           PlayerGetter
	teams             TeamGetter
	teamGroups        TeamGroupLister
	tournaments       TournamentLockReader
	leagues           LeagueReader
}

// NewTournamentPredictionService constructs a TournamentPredictionService.
func NewTournamentPredictionService(
	playerPredictions PlayerPredictionRepo,
	teamPredictions TeamPredictionRepo,
	players PlayerGetter,
	teams TeamGetter,
	teamGroups TeamGroupLister,
	tournaments TournamentLockReader,
	leagues LeagueReader,
) *TournamentPredictionService {
	return &TournamentPredictionService{
		playerPredictions: playerPredictions,
		teamPredictions:   teamPredictions,
		players:           players,
		teams:             teams,
		teamGroups:        teamGroups,
		tournaments:       tournaments,
		leagues:           leagues,
	}
}

// UpsertPlayerPrediction validates and saves a player tournament prediction.
// Returns domain.ErrNotFound if the player doesn't exist or belongs to a different tournament.
// Returns domain.ErrInvalid if group requirements are not met.
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

	// Validate group requirements.
	switch in.Category {
	case domain.PlayerHandicapCategoryGroupTopScorer:
		if in.GroupLetter == nil {
			return nil, fmt.Errorf("group is required for group_top_scorer: %w", domain.ErrInvalid)
		}
		if player.GroupLetter == nil || *player.GroupLetter != *in.GroupLetter {
			return nil, fmt.Errorf("player's team is not in group %s: %w", *in.GroupLetter, domain.ErrNotFound)
		}
	case domain.PlayerHandicapCategoryTotalTopScorer:
		if in.GroupLetter != nil {
			return nil, fmt.Errorf("group must be null for total_top_scorer: %w", domain.ErrInvalid)
		}
	}

	locked, err := s.tournaments.IsPredictionsLocked(ctx, in.TournamentID)
	if err != nil {
		return nil, fmt.Errorf("check lock: %w", err)
	}
	if locked {
		return nil, fmt.Errorf("predictions are locked for this tournament: %w", domain.ErrForbidden)
	}

	return s.playerPredictions.UpsertPlayer(ctx, in)
}

// UpsertTeamPrediction validates and saves a team tournament prediction.
// Returns domain.ErrNotFound if the team doesn't exist or belongs to a different tournament.
// Returns domain.ErrInvalid if slot/group requirements are not met.
// Returns domain.ErrForbidden if predictions are locked or wildcard limit reached.
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

	if err := s.validateTeamCategory(ctx, team, in); err != nil {
		return nil, err
	}

	if err := s.checkGroupCrossConflict(
		ctx,
		in.TournamentID,
		in.UserID,
		in.Category,
		in.GroupLetter,
		in.Pick,
	); err != nil {
		return nil, err
	}

	locked, err := s.tournaments.IsPredictionsLocked(ctx, in.TournamentID)
	if err != nil {
		return nil, fmt.Errorf("check lock: %w", err)
	}
	if locked {
		return nil, fmt.Errorf("predictions are locked for this tournament: %w", domain.ErrForbidden)
	}

	return s.teamPredictions.UpsertTeam(ctx, in)
}

func (s *TournamentPredictionService) validateTeamCategory(
	ctx context.Context,
	team *domain.Team,
	in domain.UpsertTeamPredictionInput,
) error {
	switch in.Category {
	case domain.TeamHandicapCategoryGroupWinner:
		return validateGroupWinner(team, in)
	case domain.TeamHandicapCategoryPlayoff:
		return s.validatePlayoff(ctx, team, in)
	case domain.TeamHandicapCategorySemifinalist:
		return validateSemifinalist(in)
	case domain.TeamHandicapCategoryWinner:
		return validateWinner(in)
	}
	return nil
}

func validateGroupWinner(
	team *domain.Team,
	in domain.UpsertTeamPredictionInput,
) error {
	if in.GroupLetter == nil {
		return fmt.Errorf("group is required for group_winner: %w", domain.ErrInvalid)
	}
	if in.SlotIndex != 0 {
		return fmt.Errorf("slot must be 0 for group_winner: %w", domain.ErrInvalid)
	}
	if team.GroupLetter == nil || *team.GroupLetter != *in.GroupLetter {
		return fmt.Errorf("team is not in group %s: %w", *in.GroupLetter, domain.ErrNotFound)
	}
	return nil
}

func (s *TournamentPredictionService) validatePlayoff(
	ctx context.Context,
	team *domain.Team,
	in domain.UpsertTeamPredictionInput,
) error {
	if in.GroupLetter == nil {
		return fmt.Errorf("group is required for playoff: %w", domain.ErrInvalid)
	}
	if in.SlotIndex < 0 || in.SlotIndex > 2 {
		return fmt.Errorf("slot must be 0, 1, or 2 for playoff: %w", domain.ErrInvalid)
	}
	if team.GroupLetter == nil || *team.GroupLetter != *in.GroupLetter {
		return fmt.Errorf("team is not in group %s: %w", *in.GroupLetter, domain.ErrNotFound)
	}
	if in.SlotIndex == 2 {
		count, err := s.teamPredictions.CountPlayoffWildcards(ctx, in.TournamentID, in.UserID)
		if err != nil {
			return fmt.Errorf("count wildcards: %w", err)
		}
		if count >= 8 {
			return fmt.Errorf("maximum 8 wildcard picks reached: %w", domain.ErrForbidden)
		}
	}
	return nil
}

func validateSemifinalist(in domain.UpsertTeamPredictionInput) error {
	if in.GroupLetter != nil {
		return fmt.Errorf("group must be null for semifinalist: %w", domain.ErrInvalid)
	}
	if in.SlotIndex < 0 || in.SlotIndex > 3 {
		return fmt.Errorf("slot must be 0-3 for semifinalist: %w", domain.ErrInvalid)
	}
	return nil
}

func validateWinner(in domain.UpsertTeamPredictionInput) error {
	if in.GroupLetter != nil {
		return fmt.Errorf("group must be null for winner: %w", domain.ErrInvalid)
	}
	if in.SlotIndex != 0 {
		return fmt.Errorf("slot must be 0 for winner: %w", domain.ErrInvalid)
	}
	return nil
}

// BulkUpsertPlayerPredictions validates and saves a batch of player tournament predictions.
// Items with nil PlayerID clear (delete) the slot; items with a non-nil PlayerID upsert it.
// Returns the full player prediction view after applying the batch.
// Returns domain.ErrForbidden if predictions are locked.
// Returns domain.ErrNotFound / domain.ErrInvalid for the first invalid item encountered.
func (s *TournamentPredictionService) BulkUpsertPlayerPredictions(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	items []domain.BulkPlayerPredictionItem,
) ([]*domain.PlayerPredictionView, error) {
	if len(items) == 0 {
		_, views, err := s.ListPlayerPredictionsForUser(ctx, tournamentID, userID)
		return views, err
	}

	locked, err := s.tournaments.IsPredictionsLocked(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("check lock: %w", err)
	}
	if locked {
		return nil, fmt.Errorf("predictions are locked for this tournament: %w", domain.ErrForbidden)
	}

	for _, item := range items {
		if err := s.applyPlayerItem(ctx, tournamentID, userID, item); err != nil {
			return nil, err
		}
	}

	_, views, err := s.ListPlayerPredictionsForUser(ctx, tournamentID, userID)
	return views, err
}

// applyPlayerItem processes a single item from a bulk player prediction batch.
func (s *TournamentPredictionService) applyPlayerItem(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	item domain.BulkPlayerPredictionItem,
) error {
	if item.PlayerID == nil {
		return s.playerPredictions.DeletePlayer(ctx, userID, tournamentID, item.Category, item.GroupLetter)
	}

	player, err := s.players.GetByID(ctx, *item.PlayerID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return fmt.Errorf("player %s not found in this tournament: %w", *item.PlayerID, domain.ErrNotFound)
		}
		return fmt.Errorf("get player: %w", err)
	}
	if player.TournamentID != tournamentID {
		return fmt.Errorf("player %s not found in this tournament: %w", *item.PlayerID, domain.ErrNotFound)
	}

	switch item.Category {
	case domain.PlayerHandicapCategoryGroupTopScorer:
		if item.GroupLetter == nil {
			return fmt.Errorf("group is required for group_top_scorer: %w", domain.ErrInvalid)
		}
		if player.GroupLetter == nil || *player.GroupLetter != *item.GroupLetter {
			return fmt.Errorf("player's team is not in group %s: %w", *item.GroupLetter, domain.ErrNotFound)
		}
	case domain.PlayerHandicapCategoryTotalTopScorer:
		if item.GroupLetter != nil {
			return fmt.Errorf("group must be null for total_top_scorer: %w", domain.ErrInvalid)
		}
	}

	_, err = s.playerPredictions.UpsertPlayer(
		ctx, domain.UpsertPlayerPredictionInput{
			UserID:       userID,
			TournamentID: tournamentID,
			Category:     item.Category,
			Pick:         *item.PlayerID,
			GroupLetter:  item.GroupLetter,
		},
	)
	if err != nil {
		return fmt.Errorf("upsert player prediction: %w", err)
	}
	return nil
}

// BulkUpsertTeamPredictions validates and saves a batch of team tournament predictions.
// Items with nil TeamID clear (delete) the slot; items with a non-nil TeamID upsert it.
// Returns the full team prediction view after applying the batch.
// Returns domain.ErrForbidden if predictions are locked or the wildcard cap (8) would be exceeded.
// Returns domain.ErrNotFound / domain.ErrInvalid for the first invalid item encountered.
func (s *TournamentPredictionService) BulkUpsertTeamPredictions(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	items []domain.BulkTeamPredictionItem,
) ([]*domain.TeamPredictionView, error) {
	if len(items) == 0 {
		_, views, err := s.ListTeamPredictionsForUser(ctx, tournamentID, userID)
		return views, err
	}

	locked, err := s.tournaments.IsPredictionsLocked(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("check lock: %w", err)
	}
	if locked {
		return nil, fmt.Errorf("predictions are locked for this tournament: %w", domain.ErrForbidden)
	}

	// Check wildcard cap before persisting anything.
	if err := s.checkWildcardCap(ctx, tournamentID, userID, items); err != nil {
		return nil, err
	}

	// Check that no team ends up as both group_winner and playoff for the same group.
	existingRows, err := s.teamPredictions.ListTeamsByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return nil, fmt.Errorf("list existing team predictions: %w", err)
	}
	if err := validateProjectedTeamConflicts(existingRows, items); err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.TeamID == nil {
			if err := s.teamPredictions.DeleteTeam(
				ctx,
				userID,
				tournamentID,
				item.Category,
				item.GroupLetter,
				item.SlotIndex,
			); err != nil {
				return nil, fmt.Errorf("delete team prediction: %w", err)
			}
			continue
		}

		team, err := s.teams.GetByID(ctx, *item.TeamID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, fmt.Errorf("team %s not found in this tournament: %w", *item.TeamID, domain.ErrNotFound)
			}
			return nil, fmt.Errorf("get team: %w", err)
		}
		if team.TournamentID != tournamentID {
			return nil, fmt.Errorf("team %s not found in this tournament: %w", *item.TeamID, domain.ErrNotFound)
		}

		in := domain.UpsertTeamPredictionInput{
			UserID:       userID,
			TournamentID: tournamentID,
			Category:     item.Category,
			Pick:         *item.TeamID,
			GroupLetter:  item.GroupLetter,
			SlotIndex:    item.SlotIndex,
		}
		if err := s.validateTeamCategory(ctx, team, in); err != nil {
			return nil, err
		}

		if _, err := s.teamPredictions.UpsertTeam(ctx, in); err != nil {
			return nil, fmt.Errorf("upsert team prediction: %w", err)
		}
	}

	_, teamViews, err := s.ListTeamPredictionsForUser(ctx, tournamentID, userID)
	return teamViews, err
}

// checkGroupCrossConflict prevents picking the same team as both group_winner and playoff
// in the same group. It loads existing picks and checks for a collision with the incoming pick.
func (s *TournamentPredictionService) checkGroupCrossConflict(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	category domain.TeamHandicapCategory,
	groupLetter *string,
	teamID uuid.UUID,
) error {
	if (category != domain.TeamHandicapCategoryGroupWinner && category != domain.TeamHandicapCategoryPlayoff) ||
		groupLetter == nil {
		return nil
	}

	rows, err := s.teamPredictions.ListTeamsByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return fmt.Errorf("list team predictions: %w", err)
	}

	for _, r := range rows {
		if r.GroupLetter == nil || *r.GroupLetter != *groupLetter || r.Pick != teamID {
			continue
		}
		if category == domain.TeamHandicapCategoryGroupWinner && r.Category == domain.TeamHandicapCategoryPlayoff {
			return fmt.Errorf(
				"team is already picked as a playoff team for group %s: %w", *groupLetter, domain.ErrInvalid,
			)
		}
		if category == domain.TeamHandicapCategoryPlayoff && r.Category == domain.TeamHandicapCategoryGroupWinner {
			return fmt.Errorf(
				"team is already picked as group winner for group %s: %w", *groupLetter, domain.ErrInvalid,
			)
		}
	}
	return nil
}

// validateProjectedTeamConflicts builds the post-batch state in memory and checks that
// no team appears as both group_winner and a playoff pick for the same group.
func validateProjectedTeamConflicts(
	existing []*domain.TeamPrediction,
	items []domain.BulkTeamPredictionItem,
) error {
	type slotKey struct {
		category    domain.TeamHandicapCategory
		groupLetter string
		slotIndex   int
	}

	// Start from existing DB state.
	projected := make(map[slotKey]uuid.UUID, len(existing))
	for _, r := range existing {
		gl := ""
		if r.GroupLetter != nil {
			gl = *r.GroupLetter
		}
		projected[slotKey{r.Category, gl, r.SlotIndex}] = r.Pick
	}

	// Apply batch items.
	for _, item := range items {
		gl := ""
		if item.GroupLetter != nil {
			gl = *item.GroupLetter
		}
		k := slotKey{item.Category, gl, item.SlotIndex}
		if item.TeamID == nil {
			delete(projected, k)
		} else {
			projected[k] = *item.TeamID
		}
	}

	// Validate: no team may be both group_winner(slot 0) and any playoff slot for the same group.
	type groupKey struct {
		groupLetter string
		teamID      uuid.UUID
	}
	winners := make(map[groupKey]bool)
	playoffs := make(map[groupKey]bool)
	for k, teamID := range projected {
		if k.groupLetter == "" {
			continue
		}
		gk := groupKey{k.groupLetter, teamID}
		switch k.category {
		case domain.TeamHandicapCategoryGroupWinner:
			winners[gk] = true
		case domain.TeamHandicapCategoryPlayoff:
			playoffs[gk] = true
		}
	}
	for gk := range winners {
		if playoffs[gk] {
			return fmt.Errorf(
				"team cannot be picked as both group winner and playoff team for group %s: %w",
				gk.groupLetter, domain.ErrInvalid,
			)
		}
	}
	return nil
}

// checkWildcardCap verifies that applying the batch would not exceed 8 playoff wildcard slots.// The estimate is: dbCount + batchAdds - batchClears where adds/clears are playoff slot_index=2 items.
func (s *TournamentPredictionService) checkWildcardCap(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	items []domain.BulkTeamPredictionItem,
) error {
	var batchAdds, batchClears int
	for _, item := range items {
		if item.Category == domain.TeamHandicapCategoryPlayoff && item.SlotIndex == 2 {
			if item.TeamID != nil {
				batchAdds++
			} else {
				batchClears++
			}
		}
	}
	if batchAdds == 0 {
		return nil // no wildcards being added; nothing to check
	}

	dbCount, err := s.teamPredictions.CountPlayoffWildcards(ctx, tournamentID, userID)
	if err != nil {
		return fmt.Errorf("count wildcards: %w", err)
	}
	if dbCount+batchAdds-batchClears > 8 {
		return fmt.Errorf("maximum 8 wildcard picks reached: %w", domain.ErrForbidden)
	}
	return nil
}

// ListPlayerPredictionsForUser returns one entry per (category, group) slot.
// Slots the user has not yet predicted have Prediction == nil.
func (s *TournamentPredictionService) ListPlayerPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) (locked bool, views []*domain.PlayerPredictionView, err error) {
	locked, err = s.tournaments.IsPredictionsLocked(ctx, tournamentID)
	if err != nil {
		return false, nil, fmt.Errorf("check lock: %w", err)
	}

	groups, err := s.teamGroups.ListGroupLettersByTournament(ctx, tournamentID)
	if err != nil {
		return false, nil, fmt.Errorf("list groups: %w", err)
	}

	rows, err := s.playerPredictions.ListPlayersByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return false, nil, fmt.Errorf("list player predictions: %w", err)
	}

	// Index by (category, groupLetter).
	type key struct {
		category    domain.PlayerHandicapCategory
		groupLetter string // empty string for nil (total_top_scorer)
	}
	byKey := make(map[key]*domain.PlayerPrediction, len(rows))
	for _, r := range rows {
		gl := ""
		if r.GroupLetter != nil {
			gl = *r.GroupLetter
		}
		byKey[key{r.Category, gl}] = r
	}

	playerViews := make([]*domain.PlayerPredictionView, 0, len(groups)+1)
	// group_top_scorer: one view per group
	for _, g := range groups {
		gCopy := g
		playerViews = append(
			playerViews, &domain.PlayerPredictionView{
				Category:    domain.PlayerHandicapCategoryGroupTopScorer,
				GroupLetter: &gCopy,
				Prediction:  byKey[key{domain.PlayerHandicapCategoryGroupTopScorer, g}],
			},
		)
	}
	// total_top_scorer: one view, group=nil
	playerViews = append(
		playerViews, &domain.PlayerPredictionView{
			Category:    domain.PlayerHandicapCategoryTotalTopScorer,
			GroupLetter: nil,
			Prediction:  byKey[key{domain.PlayerHandicapCategoryTotalTopScorer, ""}],
		},
	)
	return locked, playerViews, nil
}

// ListTeamPredictionsForUser returns one entry per (category, group, slot) slot.
// Slots the user has not yet predicted have Prediction == nil.
func (s *TournamentPredictionService) ListTeamPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) (locked bool, views []*domain.TeamPredictionView, err error) {
	locked, err = s.tournaments.IsPredictionsLocked(ctx, tournamentID)
	if err != nil {
		return false, nil, fmt.Errorf("check lock: %w", err)
	}

	groups, err := s.teamGroups.ListGroupLettersByTournament(ctx, tournamentID)
	if err != nil {
		return false, nil, fmt.Errorf("list groups: %w", err)
	}

	rows, err := s.teamPredictions.ListTeamsByTournamentForUser(ctx, tournamentID, userID)
	if err != nil {
		return false, nil, fmt.Errorf("list team predictions: %w", err)
	}

	type key struct {
		category    domain.TeamHandicapCategory
		groupLetter string
		slotIndex   int
	}
	byKey := make(map[key]*domain.TeamPrediction, len(rows))
	for _, r := range rows {
		gl := ""
		if r.GroupLetter != nil {
			gl = *r.GroupLetter
		}
		byKey[key{r.Category, gl, r.SlotIndex}] = r
	}

	teamViews := make([]*domain.TeamPredictionView, 0, 4*len(groups)+5)
	for _, g := range groups {
		gCopy := g
		// group_winner: 1 view per group, slot 0
		teamViews = append(
			teamViews, &domain.TeamPredictionView{
				Category:    domain.TeamHandicapCategoryGroupWinner,
				GroupLetter: &gCopy,
				SlotIndex:   0,
				Prediction:  byKey[key{domain.TeamHandicapCategoryGroupWinner, g, 0}],
			},
		)
		// playoff: 3 views per group (slots 0, 1, 2)
		for slot := 0; slot <= 2; slot++ {
			gCopy2 := g
			teamViews = append(
				teamViews, &domain.TeamPredictionView{
					Category:    domain.TeamHandicapCategoryPlayoff,
					GroupLetter: &gCopy2,
					SlotIndex:   slot,
					Prediction:  byKey[key{domain.TeamHandicapCategoryPlayoff, g, slot}],
				},
			)
		}
	}
	// semifinalist: 4 views (slots 0-3, group=nil)
	for slot := 0; slot <= 3; slot++ {
		teamViews = append(
			teamViews, &domain.TeamPredictionView{
				Category:    domain.TeamHandicapCategorySemifinalist,
				GroupLetter: nil,
				SlotIndex:   slot,
				Prediction:  byKey[key{domain.TeamHandicapCategorySemifinalist, "", slot}],
			},
		)
	}
	// winner: 1 view (slot 0, group=nil)
	teamViews = append(
		teamViews, &domain.TeamPredictionView{
			Category:    domain.TeamHandicapCategoryWinner,
			GroupLetter: nil,
			SlotIndex:   0,
			Prediction:  byKey[key{domain.TeamHandicapCategoryWinner, "", 0}],
		},
	)
	return locked, teamViews, nil
}

// checkLeagueMembership verifies the caller is a member of the league and that
// predictions are locked. Returns (league, nil) on success.
func (s *TournamentPredictionService) checkLeagueMembership(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
) (*domain.League, error) {
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

	locked, err := s.tournaments.IsPredictionsLocked(ctx, league.TournamentID)
	if err != nil {
		return nil, fmt.Errorf("check lock: %w", err)
	}
	if !locked {
		return nil, fmt.Errorf("league predictions not available until lock time has passed: %w", domain.ErrForbidden)
	}
	return league, nil
}

// ListLeagueGroupPredictions returns group-stage predictions for all league members
// in the given group: group_winner team picks and group_top_scorer player picks.
// Requires the caller to be a league member and predictions to be locked.
func (s *TournamentPredictionService) ListLeagueGroupPredictions(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
	groupLetter string,
) (*domain.LeagueGroupPredictions, error) {
	if _, err := s.checkLeagueMembership(ctx, leagueID, requestingUserID); err != nil {
		return nil, err
	}

	members, err := s.leagues.ListMembersForPredictions(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	teamPicks, err := s.teamPredictions.ListTeamsByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list team predictions: %w", err)
	}

	playerPicks, err := s.playerPredictions.ListPlayersByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list player predictions: %w", err)
	}

	sorted := sortedMembers(members, requestingUserID)

	teamViews := buildGroupTeamViews(sorted, teamPicks, groupLetter)
	playerViews := buildGroupPlayerViews(sorted, playerPicks, groupLetter)

	return &domain.LeagueGroupPredictions{
		Group:             groupLetter,
		TeamPredictions:   teamViews,
		PlayerPredictions: playerViews,
	}, nil
}

// ListLeaguePlayoffPredictions returns knockout/outright predictions for all
// league members: semifinalist and winner team picks, and total_top_scorer player picks.
// Requires the caller to be a league member and predictions to be locked.
func (s *TournamentPredictionService) ListLeaguePlayoffPredictions(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
) (*domain.LeaguePlayoffPredictions, error) {
	if _, err := s.checkLeagueMembership(ctx, leagueID, requestingUserID); err != nil {
		return nil, err
	}

	members, err := s.leagues.ListMembersForPredictions(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}

	teamPicks, err := s.teamPredictions.ListTeamsByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list team predictions: %w", err)
	}

	playerPicks, err := s.playerPredictions.ListPlayersByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list player predictions: %w", err)
	}

	sorted := sortedMembers(members, requestingUserID)

	teamViews := buildPlayoffTeamViews(sorted, teamPicks)
	playerViews := buildPlayoffPlayerViews(sorted, playerPicks)

	return &domain.LeaguePlayoffPredictions{
		TeamPredictions:   teamViews,
		PlayerPredictions: playerViews,
	}, nil
}

// buildGroupTeamViews builds group_winner and playoff (slots 0-2) category views for a specific group.
func buildGroupTeamViews(
	sorted []*domain.LeagueMemberDisplay,
	picks []*domain.TeamLeaguePick,
	groupLetter string,
) []*domain.LeagueTeamCategoryView {
	type key struct {
		userID    uuid.UUID
		category  domain.TeamHandicapCategory
		slotIndex int
	}
	pickMap := make(map[key]*domain.TeamLeaguePick, len(picks))
	for _, p := range picks {
		if p.GroupLetter == nil || *p.GroupLetter != groupLetter {
			continue
		}
		if p.Category != domain.TeamHandicapCategoryGroupWinner && p.Category != domain.TeamHandicapCategoryPlayoff {
			continue
		}
		pickMap[key{p.UserID, p.Category, p.SlotIndex}] = p
	}

	views := make([]*domain.LeagueTeamCategoryView, 0, 4)

	// group_winner: slot 0
	{
		gCopy := groupLetter
		memberPicks := make([]domain.LeagueMemberTeamPick, 0, len(sorted))
		for _, m := range sorted {
			pick := pickMap[key{m.UserID, domain.TeamHandicapCategoryGroupWinner, 0}]
			mp := domain.LeagueMemberTeamPick{UserID: m.UserID, DisplayName: m.DisplayName}
			if pick != nil {
				mp.TeamID = &pick.TeamID
				mp.TeamName = &pick.TeamName
				mp.Points = pick.Points
			}
			memberPicks = append(memberPicks, mp)
		}
		views = append(
			views, &domain.LeagueTeamCategoryView{
				Category:    domain.TeamHandicapCategoryGroupWinner,
				GroupLetter: &gCopy,
				SlotIndex:   0,
				Predictions: memberPicks,
			},
		)
	}

	// playoff: slots 0, 1, 2
	for slot := 0; slot <= 2; slot++ {
		gCopy := groupLetter
		memberPicks := make([]domain.LeagueMemberTeamPick, 0, len(sorted))
		for _, m := range sorted {
			pick := pickMap[key{m.UserID, domain.TeamHandicapCategoryPlayoff, slot}]
			mp := domain.LeagueMemberTeamPick{UserID: m.UserID, DisplayName: m.DisplayName}
			if pick != nil {
				mp.TeamID = &pick.TeamID
				mp.TeamName = &pick.TeamName
				mp.Points = pick.Points
			}
			memberPicks = append(memberPicks, mp)
		}
		views = append(
			views, &domain.LeagueTeamCategoryView{
				Category:    domain.TeamHandicapCategoryPlayoff,
				GroupLetter: &gCopy,
				SlotIndex:   slot,
				Predictions: memberPicks,
			},
		)
	}

	return views
}

// buildGroupPlayerViews builds the group_top_scorer category view for a specific group.
func buildGroupPlayerViews(
	sorted []*domain.LeagueMemberDisplay,
	picks []*domain.PlayerLeaguePick,
	groupLetter string,
) []*domain.LeaguePlayerCategoryView {
	type key struct{ userID uuid.UUID }
	pickMap := make(map[key]*domain.PlayerLeaguePick, len(picks))
	for _, p := range picks {
		if p.Category == domain.PlayerHandicapCategoryGroupTopScorer &&
			p.GroupLetter != nil && *p.GroupLetter == groupLetter {
			pickMap[key{p.UserID}] = p
		}
	}

	gCopy := groupLetter
	memberPicks := make([]domain.LeagueMemberPlayerPick, 0, len(sorted))
	for _, m := range sorted {
		pick := pickMap[key{m.UserID}]
		mp := domain.LeagueMemberPlayerPick{UserID: m.UserID, DisplayName: m.DisplayName}
		if pick != nil {
			mp.PlayerID = &pick.PlayerID
			mp.PlayerName = &pick.PlayerName
			mp.Points = pick.Points
		}
		memberPicks = append(memberPicks, mp)
	}
	return []*domain.LeaguePlayerCategoryView{
		{
			Category:    domain.PlayerHandicapCategoryGroupTopScorer,
			GroupLetter: &gCopy,
			Predictions: memberPicks,
		},
	}
}

// buildPlayoffTeamViews builds category views for semifinalist and winner.
func buildPlayoffTeamViews(
	sorted []*domain.LeagueMemberDisplay,
	picks []*domain.TeamLeaguePick,
) []*domain.LeagueTeamCategoryView {
	type key struct {
		userID      uuid.UUID
		category    domain.TeamHandicapCategory
		groupLetter string
		slotIndex   int
	}
	pickMap := make(map[key]*domain.TeamLeaguePick, len(picks))
	for _, p := range picks {
		gl := ""
		if p.GroupLetter != nil {
			gl = *p.GroupLetter
		}
		pickMap[key{p.UserID, p.Category, gl, p.SlotIndex}] = p
	}

	views := make([]*domain.LeagueTeamCategoryView, 0, 5)

	// semifinalist: slots 0-3
	for slot := 0; slot <= 3; slot++ {
		memberPicks := make([]domain.LeagueMemberTeamPick, 0, len(sorted))
		for _, m := range sorted {
			pick := pickMap[key{m.UserID, domain.TeamHandicapCategorySemifinalist, "", slot}]
			mp := domain.LeagueMemberTeamPick{UserID: m.UserID, DisplayName: m.DisplayName}
			if pick != nil {
				mp.TeamID = &pick.TeamID
				mp.TeamName = &pick.TeamName
				mp.Points = pick.Points
			}
			memberPicks = append(memberPicks, mp)
		}
		views = append(
			views, &domain.LeagueTeamCategoryView{
				Category:    domain.TeamHandicapCategorySemifinalist,
				GroupLetter: nil,
				SlotIndex:   slot,
				Predictions: memberPicks,
			},
		)
	}

	// winner: slot 0
	{
		memberPicks := make([]domain.LeagueMemberTeamPick, 0, len(sorted))
		for _, m := range sorted {
			pick := pickMap[key{m.UserID, domain.TeamHandicapCategoryWinner, "", 0}]
			mp := domain.LeagueMemberTeamPick{UserID: m.UserID, DisplayName: m.DisplayName}
			if pick != nil {
				mp.TeamID = &pick.TeamID
				mp.TeamName = &pick.TeamName
				mp.Points = pick.Points
			}
			memberPicks = append(memberPicks, mp)
		}
		views = append(
			views, &domain.LeagueTeamCategoryView{
				Category:    domain.TeamHandicapCategoryWinner,
				GroupLetter: nil,
				SlotIndex:   0,
				Predictions: memberPicks,
			},
		)
	}

	return views
}

// buildPlayoffPlayerViews builds the total_top_scorer category view.
func buildPlayoffPlayerViews(
	sorted []*domain.LeagueMemberDisplay,
	picks []*domain.PlayerLeaguePick,
) []*domain.LeaguePlayerCategoryView {
	type key struct{ userID uuid.UUID }
	pickMap := make(map[key]*domain.PlayerLeaguePick, len(picks))
	for _, p := range picks {
		if p.Category == domain.PlayerHandicapCategoryTotalTopScorer {
			pickMap[key{p.UserID}] = p
		}
	}

	memberPicks := make([]domain.LeagueMemberPlayerPick, 0, len(sorted))
	for _, m := range sorted {
		pick := pickMap[key{m.UserID}]
		mp := domain.LeagueMemberPlayerPick{UserID: m.UserID, DisplayName: m.DisplayName}
		if pick != nil {
			mp.PlayerID = &pick.PlayerID
			mp.PlayerName = &pick.PlayerName
			mp.Points = pick.Points
		}
		memberPicks = append(memberPicks, mp)
	}
	return []*domain.LeaguePlayerCategoryView{
		{
			Category:    domain.PlayerHandicapCategoryTotalTopScorer,
			GroupLetter: nil,
			Predictions: memberPicks,
		},
	}
}

// compile-time interface checks
var (
	_ PlayerPredictionRepo = (*repository.PlayerPredictionRepository)(nil)
	_ TeamPredictionRepo   = (*repository.TeamPredictionRepository)(nil)
	_ PlayerGetter         = (*repository.PlayerRepository)(nil)
	_ TeamGetter           = (*repository.TeamRepository)(nil)
	_ TeamGroupLister      = (*repository.TeamRepository)(nil)
	_ TournamentLockReader = (*repository.TournamentRepository)(nil)
	_ LeagueReader         = (*repository.LeagueRepository)(nil)
)

// sortedMembers returns members sorted with requestingUserID first, then alphabetically.
func sortedMembers(
	members []*domain.LeagueMemberDisplay,
	requestingUserID uuid.UUID,
) []*domain.LeagueMemberDisplay {
	sorted := make([]*domain.LeagueMemberDisplay, len(members))
	copy(sorted, members)
	sort.SliceStable(
		sorted, func(i, j int) bool {
			iIsMe := sorted[i].UserID == requestingUserID
			jIsMe := sorted[j].UserID == requestingUserID
			if iIsMe != jIsMe {
				return iIsMe
			}
			return sorted[i].DisplayName < sorted[j].DisplayName
		},
	)
	return sorted
}
