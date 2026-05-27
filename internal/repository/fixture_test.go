//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestFixtureRepository_GetLockedFixtureDates_HappyPath(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(
		ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'fxdate-test', 'Fixture Date Test', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO users (id, cognito_sub, email) VALUES ($1, 'fxdate-sub-1', 'fxdate1@example.com')`,
		userID,
	)
	require.NoError(t, err)

	leagueID := uuid.Must(uuid.NewV7())
	_, err = pool.Exec(
		ctx, `
		INSERT INTO leagues (id, tournament_id, owner_id, name, code)
		VALUES ($1, $2, $3, 'FXDate League', 'FXDCODE')`,
		leagueID, tournamentID, userID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO league_members (league_id, user_id) VALUES ($1, $2)`,
		leagueID, userID,
	)
	require.NoError(t, err)

	homeTeamID := insertTeam(t, pool, ctx, "Home FXD", "", tournamentID)
	awayTeamID := insertTeam(t, pool, ctx, "Away FXD", "", tournamentID)

	// Three dates: 2 past locked, 1 future unlocked
	past1 := time.Date(2026, 5, 25, 15, 0, 0, 0, time.UTC)
	past2 := time.Date(2026, 5, 24, 15, 0, 0, 0, time.UTC)
	future := time.Date(2030, 6, 1, 15, 0, 0, 0, time.UTC)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO fixtures (tournament_id, home_team_id, away_team_id, round, kickoff_at, status, prediction_locked, external_id)
		VALUES
		($1, $2, $3, 'GS', $4, 'finished', true, 1),
		($1, $2, $3, 'GS', $5, 'finished', true, 2),
		($1, $2, $3, 'GS', $6, 'upcoming', false, 3)`,
		tournamentID, homeTeamID, awayTeamID, past1, past2, future,
	)
	require.NoError(t, err)

	repo := repository.NewFixtureRepository(pool)
	dates, err := repo.GetLockedFixtureDates(ctx, leagueID)
	require.NoError(t, err)

	// Only 2 past locked dates, newest first
	require.Len(t, dates, 2)
	require.Equal(t, past1.UTC().Truncate(24*time.Hour), dates[0].UTC())
	require.Equal(t, past2.UTC().Truncate(24*time.Hour), dates[1].UTC())
}

func TestFixtureRepository_ListLockedByLeagueAndDates_HappyPath(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(
		ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'fxbydates-test', 'FX By Dates Test', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO users (id, cognito_sub, email) VALUES ($1, 'fxbyd-sub-1', 'fxbyd1@example.com')`,
		userID,
	)
	require.NoError(t, err)

	leagueID := uuid.Must(uuid.NewV7())
	_, err = pool.Exec(
		ctx, `
		INSERT INTO leagues (id, tournament_id, owner_id, name, code)
		VALUES ($1, $2, $3, 'FXByDates League', 'FXBYDCODE')`,
		leagueID, tournamentID, userID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO league_members (league_id, user_id) VALUES ($1, $2)`,
		leagueID, userID,
	)
	require.NoError(t, err)

	homeTeamID := insertTeam(t, pool, ctx, "Home FXBD", "", tournamentID)
	awayTeamID := insertTeam(t, pool, ctx, "Away FXBD", "", tournamentID)

	d1 := time.Date(2026, 5, 25, 15, 0, 0, 0, time.UTC)
	d2 := time.Date(2026, 5, 24, 15, 0, 0, 0, time.UTC)
	d3 := time.Date(2026, 5, 23, 15, 0, 0, 0, time.UTC)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO fixtures (tournament_id, home_team_id, away_team_id, round, kickoff_at, status, prediction_locked, external_id)
		VALUES
		($1, $2, $3, 'GS', $4, 'finished', true, 1),
		($1, $2, $3, 'GS', $5, 'finished', true, 2),
		($1, $2, $3, 'GS', $6, 'finished', true, 3)`,
		tournamentID, homeTeamID, awayTeamID, d1, d2, d3,
	)
	require.NoError(t, err)

	repo := repository.NewFixtureRepository(pool)

	// Request only d1 and d2 dates
	views, err := repo.ListLockedByLeagueAndDates(
		ctx, leagueID, userID, []time.Time{
			d1.UTC().Truncate(24 * time.Hour),
			d2.UTC().Truncate(24 * time.Hour),
		},
	)
	require.NoError(t, err)
	require.Len(t, views, 2)
}

func TestFixtureRepository_ListLockedByLeagueAndDates_NoMatch(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(
		ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'fxnomatch-test', 'FX No Match Test', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO users (id, cognito_sub, email) VALUES ($1, 'fxnm-sub-1', 'fxnm1@example.com')`,
		userID,
	)
	require.NoError(t, err)

	leagueID := uuid.Must(uuid.NewV7())
	_, err = pool.Exec(
		ctx, `
		INSERT INTO leagues (id, tournament_id, owner_id, name, code)
		VALUES ($1, $2, $3, 'FXNoMatch League', 'FXNMCODE')`,
		leagueID, tournamentID, userID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO league_members (league_id, user_id) VALUES ($1, $2)`,
		leagueID, userID,
	)
	require.NoError(t, err)

	repo := repository.NewFixtureRepository(pool)

	unrelatedDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	views, err := repo.ListLockedByLeagueAndDates(ctx, leagueID, userID, []time.Time{unrelatedDate})
	require.NoError(t, err)
	require.Empty(t, views)
}
