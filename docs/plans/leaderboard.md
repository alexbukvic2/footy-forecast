# Plan: Leaderboard

## Goal

Expose per-league and global tournament leaderboards that rank users by total prediction points (score + player + team), and add a `my_position` field to each item in the `GET /leagues` response so users can see their current rank at a glance.

---

## Out of scope

- Caching or materialized views for leaderboard query performance.
- Historical leaderboard snapshots (point-in-time rank history).
- Pagination of leaderboard results.
- Tie-breaking by secondary criteria (e.g. earliest prediction submission). Ties share a rank.
- Per-fixture score-points breakdown in the response — only the aggregate `score_points` total is returned.
- Any admin endpoints for leaderboard data.

---

## Data model changes

No new migrations. The leaderboard is computed on the fly from existing tables: `score_predictions`, `player_predictions`, `team_predictions`, `league_members`, `users`, `fixtures`, `leagues`.

**Prerequisite**: The `score_predictions` table (migration `20260524000002_score_predictions.sql` from the score-predictions plan) must exist before the leaderboard query can include match points. If the leaderboard is shipped before that feature, the `score_pts` CTE must be omitted and `score_points` hard-coded to `0` for all users. This must be noted as a known limitation in the release.

---

## API Contract

### 1. GET /leagues/{leagueId}/leaderboard

Returns the ranked leaderboard for a league. The requesting user must be a league member.

- Method + path: `GET /leagues/{leagueId}/leaderboard`
- Auth: Bearer JWT

**Path parameters**

| Name | Type/Format | Constraints |
|------|-------------|-------------|
| leagueId | string, uuid | Valid UUID; must identify an existing league |

**Success response: 200**

| Field | Type | Description |
|-------|------|-------------|
| leaderboard | array of LeaderboardEntry | Ordered by `total_points DESC` |

**LeaderboardEntry**

| Field | Type | Description |
|-------|------|-------------|
| position | integer | DENSE_RANK (1-based; ties share position) |
| user_id | string (uuid) | User identifier |
| display_name | string | User display name |
| score_points | integer | Points from score predictions (0 if score-predictions not yet shipped) |
| player_points | integer | Points from player predictions |
| team_points | integer | Points from team predictions |
| total_points | integer | Sum of the three category totals |

**Error responses**

| Condition | Status | Body |
|-----------|--------|------|
| JWT missing or expired | 401 | `{"error": "unauthorized"}` |
| `leagueId` is not a valid UUID | 400 | `{"error": "invalid leagueId \"...\": invalid"}` |
| League does not exist | 404 | `{"error": "not found"}` |
| Requesting user is not a league member | 403 | `{"error": "forbidden"}` |
| Unexpected failure | 500 | `{"error": "internal server error"}` |

---

### 2. GET /tournaments/{tournamentId}/leaderboard

Returns the global ranked leaderboard for a tournament. Includes every user who has made at least one prediction for that tournament. Any authenticated user may call this.

- Method + path: `GET /tournaments/{tournamentId}/leaderboard`
- Auth: Bearer JWT

**Path parameters**

| Name | Type/Format | Constraints |
|------|-------------|-------------|
| tournamentId | string, uuid | Valid UUID; must identify an existing tournament |

**Success response: 200**

| Field | Type | Description |
|-------|------|-------------|
| leaderboard | array of LeaderboardEntry | Ordered by `total_points DESC` |

**Error responses**

| Condition | Status | Body |
|-----------|--------|------|
| JWT missing or expired | 401 | `{"error": "unauthorized"}` |
| `tournamentId` is not a valid UUID | 400 | `{"error": "invalid tournamentId \"...\": invalid"}` |
| Tournament does not exist | 404 | `{"error": "not found"}` |
| Unexpected failure | 500 | `{"error": "internal server error"}` |

---

### 3. GET /leagues (modified)

`LeagueListItem` gains a required `my_position` field. No change to path, method, auth, or error responses.

**Added field on LeagueListItem**

| Field | Type | Description |
|-------|------|-------------|
| my_position | integer | Requesting user's DENSE_RANK position in this league. Always present; equals 1 when all members are tied at 0. |

---

> Implementer: add `LeaderboardEntry`, `LeagueLeaderboardResponse`, and `TournamentLeaderboardResponse` schemas to `docs/openapi.yaml`. Add `my_position` (required, integer) to `LeagueListItem`. Add the two new paths under a `leaderboard` tag. Run `make generate` before writing any handler code.

---

## SQL approach

### sqlc vs. raw pgx

These queries use CTEs and `DENSE_RANK()` window functions. The `GetUserPositionsInLeagues` query also takes a dynamic `ANY($2)` array parameter for a variable list of league IDs. sqlc cannot cleanly generate code for these patterns. The leaderboard repository uses raw `pool.Query()` with `pgx.CollectRows()` directly. Query strings live as unexported `const` values inside `internal/repository/leaderboard.go`. No `.sql` files under `internal/repository/queries/` are needed for this feature.

### Per-league leaderboard (pseudo-SQL)

```sql
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
      COALESCE(s.pts, 0) AS score_points,
      COALESCE(p.pts, 0) AS player_points,
      COALESCE(t.pts, 0) AS team_points,
      COALESCE(s.pts, 0) + COALESCE(p.pts, 0) + COALESCE(t.pts, 0) AS total_points
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
ORDER BY total_points DESC;
```

If `score_predictions` is not yet available, drop the `score_pts` CTE and replace `COALESCE(s.pts, 0)` with `0`, removing the `LEFT JOIN score_pts` line.

### Global tournament leaderboard (pseudo-SQL)

```sql
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
      COALESCE(s.pts, 0) AS score_points,
      COALESCE(p.pts, 0) AS player_points,
      COALESCE(t.pts, 0) AS team_points,
      COALESCE(s.pts, 0) + COALESCE(p.pts, 0) + COALESCE(t.pts, 0) AS total_points
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
ORDER BY total_points DESC;
```

### User positions across all leagues for GET /leagues (pseudo-SQL)

A single round trip computes the requesting user's rank in each of their leagues using `PARTITION BY league_id`. The outer query filters to the requesting user after ranking.

```sql
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
WHERE user_id = $1;
```

`$1` = requesting user UUID, `$2` = `[]uuid.UUID` of their league IDs. Result is scanned into a `map[uuid.UUID]int`. The service merges this into the `[]*domain.LeagueSummary` slice returned to the handler. If `leagueIDs` is empty, skip the DB call and return an empty map.

---

## Touched files

| File | Reason |
|------|--------|
| `internal/domain/leaderboard.go` | New file. Define `LeaderboardEntry` struct. |
| `internal/domain/league.go` | Add `LeagueSummary` struct (embeds `*League`, adds `MyPosition int`). |
| `internal/repository/leaderboard.go` | New file. `LeaderboardRepository` with `GetForLeague`, `GetForTournament`, `GetUserPositionsInLeagues`. Raw pgx queries, no sqlc. |
| `internal/repository/leaderboard_test.go` | New file. Integration tests using testcontainers. |
| `internal/service/leaderboard.go` | New file. `LeaderboardService` with `GetLeagueLeaderboard` and `GetTournamentLeaderboard`; membership check. |
| `internal/service/leaderboard_test.go` | New file. Unit tests with hand-written fakes. |
| `internal/service/league.go` | Add `LeaderboardRepo` interface and dependency to `LeagueService`; change `ListLeaguesForUser` return type to `[]*domain.LeagueSummary`; call `GetUserPositionsInLeagues` and attach positions. |
| `internal/service/league_test.go` | Update `ListLeaguesForUser` tests for new return type and leaderboard repo fake. |
| `internal/server/handler/leaderboard.go` | New file. `LeaderboardHandler` with `GetLeagueLeaderboard` and `GetTournamentLeaderboard` methods. |
| `internal/server/handler/leaderboard_test.go` | New file. Handler tests. |
| `internal/server/handler/league.go` | Update `leagueListItemResponse` to add `MyPosition int \`json:"my_position"\``; update `List` handler mapping; update `LeagueService` interface to use `[]*domain.LeagueSummary`. |
| `internal/server/handler/league_test.go` | Update `List` handler test to assert `my_position` field. |
| `docs/openapi.yaml` | Add `LeaderboardEntry`, `LeagueLeaderboardResponse`, `TournamentLeaderboardResponse` schemas; extend `LeagueListItem` with `my_position`; add `GET /leagues/{leagueId}/leaderboard` and `GET /tournaments/{tournamentId}/leaderboard` paths; add `leaderboard` tag. |
| `docs/openapi.json` | Regenerated by `make generate` — do not edit manually. |
| `internal/server/oapi/models.gen.go` | Regenerated by `make generate` — do not edit manually. |
| `postman/collections/footy-forecast/leaderboard/GetLeagueLeaderboard.request.yaml` | New Postman request for league leaderboard. |
| `postman/collections/footy-forecast/leaderboard/GetTournamentLeaderboard.request.yaml` | New Postman request for global leaderboard. |
| `postman/collections/footy-forecast/leagues/List-Leagues.request.yaml` | Add `my_position` to expected response example. |

---

## Edge cases

| Case | Expected behaviour |
|------|--------------------|
| League member with no predictions | Included in results; all category points = 0; position determined by DENSE_RANK relative to others. |
| All league members at 0 points (tournament not started) | All receive position 1. |
| Empty league (no members) | `leaderboard` is an empty array, no error. |
| Single-member league | That member is at position 1. |
| Two users tied at same total | Both receive the same position; the next distinct total gets the next consecutive DENSE_RANK value (e.g. 1, 1, 2). |
| Global leaderboard when tournament has no predictions yet | Returns empty array. |
| `score_predictions` table not yet implemented | `score_points` = 0 for all users; `score_pts` CTE omitted. Known limitation. |
| `points` column is NULL on a prediction row (not yet scored) | `COALESCE(SUM(points), 0)` treats NULL contributions as 0. |
| Non-existent `leagueId` | Membership check returns `domain.ErrNotFound`; handler returns 404. |
| Non-existent `tournamentId` | Service tournament existence check returns `domain.ErrNotFound`; handler returns 404. |
| `leagueIDs` slice is empty in `GetUserPositionsInLeagues` | Short-circuit: return empty map, skip DB query. `ANY($2)` with an empty array is undefined across drivers. |
| User queries league leaderboard while not a member | Service returns `domain.ErrForbidden`; handler returns 403. |

---

## Test plan

### Repository tests (real Postgres via testcontainers)

File: `internal/repository/leaderboard_test.go`

- `GetForLeague` — 3-member league with varied points across all three categories: assert correct per-category values and positions.
- `GetForLeague` — tie scenario: two members with equal total; both get the same DENSE_RANK position; third gets next consecutive position.
- `GetForLeague` — one member with no predictions: included with all zeros; ranks last when others have points.
- `GetForLeague` — empty league: returns empty slice, nil error.
- `GetForTournament` — users spread across all three prediction sources: correct aggregation and ranking.
- `GetForTournament` — user with only player predictions, another with only team predictions: both appear; correct zeros for unused categories.
- `GetForTournament` — tournament with zero predictions: returns empty slice.
- `GetUserPositionsInLeagues` — user in two leagues with distinct positions: returned map contains correct position for each.
- `GetUserPositionsInLeagues` — empty `leagueIDs` slice: returns empty map, nil error, no DB round trip.

### Service tests (hand-written fakes)

File: `internal/service/leaderboard_test.go`

- `GetLeagueLeaderboard` — requesting user is not a member: fake `GetMember` returns `domain.ErrNotFound`; service returns `domain.ErrForbidden`.
- `GetLeagueLeaderboard` — requesting user is a member: delegates to leaderboard repo fake; returns result unchanged.
- `GetLeagueLeaderboard` — league does not exist: propagates `domain.ErrNotFound`.
- `GetTournamentLeaderboard` — tournament not found: returns `domain.ErrNotFound`.
- `GetTournamentLeaderboard` — valid tournament: delegates to repo; returns result unchanged.

File: `internal/service/league_test.go` (updates)

- `ListLeaguesForUser` — assert returned `LeagueSummary` slice has `MyPosition` populated from leaderboard repo.
- `ListLeaguesForUser` — leaderboard repo returns error: service wraps and returns error.

### Handler tests (httptest)

File: `internal/server/handler/leaderboard_test.go`

- `GetLeagueLeaderboard` happy path: 200; response JSON contains `leaderboard` array with `position`, `user_id`, `display_name`, `score_points`, `player_points`, `team_points`, `total_points`.
- `GetLeagueLeaderboard` — non-member: service fake returns `ErrForbidden`; 403.
- `GetLeagueLeaderboard` — invalid UUID path param: 400.
- `GetLeagueLeaderboard` — no JWT: 401.
- `GetTournamentLeaderboard` happy path: 200; response shape correct.
- `GetTournamentLeaderboard` — invalid UUID path param: 400.
- `GetTournamentLeaderboard` — no JWT: 401.

File: `internal/server/handler/league_test.go` (updates)

- `List` handler: assert `my_position` field is present and correct on each list item.

---

## Open questions

1. ~~**`LeagueSummary` vs. enriching `domain.League`**~~ **Decided: use `domain.LeagueSummary`**. Introduce `domain.LeagueSummary` embedding `*League` with `MyPosition int`. `ListLeaguesForUser` returns `[]*domain.LeagueSummary`; update the handler's `LeagueService` interface accordingly.

2. **Membership check round trip in `GetLeagueLeaderboard`**: the service calls `GetMember` (one DB round trip) then `GetForLeague` (a second round trip). An alternative is to fold the membership check into the leaderboard query — if the requesting user is not in `league_members` for that `league_id`, return an empty result and let the service infer non-membership. This is a minor optimisation and complicates the service logic; the separate `GetMember` call is the safe default.

3. **Tournament existence check for global leaderboard**: if the tournament does not exist, the leaderboard query returns an empty result rather than an error. The service should call `TournamentGetter.GetByID` first so the client receives a meaningful 404. `LeaderboardService` should accept `TournamentGetter` as a constructor dependency, consistent with `LeagueService`.

4. **Score predictions deployment order**: confirm whether score-predictions and leaderboard ship in the same release or sequentially. If sequential, the leaderboard repository must conditionally omit the `score_pts` CTE. The simplest path is to ship both in the same release.

5. **`ANY($2)` with empty array**: passing an empty `[]uuid.UUID` to `ANY($2)` in pgx is legal (results in no rows) but the behaviour should be verified in a test. The implementation must short-circuit before issuing the query when `leagueIDs` is empty to avoid any ambiguity.
