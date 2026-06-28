// Package worker runs the background match-polling and scoring loop.
package worker

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/google/uuid"
)

// Worker polls live match data and scores predictions in real time.
type Worker struct {
	repo            Repo
	api             MatchAPI
	clock           Clock
	logger          *slog.Logger
	pollInterval    time.Duration
	lockLeadMinutes int
}

// New constructs a Worker. pollInterval controls how often the polling loop runs.
// lockLeadMinutes is how many minutes before kickoff predictions are locked.
func New(
	repo Repo,
	api MatchAPI,
	clock Clock,
	logger *slog.Logger,
	pollInterval time.Duration,
	lockLeadMinutes int,
) *Worker {
	return &Worker{
		repo:            repo,
		api:             api,
		clock:           clock,
		logger:          logger,
		pollInterval:    pollInterval,
		lockLeadMinutes: lockLeadMinutes,
	}
}

// Run starts the polling loop. It blocks until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	w.refreshAllActiveTournamentFixtures(ctx)
	tournamentID, _ := uuid.Parse("6c4ef9ed-a44a-4dc6-8429-4ee0f191c27e")
	if err := w.repo.SettlePlayoffWildcardPredictions(ctx, tournamentID); err != nil {
		w.logger.Error("worker: settle playoff wildcard predictions", "fixture_id", tournamentID, "err", err)
	}

	groupLetter := "J"
	if err := w.repo.SettleGroupWinnerPredictions(ctx, tournamentID, groupLetter); err != nil {
		w.logger.Error("worker: settle group winner predictions", "err", err)
	}
	if err := w.repo.SettlePlayoffGroupPredictions(ctx, tournamentID, groupLetter); err != nil {
		w.logger.Error("worker: settle playoff group predictions", "err", err)
	}

	playerIDs, err := w.repo.GetGroupTopScorerPlayerIDs(ctx, tournamentID, groupLetter)
	if err != nil {
		w.logger.Warn("worker: get group top scorer player ids", "group", groupLetter, "err", err)
	} else if len(playerIDs) > 0 {
		if err := w.repo.SettleGroupTopScorerPredictions(
			ctx,
			tournamentID,
			groupLetter,
			playerIDs,
		); err != nil {
			w.logger.Error("worker: settle group top scorer predictions", "err", err)
		}
	}

	for {
		w.tick(ctx)

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(w.pollInterval):
		}
	}
}

func (w *Worker) tick(ctx context.Context) {
	if err := w.repo.LockImminentFixtures(ctx, w.lockLeadMinutes); err != nil {
		w.logger.Warn("worker: lock imminent fixtures", "err", err)
	}

	fixtures, err := w.repo.ListPollableMatches(ctx)
	if err != nil {
		w.logger.Warn("worker: list pollable matches", "err", err)
		return
	}
	if len(fixtures) == 0 {
		w.logger.Info("worker idle")
		return
	}

	sem := newSemaphore(5)
	for _, f := range fixtures {
		if err := sem.acquire(ctx); err != nil {
			return
		}
		go func(f domain.PollableFixture) {
			defer sem.release()
			w.processSingleFixture(ctx, f)
		}(f)
	}
	sem.drain()
}

func (w *Worker) processSingleFixture(
	ctx context.Context,
	f domain.PollableFixture,
) {
	result, err := w.api.GetFixture(ctx, f.ExternalID)
	if err != nil {
		w.logger.Warn("worker: get fixture", "fixture_id", f.ID, "external_id", f.ExternalID, "err", err)
		return
	}

	newStatus := MapAPIStatus(result.StatusShort)

	if isUnchanged(f, result, newStatus) {
		return
	}

	if err := w.repo.UpdateMatchAndRescoreLivePredictions(ctx, f, result); err != nil {
		w.logger.Error("worker: update match and rescore", "fixture_id", f.ID, "err", err)
		return
	}

	if f.GroupLetter != nil {
		w.updateGroupStandings(ctx, f)
	}

	w.runSettlement(ctx, f, result, newStatus)
}

// MapAPIStatus converts an api-sports.io status.short string to a domain FixtureStatus.
func MapAPIStatus(short string) domain.FixtureStatus {
	switch short {
	case "1H", "HT", "2H", "ET", "BT", "P", "SUSP", "INT":
		return domain.FixtureStatusInProgress
	case "FT", "AET", "PEN", "AWD", "WO":
		return domain.FixtureStatusFinished
	case "CANC", "ABD":
		return domain.FixtureStatusCancelled
	default:
		return domain.FixtureStatusUpcoming
	}
}

func isUnchanged(
	f domain.PollableFixture,
	result APIFixtureResult,
	newStatus domain.FixtureStatus,
) bool {
	if newStatus == domain.FixtureStatusCancelled {
		return false
	}
	statusSame := newStatus == f.Status
	goalsSame := intPtrEq(result.GoalsHome, f.GoalsHome) && intPtrEq(result.GoalsAway, f.GoalsAway)
	return statusSame && goalsSame
}

func intPtrEq(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// isQuarterfinalRound returns true if the round string indicates a quarterfinal.
func isQuarterfinalRound(round string) bool {
	return strings.Contains(strings.ToLower(round), "quarter")
}

// isFinalRound returns true if the round string indicates the final (not semi/quarter).
func isFinalRound(round string) bool {
	lower := strings.ToLower(round)
	return strings.HasSuffix(lower, "final") &&
		!strings.Contains(lower, "semi") &&
		!strings.Contains(lower, "quarter")
}
