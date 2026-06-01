//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestOutcomesRepository(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	team1ID := uuid.Must(uuid.NewV7())
	team2ID := uuid.Must(uuid.NewV7())
	player1ID := uuid.Must(uuid.NewV7())
	player2ID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(
		ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'outcomes-test-cup', 'Outcomes Test Cup', '2026-06-01', '2026-07-01')`,
		tournamentID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO teams (id, name, tournament_id) VALUES
		($1, 'Brazil',    $3),
		($2, 'Argentina', $3)`,
		team1ID, team2ID, tournamentID,
	)
	require.NoError(t, err)

	_, err = pool.Exec(
		ctx, `
		INSERT INTO players (id, name, tournament_id, team_id, external_id) VALUES
		($1, 'Neymar',  $3, $4, $6),
		($2, 'Messi',   $3, $5, $7)`,
		player1ID, player2ID, tournamentID, team1ID, team2ID, "ext-1", "ext-2",
	)
	require.NoError(t, err)

	repo := repository.NewOutcomesRepository(pool)

	t.Run(
		"empty slices when no outcomes", func(t *testing.T) {
			players, err := repo.ListPlayerOutcomes(ctx, tournamentID)
			require.NoError(t, err)
			require.Empty(t, players)

			teams, err := repo.ListTeamOutcomes(ctx, tournamentID)
			require.NoError(t, err)
			require.Empty(t, teams)
		},
	)

	t.Run(
		"empty slices for unknown tournament", func(t *testing.T) {
			unknown := uuid.Must(uuid.NewV7())

			players, err := repo.ListPlayerOutcomes(ctx, unknown)
			require.NoError(t, err)
			require.Empty(t, players)

			teams, err := repo.ListTeamOutcomes(ctx, unknown)
			require.NoError(t, err)
			require.Empty(t, teams)
		},
	)

	t.Run(
		"returns enriched player and team outcomes", func(t *testing.T) {
			_, err = pool.Exec(
				ctx, `
			INSERT INTO player_outcomes (tournament_id, category, player_id) VALUES
			($1, 'total_top_scorer',  $2),
			($1, 'group_top_scorer',  $3)`,
				tournamentID, player1ID, player2ID,
			)
			require.NoError(t, err)

			_, err = pool.Exec(
				ctx, `
			INSERT INTO team_outcomes (tournament_id, category, team_id) VALUES
			($1, 'winner', $2)`,
				tournamentID, team1ID,
			)
			require.NoError(t, err)

			players, err := repo.ListPlayerOutcomes(ctx, tournamentID)
			require.NoError(t, err)
			require.Len(t, players, 2)

			// ordered by category ASC, name ASC: group_top_scorer < total_top_scorer
			require.Equal(t, "group_top_scorer", string(players[0].Category))
			require.Equal(t, "Messi", players[0].PlayerName)
			require.Equal(t, "Argentina", players[0].TeamName)

			require.Equal(t, "total_top_scorer", string(players[1].Category))
			require.Equal(t, "Neymar", players[1].PlayerName)
			require.Equal(t, "Brazil", players[1].TeamName)

			teams, err := repo.ListTeamOutcomes(ctx, tournamentID)
			require.NoError(t, err)
			require.Len(t, teams, 1)
			require.Equal(t, "winner", string(teams[0].Category))
			require.Equal(t, "Brazil", teams[0].TeamName)
		},
	)
}
