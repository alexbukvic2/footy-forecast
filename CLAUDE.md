# CLAUDE.md — Conventions for AI agents working on footy-forecast

You are working on **footy-forecast**, a football (soccer) prediction game.
Read this file fully before making any change. When in doubt, ask before
acting; do not silently invent conventions.

## What this app is

A prediction pool / pick'em game (NOT fantasy football). Users predict:
- Individual match results (exact scores)
- Group-stage standings
- Knockout bracket outcomes
- Tournament-level questions (top scorer, eventual winner)

Players score points per scoring rules. They can join private "leagues" and
compete on per-league and global leaderboards.

Predictions lock at kickoff. Scoring is deterministic and will be updated in near real-time as matches progress (e.g. points for a correct match winner are awarded as soon as the result is known, even if the exact score isn't final yet).

## Stack
- Go 1.26
- pgx/v5 + sqlc for DB access (no ORM)
- goose for migrations
- Postgres 16 on AWS RDS (free tier: db.t4g.micro)
- AWS Cognito for auth (added in a later phase; not yet wired up)
- Caddy for TLS termination on EC2
- systemd to run the binary
- GitHub Actions for CI/CD
- Single EC2 t4g.micro (Graviton/ARM64) in eu-central-1

## Architecture
Layered, not hexagonal — keep it simple.

**Strict rules** (enforce in review; refuse to violate):

1. **`handler` may call `service` only.** Never call `repository` directly.
2. **`service` may call `repository` and other `service`s.** Never import `http` types.
3. **`repository` returns `domain` types.** Never returns `*sql.Rows` or raw DB rows to callers.
4. **`domain` imports nothing from this project.** Pure types and errors only.
5. **No global state.** All dependencies passed via constructors.
6. **No `init()` functions** except for trivial registration (none allowed for now).
7. **Context is the first arg** of every function that does I/O or might block.
8. **Errors are wrapped with `fmt.Errorf("...: %w", err)`** at every layer boundary.
9. **No `panic` outside `main`.** Return errors.
10. **Domain errors are sentinel values** in `internal/domain/errors.go`
    (e.g. `domain.ErrNotFound`), checked with `errors.Is`.

## Error handling

- Repositories return `domain.ErrNotFound` for missing rows, not `sql.ErrNoRows`.
- Services translate repository errors into domain errors when crossing boundaries.
- Handlers translate domain errors into HTTP status codes in one place
  (see `internal/http/handler/errors.go` once it exists).
- Never expose internal error messages to clients in production.
- Logs include `request_id` from context for every error.

## Logging

- Use `slog` with the default logger configured in `main`.
- JSON output. Levels: `Info` for normal events, `Warn` for unexpected-but-handled,
  `Error` for failures requiring attention.
- Always include `request_id` for request-scoped logs.
- Never log secrets, passwords, JWTs, or full request bodies.

## Testing

- Unit tests sit next to code: `foo.go` → `foo_test.go`.
- Use the standard `testing` package + `testify/assert` only if it materially helps
  (avoid for trivial assertions).
- Table-driven tests preferred.
- Repository tests run against a real Postgres (testcontainers-go) — never mock
  the DB driver.
- Service tests use repository interfaces with hand-written fakes (not mockgen).
- Aim for meaningful coverage on `service/` and `repository/`. Handlers should
  have at least one happy-path and one error-path test.

## Commands you must run before declaring a task done

```bash
make fmt    # gofmt + go mod tidy
make lint   # golangci-lint
make test   # tests with race detector
```

All three must pass. Do not skip `make lint` even if "the change is small."

## Working with planning artifacts

- All non-trivial features get a plan first, written to `docs/plans/<slug>.md`
  by the planner. Plans include: goal, data model changes, API endpoints,
  edge cases, test plan, acceptance criteria.
- Implementer reads the plan and implements. If the plan is ambiguous or wrong,
  stop and ask — don't paper over it.
- Reviewer reads the plan + diff and produces a review file. Blocking comments
  mean the plan iterates, not just the code.
- Architectural decisions go in `docs/decisions/NNNN-title.md` (ADR-style).

## Things to never do

- Never commit secrets, even encrypted, even "temporarily."
- Never introduce a new dependency without justifying it in a comment or ADR.
  The Go stdlib is usually enough; prefer it.
- Never silently swallow an error. Either handle it or wrap-and-return.
- Never write code that depends on time without taking a `Clock` interface or
  using a passed-in `time.Now` source — makes testing miserable.
- Never use `context.Background()` inside request handlers — always derive from
  `r.Context()`.

## When you're unsure

Ask. A clarifying question is cheaper than a wrong implementation.
