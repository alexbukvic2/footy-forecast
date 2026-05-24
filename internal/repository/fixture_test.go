//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestFixtureRepository(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	// Seed tournament, teams, users, league, and league member.
	tournamentID := uuid.Must(uuid.NewV7())
	homeTeamID := uuid.Must(uuid.NewV7())
	awayTeamID := uuid.Must(uuid.NewV7())
	user1ID := uuid.Must(uuid.NewV7())
	user2ID := uuid.Must(uuid.NewV7())
	leagueID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'test-cup', 'Test Cup', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, logo, tournament_id) VALUES
		($1, 'Home FC', '', $3),
		($2, 'Away FC', '', $3)`,
		homeTeamID, awayTeamID, tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, cognito_sub, email, display_name) VALUES
		($1, 'sub-1', 'user1@example.com', 'User One'),
		($2, 'sub-2', 'user2@example.com', 'User Two')`,
		user1ID, user2ID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES
		($1, $2, $3, 'Test League', 'TESTCODE')`,
		leagueID, tournamentID, user1ID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO league_members (league_id, user_id, role) VALUES
		($1, $2, 'owner'),
		($1, $3, 'member')`,
		leagueID, user1ID, user2ID)
	require.NoError(t, err)

	// Insert fixtures: one upcoming, one finished.
	upcomingID := uuid.Must(uuid.NewV7())
	finishedID := uuid.Must(uuid.NewV7())
	kickoff := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)

	_, err = pool.Exec(ctx, `
		INSERT INTO fixtures (id, external_id, tournament_id, home_team_id, away_team_id, kickoff_at, status)
		VALUES ($1, 1001, $2, $3, $4, $5, 'upcoming')`,
		upcomingID, tournamentID, homeTeamID, awayTeamID, kickoff)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO fixtures (id, external_id, tournament_id, home_team_id, away_team_id, kickoff_at, status, goals_home, goals_away)
		VALUES ($1, 1002, $2, $3, $4, $5, 'finished', 2, 1)`,
		finishedID, tournamentID, homeTeamID, awayTeamID, kickoff.Add(-24*time.Hour))
	require.NoError(t, err)

	fixtureRepo := repository.NewFixtureRepository(pool)

	t.Run("ListByTournamentForUser returns all fixtures for tournament", func(t *testing.T) {
		views, err := fixtureRepo.ListByTournamentForUser(ctx, tournamentID, user1ID)
		require.NoError(t, err)
		require.Len(t, views, 2)
	})

	t.Run("ListByTournamentForUser returns nil prediction when none exists", func(t *testing.T) {
		views, err := fixtureRepo.ListByTournamentForUser(ctx, tournamentID, user1ID)
		require.NoError(t, err)
		for _, v := range views {
			require.Nil(t, v.Prediction)
		}
	})

	t.Run("ListByTournamentForUser returns prediction when one exists", func(t *testing.T) {
		predRepo := repository.NewPredictionRepository(pool)
		_, err := predRepo.Upsert(ctx, domain.UpsertScorePredictionInput{
			UserID:    user1ID,
			FixtureID: upcomingID,
			GoalsHome: 1,
			GoalsAway: 0,
		})
		require.NoError(t, err)

		views, err := fixtureRepo.ListByTournamentForUser(ctx, tournamentID, user1ID)
		require.NoError(t, err)

		var found *domain.UserFixtureView
		for _, v := range views {
			if v.Fixture.ID == upcomingID {
				found = v
			}
		}
		require.NotNil(t, found)
		require.NotNil(t, found.Prediction)
		require.Equal(t, 1, found.Prediction.GoalsHome)
		require.Equal(t, 0, found.Prediction.GoalsAway)
	})

	t.Run("ListLockedByLeague returns only finished/in-progress fixtures", func(t *testing.T) {
		views, err := fixtureRepo.ListLockedByLeague(ctx, leagueID, user1ID)
		require.NoError(t, err)
		require.Len(t, views, 1)
		require.Equal(t, finishedID, views[0].Fixture.ID)
	})

	t.Run("ListLockedByLeague includes all league members in predictions", func(t *testing.T) {
		views, err := fixtureRepo.ListLockedByLeague(ctx, leagueID, user1ID)
		require.NoError(t, err)
		require.Len(t, views[0].Predictions, 2)
	})

	t.Run("ListLockedByLeague puts requesting user first", func(t *testing.T) {
		views, err := fixtureRepo.ListLockedByLeague(ctx, leagueID, user1ID)
		require.NoError(t, err)
		require.Equal(t, user1ID, views[0].Predictions[0].UserID)
	})

	t.Run("ListLockedByLeague includes display_name", func(t *testing.T) {
		views, err := fixtureRepo.ListLockedByLeague(ctx, leagueID, user1ID)
		require.NoError(t, err)
		require.Equal(t, "User One", views[0].Predictions[0].DisplayName)
	})
}
