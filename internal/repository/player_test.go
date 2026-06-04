//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func insertTeam(
	t *testing.T,
	pool *db.Pool,
	ctx context.Context,
	name, logo string,
	tournamentID uuid.UUID,
) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	row := pool.QueryRow(
		ctx,
		`INSERT INTO teams (name, logo, tournament_id) VALUES ($1, $2, $3) RETURNING id`,
		name,
		logo,
		tournamentID,
	)
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

	argentinaID := insertTeam(t, pool, ctx, "Argentina", "<svg>arg</svg>", tournamentAID)
	portugalID := insertTeam(t, pool, ctx, "Portugal", "<svg>por</svg>", tournamentAID)
	englandID := insertTeam(t, pool, ctx, "England", "<svg>eng</svg>", tournamentAID)
	brazilID := insertTeam(t, pool, ctx, "Brazil", "<svg>bra</svg>", tournamentBID)
	usaID := insertTeam(t, pool, ctx, "USA", "<svg>usa</svg>", tournamentAID)

	insertPlayer := func(
		externalID int64, name string,
		tournamentID, teamID uuid.UUID,
	) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		row := pool.QueryRow(
			ctx,
			`INSERT INTO players (external_id, name, tournament_id, team_id) VALUES ($1, $2, $3, $4) RETURNING id`,
			externalID, name, tournamentID, teamID,
		)
		require.NoError(t, row.Scan(&id))
		return id
	}

	insertPlayerHandicap := func(
		playerID uuid.UUID,
		category domain.PlayerHandicapCategory,
		points int,
	) {
		t.Helper()
		_, err := pool.Exec(
			ctx,
			`INSERT INTO player_handicap (player_id, category, points) VALUES ($1, $2::player_handicap_category, $3)`,
			playerID, string(category), points,
		)
		require.NoError(t, err)
	}

	repo := repository.NewPlayerRepository(pool)

	t.Run(
		"returns matching players", func(t *testing.T) {
			insertPlayer(1, "Lionel Messi", tournamentAID, argentinaID)
			insertPlayer(2, "Ronaldo", tournamentAID, portugalID)
			insertPlayer(3, "Mason Mount", tournamentAID, englandID)

			players, err := repo.Search(ctx, tournamentAID, "ess", "ess", nil, false)
			require.NoError(t, err)
			require.Len(t, players, 1)
			require.Equal(t, "Lionel Messi", players[0].Name)
			require.Equal(t, "Argentina", players[0].TeamName)
			require.Equal(t, "<svg>arg</svg>", players[0].TeamLogo)
		},
	)

	t.Run(
		"search is case-insensitive", func(t *testing.T) {
			players, err := repo.Search(ctx, tournamentAID, "MESSI", "MESSI", nil, false)
			require.NoError(t, err)
			require.Len(t, players, 1)
			require.Equal(t, "Lionel Messi", players[0].Name)
			require.Equal(t, "Argentina", players[0].TeamName)
			require.Equal(t, "<svg>arg</svg>", players[0].TeamLogo)
		},
	)

	t.Run(
		"search is tournament-scoped", func(t *testing.T) {
			insertPlayer(101, "John Doe", tournamentBID, brazilID)
			insertPlayer(4, "John Smith", tournamentAID, usaID)

			playersA, err := repo.Search(ctx, tournamentAID, "john", "john", nil, false)
			require.NoError(t, err)
			require.Len(t, playersA, 1)
			require.Equal(t, "John Smith", playersA[0].Name)
			require.Equal(t, "USA", playersA[0].TeamName)
			require.Equal(t, "<svg>usa</svg>", playersA[0].TeamLogo)
		},
	)

	t.Run(
		"returns empty slice on no match", func(t *testing.T) {
			players, err := repo.Search(ctx, tournamentAID, "zzznomatch", "zzznomatch", nil, false)
			require.NoError(t, err)
			require.NotNil(t, players)
			require.Empty(t, players)
		},
	)

	t.Run(
		"returns empty slice for nonexistent tournament", func(t *testing.T) {
			players, err := repo.Search(ctx, uuid.New(), "Messi", "Messi", nil, false)
			require.NoError(t, err)
			require.NotNil(t, players)
			require.Empty(t, players)
		},
	)

	t.Run(
		"JOIN returns correct team_name and team_logo", func(t *testing.T) {
			franceID := insertTeam(t, pool, ctx, "France", "<svg>flag</svg>", tournamentAID)
			insertPlayer(5, "Kylian Mbappé", tournamentAID, franceID)

			players, err := repo.Search(ctx, tournamentAID, "mbappe", "mbappe", nil, false)
			require.NoError(t, err)
			require.Len(t, players, 1)
			require.Equal(t, "Kylian Mbappé", players[0].Name)
			require.Equal(t, "France", players[0].TeamName)
			require.Equal(t, "<svg>flag</svg>", players[0].TeamLogo)
		},
	)

	t.Run(
		"handicaps: both categories returned when rows exist", func(t *testing.T) {
			spainID := insertTeam(t, pool, ctx, "Spain", "<svg>esp</svg>", tournamentAID)
			playerID := insertPlayer(6, "Lamine Yamal", tournamentAID, spainID)
			insertPlayerHandicap(playerID, domain.PlayerHandicapCategoryGroupTopScorer, 7)
			insertPlayerHandicap(playerID, domain.PlayerHandicapCategoryTotalTopScorer, 15)

			players, err := repo.Search(ctx, tournamentAID, "yamal", "yamal", nil, false)
			require.NoError(t, err)
			require.Len(t, players, 1)
			require.Equal(t, 7, players[0].Handicaps[domain.PlayerHandicapCategoryGroupTopScorer])
			require.Equal(t, 15, players[0].Handicaps[domain.PlayerHandicapCategoryTotalTopScorer])
		},
	)

	t.Run(
		"handicaps: single category returned when only one row exists", func(t *testing.T) {
			germanyID := insertTeam(t, pool, ctx, "Germany", "<svg>ger</svg>", tournamentAID)
			playerID := insertPlayer(7, "Jamal Musiala", tournamentAID, germanyID)
			insertPlayerHandicap(playerID, domain.PlayerHandicapCategoryGroupTopScorer, 3)

			players, err := repo.Search(ctx, tournamentAID, "musiala", "musiala", nil, false)
			require.NoError(t, err)
			require.Len(t, players, 1)
			require.Equal(t, 3, players[0].Handicaps[domain.PlayerHandicapCategoryGroupTopScorer])
			_, hasTotal := players[0].Handicaps[domain.PlayerHandicapCategoryTotalTopScorer]
			require.False(t, hasTotal, "repo must not invent missing categories — that is the service's job")
		},
	)

	t.Run(
		"handicaps: empty map returned when no handicap rows exist", func(t *testing.T) {
			italyID := insertTeam(t, pool, ctx, "Italy", "<svg>ita</svg>", tournamentAID)
			insertPlayer(8, "Nicolo Barella", tournamentAID, italyID)

			players, err := repo.Search(ctx, tournamentAID, "barella", "barella", nil, false)
			require.NoError(t, err)
			require.Len(t, players, 1)
			require.Empty(
				t,
				players[0].Handicaps,
				"repo must return empty map, not nil, for players with no handicap rows",
			)
		},
	)
}
