package worker

import (
	"context"
	"log/slog"

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

// updateGroupStandings fetches live standings from the API and writes them to
// the DB. Called on every result change for group-stage fixtures so the group
// table stays current during a match, not just at full time.
func (w *Worker) updateGroupStandings(ctx context.Context, f domain.PollableFixture) {
	standings, err := w.api.GetStandings(ctx, f.TournamentExternalID, f.TournamentSeason)
	if err != nil {
		w.logger.Warn("worker: get standings", "fixture_id", f.ID, "err", err)
		return
	}
	entries := w.resolveTeamIDs(ctx, standings, f.TournamentID)
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

		topScorer, err := w.api.GetGroupTopScorer(ctx, f.TournamentExternalID, f.TournamentSeason, *f.GroupLetter)
		if err != nil {
			w.logger.Warn("worker: get group top scorer", "fixture_id", f.ID, "group", *f.GroupLetter, "err", err)
		} else {
			playerID, err := w.repo.GetPlayerByExternalID(ctx, topScorer.PlayerExternalID, f.TournamentID)
			if err != nil {
				w.logger.Warn(
					"worker: resolve group top scorer player",
					"external_id",
					topScorer.PlayerExternalID,
					"err",
					err,
				)
			} else {
				if err := w.repo.SettleGroupTopScorerPredictions(
					ctx,
					f.TournamentID,
					*f.GroupLetter,
					playerID,
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

		topScorer, err := w.api.GetTournamentTopScorer(ctx, f.TournamentExternalID, f.TournamentSeason)
		if err != nil {
			w.logger.Warn("worker: get tournament top scorer", "fixture_id", f.ID, "err", err)
		} else {
			playerID, err := w.repo.GetPlayerByExternalID(ctx, topScorer.PlayerExternalID, f.TournamentID)
			if err != nil {
				w.logger.Warn(
					"worker: resolve tournament top scorer player",
					"external_id",
					topScorer.PlayerExternalID,
					"err",
					err,
				)
			} else {
				if err := w.repo.SettleTopScorerPredictions(ctx, f.TournamentID, playerID); err != nil {
					w.logger.Error("worker: settle top scorer predictions", "fixture_id", f.ID, "err", err)
				}
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

// resolveTeamIDs converts a slice of API standings entries to domain entries by
// looking up each team's internal UUID. Entries whose external ID is unknown are skipped.
func (w *Worker) resolveTeamIDs(
	ctx context.Context,
	standings []APIStandingsEntry,
	tournamentID uuid.UUID,
) []domain.StandingsEntry {
	out := make([]domain.StandingsEntry, 0, len(standings))
	for _, s := range standings {
		teamID, err := w.repo.GetTeamByExternalID(ctx, s.TeamExternalID, tournamentID)
		if err != nil {
			w.logger.Error(
				"worker: resolve team external id", "external_id", s.TeamExternalID, "err", err,
				slog.String("hint", "populate teams.external_id"),
			)
			continue
		}
		out = append(
			out, domain.StandingsEntry{
				TeamID:       teamID,
				Position:     s.Position,
				Points:       s.Points,
				Played:       s.Played,
				Won:          s.Won,
				Drawn:        s.Drawn,
				Lost:         s.Lost,
				GoalsFor:     s.GoalsFor,
				GoalsAgainst: s.GoalsAgainst,
				Description:  s.Description,
			},
		)
	}
	return out
}
