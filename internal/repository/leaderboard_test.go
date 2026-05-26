//go:build integration

package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/repository"
)

func TestLeaderboardRepository_GetForLeague(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	t.Run("empty league returns empty slice", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		ownerID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'el-cup','EL Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'el-owner','el-owner@x.com','Owner')", ownerID)
		emptyLeagueID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES ($1,$2,$3,'Empty League','EMPTYLB')", emptyLeagueID, tournID, ownerID)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForLeague(ctx, emptyLeagueID)
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("all members at 0 points get position 1", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		uA := uuid.Must(uuid.NewV7())
		uB := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'zero-cup','Zero Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'za','za@x.com','Za'),($2,'zb','zb@x.com','Zb')", uA, uB)
		lgID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES ($1,$2,$3,'Zero League','ZEROLB')", lgID, tournID, uA)
		_, _ = pool.Exec(ctx, "INSERT INTO league_members (league_id, user_id, role) VALUES ($1,$2,'owner'),($1,$3,'member')", lgID, uA, uB)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForLeague(ctx, lgID)
		require.NoError(t, err)
		require.Len(t, entries, 2)
		for _, e := range entries {
			require.Equal(t, 1, e.Position)
			require.Equal(t, 0, e.TotalPoints)
		}
	})

	t.Run("single member gets position 1", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		uA := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'single-cup','Single Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'sa','sa@x.com','Sa')", uA)
		lgID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES ($1,$2,$3,'Single League','SINGLB')", lgID, tournID, uA)
		_, _ = pool.Exec(ctx, "INSERT INTO league_members (league_id, user_id, role) VALUES ($1,$2,'owner')", lgID, uA)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForLeague(ctx, lgID)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		require.Equal(t, 1, entries[0].Position)
	})

	t.Run("members with different team points ranked correctly with tie", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		uA := uuid.Must(uuid.NewV7())
		uB := uuid.Must(uuid.NewV7())
		uC := uuid.Must(uuid.NewV7())
		teamID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'rank-cup','Rank Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'ra','ra@x.com','Ra'),($2,'rb','rb@x.com','Rb'),($3,'rc','rc@x.com','Rc')", uA, uB, uC)
		_, _ = pool.Exec(ctx, "INSERT INTO teams (id, tournament_id, name, logo) VALUES ($1,$2,'Team A','')", teamID, tournID)
		lgID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES ($1,$2,$3,'Rank League','RANKLB')", lgID, tournID, uA)
		_, _ = pool.Exec(ctx, "INSERT INTO league_members (league_id, user_id, role) VALUES ($1,$2,'owner'),($1,$3,'member'),($1,$4,'member')", lgID, uA, uB, uC)
		// uA:10, uB:10 (tie), uC:5
		_, err := pool.Exec(ctx,
			"INSERT INTO team_predictions (id, user_id, tournament_id, category, pick, points) VALUES ($1,$2,$3,'winner',$4,10),($5,$6,$3,'winner',$4,10),($7,$8,$3,'winner',$4,5)",
			uuid.Must(uuid.NewV7()), uA, tournID, teamID,
			uuid.Must(uuid.NewV7()), uB, tournID, teamID,
			uuid.Must(uuid.NewV7()), uC, tournID, teamID)
		require.NoError(t, err)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForLeague(ctx, lgID)
		require.NoError(t, err)
		require.Len(t, entries, 3)
		require.Equal(t, 1, entries[0].Position)
		require.Equal(t, 1, entries[1].Position)
		require.Equal(t, 2, entries[2].Position)
		require.Equal(t, 5, entries[2].TotalPoints)
	})

	t.Run("member with no predictions included with zeros ranks last", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		uA := uuid.Must(uuid.NewV7())
		uB := uuid.Must(uuid.NewV7())
		teamID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'zero2-cup','Zero2 Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'z2a','z2a@x.com','Z2A'),($2,'z2b','z2b@x.com','Z2B')", uA, uB)
		_, _ = pool.Exec(ctx, "INSERT INTO teams (id, tournament_id, name, logo) VALUES ($1,$2,'Z2 Team','')", teamID, tournID)
		lgID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES ($1,$2,$3,'Zero2 League','ZERO2LB')", lgID, tournID, uA)
		_, _ = pool.Exec(ctx, "INSERT INTO league_members (league_id, user_id, role) VALUES ($1,$2,'owner'),($1,$3,'member')", lgID, uA, uB)
		_, err := pool.Exec(ctx, "INSERT INTO team_predictions (id, user_id, tournament_id, category, pick, points) VALUES ($1,$2,$3,'winner',$4,7)",
			uuid.Must(uuid.NewV7()), uA, tournID, teamID)
		require.NoError(t, err)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForLeague(ctx, lgID)
		require.NoError(t, err)
		require.Len(t, entries, 2)
		require.Equal(t, 1, entries[0].Position)
		require.Equal(t, 7, entries[0].TotalPoints)
		require.Equal(t, 2, entries[1].Position)
		require.Equal(t, 0, entries[1].TotalPoints)
	})
}

func TestLeaderboardRepository_GetForTournament(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()

	t.Run("tournament with no predictions returns empty slice", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'tourn-lb','Tourn LB','2026-06-01','2026-07-01')", tournID)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForTournament(ctx, tournID)
		require.NoError(t, err)
		require.Empty(t, entries)
	})

	t.Run("user with team prediction appears with correct category zeros", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		uA := uuid.Must(uuid.NewV7())
		teamID := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'tl-cup','TL Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'tla','tla@x.com','TLA')", uA)
		_, _ = pool.Exec(ctx, "INSERT INTO teams (id, tournament_id, name, logo) VALUES ($1,$2,'TL Team','')", teamID, tournID)
		_, err := pool.Exec(ctx, "INSERT INTO team_predictions (id, user_id, tournament_id, category, pick, points) VALUES ($1,$2,$3,'winner',$4,8)",
			uuid.Must(uuid.NewV7()), uA, tournID, teamID)
		require.NoError(t, err)

		repo := repository.NewLeaderboardRepository(pool)
		entries, err := repo.GetForTournament(ctx, tournID)
		require.NoError(t, err)
		require.Len(t, entries, 1)
		require.Equal(t, 1, entries[0].Position)
		require.Equal(t, 8, entries[0].TeamPoints)
		require.Equal(t, 0, entries[0].ScorePoints)
		require.Equal(t, 0, entries[0].PlayerPoints)
		require.Equal(t, 8, entries[0].TotalPoints)
	})
}

func TestLeaderboardRepository_GetUserPositionsInLeagues(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewLeaderboardRepository(pool)

	t.Run("empty leagueIDs returns empty map without DB round trip", func(t *testing.T) {
		t.Parallel()
		result, err := repo.GetUserPositionsInLeagues(ctx, uuid.Must(uuid.NewV7()), []uuid.UUID{})
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("user in two leagues returns correct positions", func(t *testing.T) {
		t.Parallel()
		tournID := uuid.Must(uuid.NewV7())
		uA := uuid.Must(uuid.NewV7())
		uB := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO tournaments (id, slug, name, starts_at, ends_at) VALUES ($1,'pos-cup','Pos Cup','2026-06-01','2026-07-01')", tournID)
		_, _ = pool.Exec(ctx, "INSERT INTO users (id, cognito_sub, email, display_name) VALUES ($1,'posa','posa@x.com','PosA'),($2,'posb','posb@x.com','PosB')", uA, uB)
		lg1 := uuid.Must(uuid.NewV7())
		lg2 := uuid.Must(uuid.NewV7())
		_, _ = pool.Exec(ctx, "INSERT INTO leagues (id, tournament_id, owner_id, name, code) VALUES ($1,$2,$3,'Pos L1','POSLB01'),($4,$2,$3,'Pos L2','POSLB02')", lg1, tournID, uA, lg2)
		_, _ = pool.Exec(ctx, "INSERT INTO league_members (league_id, user_id, role) VALUES ($1,$2,'owner'),($1,$3,'member'),($4,$2,'owner')", lg1, uA, uB, lg2)

		result, err := repo.GetUserPositionsInLeagues(ctx, uA, []uuid.UUID{lg1, lg2})
		require.NoError(t, err)
		require.Equal(t, 1, result[lg1])
		require.Equal(t, 1, result[lg2])
	})
}
