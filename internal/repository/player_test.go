//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func insertTeam(
	t *testing.T,
	pool *db.Pool,
	ctx context.Context,
	name, logo string,
) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	row := pool.QueryRow(ctx, `INSERT INTO teams (name, logo) VALUES ($1, $2) RETURNING id`, name, logo)
	require.NoError(t, row.Scan(&id))
	return id
}

func TestPlayerRepository_Search(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	tournamentAID := uuid.Must(uuid.NewV7())
	tournamentBID := uuid.Must(uuid.NewV7())

	_, err := pool.Exec(
		ctx, `
		INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES
		($1, 'player-test-a', 'Player Test A', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z'),
		($2, 'player-test-b', 'Player Test B', '2026-06-01T00:00:00Z', '2026-07-01T00:00:00Z')`,
		tournamentAID, tournamentBID,
	)
	require.NoError(t, err)

	argentinaID := insertTeam(t, pool, ctx, "Argentina", "<svg>arg</svg>")
	portugalID := insertTeam(t, pool, ctx, "Portugal", "<svg>por</svg>")
	englandID := insertTeam(t, pool, ctx, "England", "<svg>eng</svg>")
	brazilID := insertTeam(t, pool, ctx, "Brazil", "<svg>bra</svg>")
	usaID := insertTeam(t, pool, ctx, "USA", "<svg>usa</svg>")
	genericID := insertTeam(t, pool, ctx, "Generic", "")

	insertPlayer := func(
		externalID, name string,
		tournamentID, teamID uuid.UUID,
	) {
		t.Helper()
		_, err := pool.Exec(
			ctx,
			`INSERT INTO players (external_id, name, tournament_id, team_id) VALUES ($1, $2, $3, $4)`,
			externalID, name, tournamentID, teamID,
		)
		require.NoError(t, err)
	}

	repo := repository.NewPlayerRepository(pool)

	t.Run("returns matching players", func(t *testing.T) {
		insertPlayer("ext-1", "Lionel Messi", tournamentAID, argentinaID)
		insertPlayer("ext-2", "Ronaldo", tournamentAID, portugalID)
		insertPlayer("ext-3", "Mason Mount", tournamentAID, englandID)

		players, err := repo.Search(ctx, tournamentAID, "ess", "ess")
		require.NoError(t, err)
		require.Len(t, players, 1)
		require.Equal(t, "Lionel Messi", players[0].Name)
		require.Equal(t, "Argentina", players[0].TeamName)
		require.Equal(t, "<svg>arg</svg>", players[0].TeamLogo)
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		players, err := repo.Search(ctx, tournamentAID, "MESSI", "MESSI")
		require.NoError(t, err)
		require.Len(t, players, 1)
		require.Equal(t, "Lionel Messi", players[0].Name)
		require.Equal(t, "Argentina", players[0].TeamName)
		require.Equal(t, "<svg>arg</svg>", players[0].TeamLogo)
	})

	t.Run("search is tournament-scoped", func(t *testing.T) {
		insertPlayer("ext-b-1", "John Doe", tournamentBID, brazilID)
		insertPlayer("ext-a-4", "John Smith", tournamentAID, usaID)

		playersA, err := repo.Search(ctx, tournamentAID, "john", "john")
		require.NoError(t, err)
		require.Len(t, playersA, 1)
		require.Equal(t, "John Smith", playersA[0].Name)
		require.Equal(t, "USA", playersA[0].TeamName)
		require.Equal(t, "<svg>usa</svg>", playersA[0].TeamLogo)
	})

	t.Run("respects LIMIT 5", func(t *testing.T) {
		for i := 0; i < 8; i++ {
			insertPlayer(uuid.New().String(), "Player Son "+uuid.New().String(), tournamentAID, genericID)
		}

		players, err := repo.Search(ctx, tournamentAID, "Son", "Son")
		require.NoError(t, err)
		require.Len(t, players, 5)
		// Each result must have a non-empty TeamName from the JOIN — verifies the
		// join mapping is populated even when we only care about the count cap.
		for _, p := range players {
			require.NotEmpty(t, p.TeamName)
		}
	})

	t.Run("returns empty slice on no match", func(t *testing.T) {
		players, err := repo.Search(ctx, tournamentAID, "zzznomatch", "zzznomatch")
		require.NoError(t, err)
		require.NotNil(t, players)
		require.Empty(t, players)
	})

	t.Run("returns empty slice for nonexistent tournament", func(t *testing.T) {
		players, err := repo.Search(ctx, uuid.New(), "Messi", "Messi")
		require.NoError(t, err)
		require.NotNil(t, players)
		require.Empty(t, players)
	})

	t.Run("JOIN returns correct team_name and team_logo", func(t *testing.T) {
		franceID := insertTeam(t, pool, ctx, "France", "<svg>flag</svg>")
		insertPlayer("ext-mbappe", "Kylian Mbappé", tournamentAID, franceID)

		// unaccent_immutable maps "mbappe" → matches "Mbappé"; verifies the JOIN
		// correctly returns team_name and team_logo from the teams row.
		players, err := repo.Search(ctx, tournamentAID, "mbappe", "mbappe")
		require.NoError(t, err)
		require.Len(t, players, 1)
		require.Equal(t, "Kylian Mbappé", players[0].Name)
		require.Equal(t, "France", players[0].TeamName)
		require.Equal(t, "<svg>flag</svg>", players[0].TeamLogo)
	})
}
