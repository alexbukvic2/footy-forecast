//go:build integration

package repository_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
	"github.com/alexbukvic2/footy-forecast/internal/worker"
)

// wExtIDSeq is an atomic counter that guarantees unique external IDs across all
// subtests sharing the same Postgres container.
var wExtIDSeq atomic.Int64

// ---- shared seeding helpers ----

func wInsertTournament(t *testing.T, pool *db.Pool, ctx context.Context, withExtID bool) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	slug := "wt-" + id.String()
	if withExtID {
		_, err := pool.Exec(ctx,
			`INSERT INTO tournaments (id, slug, name, starts_at, ends_at, external_id, season)
			 VALUES ($1, $2, 'Test', now(), now()+interval '30 days', 200, 2026)`,
			id, slug)
		require.NoError(t, err)
	} else {
		_, err := pool.Exec(ctx,
			`INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
			 VALUES ($1, $2, 'Test', now(), now()+interval '30 days')`,
			id, slug)
		require.NoError(t, err)
	}
	return id
}

func wInsertTeam(t *testing.T, pool *db.Pool, ctx context.Context, tournamentID uuid.UUID, groupLetter string, extID int64) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	name := "team-" + id.String()
	var gl interface{}
	if groupLetter != "" {
		gl = groupLetter
	}
	var ext interface{}
	if extID != 0 {
		ext = extID
	}
	_, err := pool.Exec(ctx,
		`INSERT INTO teams (id, name, logo, tournament_id, group_letter, external_id)
		 VALUES ($1, $2, '', $3, $4, $5)`,
		id, name, tournamentID, gl, ext)
	require.NoError(t, err)
	return id
}

func wInsertFixture(t *testing.T, pool *db.Pool, ctx context.Context,
	tournamentID, homeID, awayID uuid.UUID,
	kickoff time.Time, status, round string,
	isDemo bool,
) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	extID := wExtIDSeq.Add(1)
	_, err := pool.Exec(ctx,
		`INSERT INTO fixtures (id, external_id, tournament_id, home_team_id, away_team_id,
		                       kickoff_at, status, round, is_demo)
		 VALUES ($1, $2, $3, $4, $5, $6, $7::fixture_status, $8, $9)`,
		id, extID, tournamentID, homeID, awayID, kickoff, status, round, isDemo)
	require.NoError(t, err)
	return id
}

func wInsertUser(t *testing.T, pool *db.Pool, ctx context.Context) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	_, err := pool.Exec(ctx,
		`INSERT INTO users (id, cognito_sub, email, display_name)
		 VALUES ($1, $2, $3, 'Tester')`,
		id, id.String(), id.String()+"@t.com")
	require.NoError(t, err)
	return id
}

func wInsertScorePrediction(t *testing.T, pool *db.Pool, ctx context.Context, userID, fixtureID uuid.UUID, gh, ga int) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO score_predictions (user_id, fixture_id, goals_home, goals_away)
		 VALUES ($1, $2, $3, $4)`,
		userID, fixtureID, gh, ga)
	require.NoError(t, err)
}

func wQueryPoints(t *testing.T, pool *db.Pool, ctx context.Context, fixtureID, userID uuid.UUID) int {
	t.Helper()
	var pts int
	err := pool.QueryRow(ctx,
		`SELECT COALESCE(points, -1) FROM score_predictions WHERE fixture_id=$1 AND user_id=$2`,
		fixtureID, userID).Scan(&pts)
	require.NoError(t, err)
	return pts
}

func wSetWinnerTeamID(t *testing.T, pool *db.Pool, ctx context.Context, fixtureID, teamID uuid.UUID) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`UPDATE fixtures SET winner_team_id=$1, status='finished' WHERE id=$2`,
		teamID, fixtureID)
	require.NoError(t, err)
}

func wSetUpdatedAt(t *testing.T, pool *db.Pool, ctx context.Context, fixtureID uuid.UUID, at time.Time) {
	t.Helper()
	// The fixtures_set_updated_at trigger overrides any manual updated_at on every UPDATE,
	// so we must disable it temporarily. The test user owns the table, so ALTER TABLE is allowed.
	_, err := pool.Exec(ctx, `ALTER TABLE fixtures DISABLE TRIGGER fixtures_set_updated_at`)
	require.NoError(t, err)
	defer func() {
		_, _ = pool.Exec(ctx, `ALTER TABLE fixtures ENABLE TRIGGER fixtures_set_updated_at`)
	}()
	_, err = pool.Exec(ctx, `UPDATE fixtures SET updated_at=$1 WHERE id=$2`, at, fixtureID)
	require.NoError(t, err)
}

func wInsertTeamHandicap(t *testing.T, pool *db.Pool, ctx context.Context, teamID uuid.UUID, category string, pts int) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO team_handicap (team_id, category, points) VALUES ($1, $2::team_handicap_category, $3)`,
		teamID, category, pts)
	require.NoError(t, err)
}

func wInsertPlayerHandicap(t *testing.T, pool *db.Pool, ctx context.Context, playerID uuid.UUID, category string, pts int) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO player_handicap (player_id, category, points) VALUES ($1, $2::player_handicap_category, $3)`,
		playerID, category, pts)
	require.NoError(t, err)
}

func wInsertGroupTableEntry(t *testing.T, pool *db.Pool, ctx context.Context, tournamentID, teamID uuid.UUID, group string, pos int) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO tournament_group_table (tournament_id, team_id, group_letter, position, points, played)
		 VALUES ($1, $2, $3, $4, 0, 0)
		 ON CONFLICT (tournament_id, team_id) DO UPDATE SET position=$4`,
		tournamentID, teamID, group, pos)
	require.NoError(t, err)
}

// ---- ListPollableMatches ----

func TestWorkerRepository_ListPollableMatches(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)

	tourID := wInsertTournament(t, pool, ctx, true)
	home := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	away := wInsertTeam(t, pool, ctx, tourID, "A", 0)

	t.Run("live fixture returned", func(t *testing.T) {
		id := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(-30*time.Minute), "in_progress", "GS", false)
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.True(t, containsFixture(fixtures, id))
	})

	t.Run("kickoff in 4 minutes returned", func(t *testing.T) {
		id := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(4*time.Minute), "upcoming", "GS", false)
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.True(t, containsFixture(fixtures, id))
	})

	t.Run("kickoff in 60 minutes not returned", func(t *testing.T) {
		id := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(60*time.Minute), "upcoming", "GS", false)
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.False(t, containsFixture(fixtures, id))
	})

	t.Run("finished within 24h returned", func(t *testing.T) {
		id := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(-23*time.Hour), "finished", "GS", false)
		wSetUpdatedAt(t, pool, ctx, id, time.Now().Add(-23*time.Hour))
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.True(t, containsFixture(fixtures, id))
	})

	t.Run("finished 25h ago not returned", func(t *testing.T) {
		id := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(-26*time.Hour), "finished", "GS", false)
		wSetUpdatedAt(t, pool, ctx, id, time.Now().Add(-26*time.Hour))
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.False(t, containsFixture(fixtures, id))
	})

	t.Run("demo fixture not returned", func(t *testing.T) {
		id := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(-10*time.Minute), "in_progress", "GS", true)
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.False(t, containsFixture(fixtures, id))
	})

	t.Run("tournament without external_id not returned", func(t *testing.T) {
		noExtTour := wInsertTournament(t, pool, ctx, false)
		h2 := wInsertTeam(t, pool, ctx, noExtTour, "B", 0)
		a2 := wInsertTeam(t, pool, ctx, noExtTour, "B", 0)
		id := wInsertFixture(t, pool, ctx, noExtTour, h2, a2, time.Now().Add(-10*time.Minute), "in_progress", "GS", false)
		fixtures, err := repo.ListPollableMatches(ctx)
		require.NoError(t, err)
		assert.False(t, containsFixture(fixtures, id))
	})
}

// ---- UpdateMatchAndRescoreLivePredictions ----

func TestWorkerRepository_UpdateMatchAndRescore(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)

	tourID := wInsertTournament(t, pool, ctx, true)
	home := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	away := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	userID := wInsertUser(t, pool, ctx)

	newF := func() domain.PollableFixture {
		fID := wInsertFixture(t, pool, ctx, tourID, home, away, time.Now().Add(-1*time.Hour), "in_progress", "GS", false)
		return domain.PollableFixture{ID: fID, HomeTeamID: home, AwayTeamID: away}
	}

	t.Run("exact score = 6 pts", func(t *testing.T) {
		// home correct (+1) + away correct (+1) + outcome correct (+2) + GD correct (+2) = 6
		f := newF()
		wInsertScorePrediction(t, pool, ctx, userID, f.ID, 2, 1)
		gh, ga := 2, 1
		require.NoError(t, repo.UpdateMatchAndRescoreLivePredictions(ctx, f, worker.APIFixtureResult{StatusShort: "FT", GoalsHome: &gh, GoalsAway: &ga}))
		assert.Equal(t, 6, wQueryPoints(t, pool, ctx, f.ID, userID))
	})

	t.Run("correct outcome + away goals = 3 pts", func(t *testing.T) {
		// predict 2-0, result 3-0: away correct (+1) + outcome correct (+2) = 3
		f := newF()
		wInsertScorePrediction(t, pool, ctx, userID, f.ID, 2, 0)
		gh, ga := 3, 0
		require.NoError(t, repo.UpdateMatchAndRescoreLivePredictions(ctx, f, worker.APIFixtureResult{StatusShort: "FT", GoalsHome: &gh, GoalsAway: &ga}))
		assert.Equal(t, 3, wQueryPoints(t, pool, ctx, f.ID, userID))
	})

	t.Run("wrong result = 0 pts", func(t *testing.T) {
		// predict 2-0, result 0-1: no component correct
		f := newF()
		wInsertScorePrediction(t, pool, ctx, userID, f.ID, 2, 0)
		gh, ga := 0, 1
		require.NoError(t, repo.UpdateMatchAndRescoreLivePredictions(ctx, f, worker.APIFixtureResult{StatusShort: "FT", GoalsHome: &gh, GoalsAway: &ga}))
		assert.Equal(t, 0, wQueryPoints(t, pool, ctx, f.ID, userID))
	})

	t.Run("cancelled = 0 pts", func(t *testing.T) {
		f := newF()
		wInsertScorePrediction(t, pool, ctx, userID, f.ID, 2, 0)
		require.NoError(t, repo.UpdateMatchAndRescoreLivePredictions(ctx, f, worker.APIFixtureResult{StatusShort: "CANC"}))
		assert.Equal(t, 0, wQueryPoints(t, pool, ctx, f.ID, userID))
	})

	t.Run("rescore on second call", func(t *testing.T) {
		// predict 1-0; at HT score is 1-0: all components match → 6 pts
		// at FT score is 2-0: away correct (+1) + outcome correct (+2) = 3 pts
		f := newF()
		wInsertScorePrediction(t, pool, ctx, userID, f.ID, 1, 0)
		gh, ga := 1, 0
		require.NoError(t, repo.UpdateMatchAndRescoreLivePredictions(ctx, f, worker.APIFixtureResult{StatusShort: "1H", GoalsHome: &gh, GoalsAway: &ga}))
		assert.Equal(t, 6, wQueryPoints(t, pool, ctx, f.ID, userID))
		gh2, ga2 := 2, 0
		require.NoError(t, repo.UpdateMatchAndRescoreLivePredictions(ctx, f, worker.APIFixtureResult{StatusShort: "FT", GoalsHome: &gh2, GoalsAway: &ga2}))
		assert.Equal(t, 3, wQueryPoints(t, pool, ctx, f.ID, userID))
	})
}

// ---- Completion checks ----

func TestWorkerRepository_IsGroupComplete(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	h := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	a := wInsertTeam(t, pool, ctx, tourID, "A", 0)

	t.Run("all finished", func(t *testing.T) {
		wInsertFixture(t, pool, ctx, tourID, h, a, time.Now().Add(-3*time.Hour), "finished", "GS", false)
		done, err := repo.IsGroupComplete(ctx, tourID, "A")
		require.NoError(t, err)
		assert.True(t, done)
	})
	t.Run("one in_progress", func(t *testing.T) {
		wInsertFixture(t, pool, ctx, tourID, h, a, time.Now().Add(-1*time.Hour), "in_progress", "GS", false)
		done, err := repo.IsGroupComplete(ctx, tourID, "A")
		require.NoError(t, err)
		assert.False(t, done)
	})
	t.Run("all cancelled counts as done", func(t *testing.T) {
		// all existing fixtures are finished or in_progress; add one more cancelled
		tourID2 := wInsertTournament(t, pool, ctx, true)
		h2 := wInsertTeam(t, pool, ctx, tourID2, "C", 0)
		a2 := wInsertTeam(t, pool, ctx, tourID2, "C", 0)
		wInsertFixture(t, pool, ctx, tourID2, h2, a2, time.Now().Add(-3*time.Hour), "cancelled", "GS", false)
		done, err := repo.IsGroupComplete(ctx, tourID2, "C")
		require.NoError(t, err)
		assert.True(t, done)
	})
}

func TestWorkerRepository_IsGroupStageComplete(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	h := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	a := wInsertTeam(t, pool, ctx, tourID, "A", 0)

	wInsertFixture(t, pool, ctx, tourID, h, a, time.Now().Add(-3*time.Hour), "finished", "GS", false)
	done, err := repo.IsGroupStageComplete(ctx, tourID)
	require.NoError(t, err)
	assert.True(t, done)

	wInsertFixture(t, pool, ctx, tourID, h, a, time.Now().Add(2*time.Hour), "upcoming", "GS", false)
	done, err = repo.IsGroupStageComplete(ctx, tourID)
	require.NoError(t, err)
	assert.False(t, done)
}

func TestWorkerRepository_IsRoundComplete(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	h := wInsertTeam(t, pool, ctx, tourID, "", 0)
	a := wInsertTeam(t, pool, ctx, tourID, "", 0)
	winnerID := wInsertTeam(t, pool, ctx, tourID, "", 0)

	t.Run("all QFs have winner_team_id", func(t *testing.T) {
		fID := wInsertFixture(t, pool, ctx, tourID, h, a, time.Now().Add(-3*time.Hour), "finished", "Quarter-finals", false)
		wSetWinnerTeamID(t, pool, ctx, fID, winnerID)
		done, err := repo.IsRoundComplete(ctx, tourID, "Quarter-finals")
		require.NoError(t, err)
		assert.True(t, done)
	})
	t.Run("cancelled QF has no winner_team_id → false", func(t *testing.T) {
		wInsertFixture(t, pool, ctx, tourID, h, a, time.Now().Add(-3*time.Hour), "cancelled", "Quarter-finals", false)
		done, err := repo.IsRoundComplete(ctx, tourID, "Quarter-finals")
		require.NoError(t, err)
		assert.False(t, done)
	})
}

// ---- UpdateGroupStandings ----

func TestWorkerRepository_UpdateGroupStandings(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	t1 := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	t2 := wInsertTeam(t, pool, ctx, tourID, "A", 0)

	entries := []domain.StandingsEntry{
		{TeamID: t1, Position: 1, Points: 6, Played: 2, Won: 2, Drawn: 0, Lost: 0, GoalsFor: 4, GoalsAgainst: 1},
		{TeamID: t2, Position: 2, Points: 3, Played: 2, Won: 1, Drawn: 0, Lost: 1, GoalsFor: 2, GoalsAgainst: 3},
	}
	require.NoError(t, repo.UpdateGroupStandings(ctx, tourID, "A", entries))

	var pos, pts int
	require.NoError(t, pool.QueryRow(ctx, `SELECT position, points FROM tournament_group_table WHERE tournament_id=$1 AND team_id=$2`, tourID, t1).Scan(&pos, &pts))
	assert.Equal(t, 1, pos)
	assert.Equal(t, 6, pts)

	// Idempotent second call updates values.
	entries[0].Points = 9
	require.NoError(t, repo.UpdateGroupStandings(ctx, tourID, "A", entries))
	require.NoError(t, pool.QueryRow(ctx, `SELECT points FROM tournament_group_table WHERE tournament_id=$1 AND team_id=$2`, tourID, t1).Scan(&pts))
	assert.Equal(t, 9, pts)
}

// ---- GetTeamByExternalID / GetPlayerByExternalID ----

func TestWorkerRepository_GetTeamByExternalID(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	teamID := wInsertTeam(t, pool, ctx, tourID, "A", 789)

	got, err := repo.GetTeamByExternalID(ctx, 789, tourID)
	require.NoError(t, err)
	assert.Equal(t, teamID, got)
}

func TestWorkerRepository_GetPlayerByExternalID(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	teamID := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	playerID := wInsertPlayer(t, pool, ctx, tourID, teamID, "ext-42")

	got, err := repo.GetPlayerByExternalID(ctx, "ext-42", tourID)
	require.NoError(t, err)
	assert.Equal(t, playerID, got)
}

// ---- SettleGroupWinnerPredictions ----

func TestWorkerRepository_SettleGroupWinnerPredictions(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	t1 := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	t2 := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	t3 := wInsertTeam(t, pool, ctx, tourID, "A", 0)
	userID := wInsertUser(t, pool, ctx)
	wInsertGroupTableEntry(t, pool, ctx, tourID, t1, "A", 1)
	wInsertGroupTableEntry(t, pool, ctx, tourID, t2, "A", 2)
	wInsertGroupTableEntry(t, pool, ctx, tourID, t3, "A", 3)
	wInsertTeamHandicap(t, pool, ctx, t1, "group_winner", 5)
	wInsertTeamHandicap(t, pool, ctx, t2, "group_winner", 3)

	// Correct slot-0 pick (team t1 at pos 1)
	wInsertTeamPrediction(t, pool, ctx, userID, tourID, "group_winner", t1, "A", 0)
	// Wrong position (t3 placed 3rd)
	wInsertTeamPrediction(t, pool, ctx, userID, tourID, "group_winner", t3, "A", 1)

	require.NoError(t, repo.SettleGroupWinnerPredictions(ctx, tourID, "A"))

	assert.Equal(t, 5, queryTeamPredictionPoints(t, pool, ctx, userID, tourID, "group_winner", t1, "A", 0))
	assert.Equal(t, 0, queryTeamPredictionPoints(t, pool, ctx, userID, tourID, "group_winner", t3, "A", 1))

	// Idempotent second call is no-op.
	require.NoError(t, repo.SettleGroupWinnerPredictions(ctx, tourID, "A"))
	assert.Equal(t, 5, queryTeamPredictionPoints(t, pool, ctx, userID, tourID, "group_winner", t1, "A", 0))
}

// ---- SettlePlayoffGroupPredictions ----

func TestWorkerRepository_SettlePlayoffGroupPredictions(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	t1 := wInsertTeam(t, pool, ctx, tourID, "B", 0)
	t2 := wInsertTeam(t, pool, ctx, tourID, "B", 0)
	t3 := wInsertTeam(t, pool, ctx, tourID, "B", 0)
	userID := wInsertUser(t, pool, ctx)
	wInsertGroupTableEntry(t, pool, ctx, tourID, t1, "B", 1)
	wInsertGroupTableEntry(t, pool, ctx, tourID, t2, "B", 2)
	wInsertGroupTableEntry(t, pool, ctx, tourID, t3, "B", 3)
	wInsertTeamHandicap(t, pool, ctx, t1, "playoff", 4)
	wInsertTeamHandicap(t, pool, ctx, t2, "playoff", 2)

	wInsertTeamPrediction(t, pool, ctx, userID, tourID, "playoff", t1, "B", 0)
	wInsertTeamPrediction(t, pool, ctx, userID, tourID, "playoff", t3, "B", 1)

	require.NoError(t, repo.SettlePlayoffGroupPredictions(ctx, tourID, "B"))
	assert.Equal(t, 4, queryTeamPredictionPoints(t, pool, ctx, userID, tourID, "playoff", t1, "B", 0))
	assert.Equal(t, 0, queryTeamPredictionPoints(t, pool, ctx, userID, tourID, "playoff", t3, "B", 1))
}

// ---- SettlePlayoffWildcardPredictions ----

func TestWorkerRepository_SettlePlayoffWildcardPredictions(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	t1 := wInsertTeam(t, pool, ctx, tourID, "", 0)
	t2 := wInsertTeam(t, pool, ctx, tourID, "", 0)
	// Use separate users: one picks the team that advances, one picks the eliminated team.
	userA := wInsertUser(t, pool, ctx)
	userB := wInsertUser(t, pool, ctx)
	wInsertTeamHandicap(t, pool, ctx, t1, "playoff", 6)

	// t1 advanced (in team_outcomes); t2 eliminated
	_, err := pool.Exec(ctx, `INSERT INTO team_outcomes (tournament_id, category, team_id) VALUES ($1,'playoff',$2)`, tourID, t1)
	require.NoError(t, err)

	wInsertTeamPredictionWildcard(t, pool, ctx, userA, tourID, "playoff", t1)
	wInsertTeamPredictionWildcard(t, pool, ctx, userB, tourID, "playoff", t2)

	require.NoError(t, repo.SettlePlayoffWildcardPredictions(ctx, tourID))
	assert.Equal(t, 6, queryTeamPredictionPointsWildcard(t, pool, ctx, userA, tourID, t1))
	assert.Equal(t, 0, queryTeamPredictionPointsWildcard(t, pool, ctx, userB, tourID, t2))
}

// ---- SettleSemifinalistPredictions ----

func TestWorkerRepository_SettleSemifinalistPredictions(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	// userA picks all 4 correct semifinalists (slots 0-3).
	// userB picks a team that is eliminated (slot 0).
	userA := wInsertUser(t, pool, ctx)
	userB := wInsertUser(t, pool, ctx)

	// Insert 4 QF winners.
	sfTeams := make([]uuid.UUID, 4)
	for i := range sfTeams {
		sfTeams[i] = wInsertTeam(t, pool, ctx, tourID, "", 0)
		other := wInsertTeam(t, pool, ctx, tourID, "", 0)
		fID := wInsertFixture(t, pool, ctx, tourID, sfTeams[i], other, time.Now().Add(-3*time.Hour), "finished", "Quarter-finals", false)
		wSetWinnerTeamID(t, pool, ctx, fID, sfTeams[i])
		wInsertTeamHandicap(t, pool, ctx, sfTeams[i], "semifinalist", 5)
	}
	eliminated := wInsertTeam(t, pool, ctx, tourID, "", 0)

	for _, id := range sfTeams {
		wInsertTeamPredictionSemi(t, pool, ctx, userA, tourID, id)
	}
	wInsertTeamPredictionSemi(t, pool, ctx, userB, tourID, eliminated)

	require.NoError(t, repo.SettleSemifinalistPredictions(ctx, tourID))

	for _, id := range sfTeams {
		assert.Equal(t, 5, queryTeamPredictionPointsSemi(t, pool, ctx, userA, tourID, id))
	}
	assert.Equal(t, 0, queryTeamPredictionPointsSemi(t, pool, ctx, userB, tourID, eliminated))
}

// ---- SettleTournamentWinnerPredictions ----

func TestWorkerRepository_SettleTournamentWinnerPredictions(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	winner := wInsertTeam(t, pool, ctx, tourID, "", 0)
	loser := wInsertTeam(t, pool, ctx, tourID, "", 0)
	// Use separate users: the constraint allows only one winner pick per user.
	userA := wInsertUser(t, pool, ctx)
	userB := wInsertUser(t, pool, ctx)
	wInsertTeamHandicap(t, pool, ctx, winner, "winner", 10)

	wInsertTeamPredictionWinner(t, pool, ctx, userA, tourID, winner)
	wInsertTeamPredictionWinner(t, pool, ctx, userB, tourID, loser)

	require.NoError(t, repo.SettleTournamentWinnerPredictions(ctx, tourID, winner))
	assert.Equal(t, 10, queryTeamPredictionPointsWinner(t, pool, ctx, userA, tourID, winner))
	assert.Equal(t, 0, queryTeamPredictionPointsWinner(t, pool, ctx, userB, tourID, loser))
}

// ---- SettleTopScorerPredictions ----

func TestWorkerRepository_SettleTopScorerPredictions(t *testing.T) {
	t.Parallel()
	pool := startPostgres(t)
	ctx := context.Background()
	repo := repository.NewWorkerRepository(pool)
	tourID := wInsertTournament(t, pool, ctx, true)
	teamID := wInsertTeam(t, pool, ctx, tourID, "", 0)
	topScorer := wInsertPlayer(t, pool, ctx, tourID, teamID, "top-1")
	other := wInsertPlayer(t, pool, ctx, tourID, teamID, "other-1")
	// Use separate users: the constraint allows only one top-scorer pick per user per category.
	userA := wInsertUser(t, pool, ctx)
	userB := wInsertUser(t, pool, ctx)
	wInsertPlayerHandicap(t, pool, ctx, topScorer, "total_top_scorer", 8)

	wInsertPlayerPrediction(t, pool, ctx, userA, tourID, "total_top_scorer", topScorer)
	wInsertPlayerPrediction(t, pool, ctx, userB, tourID, "total_top_scorer", other)

	require.NoError(t, repo.SettleTopScorerPredictions(ctx, tourID, topScorer))
	var pts int
	require.NoError(t, pool.QueryRow(ctx, `SELECT COALESCE(points,-1) FROM player_predictions WHERE tournament_id=$1 AND user_id=$2 AND pick=$3`, tourID, userA, topScorer).Scan(&pts))
	assert.Equal(t, 8, pts)
	require.NoError(t, pool.QueryRow(ctx, `SELECT COALESCE(points,-1) FROM player_predictions WHERE tournament_id=$1 AND user_id=$2 AND pick=$3`, tourID, userB, other).Scan(&pts))
	assert.Equal(t, 0, pts)
}

// ---- low-level helpers ----

func containsFixture(fixtures []domain.PollableFixture, id uuid.UUID) bool {
	for _, f := range fixtures {
		if f.ID == id {
			return true
		}
	}
	return false
}

func wInsertPlayer(t *testing.T, pool *db.Pool, ctx context.Context, tourID, teamID uuid.UUID, extID string) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV7())
	_, err := pool.Exec(ctx,
		`INSERT INTO players (id, external_id, name, tournament_id, team_id) VALUES ($1,$2,$3,$4,$5)`,
		id, extID, "Player "+extID, tourID, teamID)
	require.NoError(t, err)
	return id
}

func wInsertTeamPrediction(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID uuid.UUID, category string, pick uuid.UUID, group string, slot int) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO team_predictions (user_id, tournament_id, category, pick, group_letter, slot_index)
		 VALUES ($1,$2,$3::team_handicap_category,$4,$5,$6)`,
		userID, tourID, category, pick, group, slot)
	require.NoError(t, err)
}

func wInsertTeamPredictionWildcard(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID uuid.UUID, category string, pick uuid.UUID) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO team_predictions (user_id, tournament_id, category, pick, slot_index)
		 VALUES ($1,$2,$3::team_handicap_category,$4,2)`,
		userID, tourID, category, pick)
	require.NoError(t, err)
}

func wInsertTeamPredictionSemi(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID, pick uuid.UUID) {
	t.Helper()
	slotIDx := 0
	_, err := pool.Exec(ctx,
		`INSERT INTO team_predictions (user_id, tournament_id, category, pick, slot_index)
		 VALUES ($1,$2,'semifinalist',$3,$4)`,
		userID, tourID, pick, slotIDx)
	if err != nil {
		// try slot 1, 2, 3 if 0 already used
		for slot := 1; slot <= 3; slot++ {
			_, err = pool.Exec(ctx,
				`INSERT INTO team_predictions (user_id, tournament_id, category, pick, slot_index)
				 VALUES ($1,$2,'semifinalist',$3,$4)`,
				userID, tourID, pick, slot)
			if err == nil {
				return
			}
		}
		require.NoError(t, err)
	}
}

func wInsertTeamPredictionWinner(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID, pick uuid.UUID) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO team_predictions (user_id, tournament_id, category, pick, slot_index)
		 VALUES ($1,$2,'winner',$3,0)`,
		userID, tourID, pick)
	require.NoError(t, err)
}

func wInsertPlayerPrediction(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID uuid.UUID, category string, pick uuid.UUID) {
	t.Helper()
	_, err := pool.Exec(ctx,
		`INSERT INTO player_predictions (user_id, tournament_id, category, pick)
		 VALUES ($1,$2,$3::player_handicap_category,$4)`,
		userID, tourID, category, pick)
	require.NoError(t, err)
}

func queryTeamPredictionPoints(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID uuid.UUID, category string, pick uuid.UUID, group string, slot int) int {
	t.Helper()
	var pts int
	err := pool.QueryRow(ctx,
		`SELECT COALESCE(points,-999) FROM team_predictions WHERE user_id=$1 AND tournament_id=$2 AND category=$3::team_handicap_category AND pick=$4 AND group_letter=$5 AND slot_index=$6`,
		userID, tourID, category, pick, group, slot).Scan(&pts)
	require.NoError(t, err)
	return pts
}

func queryTeamPredictionPointsWildcard(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID, pick uuid.UUID) int {
	t.Helper()
	var pts int
	err := pool.QueryRow(ctx,
		`SELECT COALESCE(points,-999) FROM team_predictions WHERE user_id=$1 AND tournament_id=$2 AND category='playoff' AND pick=$3 AND slot_index=2`,
		userID, tourID, pick).Scan(&pts)
	require.NoError(t, err)
	return pts
}

func queryTeamPredictionPointsSemi(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID, pick uuid.UUID) int {
	t.Helper()
	var pts int
	err := pool.QueryRow(ctx,
		`SELECT COALESCE(points,-999) FROM team_predictions WHERE user_id=$1 AND tournament_id=$2 AND category='semifinalist' AND pick=$3`,
		userID, tourID, pick).Scan(&pts)
	require.NoError(t, err)
	return pts
}

func queryTeamPredictionPointsWinner(t *testing.T, pool *db.Pool, ctx context.Context, userID, tourID, pick uuid.UUID) int {
	t.Helper()
	var pts int
	err := pool.QueryRow(ctx,
		`SELECT COALESCE(points,-999) FROM team_predictions WHERE user_id=$1 AND tournament_id=$2 AND category='winner' AND pick=$3`,
		userID, tourID, pick).Scan(&pts)
	require.NoError(t, err)
	return pts
}
