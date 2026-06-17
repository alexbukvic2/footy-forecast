# Plan: Fixture AI Analysis via AWS Bedrock

## Goal

After a fixture finishes, fire an asynchronous background job (non-blocking to `runSettlement`) that:

1. Fetches the fixture result and all player score predictions for that fixture.
2. Calls AWS Bedrock (Claude Sonnet 4.6) with a prompt asking for 2–3 interesting stats or fun facts.
3. Upserts the generated text into a new `score_ai_analysis` table.
4. Surfaces the analysis on `GET /leagues/{leagueId}/predictions` as a nullable field on each fixture.

---

## Data Model

### New table: `score_ai_analysis`

```sql
CREATE TABLE score_ai_analysis (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    fixture_id UUID        NOT NULL REFERENCES fixtures(id),
    analysis   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (fixture_id)
);
```

- One row per fixture (enforced by `UNIQUE(fixture_id)`).
- `analysis` stores the raw text returned by the LLM (markdown bullets are fine).
- Upsert on `fixture_id` so a score correction triggers a fresh analysis.

---

## API Contract

### `GET /leagues/{leagueId}/predictions`

`LeagueFixtureViewResponse` gains one new nullable field:

```yaml
ai_analysis:
  type: string
  nullable: true
  description: >
    AI-generated match trivia (2–3 bullet points). Null if the analysis
    job has not completed yet or the fixture is not finished.
```

No other endpoints change. The response is `null` when the row is absent or the fixture is not in a finished state.

---

## New configuration

Add two fields to `internal/config/config.go`:

| Env var | Required | Default | Notes |
|---|---|---|---|
| `BEDROCK_REGION` | no | `eu-central-1` | AWS region for Bedrock API calls |
| `BEDROCK_MODEL_ID` | yes | — | Full inference profile ID; verify in AWS console (cross-region profile for Claude Sonnet 4.6 in eu-central-1 follows the pattern `eu.anthropic.claude-sonnet-4-6-<date>-v1:0`) |

`BEDROCK_MODEL_ID` is marked `required` so the app fails fast on startup if it is not set.

---

## New dependency

Add to `go.mod`:

```
github.com/aws/aws-sdk-go-v2/config
github.com/aws/aws-sdk-go-v2/service/bedrockruntime
```

No other third-party dependency is needed. AWS credentials come from the EC2 instance profile (default credential chain; no key pair required).

---

## IAM requirement (deployment)

The EC2 instance role needs one additional statement:

```json
{
  "Effect": "Allow",
  "Action": "bedrock:InvokeModel",
  "Resource": "arn:aws:bedrock:*::foundation-model/anthropic.claude-sonnet-4-6*"
}
```

This must be applied in the AWS console / Terraform before deploying. The app will start fine without it but all analysis jobs will log an error and produce no rows.

---

## Implementation steps

### 1. Migration

File: `migrations/<timestamp>_score_ai_analysis.sql`

Create `score_ai_analysis` with the schema above. No seed data.

### 2. New sqlc query file: `internal/repository/queries/ai_analysis.sql`

Two queries:

**`GetFixtureAnalysisInput :one`** — fetches everything the LLM prompt needs in a single round-trip. JOINs `fixtures` → `teams` (for home/away names and `group_letter`) and LEFT JOINs `score_predictions` → `users` (for each player's prediction). Aggregates the predictions into a JSON text column using `json_agg`, following the existing pattern in `fixtures.sql`.

**`UpsertFixtureAnalysis :exec`** — inserts a new row or updates `analysis` and `updated_at` on conflict with `fixture_id`.

Run `make generate` after adding these.

### 3. New domain types: `internal/domain/ai_analysis.go`

```go
// FixtureAnalysisInput is the data the AI job reads before calling Bedrock.
type FixtureAnalysisInput struct {
    HomeTeamName string
    AwayTeamName string
    Round        string
    GroupLetter  *string
    GoalsHome    *int
    GoalsAway    *int
    Predictions  []AnalysisPrediction
}

// AnalysisPrediction is one player's score prediction, used in the AI prompt.
type AnalysisPrediction struct {
    DisplayName string
    GoalsHome   *int
    GoalsAway   *int
    Points      *int
}
```

### 4. New repository file: `internal/repository/ai_analysis.go`

`AIAnalysisRepository` with:
- `GetFixtureAnalysisInput(ctx, fixtureID uuid.UUID) (domain.FixtureAnalysisInput, error)`
  - Calls `dbgen.GetFixtureAnalysisInput`, unmarshals the JSON predictions column (mirrors `memberPredictionJSON` pattern in `fixture.go`).
  - Returns `domain.ErrNotFound` if the fixture does not exist.
- `UpsertFixtureAnalysis(ctx, fixtureID uuid.UUID, analysis string) error`
  - Calls `dbgen.UpsertFixtureAnalysis`.

### 5. New Bedrock client: `internal/bedrock/client.go`

Defines an interface and a real implementation:

```go
// Analyser generates AI analysis text given a prompt.
type Analyser interface {
    Analyse(ctx context.Context, prompt string) (string, error)
}
```

The real implementation (`BedrockAnalyser`):
- Loads AWS config using `config.LoadDefaultConfig` (picks up EC2 instance profile automatically).
- Uses `bedrockruntime.InvokeModel` with the Anthropic Messages API format:
  ```json
  {
    "anthropic_version": "bedrock-2023-05-31",
    "max_tokens": 512,
    "messages": [{"role": "user", "content": "<prompt>"}]
  }
  ```
- Parses `content[0].text` from the response JSON.
- Returns a plain `error` on non-2xx or unmarshal failure; **does not retry internally** (retry is handled in the caller).

Constructor: `NewBedrockAnalyser(region, modelID string) (*BedrockAnalyser, error)`

### 6. Update worker

**`internal/worker/repository.go`** — add two methods to the `Repo` interface:
```go
GetFixtureAnalysisInput(ctx context.Context, fixtureID uuid.UUID) (domain.FixtureAnalysisInput, error)
UpsertFixtureAnalysis(ctx context.Context, fixtureID uuid.UUID, analysis string) error
```

**`internal/worker/worker.go`** — add `analyser bedrock.Analyser` field to `Worker` struct and update `New(...)` constructor signature.

**New file: `internal/worker/ai_analysis.go`**

```go
// runAIAnalysis fires in a goroutine. It does not block the caller.
func (w *Worker) runAIAnalysis(fixtureID uuid.UUID) {
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    input, err := w.repo.GetFixtureAnalysisInput(ctx, fixtureID)
    // ... handle error

    prompt := buildAnalysisPrompt(input)

    var analysis string
    backoff := []time.Duration{2 * time.Second, 4 * time.Second}
    for attempt := 0; attempt <= 2; attempt++ {
        if attempt > 0 {
            select {
            case <-ctx.Done():
                // log, return
            case <-time.After(backoff[attempt-1]):
            }
        }
        analysis, err = w.analyser.Analyse(ctx, prompt)
        if err == nil {
            break
        }
        // log warn with attempt number
    }
    if err != nil {
        // log error
        return
    }

    if err := w.repo.UpsertFixtureAnalysis(ctx, fixtureID, analysis); err != nil {
        // log error
    }
}
```

`buildAnalysisPrompt(input domain.FixtureAnalysisInput) string` — a pure function in the same file. Formats the fixture result and each player prediction into English prose, then asks:

> "You are a football prediction game assistant. Generate 2–3 short, interesting bullet points about this match and the players' predictions. Draw on both match facts (e.g., surprising scoreline, lopsided result) and prediction statistics (e.g., how many players predicted the correct result, whether nobody or only one person got it right, whether the majority agreed on a scoreline that turned out wrong, who was closest). Keep each bullet to one short sentence. Focus on stats and numbers, not narrative. Do not use em dashes."

**`internal/worker/worker.go` — `processSingleFixture`** — after `w.runSettlement(...)`, add:

```go
if newStatus == domain.FixtureStatusFinished {
    go w.runAIAnalysis(f.ID)
}
```

The goroutine is intentionally detached; the 60 s context timeout prevents it from leaking indefinitely.

### 7. Update `LeagueFixtureView` domain type

`internal/domain/fixture.go` — add one field:

```go
type LeagueFixtureView struct {
    Fixture     Fixture
    Predictions []LeagueMemberPrediction
    AIAnalysis  *string  // nil when analysis is not yet available
}
```

### 8. Update fixture SQL queries

`internal/repository/queries/fixtures.sql` — update both `ListLockedFixturesByLeague` and `ListLockedFixturesByLeagueAndDates` to LEFT JOIN `score_ai_analysis` and select `saa.analysis AS ai_analysis`:

```sql
LEFT JOIN score_ai_analysis saa ON saa.fixture_id = f.id
```

Add `saa.analysis` (nullable) to the SELECT list and `saa.analysis` to the GROUP BY (or use `MAX(saa.analysis)` since there is at most one row per fixture — but simply adding it to GROUP BY is cleaner since `fixture_id` is unique in that table).

Run `make generate` after editing.

Update `internal/repository/fixture.go` — `ListLockedByLeague` and `ListLockedByLeagueAndDates` now read `row.AiAnalysis` (nullable `pgtype.Text`) and set `view.AIAnalysis` accordingly.

### 9. Update the handler

`internal/server/handler/score_prediction.go`:

```go
type leagueFixtureViewResponse struct {
    fixtureResponse
    Predictions []leagueMemberScorePrediction `json:"predictions"`
    AIAnalysis  *string                        `json:"ai_analysis"`
}
```

Update `toLeagueFixtureViewResponse` to copy `v.AIAnalysis` into the response.

### 10. Update OpenAPI spec

`docs/openapi.yaml`:

1. Add `ai_analysis` to `LeagueFixtureViewResponse` schema (nullable string, not required).
2. Update the example to show a non-null value for a finished fixture.

Run `make generate` after editing.

### 11. Update router wiring

`internal/server/router.go`:

- Construct `bedrock.NewBedrockAnalyser(cfg.BedrockRegion, cfg.BedrockModelID)` (returns error — fail fast in main if Bedrock is misconfigured).
- Pass the analyser as a new argument to `worker.New(...)`.
- Update the worker's `Repo` implementation (the concrete repository struct passed to `worker.New`) to also implement the two new `Repo` interface methods; this means wiring in `AIAnalysisRepository`.

---

## Edge cases

| Scenario | Behaviour |
|---|---|
| Fixture is cancelled (not finished) | `runAIAnalysis` is never called; `ai_analysis` stays `null`. |
| Score corrected after analysis exists | `runAIAnalysis` fires again; `UpsertFixtureAnalysis` updates the row. |
| Bedrock call fails all 3 attempts | Error is logged; no row written; `ai_analysis` stays `null` on API. |
| `GetFixtureAnalysisInput` returns no predictions | Prompt is built with "no predictions were submitted"; LLM comments on the scoreline alone. |
| Analysis job still in progress when client polls | Field is `null`; client shows nothing and can re-poll. |
| `BEDROCK_MODEL_ID` not set | App exits at startup with a config error (marked `required`). |
| IAM permission not granted | `Analyse` returns a 403 error from Bedrock; logged as error; no row. |

---

## Test plan

- **`internal/bedrock/client_test.go`** — unit test `buildAnalysisPrompt` (pure function).
- **`internal/worker/ai_analysis_test.go`** — table-driven test for `runAIAnalysis` using a hand-written fake `Analyser` (returns canned text or error). Cover: happy path, Bedrock failure all attempts, repo error on upsert.
- **`internal/repository/ai_analysis_test.go`** — integration test against testcontainers Postgres: `UpsertFixtureAnalysis` (insert + update), `GetFixtureAnalysisInput` (with and without predictions).
- **Handler test** — add an assertion to the existing `ListForLeague` happy-path test that `ai_analysis` is marshalled correctly (both null and non-null).

---

## Acceptance criteria

1. Finishing a fixture triggers a background Bedrock call; the caller (`processSingleFixture`) returns without waiting.
2. A row appears in `score_ai_analysis` for the finished fixture within ~30 s in a real environment.
3. Score correction: the `updated_at` timestamp and `analysis` text change if the job re-runs.
4. `GET /leagues/{leagueId}/predictions` returns `"ai_analysis": "..."` for finished fixtures with a stored analysis and `"ai_analysis": null` for fixtures without one.
5. `make fmt && make lint && make test` all pass.
6. `make generate` produces in-sync `openapi.json` and `models.gen.go`.
7. App fails to start (config error) if `BEDROCK_MODEL_ID` is not set in the environment.
