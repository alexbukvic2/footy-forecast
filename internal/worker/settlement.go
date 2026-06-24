package worker

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// runSettlement fires outright-prediction settlement after a fixture reaches a terminal state.
// It only acts on the first transition to terminal (non-terminal → terminal). Subsequent calls
// are safe because all Settle* functions use WHERE points IS NULL for idempotency.
func (w *Worker) runSettlement(
	ctx context.Context,
	f domain.PollableFixture,
	result APIFixtureResult,
	newStatus domain.FixtureStatus,
) {
	nowTerminal := newStatus == domain.FixtureStatusFinished || newStatus == domain.FixtureStatusCancelled

	if !nowTerminal {
		return
	}

	if newStatus == domain.FixtureStatusCancelled {
		// Score predictions already zeroed by UpdateMatchAndRescoreLivePredictions.
		// No outright settlement for cancelled matches.
		return
	}

	w.updatePlayerGoals(ctx, f, result.GoalScorerExternalIDs)

	winnerTeamID := resolveWinnerTeamID(result, f)

	if f.GroupLetter != nil {
		w.settleGroupMatch(ctx, f, winnerTeamID)
	} else {
		w.settleKnockoutMatch(ctx, f, winnerTeamID)
	}
}

// updateGroupStandings computes standings from our own fixtures and writes them to the DB.
// Called on every result change for group-stage fixtures so the table stays current during a match.
func (w *Worker) updateGroupStandings(
	ctx context.Context,
	f domain.PollableFixture,
) {
	teams, err := w.repo.ListGroupTeams(ctx, f.TournamentID, *f.GroupLetter)
	if err != nil {
		w.logger.Warn("worker: list group teams", "fixture_id", f.ID, "err", err)
		return
	}
	fixtures, err := w.repo.ListGroupFixtures(ctx, f.TournamentID, *f.GroupLetter)
	if err != nil {
		w.logger.Warn("worker: list group fixtures", "fixture_id", f.ID, "err", err)
		return
	}
	entries := computeGroupStandings(*f.GroupLetter, teams, fixtures)
	if err := w.repo.UpdateGroupStandings(ctx, f.TournamentID, *f.GroupLetter, entries); err != nil {
		w.logger.Error("worker: update group standings", "fixture_id", f.ID, "err", err)
	}
}

func (w *Worker) settleGroupMatch(
	ctx context.Context,
	f domain.PollableFixture,
	_ *uuid.UUID,
) {
	groupDone, err := w.repo.IsGroupComplete(ctx, f.TournamentID, *f.GroupLetter)
	if err != nil {
		w.logger.Error("worker: is group complete", "fixture_id", f.ID, "err", err)
		return
	}

	if groupDone {
		if err := w.repo.SettleGroupWinnerPredictions(ctx, f.TournamentID, *f.GroupLetter); err != nil {
			w.logger.Error("worker: settle group winner predictions", "fixture_id", f.ID, "err", err)
		}
		if err := w.repo.SettlePlayoffGroupPredictions(ctx, f.TournamentID, *f.GroupLetter); err != nil {
			w.logger.Error("worker: settle playoff group predictions", "fixture_id", f.ID, "err", err)
		}

		playerIDs, err := w.repo.GetGroupTopScorerPlayerIDs(ctx, f.TournamentID, *f.GroupLetter)
		if err != nil {
			w.logger.Warn("worker: get group top scorer player ids", "fixture_id", f.ID, "group", *f.GroupLetter, "err", err)
		} else if len(playerIDs) > 0 {
			if err := w.repo.SettleGroupTopScorerPredictions(
				ctx,
				f.TournamentID,
				*f.GroupLetter,
				playerIDs,
			); err != nil {
				w.logger.Error("worker: settle group top scorer predictions", "fixture_id", f.ID, "err", err)
			}
		}

		allGroupsDone, err := w.repo.IsGroupStageComplete(ctx, f.TournamentID)
		if err != nil {
			w.logger.Error("worker: is group stage complete", "fixture_id", f.ID, "err", err)
			return
		}
		if allGroupsDone {
			if err := w.repo.SettlePlayoffWildcardPredictions(ctx, f.TournamentID); err != nil {
				w.logger.Error("worker: settle playoff wildcard predictions", "fixture_id", f.ID, "err", err)
			}
		}

		// Refresh fixtures so that knockout-draw matches added by the API are picked up.
		w.refreshFixturesForTournament(ctx, f.TournamentExternalID, f.TournamentSeason, f.TournamentID)
	}
}

func (w *Worker) settleKnockoutMatch(
	ctx context.Context,
	f domain.PollableFixture,
	winnerTeamID *uuid.UUID,
) {
	if isQuarterfinalRound(f.Round) {
		// Award points for the newly known semifinalist immediately after each QF.
		if err := w.repo.SettleSemifinalistPredictions(ctx, f.TournamentID); err != nil {
			w.logger.Error("worker: settle semifinalist predictions", "fixture_id", f.ID, "err", err)
		}

		// Zero remaining unsettled predictions only once all QFs are complete.
		roundDone, err := w.repo.IsRoundComplete(ctx, f.TournamentID, f.Round)
		if err != nil {
			w.logger.Error("worker: is round complete", "fixture_id", f.ID, "round", f.Round, "err", err)
		} else if roundDone {
			if err := w.repo.ZeroRemainingSemifinalistPredictions(ctx, f.TournamentID); err != nil {
				w.logger.Error("worker: zero remaining semifinalist predictions", "fixture_id", f.ID, "err", err)
			}
		}
	}

	if isFinalRound(f.Round) && winnerTeamID != nil {
		if err := w.repo.SettleTournamentWinnerPredictions(ctx, f.TournamentID, *winnerTeamID); err != nil {
			w.logger.Error("worker: settle tournament winner predictions", "fixture_id", f.ID, "err", err)
		}

		playerIDs, err := w.repo.GetTournamentTopScorerPlayerIDs(ctx, f.TournamentID)
		if err != nil {
			w.logger.Warn("worker: get tournament top scorer player ids", "fixture_id", f.ID, "err", err)
		} else if len(playerIDs) > 0 {
			if err := w.repo.SettleTopScorerPredictions(ctx, f.TournamentID, playerIDs); err != nil {
				w.logger.Error("worker: settle top scorer predictions", "fixture_id", f.ID, "err", err)
			}
		}
	}

	// Any playoff fixture finishing may reveal new fixtures in the API (e.g. next-round draw).
	w.refreshFixturesForTournament(ctx, f.TournamentExternalID, f.TournamentSeason, f.TournamentID)
}

// updatePlayerGoals records goals scored in a finished fixture into players_stats.
// Each entry in goalScorerExternalIDs is one goal (a player appears N times if they scored N goals).
func (w *Worker) updatePlayerGoals(
	ctx context.Context,
	f domain.PollableFixture,
	goalScorerExternalIDs []int64,
) {
	if len(goalScorerExternalIDs) == 0 {
		return
	}
	counts := make(map[int64]int, len(goalScorerExternalIDs))
	for _, id := range goalScorerExternalIDs {
		counts[id]++
	}
	for externalID, goals := range counts {
		err := w.repo.UpsertPlayerGoalsByExternalID(ctx, externalID, f.TournamentID, goals)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				w.logger.Warn("worker: goal scorer not found", "external_id", externalID, "tournament_id", f.TournamentID)
			} else {
				w.logger.Error("worker: upsert player goals", "external_id", externalID, "err", err)
			}
		}
	}
}

// resolveWinnerTeamID returns the winning team's UUID based on the API result.
func resolveWinnerTeamID(
	result APIFixtureResult,
	f domain.PollableFixture,
) *uuid.UUID {
	if result.HomeWinner != nil && *result.HomeWinner {
		id := f.HomeTeamID
		return &id
	}
	if result.AwayWinner != nil && *result.AwayWinner {
		id := f.AwayTeamID
		return &id
	}
	return nil
}
