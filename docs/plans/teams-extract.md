# Plan: Teams Extract

**File:** `docs/plans/teams-extract.md`
**Date:** 2026-05-22
**Status:** Draft
**Supersedes (team columns only):** `docs/plans/players-search.md`

---

## 1. Goal

Extract team identity data (`name`, logo SVG) from the `players` table into a global `teams` table, making teams first-class entities that can be shared across tournaments, and update the player search query and response to JOIN against the new table.

---

## 2. Out of Scope

- Any REST endpoints for `teams` (no create, read, update, delete via API). Teams are populated out-of-band by a data import script.
- A `TeamRepository` or `TeamService` — no application code queries teams independently.
- Timestamps (`created_at` / `updated_at`) on the `teams` table.
- A `set_updated_at` trigger on `teams`.
- Cross-tournament deduplication of players. "France" as a `teams` row is reused, but "Kylian Mbappé at World Cup 2026" and "Kylian Mbappé at Euro 2028" remain separate `players` rows.
- Pagination or filtering by team.
- Any changes to the `tournaments`, `leagues`, or `users` tables.
- A standalone `domain.Team` type — team data reaches the handler via the JOIN result in `domain.PlayerSearchResult`.

---

## 3. Data Model Changes

### 3a. Migration strategy

The `players` table defined in `migrations/20260522000002_players.sql` has **never been applied to production** — the players-search plan is unimplemented. The correct approach is to **replace the migration file in place** rather than writing an `ALTER TABLE` migration.

`migrations/20260522000002_players.sql` is rewritten to create `teams` first, then `players` referencing `teams(id)`, omitting `team_name` and `team_flag` from `players`. The filename slug and timestamp remain the same because the migration has never been run.

### 3b. New table: `teams`

```
teams
  id         UUID        PRIMARY KEY DEFAULT gen_random_uuid()
  name       TEXT        NOT NULL
  logo       TEXT        NOT NULL DEFAULT ''
```

- `id`: UUID generated server-side by Postgres.
- `name`: the team name (e.g. `"France"`). See Open Questions on uniqueness.
- `logo`: raw SVG markup. Empty string is the sentinel for "no logo yet".
- No timestamps; no trigger.
- No index needed — `teams` is never queried independently by the application.

### 3c. Updated table: `players`

Columns removed compared to the prior plan:
- `team_name TEXT NOT NULL DEFAULT ''`
- `team_flag TEXT NOT NULL DEFAULT ''`

Column added:
- `team_id UUID NOT NULL REFERENCES teams(id)`

All other columns and constraints from `players-search.md` are unchanged.

### 3d. Replacement migration: `migrations/20260522000002_players.sql`

```sql
-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE OR REPLACE FUNCTION unaccent_immutable(text)
    RETURNS text
    LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT
    AS $$ SELECT unaccent($1) $$;

CREATE TABLE teams (
    id         UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT  NOT NULL,
    logo  TEXT  NOT NULL DEFAULT ''
);

CREATE TABLE players (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id   TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    tournament_id UUID        NOT NULL REFERENCES tournaments(id),
    team_id       UUID        NOT NULL REFERENCES teams(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT players_external_id_tournament_uq UNIQUE (external_id, tournament_id),
    CONSTRAINT players_name_length_chk CHECK (char_length(name) BETWEEN 1 AND 200)
);

CREATE INDEX players_tournament_id_idx ON players (tournament_id);
CREATE INDEX players_name_trgm_idx     ON players USING GIN (unaccent_immutable(name) gin_trgm_ops);

CREATE TRIGGER players_set_updated_at
    BEFORE UPDATE ON players
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS players_set_updated_at ON players;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS teams;
DROP FUNCTION IF EXISTS unaccent_immutable(text);
-- Extensions left in place intentionally.

-- +goose StatementEnd
```

Notes:
- `unaccent_immutable` is an `IMMUTABLE` wrapper required because `unaccent()` is `STABLE` and cannot be used directly in a GIN index expression.
- `set_updated_at()` is defined in `20260516101839_create_tournaments.sql` — do not redeclare it here.
- `teams` must be created before `players` due to the FK dependency.

---

## 4. API Surface

No new endpoints. The existing endpoint URL and HTTP semantics are unchanged:

```
GET /tournaments/{tournament_id}/players/search?q=<string>
```

**Response shape change** — `team_flag` → `team_logo`:

```json
{
  "players": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Kylian Mbappé",
      "team_name": "France",
      "team_logo": "<svg>...</svg>"
    }
  ]
}
```

All error responses are unchanged.

---

## 5. sqlc Query

File: `internal/repository/queries/player.sql`

```sql
-- name: SearchPlayers :many
SELECT p.id, p.name, t.name AS team_name, t.logo as team_logo, p.tournament_id
FROM players p
JOIN teams t ON t.id = p.team_id
WHERE p.tournament_id = $1
  AND unaccent_immutable(p.name) ILIKE '%' || unaccent_immutable($2) || '%' ESCAPE '\'
ORDER BY similarity(unaccent_immutable(p.name), unaccent_immutable($3)) DESC
LIMIT 5;
```

- `$1` — tournament UUID for scoping
- `$2` — wildcard-escaped query string for `ILIKE` filtering
- `$3` — raw trimmed query string for `similarity()` ranking (must not be escaped)
- `t.name AS team_name` preserves the `TeamName` field name in the generated `SearchPlayersRow` struct
- Run `make sqlc` after editing this file to regenerate `dbgen/`.

---

## 6. Touched Files

| File | Change                                                                                                      |
|---|-------------------------------------------------------------------------------------------------------------|
| `migrations/20260522000002_players.sql` | Replaced in full: adds `teams` table, `team_id UUID FK` on `players`, removes `team_name`/`team_flag`       |
| `internal/repository/queries/player.sql` | `SearchPlayers` query: JOIN to `teams`, `t.name AS team_name`, `t.logo as team_logo`, three bind params     |
| `internal/repository/dbgen/player.sql.go` | Regenerated: `SearchPlayersRow` gains `TeamLogo`, loses `TeamFlag`; `SearchPlayersParams` gains third param |
| `internal/repository/dbgen/models.go` | Regenerated: `Player` struct gains `TeamID`, loses `TeamName`/`TeamFlag`; new `Team` struct                 |
| `internal/domain/player.go` | `Player` struct: remove `TeamName`, `TeamFlag`; new `PlayerSearchResult` type (see §7)                      |
| `internal/repository/player.go` | `Search` return type → `[]*domain.PlayerSearchResult`; `playerFromSearchRow` maps `TeamLogo`                |
| `internal/service/player.go` | `PlayerRepo.Search` and `PlayerService.Search` return type → `[]*domain.PlayerSearchResult`                 |
| `internal/server/handler/player.go` | `PlayerService` interface updated; `playerResponse.TeamFlag` → `TeamLogo json:"team_logo"`                  |
| `internal/repository/player_test.go` | Fixtures: `insertTeam` helper added; `insertPlayer` uses `team_id`; assertions updated                      |
| `internal/service/player_test.go` | `fakePlayerRepo.Search` return type updated; test literals use `PlayerSearchResult`                         |
| `internal/server/handler/player_test.go` | Fake service return type updated; response assertion `team_flag` → `team_logo`                              |

---

## 7. Domain Type Changes

### `internal/domain/player.go`

```
// Before
Player { ID uuid.UUID, Name string, TournamentID uuid.UUID, TeamName string, TeamFlag string }

// After
Player { ID uuid.UUID, Name string, TournamentID uuid.UUID, TeamID uuid.UUID }

PlayerSearchResult { ID uuid.UUID, Name string, TournamentID uuid.UUID, TeamName string, TeamLogo string }
```

`PlayerSearchResult` carries the JOIN columns. It is the return type of `PlayerRepository.Search`, `PlayerService.Search`, and the input to `toPlayerResponse` in the handler.

`SearchPlayersInput` is unchanged.

No `domain.Team` type is introduced — no application code fetches a standalone team.

---

## 8. Repository Changes

File: `internal/repository/player.go`

- `Search` signature: returns `([]*domain.PlayerSearchResult, error)`.
- `playerFromSearchRow` (or equivalent mapping function): maps `row.TeamName` and `row.TeamLogo` into `PlayerSearchResult`.
- Compile-time check: `var _ PlayerRepo = (*PlayerRepository)(nil)` still applies; `PlayerRepo` in `service/player.go` must match.

---

## 9. Service Changes

File: `internal/service/player.go`

- `PlayerRepo` interface: `Search` return type → `[]*domain.PlayerSearchResult`.
- `PlayerService.Search` return type → `[]*domain.PlayerSearchResult`.
- Wildcard escaping (`%` → `\%`, `_` → `\_`) applied before passing `$2`; raw trimmed query passed as `$3`. This is unchanged from the prior plan.
- Tournament existence check via `TournamentGetter.GetByID` is unchanged.

---

## 10. Handler Changes

File: `internal/server/handler/player.go`

- `PlayerService` interface: `Search` return type → `[]*domain.PlayerSearchResult`.
- `playerResponse`: rename `TeamFlag string \`json:"team_flag"\`` → `TeamLogo string \`json:"team_logo"\``.
- `toPlayerResponse`: source type → `*domain.PlayerSearchResult`; maps `.TeamLogo` to `TeamLogo`.
- Routing, middleware, and error handling are unchanged.

---

## 11. Edge Cases

| Scenario | Handling |
|---|---|
| `team_id` FK violated on player insert | DB enforces FK; insert fails. Import script must insert `teams` row first. No application-layer handling needed (no insert endpoint). |
| `teams` row deleted while players reference it | FK with no `ON DELETE` clause; Postgres rejects the deletion. Teams must not be deleted before their players. |
| `team_logo` is empty string | `team_logo` defaults to `''`. Response includes `"team_logo": ""`. Clients handle empty string as "no logo". |
| Same team appears in multiple tournaments | Intentional — one `teams` row shared across tournaments via `players.team_id`. |
| Accidental duplicate `teams.name` during import | No DB uniqueness constraint by default (see Open Questions). Import script is responsible for deduplication. |
| sqlc field name mismatch after regeneration | `t.name AS team_name` in the query ensures the generated field is still `TeamName`, not `Name`. Verify after running `make sqlc`. |
| Integration test fixture missing `teams` row | FK violation produces a clear test failure. `insertPlayer` must call `insertTeam` first and pass the returned `id`. |

---

## 12. Test Plan

### Repository — integration tests

File: `internal/repository/player_test.go`

New helper: `insertTeam(t, db, name, teamLogo string) uuid.UUID` — inserts into `teams`, returns `id`.

Updated helper: `insertPlayer(t, db, ..., teamID uuid.UUID)` — replaces `teamName`/`teamFlag` params.

| Test | What it verifies |
|---|---|
| Search returns matching players | Fixture updated to use `team_id`; semantics unchanged |
| Search is case-insensitive | Unchanged |
| Search is tournament-scoped | Unchanged |
| Search respects LIMIT 5 | Unchanged |
| Returns empty slice on no match | Unchanged |
| Returns empty slice for nonexistent tournament | Unchanged |
| JOIN returns correct `team_name` and `team_logo` | New: assert `result[0].TeamName` and `result[0].TeamLogo` match the inserted team row |

### Service — unit tests

File: `internal/service/player_test.go`

- `fakePlayerRepo.Search` return type: `[]*domain.PlayerSearchResult`.
- Test literals: construct `domain.PlayerSearchResult` where team data is asserted; `domain.Player` literals no longer have `TeamName`/`TeamFlag`.
- All existing test cases preserved. No new cases needed at this layer.

### Handler — unit tests

File: `internal/server/handler/player_test.go`

- `fakePlayerService.Search` return type: `[]*domain.PlayerSearchResult`.
- `TestPlayer_Search_HappyPath`: assert `"team_logo"` key in response JSON (not `"team_flag"`); assert `"team_name"` still present.
- All other test cases unchanged.

---

## 13. Open Questions

1. **`PlayerSearchResult` naming.** The plan uses `PlayerSearchResult`. The implementer may prefer `PlayerWithTeam` or `SearchedPlayer`. Choose one and be consistent across all layers.

2. **`teams.name` uniqueness.** Should the DB enforce `UNIQUE (name)`? A unique constraint prevents accidental duplicate rows but requires the import script to `INSERT ... ON CONFLICT DO NOTHING` or upsert. Decide before writing the migration.

3. **`team_logo` size constraint.** Should `teams.logo` have `CHECK (octet_length(logo) <= 65536)`? Deferred from the prior plan.

4. **Teams population workflow.** This plan defers the import mechanism entirely. The implementer must confirm how `teams` rows are seeded before any `players` rows are inserted.

5. **`external_id` in response.** Deferred from the prior plan. If the frontend needs the api-football player ID, add `external_id` to `PlayerSearchResult` and `playerResponse`.
