# Plan: Leagues

## Goal

Allow users to compete in private groups (leagues) scoped to a tournament. A league owner creates the league and shares an invite code; anyone with the code can join and compete on a shared leaderboard.

---

## Design Decisions

- **Invite code**: single static `code TEXT UNIQUE` on the `leagues` table. Owner can regenerate it. No expiry, no separate invite table.
- **Membership**: self-join via code. Owner can remove any member. Members can leave themselves. No owner-adds-by-email.
- **Players storage**: separate `league_members` join table (not an array column).
- **Owner**: `owner_id` denormalized on `leagues` for cheap authorization checks; also reflected as `role='owner'` in `league_members`.
- **Owner cannot leave**: must delete the league instead.

---

## Data Model

### New migration: `migrations/20260522XXXXXX_leagues.sql`

Timestamp must come after `20260521102222_users.sql`.

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TYPE league_member_role AS ENUM ('owner', 'member');

CREATE TABLE leagues (
    id            UUID        PRIMARY KEY,
    tournament_id UUID        NOT NULL REFERENCES tournaments(id),
    owner_id      UUID        NOT NULL REFERENCES users(id),
    name          TEXT        NOT NULL,
    code          TEXT        NOT NULL UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT leagues_name_length_chk CHECK (char_length(name) BETWEEN 1 AND 100)
);

CREATE INDEX leagues_tournament_id_idx ON leagues (tournament_id);
CREATE INDEX leagues_owner_id_idx      ON leagues (owner_id);
CREATE INDEX leagues_code_idx          ON leagues (code);

CREATE TABLE league_members (
    league_id UUID               NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    user_id   UUID               NOT NULL REFERENCES users(id),
    role      league_member_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ        NOT NULL DEFAULT now(),
    PRIMARY KEY (league_id, user_id)
);

CREATE INDEX league_members_user_id_idx ON league_members (user_id);

CREATE TRIGGER leagues_set_updated_at
    BEFORE UPDATE ON leagues
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS leagues_set_updated_at ON leagues;
DROP TABLE IF EXISTS league_members;
DROP TABLE IF EXISTS leagues;
DROP TYPE IF EXISTS league_member_role;

-- +goose StatementEnd
```

Notes:
- `ON DELETE CASCADE` on `league_members.league_id` — deleting a league removes all members automatically.
- `char_length` (Unicode code points) matches Go `[]rune` validation.
- Composite PK on `(league_id, user_id)` implicitly indexes membership lookup; separate `user_id` index covers reverse lookup.
- `set_updated_at()` function already exists (created in tournaments migration).

---

## Domain Types

### New file: `internal/domain/league.go`

```go
type LeagueMemberRole string

const (
    LeagueMemberRoleOwner  LeagueMemberRole = "owner"
    LeagueMemberRoleMember LeagueMemberRole = "member"
)

type League struct {
    ID           uuid.UUID
    TournamentID uuid.UUID
    OwnerID      uuid.UUID
    Name         string
    Code         string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type LeagueMember struct {
    LeagueID uuid.UUID
    UserID   uuid.UUID
    Role     LeagueMemberRole
    JoinedAt time.Time
}

type CreateLeagueInput struct {
    TournamentID uuid.UUID
    Name         string
}
```

### Update: `internal/domain/errors.go`

Add:
```go
// ErrForbidden is returned when the requester lacks permission for the operation.
var ErrForbidden = errors.New("forbidden")
```

---

## SQL Queries

### New file: `internal/repository/queries/league.sql`

| Query name | Type | Purpose |
|---|---|---|
| `CreateLeague` | `:one` | INSERT league, return full row |
| `GetLeagueByID` | `:one` | SELECT by id |
| `GetLeagueByCode` | `:one` | SELECT by code |
| `ListLeaguesForUser` | `:many` | JOIN league_members WHERE user_id = $1, ORDER BY created_at DESC |
| `UpdateLeagueName` | `:one` | UPDATE name WHERE id = $1 |
| `UpdateLeagueCode` | `:one` | UPDATE code WHERE id = $1 |
| `DeleteLeague` | `:exec` | DELETE WHERE id = $1 |
| `AddLeagueMember` | `:one` | INSERT league_members |
| `RemoveLeagueMember` | `:exec` | DELETE WHERE league_id=$1 AND user_id=$2 |
| `GetLeagueMember` | `:one` | SELECT WHERE league_id=$1 AND user_id=$2 |
| `ListLeagueMembersForLeague` | `:many` | WHERE league_id=$1 ORDER BY joined_at ASC |

Run `sqlc generate` after writing these to regenerate `internal/repository/dbgen/`.

---

## Repository

### New file: `internal/repository/league.go`

```go
type LeagueRepository struct { q *dbgen.Queries }
func NewLeagueRepository(pool *db.Pool) *LeagueRepository

type CreateLeagueParams struct {
    TournamentID uuid.UUID
    OwnerID      uuid.UUID
    Name         string
    Code         string
}

func (r *LeagueRepository) Create(ctx, CreateLeagueParams) (*domain.League, error)
func (r *LeagueRepository) GetByID(ctx, uuid.UUID) (*domain.League, error)
func (r *LeagueRepository) GetByCode(ctx, string) (*domain.League, error)
func (r *LeagueRepository) ListForUser(ctx, uuid.UUID) ([]*domain.League, error)
func (r *LeagueRepository) UpdateName(ctx, id uuid.UUID, name string) (*domain.League, error)
func (r *LeagueRepository) UpdateCode(ctx, id uuid.UUID, code string) (*domain.League, error)
func (r *LeagueRepository) Delete(ctx, uuid.UUID) error
func (r *LeagueRepository) AddMember(ctx, leagueID, userID uuid.UUID, role domain.LeagueMemberRole) (*domain.LeagueMember, error)
func (r *LeagueRepository) RemoveMember(ctx, leagueID, userID uuid.UUID) error
func (r *LeagueRepository) GetMember(ctx, leagueID, userID uuid.UUID) (*domain.LeagueMember, error)
func (r *LeagueRepository) ListMembers(ctx, leagueID uuid.UUID) ([]*domain.LeagueMember, error)
```

Error translation:
- `pgx.ErrNoRows` → `domain.ErrNotFound` in `GetByID`, `GetByCode`, `GetMember`
- `isUniqueViolation` → `domain.ErrConflict` in `Create` (code collision) and `AddMember` (duplicate membership)
- `isUniqueViolation` already exists in the `repository` package — do not redefine it

---

## Service

### New file: `internal/service/league.go`

Service defines two minimal repo interfaces (only methods it needs).

**`LeagueRepo` interface** — see Repository section above.

**`TournamentGetter` interface** — narrow interface for tournament validation:
```go
type TournamentGetter interface {
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Tournament, error)
}
```

Compile-time assertions at bottom of file:
```go
var _ LeagueRepo       = (*repository.LeagueRepository)(nil)
var _ TournamentGetter = (*repository.TournamentRepository)(nil)
```

`LeagueService` is constructed with both:
```go
func NewLeagueService(leagues LeagueRepo, tournaments TournamentGetter) *LeagueService
```

**Code generation** — 8-char uppercase alphanumeric from `crypto/rand`:
```go
const codeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() (string, error) { /* read 8 random bytes, map each via codeChars[v%36] */ }
```

**Method business rules:**

| Method | Rule |
|---|---|
| `CreateLeague(ctx, userID, CreateLeagueInput)` | Validate name (1–100 runes, trimmed) → `tournaments.GetByID` (→ `ErrNotFound` if missing; `ErrInvalid` if `status == concluded`) → generate code → `repo.Create` → on `ErrConflict` retry once → `repo.AddMember(owner)` |
| `GetLeague(ctx, leagueID, requesterID)` | `GetByID` → `GetMember` (if ErrNotFound → return ErrNotFound, not ErrForbidden, to avoid leaking existence) → `ListMembers` |
| `ListLeaguesForUser(ctx, userID)` | Delegates to repo |
| `UpdateLeagueName(ctx, leagueID, requesterID, name)` | `GetByID` → owner check (→ `ErrForbidden`) → validate name → `UpdateName` |
| `DeleteLeague(ctx, leagueID, requesterID)` | `GetByID` → owner check → `Delete` |
| `RegenerateCode(ctx, leagueID, requesterID)` | `GetByID` → owner check → generate code → `UpdateCode`, retry once on collision |
| `JoinLeague(ctx, code, userID)` | `GetByCode` → `GetMember` (if nil error → `ErrConflict` "already a member") → `AddMember(member)` |
| `RemoveMember(ctx, leagueID, targetUserID, requesterID)` | `GetByID` → if requester ≠ target and requester ≠ owner → `ErrForbidden` → if target == owner → `ErrInvalid` ("owner cannot leave; delete the league") → `RemoveMember` |

Compile-time assertion at bottom:
```go
var _ LeagueRepo = (*repository.LeagueRepository)(nil)
```

---

## Handler

### New file: `internal/server/handler/league.go`

All endpoints require auth — extract user with `ctxutil.UserFromCtx`.

Handler defines a minimal `LeagueService` interface. Compile-time assertion:
```go
var _ LeagueService = (*service.LeagueService)(nil)
```

**DTOs:**
```go
type leagueResponse struct {
    ID, TournamentID, OwnerID, Name, Code string
    CreatedAt, UpdatedAt                  time.Time
}
type leagueMemberResponse struct {
    UserID, Role string
    JoinedAt     time.Time
}
type leagueDetailResponse struct {
    leagueResponse
    Members []leagueMemberResponse `json:"members"`
}
```

**Endpoints:**

| Method | Path | Handler | Status | Notes |
|---|---|---|---|---|
| POST | `/leagues` | `Create` | 201 | `{tournament_id, name}` |
| GET | `/leagues` | `List` | 200 | `{"leagues":[…]}` |
| GET | `/leagues/{id}` | `Get` | 200 | league + members |
| PATCH | `/leagues/{id}` | `UpdateName` | 200 | `{name}` |
| DELETE | `/leagues/{id}` | `Delete` | 204 | |
| POST | `/leagues/{id}/code` | `RegenerateCode` | 200 | returns updated league |
| POST | `/leagues/join` | `Join` | 200 | `{code}` |
| DELETE | `/leagues/{id}/members/{userId}` | `RemoveMember` | 204 | |

### Update: `internal/server/handler/errors.go`

Add `ErrForbidden → 403` case before the `default` case in `writeError`.

### Update: `internal/server/router.go`

Wire bottom-up after existing wiring. All 8 routes use `authMW`.

---

## Test Plan

### Repository (`internal/repository/league_test.go`) — `//go:build integration`

Use `startPostgres(t)` from `testdb_test.go`. Insert fixture users and tournaments via raw SQL.

- Create: returned fields correct; duplicate code → ErrConflict
- GetByID: found; not found → ErrNotFound
- GetByCode: found; not found → ErrNotFound
- ListForUser: empty; only returns leagues the user belongs to
- UpdateName, UpdateCode: fields updated
- Delete: cascades to members (GetByID, GetMember both → ErrNotFound after delete)
- AddMember: persisted; duplicate → ErrConflict
- RemoveMember: removed
- GetMember: not found → ErrNotFound
- ListMembers: ordered by joined_at

### Service (`internal/service/league_test.go`)

Hand-written `fakeLeagueRepo` (same pattern as `fakeRepo` in `tournament_test.go`).

- Name validation: empty, whitespace-only, too long → ErrInvalid, repo not called
- CreateLeague: tournament not found → ErrNotFound; tournament concluded → ErrInvalid; AddMember called with owner role; code-collision retry
- GetLeague: ErrNotFound when not a member (same as league-not-found — no distinction)
- UpdateLeagueName: ErrForbidden when not owner
- DeleteLeague, RegenerateCode: ErrForbidden when not owner
- JoinLeague: ErrNotFound bad code; ErrConflict already member
- RemoveMember: ErrForbidden (non-owner removing other); ErrInvalid (owner leaving); success for self-leave and owner-removes-member

### Handler (`internal/server/handler/league_test.go`)

`fakeLeagueService` with `fn` fields. Use `ctxutil.WithUser` to inject user into request context.

At least one happy-path and one error-path per endpoint. Cover 201/200/204 successes and 400/403/404/409/500 error cases.

---

## Acceptance Criteria

1. `POST /leagues` → 201 with `code` field; creator is automatically the owner member.
2. `POST /leagues/join {code}` → 200, user added as member; same code returns 409.
3. `GET /leagues/{id}` as non-member → 404 (existence not leaked).
4. `PATCH /leagues/{id}` as non-owner → 403.
5. `DELETE /leagues/{id}/members/{ownerId}` (owner self-remove) → 400.
6. `DELETE /leagues/{id}` → 204; subsequent GET → 404.
7. `make fmt && make lint && make test` all pass.
8. Integration tests pass: `make test-integration`.

---

## Implementation Order

1. Migration file
2. `internal/domain/league.go` + update `errors.go`
3. `internal/repository/queries/league.sql` → run `sqlc generate`
4. `internal/repository/league.go` + `league_test.go`
5. `internal/service/league.go` + `league_test.go`
6. `internal/server/handler/errors.go` (add ErrForbidden case)
7. `internal/server/handler/league.go` + `league_test.go`
8. `internal/server/router.go`
9. `make fmt && make lint && make test && make test-integration`
