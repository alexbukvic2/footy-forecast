# Plan: league-outrights — tournament group table

Introduce a `tournament_group_table` DB table that stores each team's position
within their tournament group (data sourced externally, not computed). Expose
the data via a new protected `GET /tournaments/{tournamentId}/group-table`
endpoint and replace the existing `/leagues/{leagueId}/predictions/players` and
`/leagues/{leagueId}/predictions/teams` endpoints with two stage-oriented
endpoints (`/predictions/groups` and `/predictions/playoff`) that return combined
team + player predictions per stage. Follows the layered architecture throughout:
migration → domain → sqlc query → repository → service → handler → OpenAPI spec
→ router → Postman.

---

## 1. Goal

Allow clients to fetch the current/final group-stage standings for a tournament,
enriched with team names, sorted by group and position.

---

## 2. Out of Scope

- Computing W/D/L/GF/GA/GD — standings come from an external API.
- An ingest/upsert endpoint — data is loaded via seed scripts or admin tooling
  (planned separately if needed).
- Knockout-bracket or outright winner predictions (separate feature).

---

## 3. Data Model Changes

New table `tournament_group_table`:

| Column | Type | Notes |
|---|---|---|
| `id` | UUID PK | `gen_random_uuid()` |
| `tournament_id` | UUID FK → tournaments | NOT NULL |
| `team_id` | UUID FK → teams | NOT NULL |
| `group_letter` | CHAR(1) | NOT NULL |
| `position` | SMALLINT | NOT NULL, 1-based |
| `points` | SMALLINT | NOT NULL |
| `created_at` | TIMESTAMPTZ | NOT NULL DEFAULT now() |
| `updated_at` | TIMESTAMPTZ | NOT NULL DEFAULT now() |

Constraints / indexes:
- `UNIQUE (tournament_id, team_id)`
- `INDEX ON tournament_id`
- `set_updated_at()` trigger on `updated_at` (trigger already exists in DB from prior migrations)

Migration file: `migrations/20260527000001_tournament_group_table.sql`

---

## 4. API Contract

### 4.1 `GET /tournaments/{tournamentId}/group-table`

**Auth:** `bearerAuth` (JWT required)

**Path parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `tournamentId` | string (uuid) | yes | Tournament UUID |

**Success response — 200 OK**

Array of `TournamentGroupTableEntry`, sorted by `group_letter ASC, position ASC`.

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "tournament_id": "550e8400-e29b-41d4-a716-446655440000",
    "team_id": "550e8400-e29b-41d4-a716-446655440010",
    "team_name": "Argentina",
    "group_letter": "A",
    "position": 1,
    "points": 9
  }
]
```

**New schema (`docs/openapi.yaml`):**

```yaml
TournamentGroupTableEntry:
  type: object
  required: [id, tournament_id, team_id, team_name, group_letter, position, points]
  properties:
    id:
      type: string
      format: uuid
    tournament_id:
      type: string
      format: uuid
    team_id:
      type: string
      format: uuid
    team_name:
      type: string
    group_letter:
      type: string
      minLength: 1
      maxLength: 1
    position:
      type: integer
    points:
      type: integer
```

**Error responses** (all `$ref: '#/components/schemas/ErrorResponse'`):

| Condition | Status | Body |
|---|---|---|
| `tournamentId` is not a valid UUID | 400 | `{"error": "..."}` |
| JWT missing or expired | 401 | `{"error": "unauthorized"}` |
| Tournament not found | 404 | `{"error": "not found"}` |
| Unexpected server error | 500 | `{"error": "internal server error"}` |

> **Implementer:** add this endpoint and the `TournamentGroupTableEntry` schema
> to `docs/openapi.yaml` and run `make generate` **before** writing any handler code.

---

### 4.2 `GET /leagues/{leagueId}/predictions/groups` *(replaces `/predictions/teams` + `/predictions/players` for group stage)*

**Auth:** `bearerAuth` (JWT required)

**Path parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `leagueId` | string (uuid) | yes | League UUID |

**Query parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `group` | string (1 char, e.g. `"A"`) | yes | Group letter to filter by |

Returns all league members' group-stage predictions for the given group, combining:
- **Team** category `group_winner` (which group teams advance from this group)
- **Player** category `group_top_scorer` (who scores most in this group)

**Success response — 200 OK**

```json
{
  "group": "A",
  "team_predictions": [
    {
      "category": "group_winner",
      "predictions": [
        {
          "user_id": "...",
          "display_name": "Alice",
          "team_id": "...",
          "team_name": "Argentina",
          "points": 3
        }
      ]
    }
  ],
  "player_predictions": [
    {
      "category": "group_top_scorer",
      "predictions": [
        {
          "user_id": "...",
          "display_name": "Alice",
          "player_id": "...",
          "player_name": "Messi",
          "points": null
        }
      ]
    }
  ]
}
```

**New schemas (`docs/openapi.yaml`):** `LeagueGroupPredictions`,
`LeagueTeamCategoryPredictions`, `LeaguePlayerCategoryPredictions`,
`LeagueMemberTeamPick`, `LeagueMemberPlayerPick`.

**Error responses:**

| Condition | Status | Body |
|---|---|---|
| `leagueId` or `group` param missing/invalid | 400 | `{"error": "..."}` |
| JWT missing or expired | 401 | `{"error": "unauthorized"}` |
| Caller not a member of the league | 403 | `{"error": "forbidden"}` |
| League not found | 404 | `{"error": "not found"}` |
| Unexpected server error | 500 | `{"error": "internal server error"}` |

---

### 4.3 `GET /leagues/{leagueId}/predictions/playoff` *(replaces `/predictions/teams` + `/predictions/players` for knockout stage)*

**Auth:** `bearerAuth` (JWT required)

**Path parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `leagueId` | string (uuid) | yes | League UUID |

Returns all league members' knockout/outright predictions, combining:
- **Team** categories: `playoff`, `semifinalist`, `winner`
- **Player** category: `total_top_scorer`

**Success response — 200 OK**

```json
{
  "team_predictions": [
    {
      "category": "playoff",
      "slot_index": 0,
      "predictions": [
        {
          "user_id": "...",
          "display_name": "Alice",
          "team_id": "...",
          "team_name": "Argentina",
          "points": null
        }
      ]
    }
  ],
  "player_predictions": [
    {
      "category": "total_top_scorer",
      "predictions": [
        {
          "user_id": "...",
          "display_name": "Alice",
          "player_id": "...",
          "player_name": "Messi",
          "points": null
        }
      ]
    }
  ]
}
```

**New schema (`docs/openapi.yaml`):** `LeaguePlayoffPredictions` (reuses
`LeagueTeamCategoryPredictions` and `LeaguePlayerCategoryPredictions` from 4.2).

**Error responses:** same table as 4.2.

---

### 4.4 Routes removed

The following existing routes are **deleted** as part of this plan. The spec
entries, handler methods, router registrations, and Postman requests for both
must be removed:

| Old route | Replaced by |
|---|---|
| `GET /leagues/{leagueId}/predictions/players` | 4.2 + 4.3 |
| `GET /leagues/{leagueId}/predictions/teams` | 4.2 + 4.3 |

> **Implementer:**
> 1. Add endpoints 4.2 and 4.3 and their schemas to `docs/openapi.yaml`.
> 2. Remove the old `/predictions/players` and `/predictions/teams` path entries.
> 3. Run `make generate` **before** writing any handler code.
> 4. Delete the corresponding Postman request files and add new ones.

---

## 5. Touched Files

| File | Reason |
|---|---|
| `migrations/20260527000001_tournament_group_table.sql` | New migration — creates the table, constraints, index, trigger |
| `internal/domain/tournament_group.go` | New domain type `TournamentGroupEntry` (includes `TeamName`) |
| `internal/repository/queries/tournament_group_table.sql` | sqlc query: `ListGroupTableByTournament` with join to `teams` |
| `internal/repository/tournament_group_table.go` | New repository struct + `ListByTournament` |
| `internal/repository/tournament_group_table_test.go` | Repository integration test (testcontainers) |
| `internal/service/tournament_group_table.go` | New service struct + `ListByTournament` |
| `internal/service/tournament_group_table_test.go` | Service unit tests with fake repo |
| `internal/server/handler/tournament_group_table.go` | New handler `ListGroupTable` |
| `internal/server/handler/tournament_group_table_test.go` | Handler tests (happy + error paths) |
| `internal/service/tournament_prediction.go` | Add `ListLeagueGroupPredictions(ctx, leagueID, groupLetter)` and `ListLeaguePlayoffPredictions(ctx, leagueID)` methods; keep existing `ListLeaguePlayerPredictions` / `ListLeagueTeamPredictions` until old routes are removed |
| `internal/service/tournament_prediction_test.go` | Tests for the two new service methods |
| `internal/server/handler/tournament_prediction.go` | Add `ListLeagueGroupPredictions` and `ListLeaguePlayoffPredictions` handlers; remove `ListLeaguePlayerPredictions` and `ListLeagueTeamPredictions` |
| `internal/server/handler/tournament_prediction_test.go` | Tests for the two new handlers |
| `docs/openapi.yaml` | New paths + schemas; remove old `/predictions/players` and `/predictions/teams` paths |
| `docs/openapi.json` | Regenerated by `make generate` — do not edit directly |
| `internal/server/oapi/models.gen.go` | Regenerated by `make generate` — do not edit directly |
| `internal/server/router.go` | Register new routes; deregister old routes |
| `postman/collections/footy-forecast/tournaments/GetGroupTable.request.yaml` | New Postman request |
| `postman/collections/footy-forecast/leagues/GetLeagueGroupPredictions.request.yaml` | New Postman request (replaces old players/teams requests) |
| `postman/collections/footy-forecast/leagues/GetLeaguePlayoffPredictions.request.yaml` | New Postman request |
| `postman/collections/footy-forecast/leagues/GetLeaguePlayerPredictions.request.yaml` | **Delete** — replaced |
| `postman/collections/footy-forecast/leagues/GetLeagueTeamPredictions.request.yaml` | **Delete** — replaced |

---

## 6. Edge Cases

| Case | Handling |
|---|---|
| Tournament UUID syntactically invalid | `uuid.Parse` fails → 400 before hitting DB |
| Tournament UUID valid but does not exist | Repository returns `domain.ErrNotFound` → 404 |
| Tournament exists but group table is empty | Return `[]` (empty array), 200 — not an error |
| Team FK references a deleted team | Should not happen with FK constraint; if it does, the join silently drops the row — add a `NOT NULL` assert or log a `Warn` |
| `set_updated_at` trigger missing | Migration must be idempotent — create trigger only if not exists; check existing migrations for the pattern used |
| `group` query param missing on `/predictions/groups` | Return 400 — it is required |
| `group` value is more than one character | Return 400 — validate before hitting DB |
| League UUID invalid on prediction endpoints | `uuid.Parse` fails → 400 |
| Caller not a member of the league | Service checks membership; returns `domain.ErrForbidden` → 403 |
| League not found | Service returns `domain.ErrNotFound` → 404 |
| No predictions exist for the requested group/stage | Return empty `team_predictions` / `player_predictions` arrays, 200 |

---

## 7. Test Plan

**Repository (`tournament_group_table_test.go`, testcontainers-go):**
- Seed a tournament + teams + group-table rows; assert correct count, order
  (`group_letter ASC, position ASC`), and correct `team_name` from the join.
- Assert `domain.ErrNotFound` when tournament UUID does not exist.
- Assert empty slice (not error) when tournament exists but has no rows.

**Service — group table (`tournament_group_table_test.go`, hand-written fake repo):**
- Table-driven: invalid UUID → error; valid UUID, repo returns rows → rows
  passed through; repo returns `domain.ErrNotFound` → `domain.ErrNotFound`.

**Service — league predictions (`tournament_prediction_test.go`, hand-written fake repo):**
- `ListLeagueGroupPredictions`: valid group letter → correct team+player category views; caller not a member → `domain.ErrForbidden`; unknown league → `domain.ErrNotFound`.
- `ListLeaguePlayoffPredictions`: same membership/notfound cases; assert only playoff/semifinalist/winner team categories and total_top_scorer player category are returned.

**Handler (`tournament_group_table_test.go`, `tournament_prediction_test.go`):**
- Group table: happy path → 200; bad UUID → 400; not found → 404; server error → 500.
- League groups: happy path → 200 with correct `group`, `team_predictions`, `player_predictions` fields; missing `group` param → 400; bad UUID → 400; forbidden → 403; not found → 404; server error → 500.
- League playoff: happy path → 200 with correct shape; bad UUID → 400; forbidden → 403; not found → 404; server error → 500.

---

## 8. Open Questions

1. **Tournament existence check placement** — should the repository verify the
   tournament exists (extra `SELECT` or subquery) or should the service call
   `TournamentService.GetByID` first? Prefer the service-level check to avoid
   coupling the new repo to the tournaments table; confirm with implementer.
2. **`set_updated_at` trigger** — confirm the trigger function already exists
   in a prior migration and does not need to be recreated here.
3. **Param naming convention** — existing routes use `{tournamentId}` (camelCase)
   in the mux. Confirm this is intentional and continue the pattern.
4. **`slot_index` in playoff response** — `playoff` category supports multiple slots (e.g. 16 teams that make the round of 16). The playoff response shape includes `slot_index` on each category row; confirm this is the right grouping or whether each slot should be a separate row.

