package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/worker"
)

// UpdateMatchAndRescoreLivePredictions updates the fixture's status/goals/winner and atomically
// rescores all score_predictions for that fixture.
func (r *WorkerRepository) UpdateMatchAndRescoreLivePredictions(
	ctx context.Context,
	f domain.PollableFixture,
	result worker.APIFixtureResult,
) error {
	newStatus := mapAPIStatusForDB(result.StatusShort)
	winnerTeamID := winnerTeamIDFromResult(result, f)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// 1. Update the fixture row.
	const updateFixture = `
UPDATE fixtures
SET status            = $1,
    goals_home        = $2,
    goals_away        = $3,
    winner_team_id    = $4,
    prediction_locked = (kickoff_at <= now()),
    last_polled_at    = now()
WHERE id = $5`

	if _, err := tx.Exec(ctx, updateFixture,
		newStatus, result.GoalsHome, result.GoalsAway, winnerTeamID, f.ID,
	); err != nil {
		return fmt.Errorf("update fixture %s: %w", f.ID, err)
	}

	// 2. Rescore score_predictions.
	if newStatus == "cancelled" {
		const zeroPredictions = `
UPDATE score_predictions
SET points    = 0,
    scored_at = now()
WHERE fixture_id = $1`

		if _, err := tx.Exec(ctx, zeroPredictions, f.ID); err != nil {
			return fmt.Errorf("zero score predictions for fixture %s: %w", f.ID, err)
		}
	} else if result.GoalsHome != nil && result.GoalsAway != nil {
		gh := int32(*result.GoalsHome) //nolint:gosec
		ga := int32(*result.GoalsAway) //nolint:gosec

		const scorePredictions = `
UPDATE score_predictions
SET points =
      (CASE WHEN goals_home = $1                                           THEN 1 ELSE 0 END)
    + (CASE WHEN goals_away = $2                                           THEN 1 ELSE 0 END)
    + (CASE WHEN SIGN(goals_home::int - goals_away::int)
                 = SIGN($1::int - $2::int)                                 THEN 2 ELSE 0 END)
    + (CASE WHEN goals_home::int - goals_away::int = $1::int - $2::int     THEN 2 ELSE 0 END),
    scored_at = now()
WHERE fixture_id = $3`

		if _, err := tx.Exec(ctx, scorePredictions, gh, ga, f.ID); err != nil {
			return fmt.Errorf("score predictions for fixture %s: %w", f.ID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit update match: %w", err)
	}
	return nil
}

// SettleGroupWinnerPredictions awards points for group_winner predictions once the group concludes.
func (r *WorkerRepository) SettleGroupWinnerPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Fetch the group winner (1st place only).
	const getWinner = `
SELECT team_id FROM tournament_group_table
WHERE tournament_id = $1 AND group_letter = $2
ORDER BY position ASC
LIMIT 1`

	var pos1 uuid.UUID
	if err := tx.QueryRow(ctx, getWinner, tournamentID, groupLetter).Scan(&pos1); err != nil {
		return fmt.Errorf("get group winner team: %w", err)
	}

	// Record outcome.
	const recordOutcome = `
INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'group_winner', $2)
ON CONFLICT (tournament_id, category, team_id) DO UPDATE SET recorded_at = now()`

	if _, err := tx.Exec(ctx, recordOutcome, tournamentID, pos1); err != nil {
		return fmt.Errorf("record group winner outcome: %w", err)
	}

	// Award handicap to the correct prediction.
	const awardPoints = `
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'group_winner'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'group_winner'
  AND tp.group_letter  = $2
  AND tp.pick          = $3
  AND tp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID, groupLetter, pos1); err != nil {
		return fmt.Errorf("award group winner points: %w", err)
	}

	// Zero remaining unsettled rows.
	const zeroRest = `
UPDATE team_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'group_winner'
  AND group_letter  = $2
  AND points IS NULL`

	if _, err := tx.Exec(ctx, zeroRest, tournamentID, groupLetter); err != nil {
		return fmt.Errorf("zero remaining group winner predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle group winner: %w", err)
	}
	return nil
}

// SettlePlayoffGroupPredictions awards points for playoff predictions for teams from a specific group.
// Qualifying teams are identified by a description starting with "Promotion" in tournament_group_table,
// which is written by UpdateGroupStandings from the API standings response.
func (r *WorkerRepository) SettlePlayoffGroupPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const getQualifiers = `
SELECT team_id FROM tournament_group_table
WHERE tournament_id = $1
  AND group_letter  = $2
  AND description   LIKE 'Promotion%'
ORDER BY position ASC`

	rows, err := tx.Query(ctx, getQualifiers, tournamentID, groupLetter)
	if err != nil {
		return fmt.Errorf("get playoff qualifiers for group %s: %w", groupLetter, err)
	}
	defer rows.Close()

	var qualifiers []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan qualifier team id: %w", err)
		}
		qualifiers = append(qualifiers, id)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate playoff qualifiers: %w", err)
	}
	if len(qualifiers) == 0 {
		return fmt.Errorf("no playoff qualifiers found for group %s of tournament %s", groupLetter, tournamentID)
	}

	for _, teamID := range qualifiers {
		const recordOutcome = `
INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'playoff', $2)
ON CONFLICT (tournament_id, category, team_id) DO UPDATE SET recorded_at = now()`

		if _, err := tx.Exec(ctx, recordOutcome, tournamentID, teamID); err != nil {
			return fmt.Errorf("record playoff outcome for team %s: %w", teamID, err)
		}
	}

	const awardPoints = `
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'playoff'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'playoff'
  AND tp.group_letter  = $2
  AND tp.slot_index   IN (0, 1)
  AND tp.pick IN (SELECT team_id FROM tournament_group_table
                  WHERE tournament_id = $1 AND group_letter = $2 AND description LIKE 'Promotion%')
  AND tp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID, groupLetter); err != nil {
		return fmt.Errorf("award playoff group points: %w", err)
	}

	const zeroRest = `
UPDATE team_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'playoff'
  AND group_letter  = $2
  AND slot_index   IN (0, 1)
  AND points IS NULL`

	if _, err := tx.Exec(ctx, zeroRest, tournamentID, groupLetter); err != nil {
		return fmt.Errorf("zero remaining playoff group predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle playoff group: %w", err)
	}
	return nil
}

// SettlePlayoffWildcardPredictions awards points for wildcard (slot_index=2) playoff predictions
// after all groups have concluded. Advancing teams are already in team_outcomes(category='playoff').
func (r *WorkerRepository) SettlePlayoffWildcardPredictions(ctx context.Context, tournamentID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const awardPoints = `
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'playoff'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'playoff'
  AND tp.group_letter IS NULL
  AND tp.slot_index    = 2
  AND tp.pick IN (SELECT team_id FROM team_outcomes
                  WHERE tournament_id = $1 AND category = 'playoff')
  AND tp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID); err != nil {
		return fmt.Errorf("award wildcard playoff points: %w", err)
	}

	const zeroRest = `
UPDATE team_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'playoff'
  AND group_letter IS NULL
  AND slot_index    = 2
  AND points IS NULL`

	if _, err := tx.Exec(ctx, zeroRest, tournamentID); err != nil {
		return fmt.Errorf("zero remaining wildcard predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle wildcard: %w", err)
	}
	return nil
}

// SettleGroupTopScorerPredictions awards points for group_top_scorer predictions.
func (r *WorkerRepository) SettleGroupTopScorerPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string, topScorerPlayerID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const recordOutcome = `
INSERT INTO player_outcomes (tournament_id, category, player_id)
VALUES ($1, 'group_top_scorer', $2)
ON CONFLICT (tournament_id, category, player_id) DO UPDATE SET recorded_at = now()`

	if _, err := tx.Exec(ctx, recordOutcome, tournamentID, topScorerPlayerID); err != nil {
		return fmt.Errorf("record group top scorer outcome: %w", err)
	}

	const awardPoints = `
UPDATE player_predictions pp
SET points    = COALESCE((SELECT h.points FROM player_handicap h
                          WHERE h.player_id = pp.pick AND h.category = 'group_top_scorer'), 0),
    scored_at = now()
WHERE pp.tournament_id = $1
  AND pp.category      = 'group_top_scorer'
  AND pp.group_letter  = $2
  AND pp.pick          = $3
  AND pp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID, groupLetter, topScorerPlayerID); err != nil {
		return fmt.Errorf("award group top scorer points: %w", err)
	}

	const zeroRest = `
UPDATE player_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'group_top_scorer'
  AND group_letter  = $2
  AND pick         != $3
  AND points IS NULL`

	if _, err := tx.Exec(ctx, zeroRest, tournamentID, groupLetter, topScorerPlayerID); err != nil {
		return fmt.Errorf("zero remaining group top scorer predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle group top scorer: %w", err)
	}
	return nil
}

// SettleSemifinalistPredictions awards points once all quarterfinals have concluded.
func (r *WorkerRepository) SettleSemifinalistPredictions(ctx context.Context, tournamentID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const getSFTeams = `
SELECT winner_team_id FROM fixtures
WHERE tournament_id = $1
  AND LOWER(round) LIKE '%quarter%'
  AND winner_team_id IS NOT NULL`

	rows, err := tx.Query(ctx, getSFTeams, tournamentID)
	if err != nil {
		return fmt.Errorf("get semifinalist teams: %w", err)
	}
	defer rows.Close()

	var sfTeams []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("scan semifinalist team id: %w", err)
		}
		sfTeams = append(sfTeams, id)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate semifinalist teams: %w", err)
	}
	for _, teamID := range sfTeams {
		const recordOutcome = `
INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'semifinalist', $2)
ON CONFLICT (tournament_id, category, team_id) DO UPDATE SET recorded_at = now()`

		if _, err := tx.Exec(ctx, recordOutcome, tournamentID, teamID); err != nil {
			return fmt.Errorf("record semifinalist outcome for team %s: %w", teamID, err)
		}
	}

	const awardPoints = `
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'semifinalist'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'semifinalist'
  AND tp.pick IN (SELECT team_id FROM team_outcomes
                  WHERE tournament_id = $1 AND category = 'semifinalist')
  AND tp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID); err != nil {
		return fmt.Errorf("award semifinalist points: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle semifinalists: %w", err)
	}
	return nil
}

// ZeroRemainingSemifinalistPredictions zeroes all unsettled semifinalist predictions
// once the full quarterfinal round is complete and all winners are known.
func (r *WorkerRepository) ZeroRemainingSemifinalistPredictions(ctx context.Context, tournamentID uuid.UUID) error {
	const q = `
UPDATE team_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'semifinalist'
  AND points IS NULL`

	if _, err := r.pool.Exec(ctx, q, tournamentID); err != nil {
		return fmt.Errorf("zero remaining semifinalist predictions: %w", err)
	}
	return nil
}

// SettleTournamentWinnerPredictions awards points for the tournament winner prediction.
func (r *WorkerRepository) SettleTournamentWinnerPredictions(ctx context.Context, tournamentID uuid.UUID, winnerTeamID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const recordOutcome = `
INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'winner', $2)
ON CONFLICT (tournament_id, category, team_id) DO UPDATE SET recorded_at = now()`

	if _, err := tx.Exec(ctx, recordOutcome, tournamentID, winnerTeamID); err != nil {
		return fmt.Errorf("record winner outcome: %w", err)
	}

	const awardPoints = `
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'winner'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'winner'
  AND tp.pick          = $2
  AND tp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID, winnerTeamID); err != nil {
		return fmt.Errorf("award winner points: %w", err)
	}

	const zeroRest = `
UPDATE team_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'winner'
  AND pick         != $2
  AND points IS NULL`

	if _, err := tx.Exec(ctx, zeroRest, tournamentID, winnerTeamID); err != nil {
		return fmt.Errorf("zero remaining winner predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle tournament winner: %w", err)
	}
	return nil
}

// SettleTopScorerPredictions awards points for the total_top_scorer prediction.
func (r *WorkerRepository) SettleTopScorerPredictions(ctx context.Context, tournamentID uuid.UUID, topScorerPlayerID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const recordOutcome = `
INSERT INTO player_outcomes (tournament_id, category, player_id)
VALUES ($1, 'total_top_scorer', $2)
ON CONFLICT (tournament_id, category, player_id) DO UPDATE SET recorded_at = now()`

	if _, err := tx.Exec(ctx, recordOutcome, tournamentID, topScorerPlayerID); err != nil {
		return fmt.Errorf("record top scorer outcome: %w", err)
	}

	const awardPoints = `
UPDATE player_predictions pp
SET points    = COALESCE((SELECT h.points FROM player_handicap h
                          WHERE h.player_id = pp.pick AND h.category = 'total_top_scorer'), 0),
    scored_at = now()
WHERE pp.tournament_id = $1
  AND pp.category      = 'total_top_scorer'
  AND pp.pick          = $2
  AND pp.points IS NULL`

	if _, err := tx.Exec(ctx, awardPoints, tournamentID, topScorerPlayerID); err != nil {
		return fmt.Errorf("award top scorer points: %w", err)
	}

	const zeroRest = `
UPDATE player_predictions
SET points    = 0,
    scored_at = now()
WHERE tournament_id = $1
  AND category      = 'total_top_scorer'
  AND pick         != $2
  AND points IS NULL`

	if _, err := tx.Exec(ctx, zeroRest, tournamentID, topScorerPlayerID); err != nil {
		return fmt.Errorf("zero remaining top scorer predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit settle top scorer: %w", err)
	}
	return nil
}
