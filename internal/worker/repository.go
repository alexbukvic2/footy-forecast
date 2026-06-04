package worker

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// Repo is the data access interface consumed by the worker.
// Implemented in internal/repository.
type Repo interface {
	// Polling
	LockImminentFixtures(ctx context.Context, leadMinutes int) error
	ListPollableMatches(ctx context.Context) ([]domain.PollableFixture, error)
	UpdateMatchAndRescoreLivePredictions(ctx context.Context, f domain.PollableFixture, result APIFixtureResult) error

	// Group standings
	UpdateGroupStandings(ctx context.Context, tournamentID uuid.UUID, groupLetter string, entries []domain.StandingsEntry) error

	// Completion checks
	IsGroupComplete(ctx context.Context, tournamentID uuid.UUID, groupLetter string) (bool, error)
	IsRoundComplete(ctx context.Context, tournamentID uuid.UUID, round string) (bool, error)
	IsGroupStageComplete(ctx context.Context, tournamentID uuid.UUID) (bool, error)

	// ID resolution
	GetTeamByExternalID(ctx context.Context, externalID int64, tournamentID uuid.UUID) (uuid.UUID, error)
	GetPlayerByExternalID(ctx context.Context, externalID int64, tournamentID uuid.UUID) (uuid.UUID, error)

	// Settlement
	SettleGroupWinnerPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string) error
	SettlePlayoffGroupPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string) error
	SettlePlayoffWildcardPredictions(ctx context.Context, tournamentID uuid.UUID) error
	SettleGroupTopScorerPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string, topScorerPlayerIDs []uuid.UUID) error
	SettleSemifinalistPredictions(ctx context.Context, tournamentID uuid.UUID) error
	ZeroRemainingSemifinalistPredictions(ctx context.Context, tournamentID uuid.UUID) error
	SettleTournamentWinnerPredictions(ctx context.Context, tournamentID uuid.UUID, winnerTeamID uuid.UUID) error
	SettleTopScorerPredictions(ctx context.Context, tournamentID uuid.UUID, topScorerPlayerIDs []uuid.UUID) error
}
