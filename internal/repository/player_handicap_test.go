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

func TestPlayerHandicapRepository_Get(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	_, err := pool.Exec(ctx,
		`INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1, 'ph-test', 'PH Test', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID,
	)
	require.NoError(t, err)

	teamID := insertTeam(t, pool, ctx, "PH Team", "", tournamentID)

	playerID := uuid.Must(uuid.NewV7())
	_, err = pool.Exec(ctx,
		`INSERT INTO players (id, external_id, name, tournament_id, team_id) VALUES ($1, 1, 'PH Player', $2, $3)`,
		playerID, tournamentID, teamID,
	)
	require.NoError(t, err)

	repo := repository.NewPlayerHandicapRepository(pool)

	t.Run("returns points for existing row", func(t *testing.T) {
		_, err := pool.Exec(ctx,
			`INSERT INTO player_handicap (player_id, category, points) VALUES ($1, 'group_top_scorer', 5)`,
			playerID,
		)
		require.NoError(t, err)

		h, err := repo.Get(ctx, playerID, domain.PlayerHandicapCategoryGroupTopScorer)
		require.NoError(t, err)
		require.Equal(t, 5, h.Points)
		require.Equal(t, playerID, h.PlayerID)
		require.Equal(t, domain.PlayerHandicapCategoryGroupTopScorer, h.Category)
	})

	t.Run("returns ErrNotFound for unknown player", func(t *testing.T) {
		_, err := repo.Get(ctx, uuid.New(), domain.PlayerHandicapCategoryGroupTopScorer)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("returns ErrNotFound for wrong category", func(t *testing.T) {
		// group_top_scorer row exists for playerID; total_top_scorer does not
		_, err := repo.Get(ctx, playerID, domain.PlayerHandicapCategoryTotalTopScorer)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("uniqueness constraint prevents duplicate player+category", func(t *testing.T) {
		_, err := pool.Exec(ctx,
			`INSERT INTO player_handicap (player_id, category, points) VALUES ($1, 'group_top_scorer', 10)`,
			playerID,
		)
		require.Error(t, err)
	})
}
