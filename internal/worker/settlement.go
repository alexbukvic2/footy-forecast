package worker

import (
	"context"

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

		topScorers, err := w.api.GetGroupTopScorer(ctx, f.TournamentExternalID, f.TournamentSeason, *f.GroupLetter)
		if err != nil {
			w.logger.Warn("worker: get group top scorer", "fixture_id", f.ID, "group", *f.GroupLetter, "err", err)
		} else {
			playerIDs := w.resolveTopScorerPlayerIDs(ctx, topScorers, f.TournamentID)
			if len(playerIDs) > 0 {
				if err := w.repo.SettleGroupTopScorerPredictions(
					ctx,
					f.TournamentID,
					*f.GroupLetter,
					playerIDs,
				); err != nil {
					w.logger.Error("worker: settle group top scorer predictions", "fixture_id", f.ID, "err", err)
				}
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

		topScorers, err := w.api.GetTournamentTopScorer(ctx, f.TournamentExternalID, f.TournamentSeason)
		if err != nil {
			w.logger.Warn("worker: get tournament top scorer", "fixture_id", f.ID, "err", err)
		} else {
			playerIDs := w.resolveTopScorerPlayerIDs(ctx, topScorers, f.TournamentID)
			if len(playerIDs) > 0 {
				if err := w.repo.SettleTopScorerPredictions(ctx, f.TournamentID, playerIDs); err != nil {
					w.logger.Error("worker: settle top scorer predictions", "fixture_id", f.ID, "err", err)
				}
			}
		}
	}

	// Any playoff fixture finishing may reveal new fixtures in the API (e.g. next-round draw).
	w.refreshFixturesForTournament(ctx, f.TournamentExternalID, f.TournamentSeason, f.TournamentID)
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

// resolveTopScorerPlayerIDs maps APITopScorerResult entries to internal player UUIDs.
// Entries whose external ID is unknown are logged and skipped.
func (w *Worker) resolveTopScorerPlayerIDs(
	ctx context.Context,
	scorers []APITopScorerResult,
	tournamentID uuid.UUID,
) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(scorers))
	for _, ts := range scorers {
		playerID, err := w.repo.GetPlayerByExternalID(ctx, ts.PlayerExternalID, tournamentID)
		if err != nil {
			w.logger.Warn("worker: resolve top scorer player", "external_id", ts.PlayerExternalID, "err", err)
			continue
		}
		ids = append(ids, playerID)
	}
	return ids
}
