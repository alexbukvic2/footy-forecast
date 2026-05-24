//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestPredictionRepository_Upsert(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	// Seed prerequisites.
	tournamentID := uuid.Must(uuid.NewV7())
	homeTeamID := uuid.Must(uuid.NewV7())
	awayTeamID := uuid.Must(uuid.NewV7())
	userID := uuid.Must(uuid.NewV7())
	fixtureID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'pred-cup', 'Pred Cup', '2026-06-01', '2026-07-01')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, logo, tournament_id) VALUES ($1, 'Home', '', $3), ($2, 'Away', '', $3)`,
		homeTeamID, awayTeamID, tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, cognito_sub, email) VALUES ($1, 'sub-p', 'p@example.com')`,
		userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO fixtures (id, external_id, tournament_id, home_team_id, away_team_id, kickoff_at, status)
		VALUES ($1, 2001, $2, $3, $4, '2026-06-11T18:00:00Z', 'upcoming')`,
		fixtureID, tournamentID, homeTeamID, awayTeamID)
	require.NoError(t, err)

	repo := repository.NewPredictionRepository(pool)

	t.Run("creates a new prediction", func(t *testing.T) {
		pred, err := repo.Upsert(ctx, domain.UpsertScorePredictionInput{
			UserID:    userID,
			FixtureID: fixtureID,
			GoalsHome: 2,
			GoalsAway: 1,
		})
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, pred.ID)
		require.Equal(t, 2, pred.GoalsHome)
		require.Equal(t, 1, pred.GoalsAway)
		require.Nil(t, pred.Points)
	})

	t.Run("updates goals on second upsert and preserves points nil", func(t *testing.T) {
		_, err := repo.Upsert(ctx, domain.UpsertScorePredictionInput{
			UserID:    userID,
			FixtureID: fixtureID,
			GoalsHome: 3,
			GoalsAway: 0,
		})
		require.NoError(t, err)

		pred, err := repo.Upsert(ctx, domain.UpsertScorePredictionInput{
			UserID:    userID,
			FixtureID: fixtureID,
			GoalsHome: 1,
			GoalsAway: 1,
		})
		require.NoError(t, err)
		require.Equal(t, 1, pred.GoalsHome)
		require.Equal(t, 1, pred.GoalsAway)
		require.Nil(t, pred.Points)
	})
}
