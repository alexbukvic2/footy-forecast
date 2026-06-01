//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestTournamentGroupTableRepository(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentID := uuid.Must(uuid.NewV7())
	team1ID := uuid.Must(uuid.NewV7())
	team2ID := uuid.Must(uuid.NewV7())
	team3ID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
		VALUES ($1, 'tgt-test-cup', 'TGT Test Cup', '2026-06-01', '2026-07-01')`,
		tournamentID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO teams (id, name, tournament_id) VALUES
		($1, 'Argentina', $4),
		($2, 'France',    $4),
		($3, 'Brazil',    $4)`,
		team1ID, team2ID, team3ID, tournamentID)
	require.NoError(t, err)

	repo := repository.NewTournamentGroupTableRepository(pool)

	t.Run("empty slice when no entries", func(t *testing.T) {
		entries, err := repo.ListByTournament(ctx, tournamentID)
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("empty slice for unknown tournament", func(t *testing.T) {
		entries, err := repo.ListByTournament(ctx, uuid.Must(uuid.NewV7()))
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("returns rows in group_letter ASC position ASC order", func(t *testing.T) {
		_, err = pool.Exec(ctx, `
			INSERT INTO tournament_group_table
			    (tournament_id, team_id, group_letter, position, points) VALUES
			($1, $2, 'B', 1, 6),
			($1, $3, 'A', 1, 9),
			($1, $4, 'A', 2, 3)`,
			tournamentID, team2ID, team1ID, team3ID)
		require.NoError(t, err)

		entries, err := repo.ListByTournament(ctx, tournamentID)
		require.NoError(t, err)
		require.Len(t, entries, 3)

		// Sorted: A/1, A/2, B/1
		require.Equal(t, "A", entries[0].GroupLetter)
		require.Equal(t, 1, entries[0].Position)
		require.Equal(t, "Argentina", entries[0].TeamName)
		require.Equal(t, 9, entries[0].Points)

		require.Equal(t, "A", entries[1].GroupLetter)
		require.Equal(t, 2, entries[1].Position)
		require.Equal(t, "Brazil", entries[1].TeamName)
		require.Equal(t, 3, entries[1].Points)

		require.Equal(t, "B", entries[2].GroupLetter)
		require.Equal(t, 1, entries[2].Position)
		require.Equal(t, "France", entries[2].TeamName)
		require.Equal(t, 6, entries[2].Points)
	})
}
