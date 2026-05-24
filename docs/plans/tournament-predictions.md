# Plan: Tournament Predictions

## Goal

Allow authenticated users to submit tournament-level predictions — picking a
player for each player category (e.g. top scorer) and a team for each team
category (e.g. winner). All predictions lock 30 minutes before the first
fixture of the tournament. After lock, league members can see each other's
picks.

---

## Dependencies

This feature depends on the `fixtures` table introduced in the
**score-predictions** plan (`20260524000001_fixtures.sql`). That feature
must be implemented first so that `fixtures.kickoff_at` can be queried for
the locking check.

---

## Design decisions

### Same categories as handicaps — no new enums

`player_predictions.category` reuses `player_handicap_category`; 
`team_predictions.category` reuses `team_handicap_category`. No new DB types
needed.

### `tournament_id` is explicit on both tables

Although the picked player/team already carries a `tournament_id` FK, keeping
it explicit on the predictions tables:

- Makes the unique constraint clean: one pick per `(user, tournament, category)`.
- Enables simple listing queries without a JOIN through the pick's own table.
- Allows the service to validate that the pick belongs to the right tournament.

### `points` nullable — same reasoning as score predictions

`NULL` = not yet scored; `0` = scored, earned nothing. The scoring worker owns
that column; the upsert excludes it from the `DO UPDATE` clause.

### Single lock time per tournament

All tournament predictions lock together: 30 minutes before `MIN(kickoff_at)`
across all fixtures for the tournament. If no fixtures exist yet, predictions
remain open (the lock time is undefined, so we never block submission).

### "Get my predictions" returns all categories

The "list my predictions" endpoints enumerate every known category and return
`null` for any category the user hasn't predicted yet. This mirrors the score
predictions fixture listing (which shows all fixtures, not just predicted ones)
and saves the client from having to infer missing categories.

### League view requires lock to have passed

`GET /leagues/{leagueId}/predictions/players` and `.../teams` return `403`
until the lock time has passed. After that they return all member picks. This
prevents peeking at others' predictions before the window closes.

### Pick validated against tournament

The service verifies that the picked player/team belongs to the tournament in
the request path. A player that exists but belongs to a different tournament
returns `404` (conceptually: "not found in this tournament").

### Upsert semantics

`INSERT … ON CONFLICT … DO UPDATE` — same pick submitted twice updates the
row; points column excluded. Consistent with score predictions.

---

## Data model

### Migration — `20260524000003_tournament_predictions.sql`

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TABLE player_predictions (
    id            UUID                     PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID                     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tournament_id UUID                     NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    category      player_handicap_category NOT NULL,
    pick          UUID                     NOT NULL REFERENCES players(id),
    points        INT,
    created_at    TIMESTAMPTZ              NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ              NOT NULL DEFAULT now(),
    CONSTRAINT player_predictions_user_tournament_category_uq UNIQUE (user_id, tournament_id, category)
);

CREATE INDEX player_predictions_tournament_id_idx ON player_predictions (tournament_id);

CREATE TRIGGER set_updated_at_player_predictions
    BEFORE UPDATE ON player_predictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE team_predictions (
    id            UUID                   PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID                   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tournament_id UUID                   NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    category      team_handicap_category NOT NULL,
    pick          UUID                   NOT NULL REFERENCES teams(id),
    points        INT,
    created_at    TIMESTAMPTZ            NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ            NOT NULL DEFAULT now(),
    CONSTRAINT team_predictions_user_tournament_category_uq UNIQUE (user_id, tournament_id, category)
);

CREATE INDEX team_predictions_tournament_id_idx ON team_predictions (tournament_id);

CREATE TRIGGER set_updated_at_team_predictions
    BEFORE UPDATE ON team_predictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_team_predictions   ON team_predictions;
DROP TRIGGER IF EXISTS set_updated_at_player_predictions ON player_predictions;
DROP TABLE IF EXISTS team_predictions;
DROP TABLE IF EXISTS player_predictions;

-- +goose StatementEnd
```

---

## API Contract

### 1 — Upsert player prediction

```
PUT /tournaments/{tournamentId}/predictions/players/{category}
Authorization: Bearer {{authToken}}
```

`category` is one of: `group_top_scorer`, `total_top_scorer`.

**Request body:**
```json
{ "player_id": "<uuid>" }
```

**Response 200:**
```json
{
  "id": "<uuid>",
  "tournament_id": "<uuid>",
  "category": "group_top_scorer",
  "player_id": "<uuid>",
  "player_name": "Lionel Messi",
  "points": null
}
```

**Errors:**

| Condition | Status |
|-----------|--------|
| Not authenticated | 401 |
| `tournamentId` not a valid UUID | 400 |
| `category` not a recognised value | 400 |
| `player_id` missing or not a valid UUID | 400 |
| Player not found in this tournament | 404 |
| Predictions are locked (≤ 30 min before first fixture) | 403 |

---

### 2 — Upsert team prediction

```
PUT /tournaments/{tournamentId}/predictions/teams/{category}
Authorization: Bearer {{authToken}}
```

`category` is one of: `group_winner`, `playoff`, `semifinalist`, `winner`.

**Request body:**
```json
{ "team_id": "<uuid>" }
```

**Response 200:**
```json
{
  "id": "<uuid>",
  "tournament_id": "<uuid>",
  "category": "winner",
  "team_id": "<uuid>",
  "team_name": "Argentina",
  "points": null
}
```

**Errors:** same pattern as player upsert, substituting team for player.

---

### 3 — List my player predictions

```
GET /tournaments/{tournamentId}/predictions/players
Authorization: Bearer {{authToken}}
```

Returns one entry per category. Categories the user has not yet predicted have
`player_id: null` and `player_name: null`.

**Response 200:**
```json
[
  {
    "category": "group_top_scorer",
    "player_id": "<uuid>",
    "player_name": "Lionel Messi",
    "points": null
  },
  {
    "category": "total_top_scorer",
    "player_id": null,
    "player_name": null,
    "points": null
  }
]
```

**Errors:**

| Condition | Status |
|-----------|--------|
| Not authenticated | 401 |
| `tournamentId` not a valid UUID | 400 |

---

### 4 — List my team predictions

```
GET /tournaments/{tournamentId}/predictions/teams
Authorization: Bearer {{authToken}}
```

Returns one entry per category. Categories without a pick have `team_id: null`
and `team_name: null`.

**Response 200:**
```json
[
  {
    "category": "group_winner",
    "team_id": "<uuid>",
    "team_name": "Argentina",
    "points": null
  },
  {
    "category": "playoff",
    "team_id": null,
    "team_name": null,
    "points": null
  },
  {
    "category": "semifinalist",
    "team_id": null,
    "team_name": null,
    "points": null
  },
  {
    "category": "winner",
    "team_id": "<uuid>",
    "team_name": "France",
    "points": null
  }
]
```

**Errors:** same as list my player predictions.

---

### 5 — List league player predictions

```
GET /leagues/{leagueId}/predictions/players
Authorization: Bearer {{authToken}}
```

Only accessible after the lock time has passed. Returns one entry per
category, each with all league members' picks. Members who have not picked
appear with `player_id: null`. The requesting user's pick appears first within
each category; remaining members sorted alphabetically by `display_name`.

**Response 200:**
```json
[
  {
    "category": "group_top_scorer",
    "predictions": [
      { "user_id": "<uuid>", "display_name": "me",    "player_id": "<uuid>", "player_name": "Messi",   "points": 5 },
      { "user_id": "<uuid>", "display_name": "alice", "player_id": "<uuid>", "player_name": "Ronaldo", "points": null },
      { "user_id": "<uuid>", "display_name": "bob",   "player_id": null,     "player_name": null,      "points": null }
    ]
  },
  {
    "category": "total_top_scorer",
    "predictions": [...]
  }
]
```

**Errors:**

| Condition | Status |
|-----------|--------|
| Not authenticated | 401 |
| `leagueId` not a valid UUID | 400 |
| User not a member of the league | 403 |
| Lock time has not yet passed | 403 |

---

### 6 — List league team predictions

```
GET /leagues/{leagueId}/predictions/teams
Authorization: Bearer {{authToken}}
```

Same shape as league player predictions, substituting `team_id`/`team_name`
for `player_id`/`player_name`.

**Errors:** same as league player predictions.

---

## Layer-by-layer breakdown

### domain

**`internal/domain/tournament_prediction.go`**

```go
type PlayerPrediction struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    TournamentID uuid.UUID
    Category     PlayerHandicapCategory
    Pick         uuid.UUID  // player_id
    PickName     string     // resolved at query time
    Points       *int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type UpsertPlayerPredictionInput struct {
    UserID       uuid.UUID
    TournamentID uuid.UUID
    Category     PlayerHandicapCategory
    Pick         uuid.UUID // player_id
}

// PlayerPredictionView is one category row in the personal listing.
type PlayerPredictionView struct {
    Category   PlayerHandicapCategory
    Prediction *PlayerPrediction // nil = not yet predicted
}

// LeagueMemberPlayerPick is one member's pick within a league category view.
type LeagueMemberPlayerPick struct {
    UserID      uuid.UUID
    DisplayName string
    PlayerID    *uuid.UUID
    PlayerName  *string
    Points      *int
}

// LeaguePlayerCategoryView is one category row in the league listing.
type LeaguePlayerCategoryView struct {
    Category    PlayerHandicapCategory
    Predictions []LeagueMemberPlayerPick
}

// --- team variants mirror the above ---

type TeamPrediction struct {
    ID           uuid.UUID
    UserID       uuid.UUID
    TournamentID uuid.UUID
    Category     TeamHandicapCategory
    Pick         uuid.UUID
    PickName     string
    Points       *int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type UpsertTeamPredictionInput struct {
    UserID       uuid.UUID
    TournamentID uuid.UUID
    Category     TeamHandicapCategory
    Pick         uuid.UUID
}

type TeamPredictionView struct {
    Category   TeamHandicapCategory
    Prediction *TeamPrediction
}

type LeagueMemberTeamPick struct {
    UserID      uuid.UUID
    DisplayName string
    TeamID      *uuid.UUID
    TeamName    *string
    Points      *int
}

type LeagueTeamCategoryView struct {
    Category    TeamHandicapCategory
    Predictions []LeagueMemberTeamPick
}
```

---

### repository

**New file: `internal/repository/tournament_prediction.go`**

```go
// Interfaces consumed by services.

type PlayerPredictionRepo interface {
    UpsertPlayer(ctx context.Context, in domain.UpsertPlayerPredictionInput) (*domain.PlayerPrediction, error)
    ListPlayersByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.PlayerPrediction, error)
    ListPlayersByLeague(ctx context.Context, leagueID uuid.UUID) ([]*domain.PlayerPrediction, error)
    // ListPlayersByLeague returns raw rows (user_id, category, pick, pick_name, points).
    // Merging into per-category views happens in the service.
}

type TeamPredictionRepo interface {
    UpsertTeam(ctx context.Context, in domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error)
    ListTeamsByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.TeamPrediction, error)
    ListTeamsByLeague(ctx context.Context, leagueID uuid.UUID) ([]*domain.TeamPrediction, error)
}
```

**New file: `internal/repository/queries/tournament_predictions.sql`**

```sql
-- name: UpsertPlayerPrediction :one
INSERT INTO player_predictions (user_id, tournament_id, category, pick)
VALUES (@user_id, @tournament_id, @category, @pick)
ON CONFLICT (user_id, tournament_id, category) DO UPDATE
    SET pick       = EXCLUDED.pick,
        updated_at = now()
RETURNING
    id, user_id, tournament_id, category, pick,
    (SELECT name FROM players WHERE id = EXCLUDED.pick) AS pick_name,
    points, created_at, updated_at;

-- name: ListPlayerPredictionsByTournamentForUser :many
SELECT pp.id, pp.user_id, pp.tournament_id, pp.category, pp.pick,
       p.name AS pick_name, pp.points, pp.created_at, pp.updated_at
FROM player_predictions pp
JOIN players p ON p.id = pp.pick
WHERE pp.tournament_id = @tournament_id AND pp.user_id = @user_id;

-- name: ListPlayerPredictionsByLeague :many
SELECT pp.user_id, u.display_name, pp.tournament_id, pp.category,
       pp.pick AS player_id, p.name AS player_name, pp.points
FROM player_predictions pp
JOIN players p   ON p.id  = pp.pick
JOIN users   u   ON u.id  = pp.user_id
JOIN leagues l   ON l.tournament_id = pp.tournament_id
JOIN league_members lm ON lm.league_id = l.id AND lm.user_id = pp.user_id
WHERE l.id = @league_id;

-- name: UpsertTeamPrediction :one
INSERT INTO team_predictions (user_id, tournament_id, category, pick)
VALUES (@user_id, @tournament_id, @category, @pick)
ON CONFLICT (user_id, tournament_id, category) DO UPDATE
    SET pick       = EXCLUDED.pick,
        updated_at = now()
RETURNING
    id, user_id, tournament_id, category, pick,
    (SELECT name FROM teams WHERE id = EXCLUDED.pick) AS pick_name,
    points, created_at, updated_at;

-- name: ListTeamPredictionsByTournamentForUser :many
SELECT tp.id, tp.user_id, tp.tournament_id, tp.category, tp.pick,
       t.name AS pick_name, tp.points, tp.created_at, tp.updated_at
FROM team_predictions tp
JOIN teams t ON t.id = tp.pick
WHERE tp.tournament_id = @tournament_id AND tp.user_id = @user_id;

-- name: ListTeamPredictionsByLeague :many
SELECT tp.user_id, u.display_name, tp.tournament_id, tp.category,
       tp.pick AS team_id, t.name AS team_name, tp.points
FROM team_predictions tp
JOIN teams   t   ON t.id  = tp.pick
JOIN users   u   ON u.id  = tp.user_id
JOIN leagues l   ON l.tournament_id = tp.tournament_id
JOIN league_members lm ON lm.league_id = l.id AND lm.user_id = tp.user_id
WHERE l.id = @league_id;
```

**Existing `internal/repository/fixture.go`** — add:

```go
type FixtureFirstKickoffGetter interface {
    GetFirstKickoffByTournament(ctx context.Context, tournamentID uuid.UUID) (time.Time, error)
    // Returns domain.ErrNotFound if no fixtures exist for the tournament.
}
```

SQL addition to `internal/repository/queries/fixtures.sql`:

```sql
-- name: GetFirstKickoffByTournament :one
SELECT MIN(kickoff_at) AS kickoff_at
FROM fixtures
WHERE tournament_id = @tournament_id
HAVING COUNT(*) > 0;
```

**Player/Team existence check** — the service fetches the player/team by ID
(using existing `PlayerRepository.GetByID` / a new `TeamRepository.GetByID`)
and verifies `entity.TournamentID == tournamentID`. Return `domain.ErrNotFound`
on mismatch — the pick simply doesn't exist in this tournament from the caller's
perspective.

---

### service

**New file: `internal/service/tournament_prediction.go`**

```go
// Clock is already defined in service/prediction.go; reuse the interface, not redefine.

type TournamentPredictionService struct {
    playerPredictions PlayerPredictionRepo
    teamPredictions   TeamPredictionRepo
    players           PlayerGetter      // GetByID to validate pick belongs to tournament
    teams             TeamGetter        // GetByID to validate pick belongs to tournament
    fixtures          FixtureFirstKickoffGetter
    leagues           LeagueMemberChecker
    clock             Clock
}

func (s *TournamentPredictionService) lockAt(ctx context.Context, tournamentID uuid.UUID) (time.Time, bool, error)
// Returns (lockTime, isLocked, err).
// isLocked is false (and no error) when no fixtures exist yet.

func (s *TournamentPredictionService) UpsertPlayerPrediction(
    ctx context.Context, in domain.UpsertPlayerPredictionInput,
) (*domain.PlayerPrediction, error)
// Validates: pick player exists and belongs to tournament; predictions not locked.

func (s *TournamentPredictionService) UpsertTeamPrediction(
    ctx context.Context, in domain.UpsertTeamPredictionInput,
) (*domain.TeamPrediction, error)

func (s *TournamentPredictionService) ListPlayerPredictionsForUser(
    ctx context.Context, tournamentID, userID uuid.UUID,
) ([]*domain.PlayerPredictionView, error)
// Fetches existing predictions and merges with AllPlayerHandicapCategories.

func (s *TournamentPredictionService) ListTeamPredictionsForUser(
    ctx context.Context, tournamentID, userID uuid.UUID,
) ([]*domain.TeamPredictionView, error)

func (s *TournamentPredictionService) ListLeaguePlayerPredictions(
    ctx context.Context, leagueID, requestingUserID uuid.UUID,
) ([]*domain.LeaguePlayerCategoryView, error)
// Checks membership, then lock, then fetches and merges.

func (s *TournamentPredictionService) ListLeagueTeamPredictions(
    ctx context.Context, leagueID, requestingUserID uuid.UUID,
) ([]*domain.LeagueTeamCategoryView, error)
```

**Locking logic in service:**

```go
func (s *TournamentPredictionService) lockAt(ctx context.Context, tournamentID uuid.UUID) (time.Time, bool, error) {
    firstKickoff, err := s.fixtures.GetFirstKickoffByTournament(ctx, tournamentID)
    if errors.Is(err, domain.ErrNotFound) {
        return time.Time{}, false, nil // no fixtures yet — open
    }
    if err != nil {
        return time.Time{}, false, fmt.Errorf("get first kickoff: %w", err)
    }
    lockAt := firstKickoff.Add(-30 * time.Minute)
    return lockAt, !s.clock.Now().Before(lockAt), nil
}
```

**Merging for personal view** (service-layer, not SQL):

```go
func (s *TournamentPredictionService) ListPlayerPredictionsForUser(...) ([]*domain.PlayerPredictionView, error) {
    rows, err := s.playerPredictions.ListPlayersByTournamentForUser(ctx, tournamentID, userID)
    // ...
    byCategory := map[domain.PlayerHandicapCategory]*domain.PlayerPrediction{}
    for _, r := range rows { byCategory[r.Category] = r }
    views := make([]*domain.PlayerPredictionView, 0, len(domain.AllPlayerHandicapCategories))
    for _, cat := range domain.AllPlayerHandicapCategories {
        views = append(views, &domain.PlayerPredictionView{Category: cat, Prediction: byCategory[cat]})
    }
    return views, nil
}
```

**Merging for league view** — fetch raw rows from repo, then build:

```go
// For each member in the league, create an entry per category.
// Members with no prediction for a category get nil PlayerID/TeamID.
// Requesting user's pick sorted first within each category.
```

The league member list comes from a separate `ListLeagueMembers` query (already
used in the league service) or fetched via the existing league repo method.
Two-query merge: fetch all member predictions for the league, then fetch league
members, then cross-reference in Go.

---

### handler

**New file: `internal/server/handler/tournament_prediction.go`**

```go
type TournamentPredictionSvc interface {
    UpsertPlayerPrediction(ctx context.Context, in domain.UpsertPlayerPredictionInput) (*domain.PlayerPrediction, error)
    UpsertTeamPrediction(ctx context.Context, in domain.UpsertTeamPredictionInput) (*domain.TeamPrediction, error)
    ListPlayerPredictionsForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.PlayerPredictionView, error)
    ListTeamPredictionsForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.TeamPredictionView, error)
    ListLeaguePlayerPredictions(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeaguePlayerCategoryView, error)
    ListLeagueTeamPredictions(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeagueTeamCategoryView, error)
}

type TournamentPrediction struct {
    logger *slog.Logger
    svc    TournamentPredictionSvc
}

// Routes:
// PUT  /tournaments/{tournamentId}/predictions/players/{category}  → UpsertPlayer
// PUT  /tournaments/{tournamentId}/predictions/teams/{category}    → UpsertTeam
// GET  /tournaments/{tournamentId}/predictions/players             → ListMyPlayers
// GET  /tournaments/{tournamentId}/predictions/teams               → ListMyTeams
// GET  /leagues/{leagueId}/predictions/players                     → ListLeaguePlayers
// GET  /leagues/{leagueId}/predictions/teams                       → ListLeagueTeams
```

---

## OpenAPI additions

New paths (add to `docs/openapi.yaml`, then run `make generate`):

```
PUT  /tournaments/{tournamentId}/predictions/players/{category}
GET  /tournaments/{tournamentId}/predictions/players
PUT  /tournaments/{tournamentId}/predictions/teams/{category}
GET  /tournaments/{tournamentId}/predictions/teams
GET  /leagues/{leagueId}/predictions/players
GET  /leagues/{leagueId}/predictions/teams
```

New schemas:
- `UpsertPlayerPredictionRequest` — `{ player_id: string (uuid) }`
- `UpsertTeamPredictionRequest`   — `{ team_id: string (uuid) }`
- `PlayerPredictionResponse`      — upsert / item shape with `player_id`, `player_name`, `points`
- `TeamPredictionResponse`        — upsert / item shape with `team_id`, `team_name`, `points`
- `PlayerPredictionView`          — personal listing item (category + nullable player fields)
- `TeamPredictionView`            — personal listing item (category + nullable team fields)
- `LeagueMemberPlayerPick`        — `user_id`, `display_name`, nullable `player_id`, `player_name`, `points`
- `LeagueMemberTeamPick`          — same with team fields
- `LeaguePlayerCategoryView`      — `category` + `predictions: [LeagueMemberPlayerPick]`
- `LeagueTeamCategoryView`        — `category` + `predictions: [LeagueMemberTeamPick]`

All error responses `$ref: '#/components/responses/ErrorResponse'`.
All routes require `security: [{bearerAuth: []}]`.
Tag: `predictions`.

---

## Postman

Six new request files:

- `postman/collections/footy-forecast/predictions/players/UpsertPlayerPrediction.request.yaml`
  - `PUT {{baseUrl}}/tournaments/{{tournamentId}}/predictions/players/group_top_scorer`
  - Body: `{"player_id": "{{playerId}}"}`
  - Header: `Authorization: Bearer {{authToken}}`

- `postman/collections/footy-forecast/predictions/players/ListMyPlayerPredictions.request.yaml`
  - `GET {{baseUrl}}/tournaments/{{tournamentId}}/predictions/players`
  - Header: `Authorization: Bearer {{authToken}}`

- `postman/collections/footy-forecast/predictions/players/ListLeaguePlayerPredictions.request.yaml`
  - `GET {{baseUrl}}/leagues/{{leagueId}}/predictions/players`
  - Header: `Authorization: Bearer {{authToken}}`

- `postman/collections/footy-forecast/predictions/teams/UpsertTeamPrediction.request.yaml`
  - `PUT {{baseUrl}}/tournaments/{{tournamentId}}/predictions/teams/winner`
  - Body: `{"team_id": "{{teamId}}"}`
  - Header: `Authorization: Bearer {{authToken}}`

- `postman/collections/footy-forecast/predictions/teams/ListMyTeamPredictions.request.yaml`
  - `GET {{baseUrl}}/tournaments/{{tournamentId}}/predictions/teams`
  - Header: `Authorization: Bearer {{authToken}}`

- `postman/collections/footy-forecast/predictions/teams/ListLeagueTeamPredictions.request.yaml`
  - `GET {{baseUrl}}/leagues/{{leagueId}}/predictions/teams`
  - Header: `Authorization: Bearer {{authToken}}`

---

## Edge cases

| Case | Expected behaviour |
|------|--------------------|
| No fixtures exist for tournament | Predictions remain open; `lockAt` returns "not locked" |
| Pick player from wrong tournament | 404 — service validates `player.TournamentID == tournamentID` |
| Pick team from wrong tournament | 404 — same check |
| Submit after lock time | 403 — service checks clock |
| Same prediction submitted twice before lock | 200 — upsert returns updated row; points untouched |
| `category` param is unknown string | 400 — handler parse error via `domain.ParsePlayerHandicapCategory` |
| `tournamentId` / `leagueId` not UUID | 400 — `parseUUIDPathValue` |
| User not league member (league view) | 403 — service membership check |
| League view before lock | 403 — service lock check (after membership check) |
| League member who never predicted | Appears in league view with null `player_id` / `team_id` |
| League has only one member | Returns that member's picks across all categories |
| `points` already set; pick updated before lock | Points column untouched — excluded from upsert `DO UPDATE` |

---

## Test plan

### Repository tests (testcontainers)

**PlayerPredictionRepository**
- `UpsertPlayer`: creates a new row; second call with same `(user, tournament, category)` updates `pick`, leaves `points` unchanged.
- `ListPlayersByTournamentForUser`: returns only this user's rows; empty when none exist.
- `ListPlayersByLeague`: returns rows for all league members who have predicted; no row for members who haven't predicted yet (merge happens in service).

**TeamPredictionRepository** — mirror of player tests.

**FixtureRepository**
- `GetFirstKickoffByTournament`: returns earliest `kickoff_at`; returns `domain.ErrNotFound` when no fixtures exist.

### Service tests (fake repos)

**TournamentPredictionService**
- `UpsertPlayerPrediction`: pick belongs to wrong tournament → 404; lock passed → 403; valid → delegates to repo.
- `UpsertTeamPrediction`: same structure.
- `ListPlayerPredictionsForUser`: all categories returned; categories with no prediction have `Prediction: nil`.
- `ListTeamPredictionsForUser`: same.
- `ListLeaguePlayerPredictions`: non-member → 403; before lock → 403; after lock, member without prediction appears with nil player fields.
- `ListLeagueTeamPredictions`: same.

### Handler tests (httptest)

**TournamentPrediction handler**
- `UpsertPlayer` happy path: 200 with correct JSON shape.
- `UpsertPlayer` invalid category: 400.
- `UpsertPlayer` invalid UUID path param: 400.
- `UpsertPlayer` unauthenticated: 401.
- `UpsertPlayer` locked: 403.
- `UpsertPlayer` player not found: 404.
- `ListMyPlayers` happy path: 200 with all categories including nulls.
- `ListMyPlayers` unauthenticated: 401.
- `ListLeaguePlayers` happy path (after lock): 200.
- `ListLeaguePlayers` non-member: 403.
- `ListLeaguePlayers` before lock: 403.
- Team variants: mirror of each player test above.

---

## Acceptance criteria

1. `PUT /tournaments/{tournamentId}/predictions/players/{category}` before lock returns 200; re-submitting before lock updates the pick only.
2. Same endpoint at or after lock returns 403.
3. Pick a player belonging to a different tournament returns 404.
4. `GET /tournaments/{tournamentId}/predictions/players` always returns all categories; unpredicted categories have `player_id: null`.
5. `GET /leagues/{leagueId}/predictions/players` before lock returns 403; after lock returns all member picks grouped by category with requesting user first.
6. League members who have not predicted appear in the league view with null pick fields.
7. All team prediction endpoints behave identically to player equivalents (substituting team fields).
8. `make fmt && make lint && make test` all pass.
9. All six Postman request files committed alongside the code.
10. `make generate` run after OpenAPI spec update; `models.gen.go` and `openapi.json` committed.
