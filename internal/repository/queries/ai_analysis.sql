-- name: GetFixtureAnalysisInput :one
SELECT
    home_t.name              AS home_team_name,
    away_t.name              AS away_team_name,
    f.round,
    home_t.group_letter,
    f.goals_home,
    f.goals_away,
    coalesce(
        json_agg(json_build_object(
            'display_name', u.display_name,
            'goals_home',   sp.goals_home,
            'goals_away',   sp.goals_away,
            'points',       sp.points
        )) FILTER (WHERE u.id IS NOT NULL),
        '[]'::json
    )::text AS predictions
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
JOIN leagues l ON l.tournament_id = f.tournament_id AND l.id = @league_id
JOIN league_members lm ON lm.league_id = l.id
JOIN users u ON u.id = lm.user_id
LEFT JOIN score_predictions sp ON sp.fixture_id = f.id AND sp.user_id = lm.user_id
WHERE f.id = @fixture_id
GROUP BY f.id, f.round, f.goals_home, f.goals_away, home_t.name, away_t.name, home_t.group_letter;

-- name: UpsertFixtureAnalysis :exec
INSERT INTO score_ai_analysis (fixture_id, league_id, analysis)
VALUES (@fixture_id, @league_id, @analysis)
ON CONFLICT (fixture_id, league_id) DO UPDATE SET
    analysis   = EXCLUDED.analysis,
    updated_at = now();

-- name: ListLeaguesForFixture :many
SELECT l.id
FROM leagues l
JOIN fixtures f ON f.tournament_id = l.tournament_id
WHERE f.id = @fixture_id;
