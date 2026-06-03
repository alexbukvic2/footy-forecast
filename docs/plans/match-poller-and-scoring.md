# Plan: Match Poller and Live Scoring Worker

## Goal

Implement a single adaptive background goroutine inside the existing API binary that polls an external football API every 60 seconds when matches are live or imminent, writes live score and status changes to the database, rescores `score_predictions` in real time on every change, and triggers idempotent batch settlement of all outright prediction categories as each tournament phase concludes.

---

## Out of Scope

- Notification delivery to clients (separate plan).
- Fixture import and seeding from the external API (fixtures already in DB via a separate import flow; the worker only reads what is already present).
- Admin HTTP endpoints to manually re-trigger settlement.
- WebSocket / SSE push to clients; clients continue polling the existing REST endpoints.
- Per-provider rate-limit management beyond simple next-tick retries.
- Leaderboard query changes; the existing `SUM(points)` cross-table query automatically picks up updated `points` columns.
- Any new HTTP routes — this plan adds none; no API Contract section is needed.

---

## Data Model Changes

### Migration 1 — `migrations/20260528000001_worker_columns.sql`

**`fixtures` additions**

| Column | Type | Notes |
|--------|------|-------|
| `is_demo` | `BOOLEAN NOT NULL DEFAULT FALSE` | Worker skips `is_demo = TRUE` fixtures |
| `winner_team_id` | `UUID NULL REFERENCES teams(id)` | Set on knockout conclusion to the advancing team |
| `last_polled_at` | `TIMESTAMPTZ NULL` | Stamped on each successful API poll |

**`fixture_status` enum** — add `cancelled` value:

```sql
ALTER TYPE fixture_status ADD VALUE 'cancelled';
```

Used for API statuses `CANC` and `ABD`. Postgres 16 supports this inside a transaction.

**`tournaments` additions**

| Column | Type | Notes |
|--------|------|-------|
| `external_id` | `BIGINT NULL` | Football API league ID; worker skips tournaments where NULL |
| `season` | `SMALLINT NULL` | Season year e.g. 2026; needed for standings + topscorer API calls |

**`teams` addition**

| Column | Type | Notes |
|--------|------|-------|
| `external_id` | `BIGINT NULL` | Football API team ID; used by `GetTeamByExternalID` to resolve standings entries |

**Prediction table audit columns** (all three prediction tables)

- `score_predictions.scored_at TIMESTAMPTZ NULL`
- `player_predictions.scored_at TIMESTAMPTZ NULL`
- `team_predictions.scored_at TIMESTAMPTZ NULL`

Set/updated whenever the worker writes `points` to a row.

**New index**

```sql
CREATE INDEX score_predictions_fixture_id_idx ON score_predictions (fixture_id);
```

Required for `UPDATE score_predictions WHERE fixture_id = $1`. The existing unique index on `(user_id, fixture_id)` does not efficiently serve lookups by `fixture_id` alone.

## Package Structure

### `internal/worker/`

| File | Contents |
|------|----------|
| `worker.go` | `Worker` struct; `New(repo WorkerRepo, api MatchAPI, clock Clock, logger *slog.Logger)` constructor; `Run(ctx) error`; `tick(ctx)`; `processSingleFixture(ctx, domain.PollableFixture)` |
| `repository.go` | `WorkerRepo` interface (consumer-defined, per project convention) |
| `match_api.go` | `MatchAPI` interface; `APIFixtureResult`, `APIStandingsEntry`, `APITopScorerResult` types |
| `clock.go` | `Clock interface { Now() time.Time }`; `RealClock struct{}` |
| `semaphore.go` | Unexported `semaphore` backed by `chan struct{}`; `newSemaphore(n)`, `acquire(ctx)`, `release()` |
| `settlement.go` | `runSettlement(ctx, f, result, newStatus)`; helpers `isQuarterfinalRound`, `isFinalRound`, `resolveWinnerTeamID`, `resolveTeamIDs` |
| `worker_test.go` | Unit tests with hand-written fakes for `WorkerRepo`, `MatchAPI`, `Clock` |

### `internal/footballapi/`

| File | Contents |
|------|----------|
| `client.go` | `Client` struct implementing `worker.MatchAPI`; constructed with API key, base URL, `*http.Client` |
| `types.go` | JSON-mapped structs for football API fixture and standings responses |

### `internal/repository/`

| File | Contents |
|------|----------|
| `worker.go` | `WorkerRepository` struct; `NewWorkerRepository(pool)`; compile-time check `var _ worker.WorkerRepo = (*WorkerRepository)(nil)`; non-transactional methods using sqlc-generated queries |
| `worker_settlement.go` | `UpdateMatchAndRescoreLivePredictions` and all `Settle*` methods; raw pgx with explicit transaction control |
| `worker_test.go` | Integration tests with testcontainers-go |

### `internal/repository/queries/worker.sql`

New sqlc query file: `ListPollableMatches`, `IsGroupComplete`, `IsRoundComplete`, `IsGroupStageComplete`, `GetTeamByExternalID`, `GetPlayerByExternalID`, `UpdateFixtureStatus`. After adding this file, run `make sqlc-gen` to regenerate `internal/repository/dbgen/`.

---

## Worker Repository Interface

Defined in `internal/worker/repository.go`:

```go
type WorkerRepo interface {
    // Polling
    ListPollableMatches(ctx context.Context) ([]domain.PollableFixture, error)
    UpdateMatchAndRescoreLivePredictions(ctx context.Context, f domain.PollableFixture, result APIFixtureResult) error

    // Group standings
    UpdateGroupStandings(ctx context.Context, tournamentID uuid.UUID, groupLetter string, entries []domain.StandingsEntry) error

    // Completion checks
    IsGroupComplete(ctx context.Context, tournamentID uuid.UUID, groupLetter string) (bool, error)
    IsRoundComplete(ctx context.Context, tournamentID uuid.UUID, round string) (bool, error)
    IsGroupStageComplete(ctx context.Context, tournamentID uuid.UUID) (bool, error)

    // ID resolution
    GetTeamByExternalID(ctx context.Context, externalID int64, tournamentID uuid.UUID) (uuid.UUID, error)
    GetPlayerByExternalID(ctx context.Context, externalID string, tournamentID uuid.UUID) (uuid.UUID, error)

    // Settlement
    SettleGroupWinnerPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string) error
    SettlePlayoffGroupPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string) error
    SettlePlayoffWildcardPredictions(ctx context.Context, tournamentID uuid.UUID) error
    SettleGroupTopScorerPredictions(ctx context.Context, tournamentID uuid.UUID, groupLetter string, topScorerPlayerID uuid.UUID) error
    SettleSemifinalistPredictions(ctx context.Context, tournamentID uuid.UUID) error
    SettleTournamentWinnerPredictions(ctx context.Context, tournamentID uuid.UUID, winnerTeamID uuid.UUID) error
    SettleTopScorerPredictions(ctx context.Context, tournamentID uuid.UUID, topScorerPlayerID uuid.UUID) error
}
```

---

## MatchAPI Interface

Defined in `internal/worker/match_api.go`:

```go
type MatchAPI interface {
    GetFixture(ctx context.Context, externalFixtureID int64) (APIFixtureResult, error)
    GetStandings(ctx context.Context, externalLeagueID int64, season int) ([]APIStandingsEntry, error)
    GetGroupTopScorer(ctx context.Context, externalLeagueID int64, season int, groupLetter string) (APITopScorerResult, error)
    GetTournamentTopScorer(ctx context.Context, externalLeagueID int64, season int) (APITopScorerResult, error)
}

// APIFixtureResult is the data returned by GetFixture.
type APIFixtureResult struct {
    ExternalID  int64
    StatusShort string  // raw API status: "1H", "HT", "FT", "CANC", etc.
    GoalsHome   *int
    GoalsAway   *int
    HomeWinner  *bool   // teams.home.winner; handles ET/PEN
    AwayWinner  *bool
}

// APIStandingsEntry is one team row from GetStandings.
type APIStandingsEntry struct {
    TeamExternalID                               int64
    Position, Points, Played, Won, Drawn, Lost   int
    GoalsFor, GoalsAgainst                       int
}

// APITopScorerResult is returned by GetGroupTopScorer and GetTournamentTopScorer.
type APITopScorerResult struct {
    PlayerExternalID string // matches players.external_id
    Goals            int
}
```

---

## New Domain Types

In `internal/domain/worker.go`:

```go
// PollableFixture is returned by ListPollableMatches and carries all
// fields the worker needs to poll and settle a fixture.
type PollableFixture struct {
    ID                   uuid.UUID
    ExternalID           int64
    TournamentID         uuid.UUID
    TournamentExternalID int64
    TournamentSeason     int
    HomeTeamID           uuid.UUID
    AwayTeamID           uuid.UUID
    GroupLetter          *string   // nil = knockout match
    Round                string
    Status               FixtureStatus
    KickoffAt            time.Time
    GoalsHome            *int
    GoalsAway            *int
    WinnerTeamID         *uuid.UUID
    LastPolledAt         *time.Time
}

// StandingsEntry is one team's row as resolved to our internal team UUID.
type StandingsEntry struct {
    TeamID                             uuid.UUID
    Position, Points, Played           int
    Won, Drawn, Lost                   int
    GoalsFor, GoalsAgainst             int
}
```

In `internal/domain/fixture.go`: add `FixtureStatusCancelled FixtureStatus = "cancelled"`.

---

## Poll Loop Algorithm

```
Worker.Run(ctx):
    for {
        fixtures, err = repo.ListPollableMatches(ctx)
        if err: log.Warn; sleep 60s; continue
        if len(fixtures) == 0: log "worker idle"; sleep 5m; continue

        sem = newSemaphore(5)
        for each f in fixtures:
            sem.acquire(ctx)
            go func(f):
                defer sem.release()
                w.processSingleFixture(ctx, f)
        // drain: acquire all 5 slots to ensure all goroutines finished
        sleep 60s, or return immediately on ctx.Done()
    }

Worker.processSingleFixture(ctx, f):
    result, err = api.GetFixture(ctx, f.ExternalID)
    if err: log.Warn; return

    newStatus = mapAPIStatus(result.StatusShort)
    unchanged = newStatus == f.Status &&
                result.GoalsHome == f.GoalsHome &&
                result.GoalsAway == f.GoalsAway &&
                newStatus != cancelled
    if unchanged: return

    err = repo.UpdateMatchAndRescoreLivePredictions(ctx, f, result)
    if err: log.Error; return

    settlement.runSettlement(ctx, f, result, newStatus)
```

`mapAPIStatus` mapping:

| API `status.short` | Domain `FixtureStatus` |
|--------------------|------------------------|
| `NS` | `upcoming` |
| `1H` `HT` `2H` `ET` `BT` `P` `SUSP` `INT` | `in_progress` |
| `FT` `AET` `PEN` `AWD` `WO` | `finished` |
| `CANC` `ABD` | `cancelled` |
| `PST` | `upcoming` (update `kickoff_at` if API provides new time) |

`ListPollableMatches` query sketch (sqlc, `internal/repository/queries/worker.sql`):

```sql
-- name: ListPollableMatches :many
SELECT f.id, f.external_id, f.tournament_id,
       f.home_team_id, f.away_team_id,
       t.external_id  AS tournament_external_id,
       t.season       AS tournament_season,
       home_t.group_letter,
       f.round, f.status, f.kickoff_at,
       f.goals_home, f.goals_away,
       f.winner_team_id, f.last_polled_at
FROM fixtures f
JOIN tournaments t    ON t.id = f.tournament_id
JOIN teams     home_t ON home_t.id = f.home_team_id
WHERE f.is_demo = FALSE
  AND t.external_id IS NOT NULL
  AND (
    f.status = 'in_progress'
    OR (f.status = 'upcoming' AND f.kickoff_at <= now() + INTERVAL '5 minutes')
    OR (f.status = 'finished' AND f.updated_at >= now() - INTERVAL '24 hours')
  )
ORDER BY f.kickoff_at ASC;
```

---

## Settlement Trigger Logic

Called from `settlement.go` after each successful `UpdateMatchAndRescoreLivePredictions`. Only fires once on first transition to a terminal state.

```
runSettlement(ctx, f, result, newStatus):
    prevTerminal = f.Status IN {finished, cancelled}
    nowTerminal  = newStatus IN {finished, cancelled}
    if prevTerminal OR NOT nowTerminal: return

    if newStatus == cancelled: return
    // 0 points already written by UpdateMatchAndRescoreLivePredictions;
    // no outright settlement needed for cancelled matches.

    winnerTeamID = resolveWinnerTeamID(result, f)

    if f.GroupLetter != nil:
        // Group-stage match concluded
        standings, err = api.GetStandings(ctx, f.TournamentExternalID, f.TournamentSeason)
        if err: log.Warn; return  // standings not yet available; retry on next tick
        entries = resolveTeamIDs(ctx, standings, f.TournamentID)
        repo.UpdateGroupStandings(ctx, f.TournamentID, *f.GroupLetter, entries)

        groupDone, _ = repo.IsGroupComplete(ctx, f.TournamentID, *f.GroupLetter)
        if groupDone:
            repo.SettleGroupWinnerPredictions(ctx, f.TournamentID, *f.GroupLetter)
            repo.SettlePlayoffGroupPredictions(ctx, f.TournamentID, *f.GroupLetter)

            topScorer, err = api.GetGroupTopScorer(ctx, ...)
            if err == nil:
                playerID, err = repo.GetPlayerByExternalID(ctx, topScorer.PlayerExternalID, f.TournamentID)
                if err == nil:
                    repo.SettleGroupTopScorerPredictions(ctx, f.TournamentID, *f.GroupLetter, playerID)

            allGroupsDone, _ = repo.IsGroupStageComplete(ctx, f.TournamentID)
            if allGroupsDone:
                repo.SettlePlayoffWildcardPredictions(ctx, f.TournamentID)
    else:
        // Knockout match concluded
        if isQuarterfinalRound(f.Round):
            roundDone, _ = repo.IsRoundComplete(ctx, f.TournamentID, f.Round)
            if roundDone:
                repo.SettleSemifinalistPredictions(ctx, f.TournamentID)

        if isFinalRound(f.Round) AND winnerTeamID != nil:
            repo.SettleTournamentWinnerPredictions(ctx, f.TournamentID, winnerTeamID)
            topScorer, err = api.GetTournamentTopScorer(ctx, ...)
            if err == nil:
                playerID, err = repo.GetPlayerByExternalID(ctx, topScorer.PlayerExternalID, f.TournamentID)
                if err == nil:
                    repo.SettleTopScorerPredictions(ctx, f.TournamentID, playerID)
```

`resolveWinnerTeamID(result, f)`: returns `f.HomeTeamID` if `result.HomeWinner == true`, `f.AwayTeamID` if `result.AwayWinner == true`, nil otherwise.

`isQuarterfinalRound(round string)`: case-insensitive contains `"quarter"`.

`isFinalRound(round string)`: case-insensitive ends with `"final"` but does not contain `"semi"` or `"quarter"`. Exact API round strings are an open question; see open questions section.

**Note on settlement retry after standings failure**: when `api.GetStandings` fails, the worker returns early without settling. On the next tick `f.Status` will already be `finished` (terminal), so `prevTerminal == true` and settlement is skipped. The implementer must handle this by NOT using `prevTerminal` as a guard after a partial failure; instead, track settlement completion via a `group_settled_at` column on `tournament_group_table`, or check whether the Settle* functions have already run by querying `COUNT(*) WHERE points IS NOT NULL` in the group. See open question 13 for the recommended approach. A pragmatic option: remove the `prevTerminal` guard entirely and rely solely on `WHERE points IS NULL` idempotency within the Settle* functions; calls to standings/topscorer APIs are cheap and idempotent.

---

## Per-Settlement Algorithms

All `Settle*` functions execute inside a single pgx transaction. The `WHERE points IS NULL` predicate provides idempotency — already-settled rows are skipped.

### `UpdateMatchAndRescoreLivePredictions`

```sql
-- Run in a transaction.

-- 1. Update fixture
UPDATE fixtures
SET status            = $newStatus,
    goals_home        = $goalsHome,
    goals_away        = $goalsAway,
    winner_team_id    = $winnerTeamID,
    prediction_locked = (kickoff_at <= now()),
    last_polled_at    = now()
WHERE id = $fixtureID;

-- 2a. Cancelled: zero all score predictions (overwrite live points too)
UPDATE score_predictions
SET points = 0, scored_at = now()
WHERE fixture_id = $fixtureID;
-- (only run this branch when newStatus == 'cancelled')

-- 2b. Active result with known goals: compute points
UPDATE score_predictions
SET points = CASE
    WHEN goals_home = $gh AND goals_away = $ga             THEN 3
    WHEN SIGN(goals_home - goals_away) = SIGN($gh - $ga)   THEN 1
    ELSE 0
END,
scored_at = now()
WHERE fixture_id = $fixtureID;
-- (only run when $gh IS NOT NULL AND $ga IS NOT NULL)
```

Scoring formula: exact score = 3 pts; correct result (home win / draw / away win via `SIGN`) = 1 pt; wrong = 0 pts. Applied on every score change including during live play.

For PEN/AET: `goals_home`/`goals_away` = aggregate regulation + extra time (not penalty kicks). `winner_team_id` comes from `HomeWinner`/`AwayWinner` booleans.

### `SettleGroupWinnerPredictions(tournamentID, groupLetter)`

```sql
-- In transaction:
-- 1. Fetch pos1TeamID and pos2TeamID from tournament_group_table ORDER BY position ASC LIMIT 2

-- 2. Record outcomes
INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'group_winner', $pos1), ($1, 'group_winner', $pos2)
ON CONFLICT (tournament_id, category, team_id) DO UPDATE SET recorded_at = now();

-- 3. Award handicap to correct slot-position picks
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'group_winner'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'group_winner'
  AND tp.group_letter  = $2
  AND tp.points IS NULL
  AND ((tp.slot_index = 0 AND tp.pick = $pos1)
    OR (tp.slot_index = 1 AND tp.pick = $pos2));

-- 4. Zero all remaining unsettled rows for this group
UPDATE team_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'group_winner'
  AND group_letter = $2 AND points IS NULL;
```

Assumption: `slot_index = 0` = predicted 1st place; `slot_index = 1` = predicted 2nd. Points awarded only on exact slot-position match (see open questions for confirmation).

### `SettlePlayoffGroupPredictions(tournamentID, groupLetter)`

```sql
-- advancedTeams = {pos1TeamID, pos2TeamID}

INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'playoff', $pos1), ($1, 'playoff', $pos2)
ON CONFLICT (tournament_id, category, team_id) DO UPDATE SET recorded_at = now();

-- Award points: any pick that advanced earns handicap, regardless of slot index
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'playoff'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'playoff'
  AND tp.group_letter  = $2
  AND tp.slot_index   IN (0, 1)
  AND tp.pick         IN ($pos1, $pos2)
  AND tp.points IS NULL;

UPDATE team_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'playoff'
  AND group_letter = $2 AND slot_index IN (0, 1) AND points IS NULL;
```

### `SettlePlayoffWildcardPredictions(tournamentID)`

Called once after `IsGroupStageComplete`. All advancing teams are already in `team_outcomes(category='playoff')`.

```sql
UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'playoff'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1
  AND tp.category      = 'playoff'
  AND tp.group_letter IS NULL
  AND tp.slot_index    = 2
  AND tp.pick IN (SELECT team_id FROM team_outcomes
                  WHERE tournament_id = $1 AND category = 'playoff')
  AND tp.points IS NULL;

UPDATE team_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'playoff'
  AND group_letter IS NULL AND slot_index = 2 AND points IS NULL;
```

### `SettleGroupTopScorerPredictions(tournamentID, groupLetter, topScorerPlayerID)`

```sql
INSERT INTO player_outcomes (tournament_id, category, player_id)
VALUES ($1, 'group_top_scorer', $playerID)
ON CONFLICT (tournament_id, category, player_id) DO UPDATE SET recorded_at = now();

UPDATE player_predictions pp
SET points    = COALESCE((SELECT h.points FROM player_handicap h
                          WHERE h.player_id = pp.pick AND h.category = 'group_top_scorer'), 0),
    scored_at = now()
WHERE pp.tournament_id = $1 AND pp.category = 'group_top_scorer'
  AND pp.group_letter = $2 AND pp.pick = $playerID AND pp.points IS NULL;

UPDATE player_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'group_top_scorer'
  AND group_letter = $2 AND pick != $playerID AND points IS NULL;
```

### `SettleSemifinalistPredictions(tournamentID)`

```sql
-- sfTeams = SELECT winner_team_id FROM fixtures
--           WHERE tournament_id = $1
--             AND LOWER(round) LIKE '%quarter%'
--             AND winner_team_id IS NOT NULL
-- Return error if count != 4.

INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'semifinalist', $sf1), ($1, 'semifinalist', $sf2),
       ($1, 'semifinalist', $sf3), ($1, 'semifinalist', $sf4)
ON CONFLICT DO UPDATE SET recorded_at = now();

UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'semifinalist'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1 AND tp.category = 'semifinalist'
  AND tp.pick IN ($sf1, $sf2, $sf3, $sf4) AND tp.points IS NULL;

UPDATE team_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'semifinalist' AND points IS NULL;
```

### `SettleTournamentWinnerPredictions(tournamentID, winnerTeamID)`

```sql
INSERT INTO team_outcomes (tournament_id, category, team_id)
VALUES ($1, 'winner', $winnerTeamID)
ON CONFLICT DO UPDATE SET recorded_at = now();

UPDATE team_predictions tp
SET points    = COALESCE((SELECT h.points FROM team_handicap h
                          WHERE h.team_id = tp.pick AND h.category = 'winner'), 0),
    scored_at = now()
WHERE tp.tournament_id = $1 AND tp.category = 'winner'
  AND tp.pick = $winnerTeamID AND tp.points IS NULL;

UPDATE team_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'winner'
  AND pick != $winnerTeamID AND points IS NULL;
```

### `SettleTopScorerPredictions(tournamentID, topScorerPlayerID)`

```sql
INSERT INTO player_outcomes (tournament_id, category, player_id)
VALUES ($1, 'total_top_scorer', $playerID)
ON CONFLICT (tournament_id, category, player_id) DO UPDATE SET recorded_at = now();

UPDATE player_predictions pp
SET points    = COALESCE((SELECT h.points FROM player_handicap h
                          WHERE h.player_id = pp.pick AND h.category = 'total_top_scorer'), 0),
    scored_at = now()
WHERE pp.tournament_id = $1 AND pp.category = 'total_top_scorer'
  AND pp.pick = $playerID AND pp.points IS NULL;

UPDATE player_predictions
SET points = 0, scored_at = now()
WHERE tournament_id = $1 AND category = 'total_top_scorer'
  AND pick != $playerID AND points IS NULL;
```

---

## Cancelled/Abandoned Match Handling

| API status | Action |
|------------|--------|
| `CANC`, `ABD` | `fixtures.status = 'cancelled'`; `score_predictions.points = 0` for all rows on the fixture (unconditional, overwriting any live points). No outright settlement triggered. |
| `PST` | Update `fixtures.kickoff_at` if API provides a new time; leave `status = 'upcoming'`, `prediction_locked = FALSE`. No score change. |
| `SUSP`, `INT` | Leave `status = 'in_progress'`; continue polling. Apply cancelled handling if/when the match reaches `ABD`. |
| `WO`, `AWD` | Treat as `finished`. Score predictions on the awarded scoreline (typically 3-0 per football convention). |

`IsGroupComplete` counts cancelled fixtures as done:

```sql
WHERE status IN ('finished', 'cancelled')
```

`IsRoundComplete` requires all knockout fixtures in the round to have `winner_team_id IS NOT NULL`. A cancelled QF (no winner set) blocks semifinalist settlement indefinitely — see open question 8.

The existing `ListLockedFixturesByLeague` query in `fixtures.sql` filters `WHERE f.status IN ('in_progress', 'finished')`. This must be updated to also include `'cancelled'` so cancelled fixtures appear in the league predictions view with their final points.

---

## Integration with `main.go`

In `cmd/server/main.go`, `run()` is extended:

```go
workerRepo := repository.NewWorkerRepository(pool)
apiClient  := footballapi.NewClient(cfg.FootballAPIKey, cfg.FootballAPIBaseURL)
w          := worker.New(workerRepo, apiClient, worker.RealClock{}, logger)

workerErr := make(chan error, 1)
go func() {
    if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
        workerErr <- err
    }
    close(workerErr)
}()

select {
case <-ctx.Done():
    logger.Info("shutdown signal received")
case err := <-serverErr:
    return fmt.Errorf("http server: %w", err)
case err := <-workerErr:
    return fmt.Errorf("worker: %w", err)
}
// existing graceful HTTP shutdown follows unchanged
```

`worker.Run(ctx)` exits cleanly on context cancellation. Transient poll/DB errors are logged as `Warn` and the loop continues; `workerErr` only receives unrecoverable startup failures.

`internal/config/config.go` gains:

```go
FootballAPIKey     string `env:"FOOTBALL_API_KEY,required"`
FootballAPIBaseURL string `env:"FOOTBALL_API_BASE_URL" envDefault:"https://v3.football.api-sports.io"`
```

Add to `.env` for dev; add to the systemd unit file for prod. `server.NewRouter()` signature is unchanged.

---

## Touched Files

| File | Reason |
|------|--------|
| `migrations/20260528000001_worker_columns.sql` | New: is_demo+winner_team_id+last_polled_at on fixtures; external_id+season on tournaments; external_id on teams; scored_at on prediction tables; score_predictions(fixture_id) index; `cancelled` enum value |
| `migrations/20260528000002_group_table_stats.sql` | New: won/drawn/lost/goals_for/goals_against on tournament_group_table |
| `internal/domain/fixture.go` | Add `FixtureStatusCancelled FixtureStatus = "cancelled"` |
| `internal/domain/worker.go` | New: `PollableFixture`, `StandingsEntry` |
| `internal/worker/worker.go` | New: Worker, Run, tick, processSingleFixture |
| `internal/worker/repository.go` | New: WorkerRepo interface |
| `internal/worker/match_api.go` | New: MatchAPI interface + API result types |
| `internal/worker/clock.go` | New: Clock interface + RealClock |
| `internal/worker/semaphore.go` | New: bounded semaphore |
| `internal/worker/settlement.go` | New: runSettlement + round helpers |
| `internal/worker/worker_test.go` | New: unit tests with hand-written fakes |
| `internal/footballapi/client.go` | New: concrete HTTP MatchAPI client |
| `internal/footballapi/types.go` | New: JSON response structs |
| `internal/repository/worker.go` | New: WorkerRepository + non-transactional methods |
| `internal/repository/worker_settlement.go` | New: UpdateMatchAndRescoreLivePredictions + all Settle* (raw pgx) |
| `internal/repository/worker_test.go` | New: integration tests with testcontainers |
| `internal/repository/queries/worker.sql` | New sqlc queries; run `make sqlc-gen` after adding |
| `internal/repository/dbgen/` | Regenerated by `make sqlc-gen` — do not edit manually |
| `internal/config/config.go` | Add FootballAPIKey, FootballAPIBaseURL |
| `cmd/server/main.go` | Start worker goroutine; add workerErr to select |
| `internal/repository/queries/fixtures.sql` | Update `ListLockedFixturesByLeague` to include `'cancelled'` in status filter |

This plan adds no HTTP routes. No changes to `docs/openapi.yaml`, `docs/openapi.json`, `internal/server/oapi/models.gen.go`, or Postman collections.

---

## Edge Cases

| Case | Handling |
|------|----------|
| API returns same score/status as DB | No DB write; worker returns early from `processSingleFixture` |
| API call fails | Log Warn; skip fixture this tick; retry next tick |
| DB transaction fails mid-settlement | pgx rolls back; `WHERE points IS NULL` guard makes retry safe for outright predictions |
| `team_handicap` / `player_handicap` missing | `COALESCE(…, 0)` fallback; Warn logged; prediction awarded 0 pts |
| Two goroutines race to settle same fixture | `WHERE points IS NULL` is an optimistic lock; second transaction finds zero rows |
| `GetTeamByExternalID` returns not-found | Log Error; skip this standings entry; `teams.external_id` is not populated |
| Group-stage match cancelled | `IsGroupComplete` counts it; group Settle* functions run and award 0 to all score_predictions |
| All group matches cancelled | Group marked complete; all group-level predictions → 0 pts |
| Cancelled QF | `winner_team_id = NULL`; `IsRoundComplete = false`; semifinalist settlement blocked — see open question 8 |
| PEN match | `goals_home`/`goals_away` = aggregate (regulation + ET); `winner_team_id` from HomeWinner/AwayWinner |
| `PST` with no new kickoff | Leave kickoff unchanged; `status = upcoming`; re-polled on subsequent ticks |
| Standings API unavailable after group concludes | Log Warn; return early; next tick will re-attempt if `prevTerminal` guard removed (see settlement trigger note) |
| Context cancelled during tick | pgx propagates ctx; in-flight transactions roll back; goroutines exit cleanly |
| Fixture in 24h correction window re-polled | Score predictions ARE rescored (intentional). Outright predictions skip re-settlement (`points IS NULL` guard) |

---

## Test Plan

### Unit tests — `internal/worker/worker_test.go`

Hand-written fakes; no mockgen; table-driven.

| Test | Scenario |
|------|----------|
| `TestTick_NoFixtures` | Empty result → no API calls; idle sleep |
| `TestTick_LiveFixtures` | N fixtures → `GetFixture` called N times; semaphore limits to 5 concurrent |
| `TestMapAPIStatus` | Table: each API status string → correct domain status |
| `TestProcessSingleFixture_NoChange` | Same score/status → `UpdateMatchAndRescoreLivePredictions` not called |
| `TestProcessSingleFixture_ScoreChange` | Changed goals → `UpdateMatchAndRescoreLivePredictions` called |
| `TestProcessSingleFixture_APIError` | Error → logged; no DB call |
| `TestRunSettlement_GroupMatchGroupNotDone` | Only standings update; no Settle* calls |
| `TestRunSettlement_GroupMatchGroupDone` | SettleGroupWinner + SettlePlayoffGroup + top scorer + SettleGroupTopScorer all called |
| `TestRunSettlement_LastGroupDone` | SettlePlayoffWildcard called after `IsGroupStageComplete` returns true |
| `TestRunSettlement_QFNotAllDone` | SettleSemifinalist NOT called |
| `TestRunSettlement_AllQFsDone` | SettleSemifinalist called |
| `TestRunSettlement_FinalConcluded` | SettleTournamentWinner + top scorer fetch + SettleTopScorer called |
| `TestRunSettlement_Cancelled` | No outright Settle* calls |
| `TestRunSettlement_AlreadyTerminal` | No settlement (prevTerminal guard) |
| `TestRun_ContextCancelled` | Run returns nil |
| `TestRun_ListPollableMatchesError` | Error logged; loop continues |
| `TestRun_GetFixtureError` | One fixture fails; others succeed |

### Integration tests — `internal/repository/worker_test.go`

testcontainers-go Postgres; full migrations applied before each test.

| Test | Scenario |
|------|----------|
| `ListPollableMatches` — live fixture | Returned |
| `ListPollableMatches` — kickoff in 4 minutes | Returned |
| `ListPollableMatches` — kickoff in 60 minutes | Not returned |
| `ListPollableMatches` — finished within 24h | Returned |
| `ListPollableMatches` — finished 25h ago | Not returned |
| `ListPollableMatches` — is_demo=TRUE | Not returned |
| `ListPollableMatches` — tournament external_id IS NULL | Not returned |
| `UpdateMatchAndRescoreLivePredictions` — exact score | points = 3 |
| `UpdateMatchAndRescoreLivePredictions` — correct result | points = 1 |
| `UpdateMatchAndRescoreLivePredictions` — wrong result | points = 0 |
| `UpdateMatchAndRescoreLivePredictions` — cancelled | points = 0 |
| `UpdateMatchAndRescoreLivePredictions` — rescored (live 1-0 then 2-0) | points updated on second call |
| `IsGroupComplete` — all finished | true |
| `IsGroupComplete` — one in_progress | false |
| `IsGroupComplete` — all cancelled | true |
| `IsRoundComplete` — all QFs with winner_team_id | true |
| `IsRoundComplete` — one QF cancelled, winner_team_id NULL | false |
| `IsGroupStageComplete` — all groups concluded | true |
| `IsGroupStageComplete` — one group has upcoming fixture | false |
| `SettleGroupWinnerPredictions` — correct slot 0 pick | points = handicap value |
| `SettleGroupWinnerPredictions` — wrong position (slot 0 picks pos2 team) | points = 0 |
| `SettleGroupWinnerPredictions` — idempotent | Second call is no-op |
| `SettlePlayoffGroupPredictions` — pick advanced | points = handicap value |
| `SettlePlayoffGroupPredictions` — pick eliminated | points = 0 |
| `SettlePlayoffWildcardPredictions` — wildcard advanced | points = handicap value |
| `SettlePlayoffWildcardPredictions` — wildcard eliminated | points = 0 |
| `SettleGroupTopScorerPredictions` — correct pick | points = handicap value |
| `SettleGroupTopScorerPredictions` — wrong pick | points = 0 |
| `SettleGroupTopScorerPredictions` — missing handicap row | points = 0 (COALESCE) |
| `SettleSemifinalistPredictions` — correct pick | points = handicap value |
| `SettleSemifinalistPredictions` — non-semi team | points = 0 |
| `SettleTournamentWinnerPredictions` — correct pick | points = handicap value |
| `SettleTournamentWinnerPredictions` — wrong pick | points = 0 |
| `SettleTopScorerPredictions` — correct pick | points = handicap value |
| `SettleTopScorerPredictions` — wrong pick | points = 0 |
| `UpdateGroupStandings` — upsert 4 teams | All rows updated including W/D/L/GF/GA |
| `UpdateGroupStandings` — called twice | Second call updates positions correctly |
| `GetTeamByExternalID` — known external ID | Returns correct UUID |
| `GetPlayerByExternalID` — known external ID | Returns correct UUID |

---

## Acceptance Criteria

All require a real (non-demo) tournament with `external_id IS NOT NULL`.

1. **Live score update**: After one worker tick where the API returns a new score for a live match, `SELECT points FROM score_predictions WHERE fixture_id = $X` shows 3, 1, or 0 for each prediction row correctly.

2. **Group standings updated**: After a group-stage match is processed, `SELECT position, won, drawn, lost, goals_for, goals_against FROM tournament_group_table WHERE group_letter = 'A'` reflects the API-sourced standings.

3. **Group winner settled**: After the last group-A match concludes, all `team_predictions WHERE category = 'group_winner' AND group_letter = 'A'` rows have `points IS NOT NULL`. `GET /tournaments/{id}/outcomes` returns two `group_winner` team outcome rows for group A.

4. **Group top scorer settled**: After group A concludes, all `player_predictions WHERE category = 'group_top_scorer' AND group_letter = 'A'` rows have `points IS NOT NULL`. `GET /tournaments/{id}/outcomes` includes the `group_top_scorer` player outcome.

5. **Playoff group picks settled**: After group A concludes, all `team_predictions WHERE category = 'playoff' AND group_letter = 'A' AND slot_index IN (0, 1)` rows have `points IS NOT NULL`.

6. **Wildcard playoff settled**: After all groups conclude, all `team_predictions WHERE category = 'playoff' AND group_letter IS NULL AND slot_index = 2` rows have `points IS NOT NULL`.

7. **Semifinalist settled**: After all 4 QFs conclude, `SELECT COUNT(*) FROM team_outcomes WHERE category = 'semifinalist' AND tournament_id = $X` = 4. All `team_predictions WHERE category = 'semifinalist'` rows have `points IS NOT NULL`.

8. **Tournament winner settled**: After the final, `team_outcomes WHERE category = 'winner'` has 1 row. All `team_predictions WHERE category = 'winner'` rows have `points IS NOT NULL`.

9. **Top scorer settled**: After the final, `player_outcomes WHERE category = 'total_top_scorer'` has 1 row. All `player_predictions WHERE category = 'total_top_scorer'` rows have `points IS NOT NULL`.

10. **Leaderboard reflects all points**: `GET /tournaments/{id}/leaderboard` shows correct `total_points` summing all categories for each user.

11. **Cancelled match**: `SELECT points FROM score_predictions WHERE fixture_id = $X` all = 0. `SELECT status FROM fixtures WHERE id = $X` = `'cancelled'`.

12. **Idle worker**: With no live or imminent fixtures, the worker logs `"worker idle"` and does not call the football API.

---

## Open Questions

1. **Match scoring formula**: Plan proposes exact score = 3 pts, correct result = 1 pt, wrong = 0 pts. Confirm this is the intended formula. Is there any partial credit for one correct side of the score?

2. **`group_winner` slot semantics**: Plan assumes `slot_index = 0` = "team finishes 1st" and `slot_index = 1` = "team finishes 2nd" — points awarded only on exact slot-position match. Confirm, or clarify if any top-2 pick in either slot earns points.

3. **`playoff` slot semantics**: Plan assumes any pick that advanced (position 1 or 2) earns points regardless of which slot (0 or 1) it was placed in — slots evaluated independently. Confirm.

4. **Tie-breaking for group standings**: When two teams are tied, the football API's `position` field is treated as authoritative. Does the API break ties deterministically using goal difference, then goals scored, then head-to-head? If the API's position-1 team ever differs from what was displayed to users during the group stage, settlement will surprise users. Confirm the API is trusted as final arbiter.

5. **Tied top scorers**: If multiple players have the same goal tally, plan proposes awarding points to ALL tied players (predictions for any of them are correct). Confirm this is the intended rule, or specify a single tiebreaker (assists, alphabetical, etc.).

6. **Source of per-group player goal data**: Does the football API's `/players/topscorers` endpoint support filtering by group stage? If not, the worker must fetch all player stats per fixture and aggregate manually — a significantly more complex flow requiring a `GetFixturePlayerStats` method on `MatchAPI`. Clarify before implementing `GetGroupTopScorer`.

7. **Exact round string values**: `isQuarterfinalRound` and `isFinalRound` rely on string matching against `fixtures.round`. What exact strings does api-sports.io return for World Cup 2026? (e.g., `"Quarter-finals"`, `"Semi-finals"`, `"Final"`). The implementer must verify against the API documentation or real data before hardcoding. Consider making round identifiers configurable per tournament.

8. **Cancelled knockout fixture**: A cancelled QF leaves `winner_team_id = NULL`, permanently blocking semifinalist settlement. Decision needed: (a) admin endpoint to manually set `winner_team_id` for a cancelled fixture; (b) define a fallback rule (zero all semifinalist predictions when any QF is cancelled); or (c) treat the whole tournament's outright predictions as voided.

9. **PEN match `goals.home`/`goals.away` behaviour**: At API status `PEN`, do `goals.home`/`goals.away` reflect only the regulation + ET score (before penalty kicks), or do they count penalty goals too? Confirm against api-sports.io documentation before implementing the scoring formula for PEN matches.

10. **`teams.external_id` data population**: Adding `external_id BIGINT NULL` to teams is required for standings resolution. How will existing teams receive their external IDs? Via a one-time seed script? Via the fixture import process? This is a prerequisite for `GetTeamByExternalID` to work in production.

11. **`tournaments.external_id` and `season` population**: Same question for tournaments. The `POST /tournaments` endpoint does not accept these fields. Confirm the workflow for linking an existing tournament to the external API (seed script, new admin endpoint, or manual SQL update).

12. **Football API rate limits**: api-sports.io free tier is typically 100 requests/day. With 60-second polling ticks and 80+ matches in a World Cup group stage + knockout, the free tier is exhausted in under an hour. Confirm a paid subscription tier is in place before production deployment, and decide whether 429 rate-limit handling should be implemented in `internal/footballapi/client.go` at this stage.

13. **Settlement retry after transient failure**: The current settlement trigger uses `prevTerminal` to skip re-triggering. If standings or top-scorer API calls fail after a group concludes, this guard prevents retry. The implementer must resolve this by either: (a) removing the `prevTerminal` guard and relying solely on `WHERE points IS NULL` idempotency; (b) adding a `group_settlement_attempted_at` column; or (c) tracking which settlement steps have completed. Option (a) is simplest but means the standings and top-scorer API are called on every tick for 24 hours after conclusion.

14. **`is_demo` scope**: Plan adds `is_demo` only to `fixtures`. Should there also be an `is_demo` flag on `tournaments` to disable an entire tournament's polling with a single toggle?
