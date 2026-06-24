package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
	"github.com/alexbukvic2/footy-forecast/internal/worker"
)

// WorkerRepository implements worker.Repo.
type WorkerRepository struct {
	pool *db.Pool
	q    *dbgen.Queries
}

// NewWorkerRepository constructs a WorkerRepository.
func NewWorkerRepository(pool *db.Pool) *WorkerRepository {
	return &WorkerRepository{pool: pool, q: dbgen.New(pool)}
}

var _ worker.Repo = (*WorkerRepository)(nil)

// LockImminentFixtures sets prediction_locked = TRUE for every upcoming, non-demo fixture
// whose kickoff is within the next leadMinutes minutes. If the locked fixture is the first
// (earliest kickoff) fixture of its tournament, also sets tournaments.predictions_locked = TRUE.
func (r *WorkerRepository) LockImminentFixtures(
	ctx context.Context,
	leadMinutes int,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const lockFixtures = `
UPDATE fixtures
SET prediction_locked = TRUE
WHERE prediction_locked = FALSE
  AND is_demo          = FALSE
  AND status           = 'upcoming'
  AND kickoff_at       <= now() + ($1 * INTERVAL '1 minute')`
	if _, err := tx.Exec(ctx, lockFixtures, leadMinutes); err != nil {
		return fmt.Errorf("lock imminent fixtures: %w", err)
	}

	// Lock any tournament whose first (earliest kickoff) non-demo fixture is now locked.
	const lockTournaments = `
UPDATE tournaments
SET predictions_locked = TRUE
WHERE predictions_locked = FALSE
  AND EXISTS (
    SELECT 1
    FROM fixtures f
    WHERE f.tournament_id = tournaments.id
      AND f.is_demo = FALSE
      AND f.prediction_locked = TRUE
      AND f.kickoff_at = (
        SELECT MIN(f2.kickoff_at)
        FROM fixtures f2
        WHERE f2.tournament_id = tournaments.id
          AND f2.is_demo = FALSE
      )
  )`
	if _, err := tx.Exec(ctx, lockTournaments); err != nil {
		return fmt.Errorf("lock tournament predictions: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit lock: %w", err)
	}
	return nil
}

// ListPollableMatches returns fixtures that need polling: live or imminent (≤5 min).
// Excludes demo fixtures and tournaments without an external_id.
func (r *WorkerRepository) ListPollableMatches(ctx context.Context) ([]domain.PollableFixture, error) {
	rows, err := r.q.ListPollableMatches(ctx)
	if err != nil {
		return nil, fmt.Errorf("list pollable matches: %w", err)
	}
	out := make([]domain.PollableFixture, 0, len(rows))
	for _, row := range rows {
		f, err := pollableFixtureFromRow(row)
		if err != nil {
			return nil, fmt.Errorf("map pollable fixture: %w", err)
		}
		out = append(out, f)
	}
	return out, nil
}

// IsGroupComplete returns true when every fixture in the group has a terminal status.
func (r *WorkerRepository) IsGroupComplete(
	ctx context.Context,
	tournamentID uuid.UUID,
	groupLetter string,
) (bool, error) {
	n, err := r.q.CountIncompleteGroupFixtures(
		ctx, dbgen.CountIncompleteGroupFixturesParams{
			TournamentID: tournamentID,
			GroupLetter:  &groupLetter,
		},
	)
	if err != nil {
		return false, fmt.Errorf("count incomplete group fixtures: %w", err)
	}
	return n == 0, nil
}

// IsRoundComplete returns true when every fixture in the round has a winner_team_id set.
func (r *WorkerRepository) IsRoundComplete(
	ctx context.Context,
	tournamentID uuid.UUID,
	round string,
) (bool, error) {
	n, err := r.q.CountIncompleteRoundFixtures(
		ctx, dbgen.CountIncompleteRoundFixturesParams{
			TournamentID: tournamentID,
			Round:        round,
		},
	)
	if err != nil {
		return false, fmt.Errorf("count incomplete round fixtures: %w", err)
	}
	return n == 0, nil
}

// IsGroupStageComplete returns true when every group-stage fixture has a terminal status.
func (r *WorkerRepository) IsGroupStageComplete(
	ctx context.Context,
	tournamentID uuid.UUID,
) (bool, error) {
	n, err := r.q.CountIncompleteGroupStageFixtures(ctx, tournamentID)
	if err != nil {
		return false, fmt.Errorf("count incomplete group stage fixtures: %w", err)
	}
	return n == 0, nil
}

// GetTeamByExternalID resolves an api-sports.io team ID to our internal team UUID.
func (r *WorkerRepository) GetTeamByExternalID(
	ctx context.Context,
	externalID int64,
	tournamentID uuid.UUID,
) (uuid.UUID, error) {
	id, err := r.q.GetTeamByExternalID(
		ctx, dbgen.GetTeamByExternalIDParams{
			ExternalID:   &externalID,
			TournamentID: tournamentID,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf(
				"team external_id %d tournament %s: %w",
				externalID,
				tournamentID,
				domain.ErrNotFound,
			)
		}
		return uuid.UUID{}, fmt.Errorf("get team by external id: %w", err)
	}
	return id, nil
}

// GetPlayerByExternalID resolves an api-sports.io player ID to our internal player UUID.
func (r *WorkerRepository) GetPlayerByExternalID(
	ctx context.Context,
	externalID int64,
	tournamentID uuid.UUID,
) (uuid.UUID, error) {
	id, err := r.q.GetPlayerByExternalID(
		ctx, dbgen.GetPlayerByExternalIDParams{
			ExternalID:   externalID,
			TournamentID: tournamentID,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf(
				"player external_id %d tournament %s: %w",
				externalID,
				tournamentID,
				domain.ErrNotFound,
			)
		}
		return uuid.UUID{}, fmt.Errorf("get player by external id: %w", err)
	}
	return id, nil
}

func pollableFixtureFromRow(row dbgen.ListPollableMatchesRow) (domain.PollableFixture, error) {
	// tournament_external_id is guaranteed non-NULL by the WHERE clause, but sqlc types it nullable.
	if row.TournamentExternalID == nil {
		return domain.PollableFixture{}, fmt.Errorf("unexpected nil tournament_external_id for fixture %s", row.ID)
	}

	f := domain.PollableFixture{
		ID:                   row.ID,
		ExternalID:           row.ExternalID,
		TournamentID:         row.TournamentID,
		TournamentExternalID: *row.TournamentExternalID,
		TournamentSeason:     0,
		HomeTeamID:           row.HomeTeamID,
		AwayTeamID:           row.AwayTeamID,
		// group_letter is meaningful only for group-stage fixtures; clear it for knockout rounds.
		GroupLetter: groupLetterForRound(row.Round, row.GroupLetter),
		Round:       row.Round,
		Status:      domain.FixtureStatus(row.Status),
		KickoffAt:   row.KickoffAt,
	}

	if row.TournamentSeason != nil {
		f.TournamentSeason = int(*row.TournamentSeason)
	}

	if row.GoalsHome != nil {
		v := int(*row.GoalsHome)
		f.GoalsHome = &v
	}
	if row.GoalsAway != nil {
		v := int(*row.GoalsAway)
		f.GoalsAway = &v
	}

	// winner_team_id and last_polled_at: NULL scans to zero value for these types.
	if row.WinnerTeamID != (uuid.UUID{}) {
		id := row.WinnerTeamID
		f.WinnerTeamID = &id
	}
	if row.LastPolledAt.Valid {
		t := row.LastPolledAt.Time
		f.LastPolledAt = &t
	}

	return f, nil
}

func groupLetterForRound(
	round string,
	gl *string,
) *string {
	if strings.HasPrefix(round, "Group") {
		return gl
	}
	return nil
}

// mapAPIStatusForDB converts an api-sports.io status string to the DB enum string.
func mapAPIStatusForDB(short string) string {
	switch short {
	case "1H", "HT", "2H", "ET", "BT", "P", "SUSP", "INT":
		return "in_progress"
	case "FT", "AET", "PEN", "AWD", "WO":
		return "finished"
	case "CANC", "ABD":
		return "cancelled"
	default:
		return "upcoming"
	}
}

// winnerTeamIDFromResult derives the winner_team_id UUID (or nil) from the API result.
func winnerTeamIDFromResult(
	result worker.APIFixtureResult,
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

// ListActiveTournaments returns all tournaments that have an external API ID and season set.
func (r *WorkerRepository) ListActiveTournaments(ctx context.Context) ([]domain.ActiveTournament, error) {
	rows, err := r.q.ListActiveTournaments(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active tournaments: %w", err)
	}
	out := make([]domain.ActiveTournament, 0, len(rows))
	for _, row := range rows {
		if row.ExternalID == nil || row.Season == nil {
			continue
		}
		out = append(
			out, domain.ActiveTournament{
				ID:         row.ID,
				ExternalID: *row.ExternalID,
				Season:     int(*row.Season),
			},
		)
	}
	return out, nil
}

// InsertMissingFixtures inserts fixtures that are not yet in the DB.
// Existing fixtures (matched by external_id) are left untouched.
func (r *WorkerRepository) InsertMissingFixtures(
	ctx context.Context,
	tournamentID uuid.UUID,
	fixtures []domain.NewFixture,
) error {
	const q = `
INSERT INTO fixtures (external_id, tournament_id, home_team_id, away_team_id, kickoff_at, status, round, goals_home, goals_away)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (external_id) DO NOTHING`
	for _, f := range fixtures {
		if _, err := r.pool.Exec(
			ctx, q,
			f.ExternalID, tournamentID, f.HomeTeamID, f.AwayTeamID,
			f.KickoffAt, string(f.Status), f.Round, f.GoalsHome, f.GoalsAway,
		); err != nil {
			return fmt.Errorf("insert fixture external_id %d: %w", f.ExternalID, err)
		}
	}
	return nil
}

// UpsertPlayerGoalsByExternalID looks up the player by external ID and tournament,
// then adds goals to their running total in a single round trip.
// Returns domain.ErrNotFound if no matching player exists.
func (r *WorkerRepository) UpsertPlayerGoalsByExternalID(
	ctx context.Context,
	externalID int64,
	tournamentID uuid.UUID,
	goals int,
) error {
	const q = `
INSERT INTO players_stats (player_id, goals)
SELECT id, $3
FROM players
WHERE external_id = $1
  AND tournament_id = $2
ON CONFLICT (player_id) DO UPDATE SET goals = players_stats.goals + EXCLUDED.goals`
	tag, err := r.pool.Exec(ctx, q, externalID, tournamentID, goals)
	if err != nil {
		return fmt.Errorf("upsert player goals external_id %d: %w", externalID, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("upsert player goals external_id %d tournament %s: %w", externalID, tournamentID, domain.ErrNotFound)
	}
	return nil
}

// GetGroupTopScorerPlayerIDs returns internal player UUIDs tied at the highest goal
// count among players in the given group. Returns empty slice when no goals are recorded.
func (r *WorkerRepository) GetGroupTopScorerPlayerIDs(
	ctx context.Context,
	tournamentID uuid.UUID,
	groupLetter string,
) ([]uuid.UUID, error) {
	const q = `
SELECT ps.player_id
FROM players_stats ps
JOIN players p ON p.id = ps.player_id
JOIN teams   t ON t.id = p.team_id
WHERE p.tournament_id = $1
  AND t.group_letter  = $2
  AND ps.goals = (
    SELECT MAX(ps2.goals)
    FROM players_stats ps2
    JOIN players p2 ON p2.id = ps2.player_id
    JOIN teams   t2 ON t2.id = p2.team_id
    WHERE p2.tournament_id = $1
      AND t2.group_letter  = $2
  )
ORDER BY ps.player_id`
	rows, err := r.pool.Query(ctx, q, tournamentID, groupLetter)
	if err != nil {
		return nil, fmt.Errorf("get group top scorer player ids: %w", err)
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan group top scorer player id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// GetTournamentTopScorerPlayerIDs returns internal player UUIDs tied at the highest
// goal count in the tournament. Returns empty slice when no goals are recorded.
func (r *WorkerRepository) GetTournamentTopScorerPlayerIDs(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]uuid.UUID, error) {
	const q = `
SELECT ps.player_id
FROM players_stats ps
JOIN players p ON p.id = ps.player_id
WHERE p.tournament_id = $1
  AND ps.goals = (
    SELECT MAX(ps2.goals)
    FROM players_stats ps2
    JOIN players p2 ON p2.id = ps2.player_id
    WHERE p2.tournament_id = $1
  )
ORDER BY ps.player_id`
	rows, err := r.pool.Query(ctx, q, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("get tournament top scorer player ids: %w", err)
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan tournament top scorer player id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ListGroupTeams returns the IDs of all teams assigned to a given group letter in a tournament.
func (r *WorkerRepository) ListGroupTeams(
	ctx context.Context,
	tournamentID uuid.UUID,
	groupLetter string,
) ([]uuid.UUID, error) {
	const q = `SELECT id FROM teams WHERE tournament_id = $1 AND group_letter = $2 ORDER BY id`
	rows, err := r.pool.Query(ctx, q, tournamentID, groupLetter)
	if err != nil {
		return nil, fmt.Errorf("list group teams: %w", err)
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan group team id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ListGroupFixtures returns all group-stage fixtures for a tournament + group that have goals set.
func (r *WorkerRepository) ListGroupFixtures(
	ctx context.Context,
	tournamentID uuid.UUID,
	groupLetter string,
) ([]domain.GroupFixture, error) {
	const q = `
SELECT f.home_team_id, f.away_team_id, f.goals_home, f.goals_away
FROM fixtures f
JOIN teams t ON t.id = f.home_team_id
WHERE f.tournament_id = $1
  AND t.group_letter   = $2
  AND f.goals_home    IS NOT NULL
  AND f.goals_away    IS NOT NULL`
	rows, err := r.pool.Query(ctx, q, tournamentID, groupLetter)
	if err != nil {
		return nil, fmt.Errorf("list group fixtures: %w", err)
	}
	defer rows.Close()
	var out []domain.GroupFixture
	for rows.Next() {
		var (
			gf        domain.GroupFixture
			goalsHome int32
			goalsAway int32
		)
		if err := rows.Scan(&gf.HomeTeamID, &gf.AwayTeamID, &goalsHome, &goalsAway); err != nil {
			return nil, fmt.Errorf("scan group fixture: %w", err)
		}
		gf.GoalsHome = int(goalsHome)
		gf.GoalsAway = int(goalsAway)
		out = append(out, gf)
	}
	return out, rows.Err()
}

// UpdateGroupStandings upserts standings rows for all teams in a group.
func (r *WorkerRepository) UpdateGroupStandings(
	ctx context.Context,
	tournamentID uuid.UUID,
	groupLetter string,
	entries []domain.StandingsEntry,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const q = `
INSERT INTO tournament_group_table
    (tournament_id, team_id, group_letter, position, points, played, won, drawn, lost, goals_for, goals_against, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (tournament_id, team_id) DO UPDATE SET
    group_letter  = EXCLUDED.group_letter,
    position      = EXCLUDED.position,
    points        = EXCLUDED.points,
    played        = EXCLUDED.played,
    won           = EXCLUDED.won,
    drawn         = EXCLUDED.drawn,
    lost          = EXCLUDED.lost,
    goals_for     = EXCLUDED.goals_for,
    goals_against = EXCLUDED.goals_against,
    description   = EXCLUDED.description`

	for _, e := range entries {
		if e.Group == nil || *e.Group != groupLetter {
			continue
		}
		if _, err := tx.Exec(
			ctx, q,
			tournamentID, e.TeamID, groupLetter,
			int16(e.Position), int16(e.Points), int16(e.Played), //nolint:gosec
			int16(e.Won), int16(e.Drawn), int16(e.Lost), //nolint:gosec
			int16(e.GoalsFor), int16(e.GoalsAgainst), //nolint:gosec
			e.Description,
		); err != nil {
			return fmt.Errorf("upsert standings for team %s: %w", e.TeamID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit standings: %w", err)
	}
	return nil
}
