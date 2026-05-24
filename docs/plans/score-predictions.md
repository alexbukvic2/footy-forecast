# Plan: Score Predictions & Fixtures

## Goal

Introduce a local `fixtures` table (synced from the football API) and a
`score_predictions` table. Allow authenticated users to submit score predictions
for individual fixtures. Expose fixture listings — one personalised view and one
league-scoped view — so clients can display fixtures alongside each user's
prediction and points.

---

## Design decisions

### `points` on `score_predictions` — agreed, nullable

`points INT` (nullable). `NULL` = not yet scored; `0` = scored, earned nothing.
Never `NOT NULL DEFAULT 0` — that conflates the two states and makes it
impossible to tell whether the scoring worker has run on a row.
The `ON CONFLICT DO UPDATE` clause for the upsert deliberately excludes
`points`; the scoring worker owns that column entirely.

### Local `fixtures` table — required

Fixtures are synced from the football API by the import worker. Storing them
locally enables:
- Proper FK from `score_predictions` to `fixtures`
- Server-side kickoff locking (checking `kickoff_at` before accepting a prediction)
- Fixture listings without calling the external API at request time

`external_id BIGINT UNIQUE` carries the football-API identifier; `id UUID` is
our primary key, consistent with every other table.

### Kickoff locking — 30-minute window, enforced in this phase

Now that `fixtures.kickoff_at` is local, the service checks it on every
upsert. Predictions lock **30 minutes before kickoff**:

```go
lockAt := fixture.KickoffAt.Add(-30 * time.Minute)
if !clock.Now().Before(lockAt) { return ErrForbidden }
```

This gives users a clean cut-off that is visible ahead of time, rather than
a race window at the exact kickoff second.

### Fixture responses include team names, not team IDs

`home_team_id` / `away_team_id` are stored internally but the API returns
`home_team_name` / `away_team_name` (resolved via JOIN to `teams` at query
time). Clients have no use for the raw UUIDs; they need display text.

### League predictions endpoint renamed

`GET /leagues/{leagueId}/fixtures` became `GET /leagues/{leagueId}/predictions`.
The word "predictions" better describes the response (locked fixtures with all
member predictions), and avoids confusion with the personal fixtures endpoint
which is mounted under `/tournaments`.

### `PUT /predictions/{fixtureId}` uses our UUID

Path parameter is `fixtures.id` (UUID), not `external_id`. Consistent with the
rest of the API; decouples clients from the external data provider.

### Upsert semantics for submit

`INSERT … ON CONFLICT (user_id, fixture_id) DO UPDATE`. The caller does not
distinguish create from update; the latest valid prediction before kickoff wins.

---

## Data model

### Migration 1 — `fixtures` (`20260524000001_fixtures.sql`)

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TYPE fixture_status AS ENUM ('upcoming', 'in_progress', 'finished');

CREATE TABLE fixtures (
    id            UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id   BIGINT         NOT NULL UNIQUE,
    tournament_id UUID           NOT NULL REFERENCES tournaments(id),
    home_team_id  UUID           NOT NULL REFERENCES teams(id),
    away_team_id  UUID           NOT NULL REFERENCES teams(id),
    kickoff_at    TIMESTAMPTZ    NOT NULL,
    status        fixture_status NOT NULL DEFAULT 'upcoming',
    goals_home    INT,
    goals_away    INT,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX fixtures_tournament_id_idx ON fixtures (tournament_id);
CREATE INDEX fixtures_kickoff_at_idx    ON fixtures (kickoff_at);

CREATE TRIGGER set_updated_at_fixtures
    BEFORE UPDATE ON fixtures
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_fixtures ON fixtures;
DROP TABLE  IF EXISTS fixtures;
DROP TYPE   IF EXISTS fixture_status;

-- +goose StatementEnd
```

### Migration 2 — `score_predictions` (`20260524000002_score_predictions.sql`)

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TABLE score_predictions (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fixture_id  UUID        NOT NULL REFERENCES fixtures(id) ON DELETE CASCADE,
    goals_home  INT         NOT NULL CHECK (goals_home >= 0),
    goals_away  INT         NOT NULL CHECK (goals_away >= 0),
    points      INT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT score_predictions_user_fixture_uq UNIQUE (user_id, fixture_id)
);

CREATE TRIGGER set_updated_at_score_predictions
    BEFORE UPDATE ON score_predictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_score_predictions ON score_predictions;
DROP TABLE   IF EXISTS score_predictions;

-- +goose StatementEnd
```

---

## API endpoints

### 1 — Submit (upsert) prediction

```
PUT /predictions/{fixtureId}
Authorization: Bearer {{authToken}}
```

**Request body:**
```json
{ "goals_home": 2, "goals_away": 1 }
```

**Response 200:**
```json
{
  "id": "<uuid>",
  "fixture_id": "<uuid>",
  "goals_home": 2,
  "goals_away": 1,
  "points": null
}
```

**Errors:**

| Condition | Status |
|-----------|--------|
| Not authenticated | 401 |
| `fixtureId` not a valid UUID | 400 |
| `goals_home` or `goals_away` < 0 | 400 |
| Fixture not found | 404 |
| Within 30 minutes of kickoff (or past it) | 403 |

---

### 2 — List fixtures (personal view)

```
GET /tournaments/{tournamentId}/fixtures
Authorization: Bearer {{authToken}}
```

Returns all fixtures in the tournament. Each fixture includes the authenticated
user's own prediction and points (`null` if no prediction has been made yet).

**Response 200:**
```json
[
  {
    "id": "<uuid>",
    "external_id": 12345,
    "tournament_id": "<uuid>",
    "home_team_name": "Argentina",
    "away_team_name": "France",
    "kickoff_at": "2026-07-01T15:00:00Z",
    "status": "finished",
    "goals_home": 2,
    "goals_away": 1,
    "prediction": {
      "id": "<uuid>",
      "goals_home": 2,
      "goals_away": 0,
      "points": 3
    }
  },
  {
    "id": "<uuid>",
    "external_id": 12346,
    "tournament_id": "<uuid>",
    "home_team_name": "Brazil",
    "away_team_name": "Germany",
    "kickoff_at": "2026-07-03T18:00:00Z",
    "status": "upcoming",
    "goals_home": null,
    "goals_away": null,
    "prediction": null
  }
]
```

**Errors:**

| Condition | Status |
|-----------|--------|
| Not authenticated | 401 |
| `tournamentId` not a valid UUID | 400 |

---

### 3 — List predictions (league view)

```
GET /leagues/{leagueId}/predictions
Authorization: Bearer {{authToken}}
```

Returns only **locked** fixtures (`in_progress` or `finished`) for the
tournament the league belongs to, ordered by `kickoff_at DESC`. Each fixture
includes predictions from all league members. Requesting user must be a member
of the league. Within each fixture, the requesting user's prediction appears
first; remaining members are sorted alphabetically by `display_name`.

**Response 200:**
```json
[
  {
    "id": "<uuid>",
    "external_id": 12345,
    "tournament_id": "<uuid>",
    "home_team_name": "Argentina",
    "away_team_name": "France",
    "kickoff_at": "2026-07-01T15:00:00Z",
    "status": "finished",
    "goals_home": 2,
    "goals_away": 1,
    "predictions": [
      { "user_id": "<uuid>", "display_name": "me", "goals_home": 2, "goals_away": 0, "points": 3 },
      { "user_id": "<uuid>", "display_name": "alice", "goals_home": 1, "goals_away": 1, "points": 0 }
    ]
  }
]
```

**Errors:**

| Condition | Status |
|-----------|--------|
| Not authenticated | 401 |
| `leagueId` not a valid UUID | 400 |
| User not a member of the league | 403 |

---

## Layer-by-layer breakdown

### domain

**`internal/domain/fixture.go`**
```go
type FixtureStatus string

const (
    FixtureStatusUpcoming   FixtureStatus = "upcoming"
    FixtureStatusInProgress FixtureStatus = "in_progress"
    FixtureStatusFinished   FixtureStatus = "finished"
)

type Fixture struct {
    ID           uuid.UUID
    ExternalID   int64
    TournamentID uuid.UUID
    HomeTeamID   uuid.UUID
    AwayTeamID   uuid.UUID
    HomeTeamName string    // resolved at query time; empty for GetByID (internal use)
    AwayTeamName string
    KickoffAt    time.Time
    Status       FixtureStatus
    GoalsHome    *int
    GoalsAway    *int
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// UserFixtureView is the result for the personal fixtures listing.
type UserFixtureView struct {
    Fixture    Fixture
    Prediction *ScorePrediction // nil = no prediction made
}

// LeagueMemberPrediction is one member's prediction within the league view.
type LeagueMemberPrediction struct {
    UserID      uuid.UUID
    DisplayName string
    GoalsHome   *int // nil = no prediction
    GoalsAway   *int
    Points      *int // nil = not yet scored
}

// LeagueFixtureView is the result for the league fixtures listing.
type LeagueFixtureView struct {
    Fixture     Fixture
    Predictions []LeagueMemberPrediction
}
```

**`internal/domain/prediction.go`**
```go
type ScorePrediction struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    FixtureID uuid.UUID
    GoalsHome int
    GoalsAway int
    Points    *int
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UpsertScorePredictionInput struct {
    UserID    uuid.UUID
    FixtureID uuid.UUID
    GoalsHome int
    GoalsAway int
}
```

---

### repository

**`internal/repository/fixture.go`**

Interfaces (consumed by services):
```go
type FixtureGetter interface {
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Fixture, error)
}

type FixtureLister interface {
    ListByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error)
    ListLockedByLeague(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeagueFixtureView, error)
}
```

**Two-query merge for personal view**

A LEFT JOIN for `ListByTournamentForUser` causes nullable UUID scan issues with
pgx/v5 when sqlc generates non-pointer types. Instead, two separate queries are
issued and merged in Go: `ListFixturesByTournament` + `ListPredictionsByUserAndTournament`.

**`internal/repository/queries/fixtures.sql`**

```sql
-- name: GetFixtureByID :one
SELECT id, external_id, tournament_id, home_team_id, away_team_id,
       kickoff_at, status, goals_home, goals_away, created_at, updated_at
FROM fixtures WHERE id = @id;

-- name: ListFixturesByTournament :many
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       home_t.name AS home_team_name, away_t.name AS away_team_name,
       f.kickoff_at, f.status, f.goals_home, f.goals_away, f.created_at, f.updated_at
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
WHERE f.tournament_id = @tournament_id ORDER BY f.kickoff_at;

-- name: ListLockedFixturesByLeague :many
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       home_t.name AS home_team_name, away_t.name AS away_team_name,
       f.kickoff_at, f.status, f.goals_home, f.goals_away, f.created_at, f.updated_at,
       coalesce(json_agg(json_build_object(
           'user_id', lm.user_id, 'display_name', u.display_name,
           'goals_home', sp.goals_home, 'goals_away', sp.goals_away, 'points', sp.points
       ) ORDER BY (lm.user_id = @requesting_user_id) DESC, u.display_name ASC), '[]'::json)
       AS member_predictions
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
JOIN leagues l ON l.tournament_id = f.tournament_id
JOIN league_members lm ON lm.league_id = l.id
JOIN users u ON u.id = lm.user_id
LEFT JOIN score_predictions sp ON sp.fixture_id = f.id AND sp.user_id = lm.user_id
WHERE l.id = @league_id AND f.status IN ('in_progress', 'finished')
GROUP BY f.id, home_t.name, away_t.name
ORDER BY f.kickoff_at DESC;
```

**`internal/repository/prediction.go`**

Interface:
```go
type PredictionRepo interface {
    Upsert(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
}
```

SQL (`internal/repository/queries/predictions.sql`):
```sql
-- name: UpsertScorePrediction :one
INSERT INTO score_predictions (user_id, fixture_id, goals_home, goals_away)
VALUES (@user_id, @fixture_id, @goals_home, @goals_away)
ON CONFLICT (user_id, fixture_id) DO UPDATE
    SET goals_home = EXCLUDED.goals_home,
        goals_away = EXCLUDED.goals_away,
        updated_at = now()
RETURNING *;
```

**`internal/repository/league.go`** (existing — add method)

```go
type LeagueMemberChecker interface {
    IsMember(ctx context.Context, leagueID, userID uuid.UUID) (bool, error)
}
```

---

### service

**`internal/service/fixture.go`**

```go
type FixtureService struct {
    fixtures FixtureLister
    leagues  LeagueMemberChecker
}

func (s *FixtureService) ListForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error)

func (s *FixtureService) ListForLeague(ctx context.Context, leagueID, userID uuid.UUID) ([]*domain.LeagueFixtureView, error)
// ListForLeague checks league membership first; returns domain.ErrForbidden if not a member.
```

**`internal/service/prediction.go`**

```go
type Clock interface {
    Now() time.Time
}

type PredictionService struct {
    predictions PredictionRepo
    fixtures    FixtureGetter
    clock       Clock
}

func (s *PredictionService) UpsertScore(ctx context.Context, in domain.UpsertScorePredictionInput) (*domain.ScorePrediction, error)
// Validates goals >= 0, fetches fixture, rejects with ErrForbidden if kicked off.
```

---

### handler

**`internal/server/handler/fixture.go`**

```go
type Fixture struct {
    logger *slog.Logger
    svc    FixtureService
}

func (h *Fixture) ListForUser(w http.ResponseWriter, r *http.Request)   // GET /tournaments/{tournamentId}/fixtures
func (h *Fixture) ListForLeague(w http.ResponseWriter, r *http.Request) // GET /leagues/{leagueId}/predictions
```

**`internal/server/handler/prediction.go`**

```go
type Prediction struct {
    logger *slog.Logger
    svc    PredictionService
}

func (h *Prediction) UpsertScore(w http.ResponseWriter, r *http.Request) // PUT /predictions/{fixtureId}
```

Route registration:
```
PUT  /predictions/{fixtureId}
GET  /tournaments/{tournamentId}/fixtures
GET  /leagues/{leagueId}/predictions
```

---

## Postman

Three new request files:

**`postman/collections/footy-forecast/predictions/UpsertScorePrediction.request.yaml`**
- Method: `PUT`
- URL: `{{baseUrl}}/predictions/{{fixtureId}}`
- Header: `Authorization: Bearer {{authToken}}`
- Body: `{"goals_home": 2, "goals_away": 1}`
- Post-response script: capture `id` → `{{predictionId}}`

**`postman/collections/footy-forecast/tournaments/ListFixtures.request.yaml`**
- Method: `GET`
- URL: `{{baseUrl}}/tournaments/{{tournamentId}}/fixtures`
- Header: `Authorization: Bearer {{authToken}}`

**`postman/collections/footy-forecast/leagues/ListPredictions.request.yaml`**
- Method: `GET`
- URL: `{{baseUrl}}/leagues/{{leagueId}}/predictions`
- Header: `Authorization: Bearer {{authToken}}`

---

## Edge cases

| Case | Expected behaviour |
|------|--------------------|
| Negative goals | 400 — service validates before DB |
| `fixtureId` not a UUID | 400 — handler parse error |
| Fixture not found | 404 — repo returns `ErrNotFound` |
| Within 30 min of kickoff (or past it) | 403 — service checks `lockAt = kickoff_at - 30m` via Clock |
| Same prediction submitted twice before kickoff | 200 — upsert returns updated row |
| `points` already set, prediction updated before kickoff | `points` untouched — upsert excludes the column |
| User not a member of requested league | 403 — service membership check |
| League has members with no predictions yet | predictions array entry has `goals_home: null, goals_away: null, points: null` |
| Fixture is `upcoming` — excluded from league view | filtered server-side by `status IN ('in_progress', 'finished')` |

---

## Test plan

### Repository tests (testcontainers)

**FixtureRepository**
- `ListByTournamentForUser`: returns all fixtures; prediction fields are null when no prediction exists; populated when one does.
- `ListLockedByLeague`: returns only `in_progress`/`finished` fixtures; `member_predictions` JSON includes all members, null goals for those who haven't predicted.

**PredictionRepository**
- `Upsert` creates a new row; second call with same (user, fixture) updates goals, does not touch `points`.
- `updated_at` advances on second upsert.

### Service tests (fake repos)

**FixtureService**
- `ListForLeague`: non-member gets `ErrForbidden`; member gets results.

**PredictionService**
- `GoalsHome < 0` → `ErrInvalid`.
- `GoalsAway < 0` → `ErrInvalid`.
- Fixture not found → `ErrNotFound`.
- Within 30 minutes of kickoff → `ErrForbidden` (inject frozen Clock at `kickoff_at - 29m`).
- Valid input, 31+ minutes before kickoff → delegates to repo, returns result.

### Handler tests (httptest)

**Fixture handler**
- `ListForUser` happy path: 200 with correct shape (team names, not IDs).
- `ListForLeague` (GET /leagues/{leagueId}/predictions) happy path: 200 with predictions array including `display_name`.
- `ListForLeague` non-member: 403.
- Invalid UUID path params: 400.
- Unauthenticated: 401.

**Prediction handler**
- Happy path: 200 with `points: null`.
- Non-UUID `fixtureId`: 400.
- Negative goals: 400.
- Fixture not found: 404.
- Fixture kicked off: 403.
- Unauthenticated: 401.

---

## Acceptance criteria

1. `PUT /predictions/{fixtureId}` with valid body and 31+ minutes before kickoff returns 200; re-submitting updates goals only.
2. `PUT /predictions/{fixtureId}` within 30 minutes of kickoff (or after) returns 403.
3. `GET /tournaments/{tournamentId}/fixtures` returns all fixtures with `home_team_name`/`away_team_name`; prediction is null when none made.
4. `GET /leagues/{leagueId}/predictions` returns only locked fixtures ordered by `kickoff_at DESC`; non-member gets 403.
5. League prediction response includes an entry per member regardless of whether they predicted; requesting user appears first, others alphabetical by `display_name`.
6. Negative goals return 400; unauthenticated requests return 401.
7. `make fmt && make lint && make test` all pass.
8. All three Postman files committed alongside the code.
