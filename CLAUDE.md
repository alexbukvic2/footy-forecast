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
- oapi-codegen (build-time only) for generating Go types from the OpenAPI spec
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
  (see `internal/server/handler/errors.go` once it exists).
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

If you changed `docs/openapi.yaml`, also run:

```bash
make generate   # regenerates docs/openapi.json and internal/server/oapi/models.gen.go
```

Commit both regenerated files alongside the spec change.

## OpenAPI

The canonical API contract lives at `docs/openapi.yaml` (OpenAPI 3.1).
**This is the single source of truth for all HTTP routes.** If there is ever a
conflict between the spec and Postman, the spec wins.

Spec-first workflow:

1. Update `docs/openapi.yaml` **before** writing any handler code.
2. Run `make generate` to regenerate `docs/openapi.json` and `internal/server/oapi/models.gen.go`.
3. Implement the handler using the generated types.
4. Never edit `models.gen.go` or `openapi.json` directly — both are overwritten on every `make generate`.

Rules:
- **New route** → add the path and operation to the spec; run `make generate`.
- **Route removed** → delete the spec entry; run `make generate`.
- **Request/response shape changed** → update the spec first; run `make generate`.
- All error responses must `$ref` the shared `ErrorResponse` schema.
- Public routes must carry `security: []`; protected routes `security: [{bearerAuth: []}]`.
- All request/response schemas must have explicit types — no `{}` or `any`.

Special public routes (no auth):
- `GET /openapi.json` — serves `docs/openapi.json` for tooling and mobile SDK generation.
- `GET /docs` — serves Swagger UI (CDN-loaded) pointed at `/openapi.json`.

Generated files (`docs/openapi.json`, `internal/server/oapi/models.gen.go`) are committed to the repository (same convention as `dbgen/`). CI enforces they are in sync with the spec via a drift check.

## Postman

The Postman collection lives in `postman/collections/footy-forecast/`.
Postman is the **test runner**; the OpenAPI spec (above) is the contract.
Any change that adds, removes, or alters an HTTP route **must** include a
corresponding Postman update in the same commit:

- **New route** → add a `<Name>.request.yaml` in the matching subfolder (or create
  the subfolder if it is a new resource).
- **Route removed** → delete the corresponding `.request.yaml`.
- **URL, method, or request body changed** → update the `.request.yaml` in place.

Authenticated routes must include `Authorization: Bearer {{authToken}}` in the
headers. Use `{{baseUrl}}` for the host and existing env vars (`{{tournamentId}}`,
`{{leagueId}}`, etc.) for path parameters. Add a post-response script to capture
newly returned IDs into env vars when that ID will be needed by downstream requests.

## Working with planning artifacts

- All non-trivial features get a plan first, written to `docs/plans/<slug>.md`
  by the planner. Plans include: goal, data model changes, API contract (for any
  route change), edge cases, test plan, acceptance criteria.
- **Plans that add, modify, or remove HTTP routes must include an "API Contract"
  section** per the template in `.claude/agents/planner.md`. The implementer
  translates this section into `docs/openapi.yaml` before writing handler code.
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

## CI quirks

- `golangci-lint-action` uses `install-mode: goinstall` because the prebuilt
  binaries (as of golangci-lint v2.7.2) are built with Go 1.25 and refuse to
  analyze code targeting Go 1.26+. Building from source uses the local Go
  toolchain instead. Revisit once golangci-lint ships 1.26-built binaries.

## Deployment

Deploys to production happen automatically on push to `main`:

1. `.github/workflows/deploy.yml` triggers on push to main
2. It calls `.github/workflows/ci.yml` as a sub-workflow to verify the build
3. After CI passes, it builds an ARM64 Linux binary
4. Assumes the `footy-forecast-deploy` IAM role via OIDC
  - Trust policy keyed on `environment:production`, not branch ref
  - Required because the deploy job declares `environment: production`
5. Uploads binary + migrations to `s3://footy-forecast-deploy-.../`
6. Triggers `/usr/local/bin/footy-forecast-deploy` on the EC2 box via SSM Run Command
  - The script downloads, validates, migrates, installs atomically, restarts
  - Waits for `/health/ready` to confirm before exiting

Manual deploy from laptop (debugging only):

```bash
# Build + upload as if you were the workflow
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o /tmp/server ./cmd/server
aws s3 cp /tmp/server "s3://${BUCKET}/builds/server-$(git rev-parse --short HEAD)" --profile hexa
aws s3 sync ./migrations/ "s3://${BUCKET}/migrations/" --profile hexa

# Trigger the deploy script
aws ssm send-command --instance-ids "$INSTANCE_ID" \
  --document-name "AWS-RunShellScript" \
  --parameters "commands=[\"bash /usr/local/bin/footy-forecast-deploy builds/server-$(git rev-parse --short HEAD)\"]" \
  --profile hexa
```

## Backup strategy

- Postgres `pg_dump` every 6 hours, gzipped, uploaded to s3://footy-forecast-backups/daily/
- 90-day retention via S3 lifecycle policy
- Worst-case data loss: 6 hours

**Point-in-time recovery (PITR) is intentionally deferred.** When user base grows
or data becomes harder to recreate (real-money transactions, irrecoverable user
data), set up WAL-G or pgBackRest for continuous WAL archiving to S3.

Restore: see docs/runbooks/restore.md (to be written).

# On the box:
aws s3 cp "s3://footy-forecast-backups/daily/<date>/<key>" /tmp/restore.sql.gz --region eu-central-1

# Stop the app so it doesn't see partial state mid-restore
sudo systemctl stop footy-forecast

# Restore (clean + if-exists means it drops existing tables first)
gunzip -c /tmp/restore.sql.gz | PGPASSWORD="$(aws ssm get-parameter --name /footy-forecast/prod/postgres-password --with-decryption --region eu-central-1 --query 'Parameter.Value' --output text)" psql -h 127.0.0.1 -U footy_app -d footy_forecast

# Restart the app
sudo systemctl start footy-forecast

