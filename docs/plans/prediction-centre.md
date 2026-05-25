# Plan: Prediction Centre

## Goal

Implement the prediction centre: a set of endpoints that return **all possible
prediction slots** for a tournament — not only the ones a user has already
filled. Separate tabs on the client show score predictions (all fixtures),
player predictions (all group/category combos), and team predictions (all
group/category/slot combos), each with the user's pick or `null`.

This plan also:
- Refactors `player_predictions` and `team_predictions` to support per-group
  and per-slot picks (the current schema assumes one pick per category, which is wrong).
- Implements score predictions (planned but never built in commit c46ec5d).
- Adds a `tournament_outcomes` table for the scoring worker.

---

## Design decisions

### Group letter on `teams`, not a separate table

`teams.group_letter VARCHAR(1)` is sufficient. The service queries
`SELECT DISTINCT group_letter FROM teams WHERE tournament_id = $1` to enumerate
groups. No `tournament_groups` table needed.

### Slot model for multi-pick categories

| Category | Scoped | Slots per scope |
|---|---|---|
| `group_top_scorer` | per group | 1 (slot 0) |
| `total_top_scorer` | tournament | 1 (slot 0) |
| `group_winner` | per group | 1 (slot 0) |
| `playoff` | per group | 3 (slots 0, 1 = direct qualifiers; slot 2 = wildcard 3rd-place) |
| `semifinalist` | tournament | 4 (slots 0–3, unordered) |
| `winner` | tournament | 1 (slot 0) |

`slot_index` is added to `team_predictions` only. Player predictions never need
more than one slot per group, so `player_predictions` only gains `group_letter`.

For `playoff` wildcards: max 8 slot-2 picks per user per tournament (12 groups,
choose 8). The "max 8" rule is enforced in the service, not the DB.
"At most 1 per group" is enforced by the unique index (slot 2 can appear only
once per group per user).

### Unique constraint replacement

The old `UNIQUE (user_id, tournament_id, category)` is replaced by two partial
unique indexes on each prediction table, avoiding NULL equality issues:

```sql
-- group-scoped rows
UNIQUE (user_id, tournament_id, category, group_letter, <slot_index for teams>)
WHERE group_letter IS NOT NULL

-- tournament-scoped rows
UNIQUE (user_id, tournament_id, category, <slot_index for teams>)
WHERE group_letter IS NULL
```

Because of these partial indexes, the sqlc upsert queries must also become
partial (`ON CONFLICT (...) WHERE group_letter IS NOT NULL/NULL`). Two named
queries per table (group-upsert and no-group-upsert); the repository routes
based on whether GroupLetter is set.

### Group letter in predictions is for UI/uniqueness only

The scoring worker joins on `pick = outcome.entity_id` — it never needs to
match `group_letter`. A user can only pick Germany as `group_winner` for
Germany's actual group (validated at submission), so the pick alone is
sufficient for scoring.

### Outcome storage

A `tournament_outcomes` table (two variants: player and team) is populated by
a cron worker after each stage ends. The scoring worker reads from this table —
no API calls, no standings re-computation at scoring time.

```
player_outcomes: (tournament_id, category, player_id) — one row per fact
team_outcomes:   (tournament_id, category, team_id)   — one row per fact
```

Group letter is not stored on outcomes; it is derivable from `teams.group_letter`
and is never needed for scoring.

### Score predictions

`score_predictions` is implemented per `docs/plans/score-predictions.md`.
Create the file as `20260524000004_score_predictions.sql` (the original plan
picked `000002`, which is already taken by `tournament_predictions`).

---

## Migrations

### `20260524000003_prediction_groups_slots.sql`

```sql
-- +goose Up
-- +goose StatementBegin

-- 1. Group letter on teams
ALTER TABLE teams
    ADD COLUMN group_letter VARCHAR(1);

-- 2. player_predictions: add group_letter, replace unique constraint
ALTER TABLE player_predictions
    ADD COLUMN group_letter VARCHAR(1);

ALTER TABLE player_predictions
    DROP CONSTRAINT player_predictions_user_tournament_category_uq;

CREATE UNIQUE INDEX player_predictions_group_uq
    ON player_predictions (user_id, tournament_id, category, group_letter)
    WHERE group_letter IS NOT NULL;

CREATE UNIQUE INDEX player_predictions_no_group_uq
    ON player_predictions (user_id, tournament_id, category)
    WHERE group_letter IS NULL;

-- 3. team_predictions: add group_letter + slot_index, replace unique constraint
ALTER TABLE team_predictions
    ADD COLUMN group_letter VARCHAR(1),
    ADD COLUMN slot_index   SMALLINT NOT NULL DEFAULT 0;

ALTER TABLE team_predictions
    DROP CONSTRAINT team_predictions_user_tournament_category_uq;

-- group-scoped (group_winner, playoff): unique per (category, group, slot)
CREATE UNIQUE INDEX team_predictions_group_uq
    ON team_predictions (user_id, tournament_id, category, group_letter, slot_index)
    WHERE group_letter IS NOT NULL;

-- tournament-scoped (semifinalist, winner): unique per (category, slot)
CREATE UNIQUE INDEX team_predictions_no_group_uq
    ON team_predictions (user_id, tournament_id, category, slot_index)
    WHERE group_letter IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS team_predictions_no_group_uq;
DROP INDEX IF EXISTS team_predictions_group_uq;
ALTER TABLE team_predictions
    DROP COLUMN IF EXISTS slot_index,
    DROP COLUMN IF EXISTS group_letter;
ALTER TABLE team_predictions
    ADD CONSTRAINT team_predictions_user_tournament_category_uq
        UNIQUE (user_id, tournament_id, category);

DROP INDEX IF EXISTS player_predictions_no_group_uq;
DROP INDEX IF EXISTS player_predictions_group_uq;
ALTER TABLE player_predictions
    DROP COLUMN IF EXISTS group_letter;
ALTER TABLE player_predictions
    ADD CONSTRAINT player_predictions_user_tournament_category_uq
        UNIQUE (user_id, tournament_id, category);

ALTER TABLE teams DROP COLUMN IF EXISTS group_letter;

-- +goose StatementEnd
```

### `20260524000004_score_predictions.sql`

Identical to the content in `docs/plans/score-predictions.md` → Migration 2.

### `20260524000005_tournament_outcomes.sql`

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TABLE player_outcomes (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID NOT NULL REFERENCES tournaments(id),
    category      player_handicap_category NOT NULL,
    player_id     UUID NOT NULL REFERENCES players(id),
    recorded_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT player_outcomes_tournament_category_player_uq
        UNIQUE (tournament_id, category, player_id)
);

CREATE TABLE team_outcomes (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID NOT NULL REFERENCES tournaments(id),
    category      team_handicap_category NOT NULL,
    team_id       UUID NOT NULL REFERENCES teams(id),
    recorded_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT team_outcomes_tournament_category_team_uq
        UNIQUE (tournament_id, category, team_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS team_outcomes;
DROP TABLE IF EXISTS player_outcomes;

-- +goose StatementEnd
```

---

## Domain changes

### `internal/domain/errors.go`

Add `ErrConflict` sentinel alongside the existing ones:
```go
var ErrConflict = errors.New("conflict")
```
Used when a unique slot is already taken (e.g. playoff slot 2 for a group the
user already has a wildcard pick for). The handler maps it to 409.

### `internal/domain/team.go`

Add `GroupLetter *string` to `Team` struct (nullable until the draw).

### `internal/domain/player.go`

Add `GroupLetter *string` to `Player` struct (derived from the player's team at
query time via JOIN — see repository changes). Nil when the team has no group
letter assigned yet.

### `internal/domain/tournament_prediction.go`

Replace the existing types with the following. The key changes are annotated.

```go
// UpsertPlayerPredictionInput — add GroupLetter
type UpsertPlayerPredictionInput struct {
    UserID       uuid.UUID
    TournamentID uuid.UUID
    Category     PlayerHandicapCategory
    Pick         uuid.UUID
    GroupLetter  *string // required for group_top_scorer; nil for total_top_scorer
}

// PlayerPrediction — add GroupLetter
type PlayerPrediction struct {
    // ... existing fields ...
    GroupLetter *string
}

// PlayerPredictionView — key is now (category, group); one view per slot
type PlayerPredictionView struct {
    Category    PlayerHandicapCategory
    GroupLetter *string            // nil for total_top_scorer
    Prediction  *PlayerPrediction  // nil = not yet predicted
}

// UpsertTeamPredictionInput — add GroupLetter, SlotIndex
type UpsertTeamPredictionInput struct {
    UserID       uuid.UUID
    TournamentID uuid.UUID
    Category     TeamHandicapCategory
    Pick         uuid.UUID
    GroupLetter  *string // required for group_winner, playoff; nil for semifinalist, winner
    SlotIndex    int     // 0 for group_winner/winner/total_top_scorer; 0–3 for semifinalist; 0–2 for playoff
}

// TeamPrediction — add GroupLetter, SlotIndex
type TeamPrediction struct {
    // ... existing fields ...
    GroupLetter *string
    SlotIndex   int
}

// TeamPredictionView — key is now (category, group, slot)
type TeamPredictionView struct {
    Category    TeamHandicapCategory
    GroupLetter *string
    SlotIndex   int
    Prediction  *TeamPrediction
}

// League pick types — add GroupLetter and SlotIndex to carry full slot identity
type PlayerLeaguePick struct {
    UserID      uuid.UUID
    Category    PlayerHandicapCategory
    GroupLetter *string
    PlayerID    uuid.UUID
    PlayerName  string
    Points      *int
}

type TeamLeaguePick struct {
    UserID      uuid.UUID
    Category    TeamHandicapCategory
    GroupLetter *string
    SlotIndex   int
    TeamID      uuid.UUID
    TeamName    string
    Points      *int
}

// League view types — key expands to (category, group, slot)
type LeaguePlayerCategoryView struct {
    Category    PlayerHandicapCategory
    GroupLetter *string
    Predictions []LeagueMemberPlayerPick
}

type LeagueTeamCategoryView struct {
    Category    TeamHandicapCategory
    GroupLetter *string
    SlotIndex   int
    Predictions []LeagueMemberTeamPick
}
```

### `internal/domain/outcome.go` (new file)

```go
type PlayerOutcome struct {
    ID           uuid.UUID
    TournamentID uuid.UUID
    Category     PlayerHandicapCategory
    PlayerID     uuid.UUID
    RecordedAt   time.Time
}

type TeamOutcome struct {
    ID           uuid.UUID
    TournamentID uuid.UUID
    Category     TeamHandicapCategory
    TeamID       uuid.UUID
    RecordedAt   time.Time
}
```

### `internal/domain/fixture.go` and `internal/domain/score_prediction.go`

Implement exactly as specified in `docs/plans/score-predictions.md` (Domain section).

---

## Repository changes

### `internal/repository/queries/tournament_predictions.sql`

Replace existing queries. The upsert now splits into two named queries per
entity (group vs no-group) because partial ON CONFLICT requires the WHERE clause
to match the index definition.

```sql
-- name: UpsertPlayerPredictionGroup :one
WITH ins AS (
    INSERT INTO player_predictions (user_id, tournament_id, category, pick, group_letter)
    VALUES (@user_id, @tournament_id, @category, @pick, @group_letter)
    ON CONFLICT (user_id, tournament_id, category, group_letter)
        WHERE group_letter IS NOT NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, points, created_at, updated_at
)
SELECT ins.*, p.name AS pick_name FROM ins JOIN players p ON p.id = ins.pick;

-- name: UpsertPlayerPredictionNoGroup :one
WITH ins AS (
    INSERT INTO player_predictions (user_id, tournament_id, category, pick)
    VALUES (@user_id, @tournament_id, @category, @pick)
    ON CONFLICT (user_id, tournament_id, category)
        WHERE group_letter IS NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, points, created_at, updated_at
)
SELECT ins.*, p.name AS pick_name FROM ins JOIN players p ON p.id = ins.pick;

-- name: ListPlayerPredictionsByTournamentForUser :many
SELECT pp.id, pp.user_id, pp.tournament_id, pp.category, pp.pick,
       pp.group_letter, p.name AS pick_name, pp.points, pp.created_at, pp.updated_at
FROM player_predictions pp
JOIN players p ON p.id = pp.pick
WHERE pp.tournament_id = @tournament_id AND pp.user_id = @user_id;

-- name: ListPlayerPredictionsByLeague :many
SELECT lm.user_id, pp.category, pp.group_letter, pp.pick AS player_id,
       p.name AS player_name, pp.points
FROM league_members lm
JOIN leagues l ON l.id = lm.league_id
JOIN player_predictions pp
    ON pp.user_id = lm.user_id AND pp.tournament_id = l.tournament_id
JOIN players p ON p.id = pp.pick
WHERE lm.league_id = @league_id;

-- name: UpsertTeamPredictionGroup :one
WITH ins AS (
    INSERT INTO team_predictions (user_id, tournament_id, category, pick, group_letter, slot_index)
    VALUES (@user_id, @tournament_id, @category, @pick, @group_letter, @slot_index)
    ON CONFLICT (user_id, tournament_id, category, group_letter, slot_index)
        WHERE group_letter IS NOT NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, slot_index, points, created_at, updated_at
)
SELECT ins.*, t.name AS pick_name FROM ins JOIN teams t ON t.id = ins.pick;

-- name: UpsertTeamPredictionNoGroup :one
WITH ins AS (
    INSERT INTO team_predictions (user_id, tournament_id, category, pick, slot_index)
    VALUES (@user_id, @tournament_id, @category, @pick, @slot_index)
    ON CONFLICT (user_id, tournament_id, category, slot_index)
        WHERE group_letter IS NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, slot_index, points, created_at, updated_at
)
SELECT ins.*, t.name AS pick_name FROM ins JOIN teams t ON t.id = ins.pick;

-- name: ListTeamPredictionsByTournamentForUser :many
SELECT tp.id, tp.user_id, tp.tournament_id, tp.category, tp.pick,
       tp.group_letter, tp.slot_index, t.name AS pick_name,
       tp.points, tp.created_at, tp.updated_at
FROM team_predictions tp
JOIN teams t ON t.id = tp.pick
WHERE tp.tournament_id = @tournament_id AND tp.user_id = @user_id;

-- name: ListTeamPredictionsByLeague :many
SELECT lm.user_id, tp.category, tp.group_letter, tp.slot_index,
       tp.pick AS team_id, t.name AS team_name, tp.points
FROM league_members lm
JOIN leagues l ON l.id = lm.league_id
JOIN team_predictions tp
    ON tp.user_id = lm.user_id AND tp.tournament_id = l.tournament_id
JOIN teams t ON t.id = tp.pick
WHERE lm.league_id = @league_id;

-- ListLeagueMembersForPredictions — unchanged
```

### `internal/repository/queries/team.sql` — add group query

```sql
-- name: ListGroupLettersByTournament :many
SELECT DISTINCT group_letter
FROM teams
WHERE tournament_id = @tournament_id
  AND group_letter IS NOT NULL
ORDER BY group_letter ASC;
```

### `internal/repository/tournament_prediction.go`

- `UpsertPlayer`: if `in.GroupLetter != nil`, call `UpsertPlayerPredictionGroup`; else `UpsertPlayerPredictionNoGroup`. Map `group_letter` in/out.
- `UpsertTeam`: same split for group vs no-group. Map `group_letter` and `slot_index` in/out.
- `ListPlayersByTournamentForUser`: map `group_letter` from rows.
- `ListTeamsByTournamentForUser`: map `group_letter` and `slot_index` from rows.
- `ListPlayersByLeague`: map `group_letter` from rows.
- `ListTeamsByLeague`: map `group_letter` and `slot_index` from rows.

### `internal/repository/queries/player.sql` — update `GetPlayerByID`

Add a `JOIN teams` to return `group_letter` alongside the existing player fields:
```sql
-- name: GetPlayerByID :one
SELECT p.id, p.external_id, p.name, p.tournament_id, p.team_id,
       t.group_letter,
       p.created_at, p.updated_at
FROM players p
JOIN teams t ON t.id = p.team_id
WHERE p.id = @id;
```
The `PlayerRepository.GetByID` method must map `group_letter` onto `domain.Player.GroupLetter`.

### `internal/repository/queries/tournament_predictions.sql` — add wildcard count

Add this query to the file alongside the other prediction queries:
```sql
-- name: CountPlayoffWildcards :one
SELECT COUNT(*)::int FROM team_predictions
WHERE tournament_id = @tournament_id
  AND user_id       = @user_id
  AND category      = 'playoff'
  AND slot_index    = 2;
```

### `internal/repository/tournament_prediction.go` — handle unique violation

When `UpsertTeamPredictionGroup` or `UpsertTeamPredictionNoGroup` returns a pgx
unique-violation error (pgconn error code `"23505"`), wrap it as `domain.ErrConflict`
so the handler can map it to 409. Pattern:
```go
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) && pgErr.Code == "23505" {
    return nil, fmt.Errorf("slot already taken: %w", domain.ErrConflict)
}
```

### `internal/repository/team.go` — add method

```go
// ListGroupLettersByTournament returns distinct group letters for a tournament's teams,
// sorted A→Z. Returns empty slice when no teams have a group_letter assigned yet.
func (r *TeamRepository) ListGroupLettersByTournament(ctx context.Context, tournamentID uuid.UUID) ([]string, error)
```

### Score prediction and fixture repositories

Implement exactly as specified in `docs/plans/score-predictions.md`.

---

## Service changes

### New interface in `tournament_prediction.go`

```go
// TeamGroupLister returns sorted distinct group letters for a tournament.
type TeamGroupLister interface {
    ListGroupLettersByTournament(ctx context.Context, tournamentID uuid.UUID) ([]string, error)
}
```

Add `teamGroups TeamGroupLister` to `TournamentPredictionService` and its constructor.

### `UpsertPlayerPrediction`

Additional validation before calling the repo:
1. If `category == group_top_scorer`: `in.GroupLetter` must be non-nil, and `player.Team.GroupLetter` must equal `*in.GroupLetter` (player belongs to that group).
2. If `category == total_top_scorer`: `in.GroupLetter` must be nil.

The `Team` embedded in `Player` gains `GroupLetter`; the existing `PlayerGetter.GetByID` already returns a `Player` with `TeamID` — add a team lookup or embed `GroupLetter` on player directly. Simplest: add `GroupLetter *string` to `domain.Player` (populated from `teams.group_letter` via JOIN in the player query).

### `UpsertTeamPrediction`

Additional validation:
1. `group_winner`: `in.GroupLetter` must be non-nil; `in.SlotIndex` must be 0; `team.GroupLetter` must equal `*in.GroupLetter`.
2. `playoff`: `in.GroupLetter` must be non-nil; `in.SlotIndex` must be 0, 1, or 2 (anything else → `ErrInvalid`); `team.GroupLetter` must equal `*in.GroupLetter`; if `in.SlotIndex == 2`, count existing wildcard picks — must be `< 8` (else `ErrForbidden`).
3. `semifinalist`: `in.GroupLetter` must be nil; `in.SlotIndex` ∈ {0,1,2,3} (else `ErrInvalid`).
4. `winner`: `in.GroupLetter` nil; `in.SlotIndex` must be 0.

If the repo returns `ErrConflict` (unique index violation on the slot), the
service returns it as-is; the handler maps `ErrConflict` → 409.

Add a repo method to count wildcard picks:
```go
// CountPlayoffWildcardsByUser returns the number of playoff slot-2 picks for a user in a tournament.
type TeamPredictionRepo interface {
    // ... existing methods ...
    CountPlayoffWildcards(ctx context.Context, tournamentID, userID uuid.UUID) (int, error)
}
```

SQL:
```sql
-- name: CountPlayoffWildcards :one
SELECT COUNT(*) FROM team_predictions
WHERE tournament_id = @tournament_id
  AND user_id = @user_id
  AND category = 'playoff'
  AND slot_index = 2;
```

### `ListPlayerPredictionsForUser`

Replace the category-only enumeration with a group-aware one:

```go
groups, err := s.teamGroups.ListGroupLettersByTournament(ctx, tournamentID)
// ...
// byKey: map[(category, groupLetter)] -> *PlayerPrediction
// enumerate: for group_top_scorer → one view per group letter
//            for total_top_scorer → one view with GroupLetter=nil
```

### `ListTeamPredictionsForUser`

Same pattern, with slot enumeration:

```go
// enumerate:
// group_winner   → 1 view per group (slot 0)
// playoff        → 3 views per group (slots 0, 1, 2)
// semifinalist   → 4 views (slots 0–3, GroupLetter=nil)
// winner         → 1 view (slot 0, GroupLetter=nil)
// byKey: map[(category, groupLetter, slotIndex)] -> *TeamPrediction
```

### `buildLeaguePlayerViews` / `buildLeagueTeamViews`

Map key expands to include `GroupLetter` and `SlotIndex`. The enumeration loop
(currently over `AllPlayerHandicapCategories`) must be replaced with the same
group-aware enumeration used in the personal list methods. Pass `groups []string`
as an argument.

Both `ListLeaguePlayerPredictions` and `ListLeagueTeamPredictions` must fetch
group letters before calling the build helpers.

### Score prediction service

Implement exactly as specified in `docs/plans/score-predictions.md` (Service section).

---

## API contract

### Modified: `PUT /tournaments/{tournamentId}/predictions/players/{category}`

Request body gains `group` (string | null):
```json
{ "player_id": "<uuid>", "group": "A" }
```
- `group` is required when `category = group_top_scorer`; must be null/absent otherwise.
- 400 if `group` is wrong type or present for `total_top_scorer`.
- 404 if player's team group letter doesn't match supplied `group`.

Response gains `group`:
```json
{ "id": "...", "tournament_id": "...", "category": "...", "group": "A",
  "player_id": "...", "player_name": "...", "points": null }
```

### Modified: `GET /tournaments/{tournamentId}/predictions/players`

Response is now a flat array with one item per `(category, group)` slot.
Each item has `group` (string | null) and `player_id` / `player_name` /
`points` as nullable (null = not yet predicted):
```json
[
  { "category": "group_top_scorer", "group": "A", "player_id": "...", "player_name": "Müller", "points": null },
  { "category": "group_top_scorer", "group": "B", "player_id": null,  "player_name": null,    "points": null },
  ...
  { "category": "total_top_scorer", "group": null, "player_id": "...", "player_name": "Mbappe", "points": 10 }
]
```

### Modified: `PUT /tournaments/{tournamentId}/predictions/teams/{category}`

Request body gains `group` and `slot`:
```json
{ "team_id": "<uuid>", "group": "A", "slot": 0 }
```

| category | `group` | `slot` |
|---|---|---|
| `group_winner` | required | omit (default 0) |
| `playoff` | required | required (0, 1, or 2) |
| `semifinalist` | omit | required (0–3) |
| `winner` | omit | omit (default 0) |

- 400 if `slot` out of range for category.
- 403 if `playoff` slot 2 and user already has 8 wildcard picks.
- 404 if team's group doesn't match supplied `group`.

Response gains `group` and `slot`.

### Modified: `GET /tournaments/{tournamentId}/predictions/teams`

Flat array, one item per `(category, group, slot)`:
```json
[
  { "category": "group_winner",  "group": "A", "slot": 0, "team_id": "...", "team_name": "Germany", "points": 3 },
  { "category": "group_winner",  "group": "B", "slot": 0, "team_id": null,  "team_name": null,      "points": null },
  ...
  { "category": "playoff", "group": "A", "slot": 0, "team_id": "...", "team_name": "Germany", "points": null },
  { "category": "playoff", "group": "A", "slot": 1, "team_id": "...", "team_name": "France",  "points": null },
  { "category": "playoff", "group": "A", "slot": 2, "team_id": null,  "team_name": null,      "points": null },
  ...
  { "category": "semifinalist", "group": null, "slot": 0, "team_id": "...", "team_name": "Brazil", "points": null },
  { "category": "semifinalist", "group": null, "slot": 1, "team_id": null,  "team_name": null,     "points": null },
  { "category": "semifinalist", "group": null, "slot": 2, "team_id": null,  "team_name": null,     "points": null },
  { "category": "semifinalist", "group": null, "slot": 3, "team_id": null,  "team_name": null,     "points": null },
  { "category": "winner",       "group": null, "slot": 0, "team_id": null,  "team_name": null,     "points": null }
]
```

### Modified: `GET /leagues/{leagueId}/predictions/players` and `GET /leagues/{leagueId}/predictions/teams`

Both gain `group` and (for teams) `slot` in the category view keys, mirroring
the personal list shapes. Each item in the `predictions` array gains
`group_letter` / `slot_index` for teams.

### New score prediction endpoints

Implement exactly as specified in `docs/plans/score-predictions.md` (API section):
- `PUT /predictions/{fixtureId}`
- `GET /tournaments/{tournamentId}/fixtures`
- `GET /leagues/{leagueId}/predictions` (score predictions league view)

---

## Handler changes

`internal/server/handler/errors.go` must map `domain.ErrConflict` → 409 alongside
the existing mappings (`ErrNotFound` → 404, `ErrForbidden` → 403, `ErrInvalid` → 400).

---

## OpenAPI updates

1. Update `PlayerPredictionViewResponse` schema — add `group: {type: string, nullable: true}`.
2. Update `TeamPredictionViewResponse` schema — add `group: {type: string, nullable: true}` and `slot: {type: integer}`.
3. Update PUT request body schemas for both player and team predictions.
4. Update league prediction response schemas accordingly.
5. Add schemas and paths from `docs/plans/score-predictions.md`.

Run `make generate` after all spec changes.

---

## Postman

Update existing requests:
- `Upsert-Player-Prediction.request.yaml` — add `group` to body
- `Upsert-Team-Prediction.request.yaml` — add `group` and `slot` to body

Add new requests per `docs/plans/score-predictions.md` (three new files).

---

## Edge cases

| Case | Expected |
|---|---|
| `group_top_scorer` submitted without `group` | 400 |
| `group_top_scorer` with group `"Z"` (no teams in that group) | 404 — team lookup fails |
| Player's team has no `group_letter` yet | 404 — group match fails |
| `playoff` slot 2 when user already has 8 wildcards | 403 |
| `playoff` slot 2 for a group where user already has a slot-2 pick | repo gets unique violation (23505) → wrapped as `ErrConflict` → 409 |
| `playoff` slot 3 (out of range) | 400 |
| `semifinalist` with `group` set | 400 |
| Same team submitted to two different `semifinalist` slots | 200 — allowed (service does not deduplicate) |
| No teams have `group_letter` set (pre-draw) | Personal list returns only tournament-scoped slots; group-scoped categories return zero rows |
| League view before lock time | 403 |

---

## Test plan

### Repository tests (testcontainers)

- `UpsertPlayerPredictionGroup`: inserts; conflict on same (user, tournament, category, group) updates pick; `points` untouched.
- `UpsertPlayerPredictionNoGroup`: same for group-less categories.
- `ListPlayerPredictionsByTournamentForUser`: returns `group_letter` in results.
- `UpsertTeamPredictionGroup` / `NoGroup`: same patterns; `slot_index` round-trips correctly.
- `ListTeamPredictionsByTournamentForUser`: returns `group_letter` and `slot_index`.
- `ListGroupLettersByTournament`: returns sorted distinct letters; empty when no teams assigned.
- `CountPlayoffWildcards`: returns 0 initially; increments correctly.

### Service tests (hand-written fakes)

- `UpsertPlayerPrediction`: group required for `group_top_scorer`; nil for `total_top_scorer`; mismatch → 404.
- `UpsertTeamPrediction`: group required for group-scoped; slot out of range → 400; 9th wildcard → 403; team in wrong group → 404.
- `ListPlayerPredictionsForUser`: when 3 groups exist and user has predicted 1 `group_top_scorer`, returns 2 nil + 1 filled + 1 `total_top_scorer` (nil or filled).
- `ListTeamPredictionsForUser`: 3 groups → 3 `group_winner` + 9 `playoff` + 4 `semifinalist` + 1 `winner` = 17 rows.
- `ListLeaguePlayerPredictions` / `ListLeagueTeamPredictions`: view key includes group and slot; requesting user first.

### Handler tests (httptest)

- PUT player with missing `group` for `group_top_scorer` → 400.
- PUT team with `slot: 5` for `semifinalist` → 400.
- GET personal player list: correct count of items when groups exist; correct count when no groups.
- GET personal team list: correct structure including playoff wildcard slots.
- Unauthenticated → 401.

---

## Acceptance criteria

1. `GET /tournaments/{id}/predictions/players` returns exactly `(num_groups × 1) + 1` items when groups are assigned; `1` item when no groups are assigned.
2. `GET /tournaments/{id}/predictions/teams` returns exactly `(num_groups × 4) + 5` items: `group_winner` (×g) + `playoff` (×3g) + `semifinalist` (×4) + `winner` (×1).
3. Slots the user hasn't filled have `team_id: null` / `player_id: null` / `points: null`.
4. `PUT` with correct group and slot stores the pick and returns it.
5. Submitting `playoff` slot 2 nine times for different groups returns 403 on the ninth.
6. Submitting `playoff` slot 2 twice for the same group returns conflict (409 or 400 per handler mapping).
7. Picking a team whose `group_letter` doesn't match the supplied `group` returns 404.
8. League endpoints reflect the same group/slot structure; requesting user appears first.
9. Score prediction endpoints work per acceptance criteria in `docs/plans/score-predictions.md`.
10. `make fmt && make lint && make test` all pass.
11. OpenAPI spec regenerated; Postman requests updated.
