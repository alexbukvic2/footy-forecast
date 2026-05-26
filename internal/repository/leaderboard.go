package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// LeaderboardRepository computes leaderboard data using raw pgx queries.
// These queries use CTEs and DENSE_RANK() window functions that sqlc cannot
// cleanly generate code for.
type LeaderboardRepository struct {
	pool *db.Pool
}

// NewLeaderboardRepository constructs a LeaderboardRepository.
func NewLeaderboardRepository(pool *db.Pool) *LeaderboardRepository {
	return &LeaderboardRepository{pool: pool}
}

const queryLeagueLeaderboard = `
WITH
  score_pts AS (
    SELECT sp.user_id, COALESCE(SUM(sp.points), 0) AS pts
    FROM score_predictions sp
    JOIN fixtures f ON f.id = sp.fixture_id
    WHERE f.tournament_id = (SELECT tournament_id FROM leagues WHERE id = $1)
    GROUP BY sp.user_id
  ),
  player_pts AS (
    SELECT user_id, COALESCE(SUM(points), 0) AS pts
    FROM player_predictions
    WHERE tournament_id = (SELECT tournament_id FROM leagues WHERE id = $1)
    GROUP BY user_id
  ),
  team_pts AS (
    SELECT user_id, COALESCE(SUM(points), 0) AS pts
    FROM team_predictions
    WHERE tournament_id = (SELECT tournament_id FROM leagues WHERE id = $1)
    GROUP BY user_id
  ),
  totals AS (
    SELECT
      lm.user_id,
      u.display_name,
      COALESCE(s.pts, 0)                                                       AS score_points,
      COALESCE(p.pts, 0)                                                       AS player_points,
      COALESCE(t.pts, 0)                                                       AS team_points,
      COALESCE(s.pts, 0) + COALESCE(p.pts, 0) + COALESCE(t.pts, 0)           AS total_points
    FROM league_members lm
    JOIN users u ON u.id = lm.user_id
    LEFT JOIN score_pts  s ON s.user_id = lm.user_id
    LEFT JOIN player_pts p ON p.user_id = lm.user_id
    LEFT JOIN team_pts   t ON t.user_id = lm.user_id
    WHERE lm.league_id = $1
  )
SELECT
  DENSE_RANK() OVER (ORDER BY total_points DESC) AS position,
  user_id, display_name, score_points, player_points, team_points, total_points
FROM totals
ORDER BY total_points DESC`

// GetForLeague returns the leaderboard for a league, ordered by total points descending.
// If the league has no members, an empty slice is returned.
func (r *LeaderboardRepository) GetForLeague(
	ctx context.Context,
	leagueID uuid.UUID,
) ([]*domain.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx, queryLeagueLeaderboard, leagueID)
	if err != nil {
		return nil, fmt.Errorf("query league leaderboard: %w", err)
	}
	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.LeaderboardEntry, error) {
		var e domain.LeaderboardEntry
		if err := row.Scan(
			&e.Position,
			&e.UserID,
			&e.DisplayName,
			&e.ScorePoints,
			&e.PlayerPoints,
			&e.TeamPoints,
			&e.TotalPoints,
		); err != nil {
			return nil, err
		}
		return &e, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect league leaderboard rows: %w", err)
	}
	return entries, nil
}

const queryTournamentLeaderboard = `
WITH
  score_pts AS (
    SELECT sp.user_id, COALESCE(SUM(sp.points), 0) AS pts
    FROM score_predictions sp
    JOIN fixtures f ON f.id = sp.fixture_id
    WHERE f.tournament_id = $1
    GROUP BY sp.user_id
  ),
  player_pts AS (
    SELECT user_id, COALESCE(SUM(points), 0) AS pts
    FROM player_predictions
    WHERE tournament_id = $1
    GROUP BY user_id
  ),
  team_pts AS (
    SELECT user_id, COALESCE(SUM(points), 0) AS pts
    FROM team_predictions
    WHERE tournament_id = $1
    GROUP BY user_id
  ),
  all_users AS (
    SELECT user_id FROM score_pts
    UNION
    SELECT user_id FROM player_pts
    UNION
    SELECT user_id FROM team_pts
  ),
  totals AS (
    SELECT
      au.user_id,
      u.display_name,
      COALESCE(s.pts, 0)                                                       AS score_points,
      COALESCE(p.pts, 0)                                                       AS player_points,
      COALESCE(t.pts, 0)                                                       AS team_points,
      COALESCE(s.pts, 0) + COALESCE(p.pts, 0) + COALESCE(t.pts, 0)           AS total_points
    FROM all_users au
    JOIN users u ON u.id = au.user_id
    LEFT JOIN score_pts  s ON s.user_id = au.user_id
    LEFT JOIN player_pts p ON p.user_id = au.user_id
    LEFT JOIN team_pts   t ON t.user_id = au.user_id
  )
SELECT
  DENSE_RANK() OVER (ORDER BY total_points DESC) AS position,
  user_id, display_name, score_points, player_points, team_points, total_points
FROM totals
ORDER BY total_points DESC`

// GetForTournament returns the global leaderboard for a tournament, ordered by
// total points descending. Only users who have made at least one prediction appear.
func (r *LeaderboardRepository) GetForTournament(
	ctx context.Context,
	tournamentID uuid.UUID,
) ([]*domain.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx, queryTournamentLeaderboard, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("query tournament leaderboard: %w", err)
	}
	entries, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*domain.LeaderboardEntry, error) {
		var e domain.LeaderboardEntry
		if err := row.Scan(
			&e.Position,
			&e.UserID,
			&e.DisplayName,
			&e.ScorePoints,
			&e.PlayerPoints,
			&e.TeamPoints,
			&e.TotalPoints,
		); err != nil {
			return nil, err
		}
		return &e, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect tournament leaderboard rows: %w", err)
	}
	return entries, nil
}

const queryUserPositionsInLeagues = `
WITH
  score_pts AS (
    SELECT sp.user_id, l.id AS league_id, COALESCE(SUM(sp.points), 0) AS pts
    FROM score_predictions sp
    JOIN fixtures f ON f.id = sp.fixture_id
    JOIN leagues l ON l.tournament_id = f.tournament_id
    WHERE l.id = ANY($2)
    GROUP BY sp.user_id, l.id
  ),
  player_pts AS (
    SELECT pp.user_id, l.id AS league_id, COALESCE(SUM(pp.points), 0) AS pts
    FROM player_predictions pp
    JOIN leagues l ON l.tournament_id = pp.tournament_id
    WHERE l.id = ANY($2)
    GROUP BY pp.user_id, l.id
  ),
  team_pts AS (
    SELECT tp.user_id, l.id AS league_id, COALESCE(SUM(tp.points), 0) AS pts
    FROM team_predictions tp
    JOIN leagues l ON l.tournament_id = tp.tournament_id
    WHERE l.id = ANY($2)
    GROUP BY tp.user_id, l.id
  ),
  ranked AS (
    SELECT
      lm.league_id,
      lm.user_id,
      DENSE_RANK() OVER (
        PARTITION BY lm.league_id
        ORDER BY (COALESCE(s.pts, 0) + COALESCE(p.pts, 0) + COALESCE(t.pts, 0)) DESC
      ) AS position
    FROM league_members lm
    LEFT JOIN score_pts  s ON s.user_id = lm.user_id AND s.league_id = lm.league_id
    LEFT JOIN player_pts p ON p.user_id = lm.user_id AND p.league_id = lm.league_id
    LEFT JOIN team_pts   t ON t.user_id = lm.user_id AND t.league_id = lm.league_id
    WHERE lm.league_id = ANY($2)
  )
SELECT league_id, position
FROM ranked
WHERE user_id = $1`

// GetUserPositionsInLeagues returns a map of leagueID → DENSE_RANK position for
// the given user across all specified leagues. If leagueIDs is empty the function
// returns an empty map without issuing a DB query.
func (r *LeaderboardRepository) GetUserPositionsInLeagues(
	ctx context.Context,
	userID uuid.UUID,
	leagueIDs []uuid.UUID,
) (map[uuid.UUID]int, error) {
	if len(leagueIDs) == 0 {
		return map[uuid.UUID]int{}, nil
	}

	rows, err := r.pool.Query(ctx, queryUserPositionsInLeagues, userID, leagueIDs)
	if err != nil {
		return nil, fmt.Errorf("query user positions in leagues: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]int, len(leagueIDs))
	for rows.Next() {
		var leagueID uuid.UUID
		var position int
		if err := rows.Scan(&leagueID, &position); err != nil {
			return nil, fmt.Errorf("scan user position row: %w", err)
		}
		result[leagueID] = position
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user position rows: %w", err)
	}
	return result, nil
}
