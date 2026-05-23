# Plan: Handicap — Player & Team Handicap Points

**File:** `docs/plans/handicap.md`
**Date:** 2026-05-23
**Status:** Implemented

---

## 1. Goal

Introduce handicap points for players and teams. A handicap assigns a points value to a player or team within a specific category (e.g. a long-shot winner gets more points than a favourite). Expose authenticated read endpoints to look up these values, and always enrich the player search response with handicap points for every category.

---

## 2. Out of Scope

- Write endpoints for handicap (creation, update, delete). Population is handled by an out-of-band migration script (planned separately).
- Handicap lookup scoped to a tournament — the category alone is the scope.
- Team search endpoint (only player search is enhanced).
- Validating that a `player_id` or `team_id` actually belongs to any particular tournament.

---

## 3. Data Model Changes

### 3a. `teams.tournament_id` — new NOT NULL column

`teams` currently has no tournament affiliation. Add a FK column so each team belongs to exactly one tournament.

**Assumption:** the `teams` table is empty at migration time (development environment; test data is reseeded). The column is `NOT NULL` with no default.

### 3b. New enum types

| Postgres type | Values |
|---|---|
| `player_handicap_category` | `group_top_scorer`, `total_top_scorer` |
| `team_handicap_category` | `group_winner`, `playoff`, `semifinalist`, `winner` |

`playoff` means "reached the knockout stage" (advanced from group stage).

### 3c. `player_handicap` table

```sql
CREATE TABLE player_handicap (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id   UUID        NOT NULL REFERENCES players(id),
    category    player_handicap_category NOT NULL,
    points      INTEGER     NOT NULL,
    CONSTRAINT player_handicap_player_category_uq UNIQUE (player_id, category)
);
```

### 3d. `team_handicap` table

```sql
CREATE TABLE team_handicap (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id     UUID        NOT NULL REFERENCES teams(id),
    category    team_handicap_category NOT NULL,
    points      INTEGER     NOT NULL,
    CONSTRAINT team_handicap_team_category_uq UNIQUE (team_id, category)
);
```

### 3e. Migrations

**Migration 1:** `migrations/20260523000001_teams_tournament_id.sql`

Add `tournament_id` to `teams`.

```sql
-- +goose Up
-- +goose StatementBegin

ALTER TABLE teams
    ADD COLUMN tournament_id UUID NOT NULL REFERENCES tournaments(id);

CREATE INDEX teams_tournament_id_idx ON teams (tournament_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS teams_tournament_id_idx;
ALTER TABLE teams DROP COLUMN IF EXISTS tournament_id;

-- +goose StatementEnd
```

**Migration 2:** `migrations/20260523000002_handicap.sql`

Enum types and handicap tables.

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TYPE player_handicap_category AS ENUM (
    'group_top_scorer',
    'total_top_scorer'
);

CREATE TYPE team_handicap_category AS ENUM (
    'group_winner',
    'playoff',
    'semifinalist',
    'winner'
);

CREATE TABLE player_handicap (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id   UUID        NOT NULL REFERENCES players(id),
    category    player_handicap_category NOT NULL,
    points      INTEGER     NOT NULL,
    CONSTRAINT player_handicap_player_category_uq UNIQUE (player_id, category)
);

CREATE TABLE team_handicap (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id     UUID        NOT NULL REFERENCES teams(id),
    category    team_handicap_category NOT NULL,
    points      INTEGER     NOT NULL,
    CONSTRAINT team_handicap_team_category_uq UNIQUE (team_id, category)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS team_handicap;
DROP TABLE IF EXISTS player_handicap;
DROP TYPE IF EXISTS team_handicap_category;
DROP TYPE IF EXISTS player_handicap_category;

-- +goose StatementEnd
```

---

## 4. API Surface

### 4a. Get player handicap

```
GET /player-handicaps/{player_id}/{category}
```

Protected by `authMW`.

| Parameter | In | Type | Description |
|---|---|---|---|
| `player_id` | path | UUID string | Player to look up |
| `category` | path | string | One of: `group_top_scorer`, `total_top_scorer` |

**200 OK**
```json
{ "points": 5 }
```

**Error responses**

| Condition | Status | Body |
|---|---|---|
| Invalid/missing JWT | 401 | `{"error": "unauthorized"}` |
| `player_id` not a valid UUID | 400 | `{"error": "invalid player_id: ..."}` |
| `category` not a recognised value | 400 | `{"error": "invalid category: invalid"}` |
| No handicap row for player + category | 404 | `{"error": "not found"}` |
| Unexpected error | 500 | `{"error": "internal server error"}` |

---

### 4b. Get team handicap

```
GET /team-handicaps/{team_id}/{category}
```

Protected by `authMW`.

| Parameter | In | Type | Description |
|---|---|---|---|
| `team_id` | path | UUID string | Team to look up |
| `category` | path | string | One of: `group_winner`, `playoff`, `semifinalist`, `winner` |

**200 OK**
```json
{ "points": 10 }
```

**Error responses** — same shape as 4a, using `team_id` and `team_handicap_category` values.

---

### 4c. Player search — always includes all handicap categories

```
GET /tournaments/{tournament_id}/players/search?q=<string>
```

Handicap data is **always** included in the response — no optional query parameter. Every player in the response carries a `handicaps` map with an entry for every known `player_handicap_category`. If a player has no row in `player_handicap` for a given category, the service fills in a default of **20 points**.

**200 OK**
```json
{
  "players": [
    {
      "id": "uuid-string",
      "name": "Kylian Mbappé",
      "team_name": "France",
      "team_logo": "<svg>...</svg>",
      "handicaps": {
        "group_top_scorer": 5,
        "total_top_scorer": 20
      }
    }
  ]
}
```

---

## 5. Design Decisions

### 5a. Handicap always returned; default of 20

The original plan had `handicap_category` as an optional query parameter that would optionally enrich the response with a single category's points (`null` when missing). This was changed to:

- **Always return all categories** — no optional parameter.
- **Default value of 20** when a player has no `player_handicap` row for a category.
- **Default lives in the service layer** (`service/player.go`, constant `defaultHandicapPoints = 20`), not the repository or DB. The repository faithfully returns only what the DB has; the service fills gaps.

### 5b. Single SQL query using CTE + `json_object_agg`

The original plan proposed two separate queries (`SearchPlayers` and `SearchPlayersWithHandicap`) and the initial implementation iterated players and fetched handicap rows separately (N+1). Both were replaced with a single query:

```sql
-- name: SearchPlayers :many
WITH top_players AS (
    SELECT p.id, p.name, t.name AS team_name, t.logo AS team_logo, p.tournament_id
    FROM players p
    JOIN teams t ON t.id = p.team_id
    WHERE p.tournament_id = @tournament_id
      AND unaccent_immutable(p.name) ILIKE '%' || unaccent_immutable(@escaped_query) || '%' ESCAPE '\'
    ORDER BY similarity(unaccent_immutable(p.name), unaccent_immutable(@raw_query)) DESC
    LIMIT 5
)
SELECT tp.id, tp.name, tp.team_name, tp.team_logo, tp.tournament_id,
       coalesce(
           json_object_agg(ph.category, ph.points) FILTER (WHERE ph.category IS NOT NULL),
           '{}'::json
       )::text AS handicaps
FROM top_players tp
LEFT JOIN player_handicap ph ON ph.player_id = tp.id
GROUP BY tp.id, tp.name, tp.team_name, tp.team_logo, tp.tournament_id
ORDER BY similarity(unaccent_immutable(tp.name), unaccent_immutable(@raw_query)) DESC;
```

Key design points:
- **CTE applies `LIMIT 5` before the JOIN** — without the CTE, the LIMIT would apply to the post-JOIN rows, returning fewer than 5 players.
- **Single unfiltered `LEFT JOIN`** on `player_id` only — no category filter in the JOIN. All handicap rows for a player come back and are collapsed by `GROUP BY`.
- **`json_object_agg` + `FILTER`** aggregates all handicap rows into one JSON object per player. The `FILTER (WHERE ph.category IS NOT NULL)` avoids a spurious `{null: null}` entry when the LEFT JOIN yields no match.
- **`::text` cast** on the `coalesce` result — sqlc cannot infer the return type of the aggregate expression; casting to `text` ensures sqlc generates `string` rather than `interface{}`, enabling clean `json.Unmarshal` in the repository.

### 5c. No Go-level row aggregation in the repository

Each player is now returned as exactly one row by Postgres (via `GROUP BY`). The repository performs a simple 1:1 mapping: unmarshal the JSON text column into `map[domain.PlayerHandicapCategory]int`.

---

## 6. Touched Files

| File | Action | Reason |
|---|---|---|
| `migrations/20260523000001_teams_tournament_id.sql` | Create | Add `tournament_id` column and index to `teams` |
| `migrations/20260523000002_handicap.sql` | Create | Enum types + `player_handicap` and `team_handicap` tables |
| `internal/domain/player_handicap.go` | Create | `PlayerHandicapCategory` enum constants, `AllPlayerHandicapCategories` slice, `PlayerHandicap` type |
| `internal/domain/team_handicap.go` | Create | `TeamHandicapCategory` enum constants, `TeamHandicap` type |
| `internal/domain/player.go` | Modify | `PlayerSearchResult.Handicaps` is `map[PlayerHandicapCategory]int`; `SearchPlayersInput` has no handicap param |
| `internal/repository/queries/player_handicap.sql` | Create | `GetPlayerHandicap` sqlc query |
| `internal/repository/queries/team_handicap.sql` | Create | `GetTeamHandicap` sqlc query |
| `internal/repository/queries/player.sql` | Modify | `SearchPlayers` uses CTE + `json_object_agg` + `GROUP BY`; `::text` cast on handicaps column |
| `internal/repository/dbgen/` | Regenerate | Run `make sqlc-gen` after SQL changes — do not edit manually |
| `internal/repository/player_handicap.go` | Create | `PlayerHandicapRepository` with `Get` method |
| `internal/repository/team_handicap.go` | Create | `TeamHandicapRepository` with `Get` method |
| `internal/repository/player.go` | Modify | Simple loop + JSON unmarshal; no Go-level row aggregation |
| `internal/service/player_handicap.go` | Create | `PlayerHandicapService` + `PlayerHandicapRepo` interface |
| `internal/service/team_handicap.go` | Create | `TeamHandicapService` + `TeamHandicapRepo` interface |
| `internal/service/player.go` | Modify | Fills default handicap (20) for missing categories after repo call |
| `internal/server/handler/player_handicap.go` | Create | `PlayerHandicap` handler, local `PlayerHandicapService` interface, `Get` method |
| `internal/server/handler/team_handicap.go` | Create | `TeamHandicap` handler, local `TeamHandicapService` interface, `Get` method |
| `internal/server/handler/player.go` | Modify | Response always includes `handicaps` map; no `handicap_category` query param |
| `internal/server/router.go` | Modify | Wire new repos, services, handlers; register two new routes |
| `internal/repository/player_handicap_test.go` | Create | Integration tests (testcontainers) |
| `internal/repository/team_handicap_test.go` | Create | Integration tests (testcontainers) |
| `internal/service/player_test.go` | Modify | Tests for default filling, all-categories-present, mixed presence |
| `internal/server/handler/player_handicap_test.go` | Create | Handler unit tests |
| `internal/server/handler/team_handicap_test.go` | Create | Handler unit tests |
| `postman/collections/footy-forecast/player-handicaps/Get Player Handicap.request.yaml` | Create | New request |
| `postman/collections/footy-forecast/team-handicaps/Get Team Handicap.request.yaml` | Create | New request |
| `postman/collections/footy-forecast/players/Search Players.request.yaml` | Modify | Document `handicaps` map in response |

---

## 7. Edge Cases

| Scenario | Handling |
|---|---|
| `player_id` path value is not a valid UUID | `parseUUIDPathValue` returns `domain.ErrInvalid` → 400 |
| `team_id` path value is not a valid UUID | Same → 400 |
| `category` path value not in the known enum | Handler validates against domain constants; returns `domain.ErrInvalid` → 400 |
| No handicap row for the given player + category | Repository returns `domain.ErrNotFound` → 404 |
| No handicap row for the given team + category | Same → 404 |
| Player in search results has no `player_handicap` row | LEFT JOIN returns no rows for that player; `json_object_agg FILTER` returns `'{}'`; service fills default 20 for all missing categories |
| Player has some but not all categories | Service fills default 20 only for missing categories; existing values are preserved |
| `teams.tournament_id` migration on non-empty table | Migration fails — expected in production; table must be empty. Document as a pre-condition. |

---

## 8. Test Plan

### Repository — integration tests (`//go:build integration`)

**`player_handicap_test.go`**

| Test | What it verifies |
|---|---|
| Get returns points for existing row | Insert row, query by player + category, assert `Points` |
| Get returns `ErrNotFound` for unknown player | No row exists → `domain.ErrNotFound` |
| Get returns `ErrNotFound` for wrong category | Row exists for other category → `domain.ErrNotFound` |
| Uniqueness constraint | Insert duplicate (player_id, category) → Postgres error |

**`team_handicap_test.go`** — same four cases for teams.

### Service — unit tests

**`player_test.go`:**

| Test | What it verifies |
|---|---|
| All categories present when repo returns empty handicaps map | Service fills all categories with default 20 |
| Default 20 used when repo returns no handicap row for a category | |
| Repo points used when handicap row exists | |
| Mixed: default applied only for missing categories | |

### Handler — unit tests

**`player_handicap_test.go`:**

| Test | What it verifies |
|---|---|
| Happy path → 200 `{"points": N}` | |
| Invalid `player_id` UUID → 400 | |
| Unknown category string → 400 | |
| Service `ErrNotFound` → 404 | |
| Service unexpected error → 500 | |

**`team_handicap_test.go`** — same shape.
