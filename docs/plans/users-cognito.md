# Plan: Users + AWS Cognito Auth

**File:** `docs/plans/users-cognito.md`
**Date:** 2026-05-20
**Status:** Draft

---

## 1. Goal

Wire up AWS Cognito (social login only, RS256 JWTs) as the identity layer, automatically provisioning a local `users` row on first authenticated request, and expose a `GET /users/me` endpoint that returns the caller's profile.

---

## 2. Out of Scope

- Native username/password auth in Cognito — social-only is the explicit constraint.
- `footy-forecast-web` App Client — the second client ID slot is reserved; no backend code is needed beyond adding the ID to the `COGNITO_ALLOWED_CLIENT_IDS` env var.
- Cognito User Pool / App Client creation and Google + Facebook IdP wiring are covered in **Phase 0** (Section 9) and must be completed before codebase phases begin.
- User profile update endpoint (e.g., `PATCH /users/me` to change display name).
- Roles, permissions, or any authorization beyond "authenticated vs unauthenticated".
- Token refresh or logout flows — the frontend owns the Cognito Hosted UI flow; the backend only validates tokens.
- Admin endpoints for managing users (suspend, delete).
- Rate limiting or brute-force protection on the auth middleware.

---

## 3. Data Model Changes

### 3a. New Postgres enum and table

Migration file: `migrations/20260520000000_create_users.sql` (use the actual `goose` timestamp at implementation time).

```sql
-- +goose Up
-- +goose StatementBegin

CREATE TYPE user_status AS ENUM ('active', 'suspended');

CREATE TABLE users (
    id            UUID        PRIMARY KEY,
    cognito_sub   TEXT        NOT NULL UNIQUE,
    email         TEXT        NOT NULL UNIQUE,
    display_name  TEXT        NOT NULL DEFAULT '',
    status        user_status NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX users_email_idx ON users (email);

-- Reuse set_updated_at() created in the tournaments migration.
CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS users_set_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_status;

-- NOTE: set_updated_at() is NOT dropped here — tournaments still depends on it.

-- +goose StatementEnd
```

Key decisions:
- `id` is a UUIDv7 generated in Go (not `gen_random_uuid()`), consistent with the tournaments pattern.
- `cognito_sub` is the stable Cognito identifier — never changes even if the user's email changes.
- `email` has a `UNIQUE` constraint. If two Cognito identities somehow share an email (Cognito misconfiguration), the upsert will surface a `domain.ErrConflict` 500 — see Edge Cases.
- `display_name` defaults to `''` (empty string, not NULL) to avoid nullable handling throughout.
- The Down migration does not drop `set_updated_at()` because tournaments still uses it.

### 3b. sqlc queries

New file: `internal/repository/queries/user.sql`

```sql
-- name: UpsertUser :one
INSERT INTO users (id, cognito_sub, email, display_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cognito_sub) DO UPDATE
    SET email        = EXCLUDED.email,
        display_name = CASE
                           WHEN users.display_name = '' THEN EXCLUDED.display_name
                           ELSE users.display_name
                       END,
        updated_at   = now()
RETURNING id, cognito_sub, email, display_name, status, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, cognito_sub, email, display_name, status, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByCognitoSub :one
SELECT id, cognito_sub, email, display_name, status, created_at, updated_at
FROM users
WHERE cognito_sub = $1;
```

Notes:
- `UpsertUser` always refreshes `email` on conflict (stays in sync if user changes email in Cognito).
- `display_name` is only overwritten if the stored value is empty — preserves future user-supplied values.
- `updated_at` is set explicitly in the `DO UPDATE` clause because the trigger fires only on `UPDATE` statements, and `ON CONFLICT DO UPDATE` counts as an update in terms of data but the trigger behavior should be verified.
- After adding this file, run `sqlc generate` from the repo root.

### 3c. Config additions

In `internal/config/config.go`, add to the `Config` struct:

```go
CognitoRegion           string   `env:"COGNITO_REGION,required"`
CognitoUserPoolID       string   `env:"COGNITO_USER_POOL_ID,required"`
CognitoAllowedClientIDs []string `env:"COGNITO_ALLOWED_CLIENT_IDS,required" envSeparator:","`
```

> **Note for implementer:** Verify the `envSeparator` tag key against `caarlos0/env/v11` release notes before writing — the tag syntax may differ from v10.

Add a validation guard in `Load()`:
```go
if len(cfg.CognitoAllowedClientIDs) == 0 {
    return nil, fmt.Errorf("COGNITO_ALLOWED_CLIENT_IDS must contain at least one value")
}
```

Derived values (computed, not stored in Config):
- JWKS URL: `https://cognito-idp.{CognitoRegion}.amazonaws.com/{CognitoUserPoolID}/.well-known/jwks.json`
- Token issuer: `https://cognito-idp.{CognitoRegion}.amazonaws.com/{CognitoUserPoolID}`

Update `.env.example` with the three new vars (empty values, with a comment).

---

## 4. API Surface

### Existing endpoints (unchanged)

All `/tournaments/*` and `/health/*` routes remain public. No `Authorization` header required.

### New endpoint

#### `GET /users/me`

Protected by `AuthMiddleware`. If the middleware passes, the user is guaranteed to be in the request context.

**Request**

```
GET /users/me HTTP/1.1
Authorization: Bearer <cognito-id-token>
```

**Response — 200 OK**

```json
{
  "id":           "0197c3e1-1234-7abc-9def-000000000001",
  "email":        "user@example.com",
  "display_name": "Alice",
  "status":       "active",
  "created_at":   "2026-05-20T10:00:00Z",
  "updated_at":   "2026-05-20T10:00:00Z"
}
```

`cognito_sub` is intentionally omitted — it is an internal implementation detail.

**Error responses (from middleware, before handler runs)**

| Condition | Status | Body |
|---|---|---|
| Missing `Authorization` header | 401 | `{"error":"unauthorized"}` |
| Malformed / expired / wrong-issuer / wrong-aud token | 401 | `{"error":"unauthorized"}` |
| User provision fails (DB error) | 500 | `{"error":"internal server error"}` |

The 401 response never exposes the specific rejection reason to avoid leaking validation details.

---

## 5. New Package Layout

```
internal/
  domain/
    user.go                         -- User type, UserStatus constants
  repository/
    user.go                         -- UserRepository interface + impl
    queries/user.sql                -- sqlc queries
  service/
    user.go                         -- UserService, ProvisionFromCognito
  server/
    cognito/
      jwks.go                       -- JWKSCache: fetch + cache with TTL, re-fetch on unknown kid
      validator.go                  -- JWTValidator interface, Claims struct, CognitoValidator
    ctxutil/
      user.go                       -- WithUser / UserFromCtx typed context helpers
    middleware/
      auth.go                       -- AuthMiddleware(validator, userService)
    handler/
      user.go                       -- GET /users/me
```

---

## 6. Key Interfaces

### `JWTValidator`

```go
// internal/server/cognito/validator.go

type Claims struct {
    Sub       string
    Email     string
    Name      string
    GivenName string
}

type JWTValidator interface {
    Validate(ctx context.Context, tokenString string) (Claims, error)
}
```

`CognitoValidator` implements `JWTValidator`. It:
1. Parses the JWT header to extract `kid`.
2. Looks up the key in the JWKS cache.
3. If `kid` is unknown, re-fetches JWKS once (handles key rotation). Fails if still absent.
4. Verifies the RS256 signature with `golang-jwt/jwt/v5`.
5. Validates `exp`, `iss` (must equal pool URL), and `aud` (must be in `allowedClientIDs`).
6. Returns populated `Claims`.

The JWKS cache uses a mutex-protected map and a TTL (1 hour). Re-fetch on unknown `kid` is done under a write lock to prevent thundering herd during key rotation.

Use `WithLeeway(30 * time.Second)` in the jwt parser to tolerate minor clock drift.

### `UserService.ProvisionFromCognito`

```go
// internal/service/user.go

func (s *UserService) ProvisionFromCognito(ctx context.Context, cognitoSub, email, rawDisplayName string) (domain.User, error)
```

Display name derivation (in order):
1. `strings.TrimSpace(rawDisplayName)` — use if non-empty.
2. Email prefix: `strings.Cut(email, "@")` → use the part before `@`.
3. Fallback: use full email string if no `@` present.
4. Truncate to 50 Unicode code points: `string([]rune(s)[:min(50, len([]rune(s)))])`.

Calls `repo.UpsertUser(ctx, id, cognitoSub, email, displayName)` with a freshly generated UUIDv7.

### `AuthMiddleware`

```go
// internal/server/middleware/auth.go

func Auth(validator cognito.JWTValidator, users UserProvisioner) func(http.Handler) http.Handler
```

Where `UserProvisioner` is a narrow interface:
```go
type UserProvisioner interface {
    ProvisionFromCognito(ctx context.Context, sub, email, displayName string) (domain.User, error)
}
```

Flow:
1. Extract `Authorization: Bearer <token>` — 401 if absent or malformed.
2. `validator.Validate(ctx, token)` → `Claims` — 401 on error.
3. `users.ProvisionFromCognito(ctx, claims.Sub, claims.Email, claims.Name || claims.GivenName)` → 500 on error.
4. `ctxutil.WithUser(r.Context(), user)` → call next handler.

The middleware writes JSON error responses directly (inline `writeJSON`) without importing the `handler` package, to avoid a circular import (`handler` → `middleware` would be circular if `middleware` imports `handler`).

### `ctxutil.UserFromCtx`

```go
// internal/server/ctxutil/user.go

func WithUser(ctx context.Context, u domain.User) context.Context
func UserFromCtx(ctx context.Context) (domain.User, bool)
```

Returns `(domain.User, bool)` — the `false` case is a programming error (handler called without middleware). The handler checks the bool and returns 500 with a logged error. This satisfies the "no panic outside main" rule.

---

## 7. Edge Cases

**Token expired mid-flight.** `exp` is validated by the JWT library. Middleware returns 401. No DB writes occur.

**Cognito key rotation (unknown `kid`).** Validator re-fetches JWKS once under a write lock. If the `kid` is still absent after re-fetch, returns error → 401. No infinite retry.

**Concurrent first-provision.** Two simultaneous requests for a new user both reach `UpsertUser`. `ON CONFLICT (cognito_sub) DO UPDATE` is atomic. Both succeed and return the same `id`.

**Email uniqueness collision.** If two Cognito identities somehow share an email (should be impossible), the upsert conflicts on the `email` UNIQUE constraint instead of `cognito_sub`. The repo sees `pgerrcode.UniqueViolation` and returns `domain.ErrConflict` → handler returns 500. This is a Cognito misconfiguration symptom; log it clearly.

**Suspended user.** `GET /users/me` returns the profile with `"status":"suspended"`. The middleware does not check status — that is future work once admin tooling exists.

**JWKS endpoint unreachable.** The JWKS is fetched lazily on first use. If Cognito is unreachable, `Validate` fails → middleware returns 500 (cannot evaluate the token). Log clearly with `slog.Error`.

**Display name truncation at multibyte boundary.** Use `[]rune` slicing, not `[]byte` slicing, to avoid splitting a multibyte UTF-8 character.

---

## 8. Open Questions

1. **`UserFromCtx` false case status code.** The plan says return 500, but consider whether a logged panic is cleaner since it's definitively a wiring bug. Decision needed before implementing the handler. This plan recommends the `(User, bool)` + 500 path to respect the "no panic outside main" rule.

2. **JWKS fetch failure status code.** When Cognito's JWKS endpoint is unreachable, should the middleware return 401 ("cannot validate token") or 503 ("upstream dependency down")? RFC 6750 suggests 401 for any token validation failure. This plan defaults to 500 (unexpected error) but the implementer should confirm.

3. **`caarlos0/env/v11` slice separator tag syntax.** Verify against v11 release notes before writing config code. The `envSeparator` tag may have changed from v10.

4. **`golang-jwt/jwt/v5` vs. `jwt/v4`.** `jwt/v4` may already be an indirect dependency. Adding `jwt/v5` is a separate module path and is valid. This plan defaults to `v5` (current release, better API for JWKS key functions). If you prefer avoiding two jwt majors, `v4` works with minor API differences.

5. **`NewRouter` signature.** Currently `NewRouter(logger, pool)`. Should it accept `*config.Config` directly, or a purpose-built `AuthConfig` struct? Passing `*config.Config` is simpler and consistent with how `main.go` already uses it. Implementer decides.

---

## 9. Phase 0: Cognito Infrastructure Setup

Run these steps before starting Phase 1. By the end you will have a real Pool ID, two real Client IDs, SSM params in place, and a way to mint a test token locally.

All `aws` commands use `--profile hexa` (the project's AWS profile) and `--region eu-central-1`. Capture output values into shell variables as you go — you'll need them in later steps.

---

### Step 0.1 — Create the User Pool

```bash
POOL_ID=$(aws cognito-idp create-user-pool \
  --pool-name footy-forecast \
  --region eu-central-1 \
  --admin-create-user-config 'AllowAdminCreateUserOnly=true' \
  --auto-verified-attributes email \
  --policies 'PasswordPolicy={MinimumLength=8,RequireUppercase=false,RequireLowercase=false,RequireNumbers=false,RequireSymbols=false,TemporaryPasswordValidityDays=1}' \
  --mfa-configuration OFF \
  --deletion-protection ACTIVE \
  --profile hexa \
  --query 'UserPool.Id' \
  --output text)

echo "Pool ID: $POOL_ID"
```

`AllowAdminCreateUserOnly=true` blocks direct API sign-ups while leaving social IdP federation fully functional — Cognito creates the user record automatically on first social sign-in. `DELETION_PROTECTION ACTIVE` prevents an accidental `delete-user-pool` from destroying prod data.

---

### Step 0.2 — Create the Hosted UI Domain

The domain prefix must be globally unique across all Cognito pools. If `footy-forecast` is taken, try `footy-forecast-prod` or append a random suffix.

```bash
COGNITO_DOMAIN=footy-forecast   # change if taken

aws cognito-idp create-user-pool-domain \
  --user-pool-id "$POOL_ID" \
  --domain "$COGNITO_DOMAIN" \
  --region eu-central-1 \
  --profile hexa

# Full base URL (used in IdP redirect URIs below):
echo "Hosted UI base: https://${COGNITO_DOMAIN}.auth.eu-central-1.amazoncognito.com"
```

---

### Step 0.3 — Configure Google Identity Provider

**First, in Google Cloud Console (out-of-band):**
1. Create a project (or reuse an existing one).
2. Navigate to **APIs & Services → Credentials → Create Credentials → OAuth 2.0 Client ID**.
3. Application type: **Web application**.
4. Add authorized redirect URI: `https://<COGNITO_DOMAIN>.auth.eu-central-1.amazoncognito.com/oauth2/idpresponse`
5. Note the **Client ID** and **Client Secret**.

**Then in AWS:**

```bash
GOOGLE_CLIENT_ID=<paste-google-client-id>
GOOGLE_CLIENT_SECRET=<paste-google-client-secret>

aws cognito-idp create-identity-provider \
  --user-pool-id "$POOL_ID" \
  --provider-name Google \
  --provider-type Google \
  --provider-details "{\"client_id\":\"${GOOGLE_CLIENT_ID}\",\"client_secret\":\"${GOOGLE_CLIENT_SECRET}\",\"authorize_scopes\":\"openid email profile\"}" \
  --attribute-mapping '{"email":"email","name":"name","given_name":"given_name"}' \
  --region eu-central-1 \
  --profile hexa
```

The attribute mapping keys are Cognito standard attributes; values are the Google claim names. `openid email profile` scopes are required to receive `email` and `name` in the ID token.

---

### Step 0.4 — Configure Facebook Identity Provider

> **Critical:** The Facebook app **must explicitly include the `email` permission**. Without it, Cognito receives no email from Facebook and JIT provisioning will fail with a missing-email error on every Facebook sign-in.

**First, in Facebook Developers (out-of-band):**
1. Go to [developers.facebook.com](https://developers.facebook.com) → **My Apps → Create App**.
2. App type: **Consumer**.
3. Add the **Facebook Login** product.
4. Under Facebook Login → **Settings**, add Valid OAuth Redirect URI: `https://<COGNITO_DOMAIN>.auth.eu-central-1.amazoncognito.com/oauth2/idpresponse`
5. Under **App Settings → Basic**, note the **App ID** and **App Secret**.
6. Confirm that `email` is listed under **Default Access** permissions (it should be; verify explicitly).

**Then in AWS:**

```bash
FB_APP_ID=<paste-facebook-app-id>
FB_APP_SECRET=<paste-facebook-app-secret>

aws cognito-idp create-identity-provider \
  --user-pool-id "$POOL_ID" \
  --provider-name Facebook \
  --provider-type Facebook \
  --provider-details "{\"client_id\":\"${FB_APP_ID}\",\"client_secret\":\"${FB_APP_SECRET}\",\"authorize_scopes\":\"public_profile,email\"}" \
  --attribute-mapping '{"email":"email","name":"name","given_name":"first_name"}' \
  --region eu-central-1 \
  --profile hexa
```

Facebook's Graph API returns `first_name` (not `given_name`), so the mapping uses `"given_name":"first_name"`.

---

### Step 0.5 — Create App Clients

```bash
# Mobile client (current)
MOBILE_CLIENT_ID=$(aws cognito-idp create-user-pool-client \
  --user-pool-id "$POOL_ID" \
  --client-name footy-forecast-mobile \
  --no-generate-secret \
  --allowed-o-auth-flows code \
  --allowed-o-auth-flows-user-pool-client \
  --allowed-o-auth-scopes openid email profile \
  --callback-urls '["footy-forecast://callback"]' \
  --logout-urls '["footy-forecast://signout"]' \
  --supported-identity-providers '["Google","Facebook"]' \
  --explicit-auth-flows '["ALLOW_REFRESH_TOKEN_AUTH"]' \
  --prevent-user-existence-errors ENABLED \
  --enable-token-revocation \
  --region eu-central-1 \
  --profile hexa \
  --query 'UserPoolClient.ClientId' \
  --output text)

echo "Mobile Client ID: $MOBILE_CLIENT_ID"

# Web client (placeholder — reserves the slot; no web app yet)
WEB_CLIENT_ID=$(aws cognito-idp create-user-pool-client \
  --user-pool-id "$POOL_ID" \
  --client-name footy-forecast-web \
  --no-generate-secret \
  --allowed-o-auth-flows code \
  --allowed-o-auth-flows-user-pool-client \
  --allowed-o-auth-scopes openid email profile \
  --callback-urls '["http://localhost:3000/callback"]' \
  --logout-urls '["http://localhost:3000"]' \
  --supported-identity-providers '["Google","Facebook"]' \
  --explicit-auth-flows '["ALLOW_REFRESH_TOKEN_AUTH"]' \
  --prevent-user-existence-errors ENABLED \
  --enable-token-revocation \
  --region eu-central-1 \
  --profile hexa \
  --query 'UserPoolClient.ClientId' \
  --output text)

echo "Web Client ID: $WEB_CLIENT_ID"
```

Both clients are public (no secret) with PKCE implicit in the authorization code flow. `ALLOW_REFRESH_TOKEN_AUTH` is the only explicit auth flow — no native password flows.

---

### Step 0.6 — Store SSM Parameters

```bash
aws ssm put-parameter \
  --name /footy-forecast/prod/cognito-region \
  --value eu-central-1 \
  --type String \
  --region eu-central-1 \
  --profile hexa

aws ssm put-parameter \
  --name /footy-forecast/prod/cognito-user-pool-id \
  --value "$POOL_ID" \
  --type String \
  --region eu-central-1 \
  --profile hexa

aws ssm put-parameter \
  --name /footy-forecast/prod/cognito-allowed-client-ids \
  --value "${MOBILE_CLIENT_ID},${WEB_CLIENT_ID}" \
  --type String \
  --region eu-central-1 \
  --profile hexa
```

---

### Step 0.7 — Update Local `.env`

```bash
# .env  (never commit)
COGNITO_REGION=eu-central-1
COGNITO_USER_POOL_ID=<POOL_ID>
COGNITO_ALLOWED_CLIENT_IDS=<MOBILE_CLIENT_ID>,<WEB_CLIENT_ID>
```

---

### Step 0.8 — Update EC2 Deploy Script

The EC2 deploy script (`/usr/local/bin/footy-forecast-deploy`) already pulls SSM params for the DB password. Add the same pattern for the three new Cognito params and export them so the systemd unit picks them up. Exact edits depend on the script's current structure — add a pull for each `/footy-forecast/prod/cognito-*` param alongside the existing DB password pull.

---

### Step 0.9 — Verify: Mint a Test Token

Use the web client (which has a `localhost` redirect) to get a real ID token without needing a mobile app.

1. Open this URL in a browser (substituting your values):

   ```
   https://<COGNITO_DOMAIN>.auth.eu-central-1.amazoncognito.com/login
     ?response_type=code
     &client_id=<WEB_CLIENT_ID>
     &scope=openid+email+profile
     &redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback
   ```

2. Sign in with Google or Facebook. The browser will redirect to `http://localhost:3000/callback?code=<AUTH_CODE>`. Copy the `code` value from the URL bar (the redirect will fail — that's expected, no server is running there).

3. Exchange the code for tokens:

   ```bash
   curl -s -X POST \
     "https://${COGNITO_DOMAIN}.auth.eu-central-1.amazoncognito.com/oauth2/token" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=authorization_code" \
     -d "client_id=${WEB_CLIENT_ID}" \
     -d "code=<AUTH_CODE>" \
     -d "redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback" \
     | jq .
   ```

4. Copy `id_token` from the response. Use it as the Bearer token for smoke-testing `GET /users/me` once the codebase phase is done:

   ```bash
   curl -H "Authorization: Bearer <id_token>" http://localhost:<PORT>/users/me
   ```

Paste the `id_token` into [jwt.io](https://jwt.io) to verify the claims (`sub`, `email`, `name`, `iss`, `aud`) look correct before writing a line of Go.

---

---

## 10. Implementation Phases

Complete Phase 0 (Section 9) first. You need a real Pool ID and Client IDs before Phase 1 touches config.

### Phase 1 — Foundation

1. Update `internal/config/config.go` with Cognito fields.
2. Update `.env.example`.
3. Add migration (`goose create users sql`).
4. Add `internal/domain/user.go`.
5. Add `internal/repository/queries/user.sql`.
6. Run `sqlc generate`.
7. Add `internal/repository/user.go`.
8. Add `internal/service/user.go`.
9. `make fmt && make lint && make test`.

### Phase 2 — Auth Middleware

1. `go get github.com/golang-jwt/jwt/v5`.
2. Add `internal/server/cognito/jwks.go`.
3. Add `internal/server/cognito/validator.go`.
4. Add `internal/server/ctxutil/user.go`.
5. Add `internal/server/middleware/auth.go`.
6. `make fmt && make lint && make test`.

### Phase 3 — Endpoint + Wiring

1. Add `internal/server/handler/user.go`.
2. Update `internal/server/router.go`.
3. Update `cmd/server/main.go`.
4. `make fmt && make lint && make test`.

---

## 11. Acceptance Criteria

1. `make fmt && make lint && make test` all pass with no suppressions.
2. `GET /users/me` with a valid Cognito ID token returns 200 with `id`, `email`, `display_name`, `status`, `created_at`.
3. First request from a new Cognito user creates a `users` row. Second identical request returns the same `id`.
4. After an email change in Cognito, the next authenticated request updates the `email` column.
5. `GET /tournaments` returns 200 without an `Authorization` header.
6. A token whose `aud` is not in `COGNITO_ALLOWED_CLIENT_IDS` is rejected with 401.
7. An expired token is rejected with 401.
8. `cognito_sub` is not present in any JSON response body.

---

## 12. Test Plan

### `domain/user_test.go`
- `UserStatus` valid values: `active` → true, `suspended` → true, `""` → false.
- Display name derivation: blank `name`+`given_name` → email prefix; long name → truncated to 50 runes; Unicode name truncated at code-point boundary; email with no `@` → whole email used.

### `repository/user_test.go` (integration, testcontainers)
- `UpsertUser` creates a new row on first call.
- Second upsert with same `cognito_sub` returns same `id`.
- Second upsert refreshes `email`.
- Second upsert does NOT overwrite a non-empty `display_name`.
- Second upsert DOES set `display_name` when stored value is empty.
- `GetByID` returns the user for a known ID; returns `domain.ErrNotFound` for unknown.
- `GetByCognitoSub` returns the user for a known sub; returns `domain.ErrNotFound` for unknown.
- `updated_at` advances after a second upsert.

### `service/user_test.go`
- `ProvisionFromCognito` calls upsert with correct arguments.
- Display name priority: non-blank `rawDisplayName` wins; blank → email prefix.
- Truncation applied at 50 runes.
- Repo errors are wrapped and propagated.

### `cognito/validator_test.go`
- Valid token, correct issuer, correct aud → returns `Claims`.
- Expired token → error.
- Wrong signing key → error.
- Wrong issuer → error.
- `aud` not in allowed list → error.
- Unknown `kid` triggers one JWKS re-fetch (mock HTTP server via `httptest.NewServer`).
- Unknown `kid` still absent after re-fetch → error, no further retries.

### `middleware/auth_test.go`
- Missing `Authorization` header → 401.
- Header not in `Bearer <token>` format → 401.
- Validator returns error → 401.
- Provision returns error → 500.
- Valid token + provision succeeds → next handler called; `ctxutil.UserFromCtx` returns the user.

### `handler/user_test.go`
- User in context → 200, JSON contains `id`, `email`, `display_name`, `status`, `created_at`; does NOT contain `cognito_sub`.
- User not in context → 500.
