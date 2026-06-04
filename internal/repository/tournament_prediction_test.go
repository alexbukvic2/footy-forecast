//go:build integration

package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestPlayerPredictionRepository(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	// Seed: tournament, teams, players, user, league
	tournamentID := uuid.Must(uuid.NewV7())
	teamID := uuid.Must(uuid.NewV7())
	playerID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	user2ID := uuid.Must(uuid.NewV7())
	leagueID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'pp-test-cup', 'PP Test Cup', '2026-06-01', '2026-07-01')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, tournament_id) VALUES ($1, 'Team A', $2)`,
		teamID, tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO players (id, external_id, name, tournament_id, team_id)
		VALUES ($1, 1, 'Lionel Messi', $2, $3)`,
		playerID, tournamentID, teamID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, cognito_sub, email, display_name) VALUES
		($1, 'sub-pp-1', 'ppuser1@example.com', 'Alice'),
		($2, 'sub-pp-2', 'ppuser2@example.com', 'Bob')`,
		userID, user2ID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO leagues (id, tournament_id, owner_id, name, code)
		VALUES ($1, $2, $3, 'Test League', 'PPTEST01')`,
		leagueID, tournamentID, userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO league_members (league_id, user_id, role) VALUES
		($1, $2, 'owner'), ($1, $3, 'member')`,
		leagueID, userID, user2ID)
	require.NoError(t, err)

	repo := repository.NewPlayerPredictionRepository(pool)

	t.Run("upsert creates new row", func(t *testing.T) {
		pred, err := repo.UpsertPlayer(ctx, domain.UpsertPlayerPredictionInput{
			UserID:       userID,
			TournamentID: tournamentID,
			Category:     domain.PlayerHandicapCategoryGroupTopScorer,
			Pick:         playerID,
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, pred.ID)
		require.Equal(t, userID, pred.UserID)
		require.Equal(t, tournamentID, pred.TournamentID)
		require.Equal(t, domain.PlayerHandicapCategoryGroupTopScorer, pred.Category)
		require.Equal(t, playerID, pred.Pick)
		require.Equal(t, "Lionel Messi", pred.PickName)
		require.Nil(t, pred.Points)
	})

	t.Run("upsert with same category updates pick, leaves points unchanged", func(t *testing.T) {
		// First set points manually to verify they're preserved.
		_, err := pool.Exec(ctx, `
			UPDATE player_predictions SET points = 10
			WHERE user_id = $1 AND tournament_id = $2 AND category = 'group_top_scorer'`,
			userID, tournamentID)
		require.NoError(t, err)

		// Second player in same tournament.
		player2ID := uuid.Must(uuid.NewV7())
		_, err = pool.Exec(ctx, `
			INSERT INTO players (id, external_id, name, tournament_id, team_id)
			VALUES ($1, 2, 'Cristiano Ronaldo', $2, $3)`,
			player2ID, tournamentID, teamID)
		require.NoError(t, err)

		pred, err := repo.UpsertPlayer(ctx, domain.UpsertPlayerPredictionInput{
			UserID:       userID,
			TournamentID: tournamentID,
			Category:     domain.PlayerHandicapCategoryGroupTopScorer,
			Pick:         player2ID,
		})
		require.NoError(t, err)
		require.Equal(t, player2ID, pred.Pick)
		require.Equal(t, "Cristiano Ronaldo", pred.PickName)
		require.NotNil(t, pred.Points)
		require.Equal(t, 10, *pred.Points, "points should be unchanged")
	})

	t.Run("list by tournament for user returns only that user's predictions", func(t *testing.T) {
		// Insert a prediction for user2 to ensure it's filtered out.
		_, err := repo.UpsertPlayer(ctx, domain.UpsertPlayerPredictionInput{
			UserID:       user2ID,
			TournamentID: tournamentID,
			Category:     domain.PlayerHandicapCategoryTotalTopScorer,
			Pick:         playerID,
		})
		require.NoError(t, err)

		rows, err := repo.ListPlayersByTournamentForUser(ctx, tournamentID, userID)
		require.NoError(t, err)
		for _, r := range rows {
			require.Equal(t, userID, r.UserID, "should only return caller's predictions")
		}
	})

	t.Run("list by tournament for user returns empty slice when none exist", func(t *testing.T) {
		unknownUserID := uuid.Must(uuid.NewV7())
		rows, err := repo.ListPlayersByTournamentForUser(ctx, tournamentID, unknownUserID)
		require.NoError(t, err)
		require.Empty(t, rows)
	})

	t.Run("list by league returns rows for members who predicted", func(t *testing.T) {
		rows, err := repo.ListPlayersByLeague(ctx, leagueID)
		require.NoError(t, err)
		require.NotEmpty(t, rows)
		for _, r := range rows {
			require.True(t, r.UserID == userID || r.UserID == user2ID)
		}
	})
}

func TestTeamPredictionRepository(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	teamID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	user2ID := uuid.Must(uuid.NewV7())
	leagueID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'tp-test-cup', 'TP Test Cup', '2026-06-01', '2026-07-01')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, tournament_id) VALUES ($1, 'Team B', $2)`,
		teamID, tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, cognito_sub, email, display_name) VALUES
		($1, 'sub-tp-1', 'tpuser1@example.com', 'Alice'),
		($2, 'sub-tp-2', 'tpuser2@example.com', 'Bob')`,
		userID, user2ID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO leagues (id, tournament_id, owner_id, name, code)
		VALUES ($1, $2, $3, 'TP League', 'TPTEST01')`,
		leagueID, tournamentID, userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO league_members (league_id, user_id, role) VALUES
		($1, $2, 'owner'), ($1, $3, 'member')`,
		leagueID, userID, user2ID)
	require.NoError(t, err)

	repo := repository.NewTeamPredictionRepository(pool)

	t.Run("upsert creates new row", func(t *testing.T) {
		pred, err := repo.UpsertTeam(ctx, domain.UpsertTeamPredictionInput{
			UserID:       userID,
			TournamentID: tournamentID,
			Category:     domain.TeamHandicapCategoryWinner,
			Pick:         teamID,
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, pred.ID)
		require.Equal(t, teamID, pred.Pick)
		require.Equal(t, "Team B", pred.PickName)
		require.Nil(t, pred.Points)
	})

	t.Run("upsert with same category updates pick, leaves points unchanged", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			UPDATE team_predictions SET points = 20
			WHERE user_id = $1 AND tournament_id = $2 AND category = 'winner'`,
			userID, tournamentID)
		require.NoError(t, err)

		team2ID := uuid.Must(uuid.NewV7())
		_, err = pool.Exec(ctx, `
			INSERT INTO teams (id, name, tournament_id) VALUES ($1, 'Team C', $2)`,
			team2ID, tournamentID)
		require.NoError(t, err)

		pred, err := repo.UpsertTeam(ctx, domain.UpsertTeamPredictionInput{
			UserID:       userID,
			TournamentID: tournamentID,
			Category:     domain.TeamHandicapCategoryWinner,
			Pick:         team2ID,
		})
		require.NoError(t, err)
		require.Equal(t, team2ID, pred.Pick)
		require.NotNil(t, pred.Points)
		require.Equal(t, 20, *pred.Points, "points should be unchanged")
	})

	t.Run("list by tournament for user returns empty when none", func(t *testing.T) {
		unknownUser := uuid.Must(uuid.NewV7())
		rows, err := repo.ListTeamsByTournamentForUser(ctx, tournamentID, unknownUser)
		require.NoError(t, err)
		require.Empty(t, rows)
	})

	t.Run("list by league returns rows for members who predicted", func(t *testing.T) {
		rows, err := repo.ListTeamsByLeague(ctx, leagueID)
		require.NoError(t, err)
		require.NotEmpty(t, rows)
		for _, r := range rows {
			require.True(t, r.UserID == userID || r.UserID == user2ID)
		}
	})
}

func TestFixtureRepository_GetFirstKickoff(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	teamA := uuid.Must(uuid.NewV7())
	teamB := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'fixture-test-cup', 'Fixture Test Cup', '2026-06-01', '2026-07-01')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, tournament_id) VALUES
		($1, 'TeamA', $3), ($2, 'TeamB', $3)`,
		teamA, teamB, tournamentID)
	require.NoError(t, err)

	repo := repository.NewFixtureRepository(pool)

	t.Run("returns ErrNotFound when no fixtures exist", func(t *testing.T) {
		_, err := repo.GetFirstKickoffByTournament(ctx, tournamentID)
		require.True(t, errors.Is(err, domain.ErrNotFound), "got: %v", err)
	})

	kickoff1 := time.Date(2026, 6, 10, 18, 0, 0, 0, time.UTC)
	kickoff2 := time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC)

	_, err = pool.Exec(ctx, `
		INSERT INTO fixtures (external_id, tournament_id, home_team_id, away_team_id, kickoff_at)
		VALUES (1001, $1, $2, $3, $4), (1002, $1, $3, $2, $5)`,
		tournamentID, teamA, teamB, kickoff1, kickoff2)
	require.NoError(t, err)

	t.Run("returns earliest kickoff", func(t *testing.T) {
		got, err := repo.GetFirstKickoffByTournament(ctx, tournamentID)
		require.NoError(t, err)
		require.True(t, got.Equal(kickoff1), "want %v got %v", kickoff1, got)
	})
}
