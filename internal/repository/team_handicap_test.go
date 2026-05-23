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

func TestTeamHandicapRepository_Get(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	_, err := pool.Exec(ctx,
		`INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1, 'th-test', 'TH Test', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentID,
	)
	require.NoError(t, err)

	teamID := insertTeam(t, pool, ctx, "TH Team", "", tournamentID)

	repo := repository.NewTeamHandicapRepository(pool)

	t.Run("returns points for existing row", func(t *testing.T) {
		_, err := pool.Exec(ctx,
			`INSERT INTO team_handicap (team_id, category, points) VALUES ($1, 'winner', 20)`,
			teamID,
		)
		require.NoError(t, err)

		h, err := repo.Get(ctx, teamID, domain.TeamHandicapCategoryWinner)
		require.NoError(t, err)
		require.Equal(t, 20, h.Points)
		require.Equal(t, teamID, h.TeamID)
		require.Equal(t, domain.TeamHandicapCategoryWinner, h.Category)
	})

	t.Run("returns ErrNotFound for unknown team", func(t *testing.T) {
		_, err := repo.Get(ctx, uuid.New(), domain.TeamHandicapCategoryWinner)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("returns ErrNotFound for wrong category", func(t *testing.T) {
		// winner row exists for teamID; group_winner does not
		_, err := repo.Get(ctx, teamID, domain.TeamHandicapCategoryGroupWinner)
		require.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("uniqueness constraint prevents duplicate team+category", func(t *testing.T) {
		_, err := pool.Exec(ctx,
			`INSERT INTO team_handicap (team_id, category, points) VALUES ($1, 'winner', 30)`,
			teamID,
		)
		require.Error(t, err)
	})
}
