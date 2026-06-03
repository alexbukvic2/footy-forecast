package repository

import (
	"context"
	"errors"
	"fmt"

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

// ListPollableMatches returns fixtures that need polling: live, imminent (≤5 min), or recently
// finished (≤24 h). Excludes demo fixtures and tournaments without an external_id.
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
func (r *WorkerRepository) IsGroupComplete(ctx context.Context, tournamentID uuid.UUID, groupLetter string) (bool, error) {
	n, err := r.q.CountIncompleteGroupFixtures(ctx, dbgen.CountIncompleteGroupFixturesParams{
		TournamentID: tournamentID,
		GroupLetter:  &groupLetter,
	})
	if err != nil {
		return false, fmt.Errorf("count incomplete group fixtures: %w", err)
	}
	return n == 0, nil
}

// IsRoundComplete returns true when every fixture in the round has a winner_team_id set.
func (r *WorkerRepository) IsRoundComplete(ctx context.Context, tournamentID uuid.UUID, round string) (bool, error) {
	n, err := r.q.CountIncompleteRoundFixtures(ctx, dbgen.CountIncompleteRoundFixturesParams{
		TournamentID: tournamentID,
		Round:        round,
	})
	if err != nil {
		return false, fmt.Errorf("count incomplete round fixtures: %w", err)
	}
	return n == 0, nil
}

// IsGroupStageComplete returns true when every group-stage fixture has a terminal status.
func (r *WorkerRepository) IsGroupStageComplete(ctx context.Context, tournamentID uuid.UUID) (bool, error) {
	n, err := r.q.CountIncompleteGroupStageFixtures(ctx, tournamentID)
	if err != nil {
		return false, fmt.Errorf("count incomplete group stage fixtures: %w", err)
	}
	return n == 0, nil
}

// GetTeamByExternalID resolves an api-sports.io team ID to our internal team UUID.
func (r *WorkerRepository) GetTeamByExternalID(ctx context.Context, externalID int64, tournamentID uuid.UUID) (uuid.UUID, error) {
	id, err := r.q.GetTeamByExternalID(ctx, dbgen.GetTeamByExternalIDParams{
		ExternalID:   &externalID,
		TournamentID: tournamentID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf("team external_id %d tournament %s: %w", externalID, tournamentID, domain.ErrNotFound)
		}
		return uuid.UUID{}, fmt.Errorf("get team by external id: %w", err)
	}
	return id, nil
}

// GetPlayerByExternalID resolves an api-sports.io player ID to our internal player UUID.
func (r *WorkerRepository) GetPlayerByExternalID(ctx context.Context, externalID string, tournamentID uuid.UUID) (uuid.UUID, error) {
	id, err := r.q.GetPlayerByExternalID(ctx, dbgen.GetPlayerByExternalIDParams{
		ExternalID:   externalID,
		TournamentID: tournamentID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf("player external_id %s tournament %s: %w", externalID, tournamentID, domain.ErrNotFound)
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
		GroupLetter:          row.GroupLetter,
		Round:                row.Round,
		Status:               domain.FixtureStatus(row.Status),
		KickoffAt:            row.KickoffAt,
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
	if !row.LastPolledAt.IsZero() {
		t := row.LastPolledAt
		f.LastPolledAt = &t
	}

	return f, nil
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
func winnerTeamIDFromResult(result worker.APIFixtureResult, f domain.PollableFixture) *uuid.UUID {
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

// UpdateGroupStandings upserts standings rows for all teams in a group.
func (r *WorkerRepository) UpdateGroupStandings(ctx context.Context, tournamentID uuid.UUID, groupLetter string, entries []domain.StandingsEntry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const q = `
INSERT INTO tournament_group_table
    (tournament_id, team_id, group_letter, position, points, played, won, drawn, lost, goals_for, goals_against)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (tournament_id, team_id) DO UPDATE SET
    group_letter  = EXCLUDED.group_letter,
    position      = EXCLUDED.position,
    points        = EXCLUDED.points,
    played        = EXCLUDED.played,
    won           = EXCLUDED.won,
    drawn         = EXCLUDED.drawn,
    lost          = EXCLUDED.lost,
    goals_for     = EXCLUDED.goals_for,
    goals_against = EXCLUDED.goals_against`

	for _, e := range entries {
		if _, err := tx.Exec(ctx, q,
			tournamentID, e.TeamID, groupLetter,
			int16(e.Position), int16(e.Points), int16(e.Played), //nolint:gosec
			int16(e.Won), int16(e.Drawn), int16(e.Lost), //nolint:gosec
			int16(e.GoalsFor), int16(e.GoalsAgainst), //nolint:gosec
		); err != nil {
			return fmt.Errorf("upsert standings for team %s: %w", e.TeamID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit standings: %w", err)
	}
	return nil
}
