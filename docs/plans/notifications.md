# Plan: Push Notifications

**File:** `docs/plans/notifications.md`
**Date:** 2026-06-07
**Status:** Draft

---

## 1. Goal

Add push notification support for mobile clients (Expo). Two initial types:

- **`matchday`** — sent at noon in the user's local timezone on any day that has fixtures, but only if the user has at least one unpredicted fixture in the next 24 hours.
- **`pre_match`** — sent 2.5 hours before kickoff for each fixture the user has not yet predicted.
- **`tournament_reminder`** — sent once per tournament, 5 hours before its earliest fixture, to users who have not yet submitted any tournament-level predictions (group standings, bracket, outright questions).

Users can enable/disable each type independently and configure a daily silent window during which no notifications are sent.

---

## 2. Out of Scope

- Email or web-push notifications.
- Per-league notification preferences.
- Live/in-match notifications (goal scored, result confirmed).
- Notification inbox / read-status tracking.
- Notification history visible in the app.
- Delivery receipts or open-rate analytics.

---

## 3. Push Provider

**Use Expo Push Notifications (EAS)** — not raw FCM.

Rationale: the mobile client is Expo, so every device already has an Expo push token (`ExponentPushToken[...]`). Expo's server-side API (`https://exp.host/--/expo-push-notification/send`) routes to FCM (Android) and APNs (iOS) automatically. This means:

- One API call + one token format on the server, regardless of platform.
- No APNs certificate management on our side.
- FCM credentials are configured in the Expo dashboard, not in our backend.

No additional Go dependencies are needed — we talk to Expo over plain HTTPS with `encoding/json`.

---

## 4. Data Model Changes

### 4a. `users` — add `timezone`, `silent_from`, `silent_until`

```sql
ALTER TABLE users ADD COLUMN timezone    TEXT NOT NULL DEFAULT 'UTC';
ALTER TABLE users ADD COLUMN silent_from TIME;
ALTER TABLE users ADD COLUMN silent_until TIME;
```

`timezone` is an IANA timezone string (e.g. `"Europe/London"`). The mobile app passes the device timezone on sign-in and whenever it changes. Defaults to `UTC`.

`silent_from` / `silent_until` are a single global quiet window applied to **all** notification types. Both must be non-null together or both null. Evaluated in the user's local timezone.

### 4b. New table: `push_tokens`

```sql
CREATE TABLE push_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (token)
);
CREATE INDEX push_tokens_user_id_idx ON push_tokens (user_id);
```

One row per device. A user may have multiple tokens (phone + tablet). `token` is unique globally because each Expo token maps to exactly one device; if two users somehow share a device, the second registration replaces the first.

On registration: upsert by token, updating `user_id` if it changed (device re-used after sign-out).

### 4c. New table: `notification_preferences`

```sql
CREATE TABLE notification_preferences (
    user_id  UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type     TEXT    NOT NULL,  -- 'matchday' | 'pre_match'
    enabled  BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (user_id, type)
);
```

A missing row means "default enabled" — equivalent to `enabled = TRUE`. The `GET` endpoint materialises defaults for types that have no row yet.

### 4d. New table: `notification_log`

```sql
CREATE TABLE notification_log (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         TEXT        NOT NULL,
    reference_id TEXT        NOT NULL,  -- date string for matchday; fixture_id for pre_match
    sent_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, type, reference_id)
);
CREATE INDEX notification_log_sent_at_idx ON notification_log (sent_at);
```

`reference_id` for `matchday` is the calendar date string `"2026-06-07"` (in the user's timezone). For `pre_match` it is the fixture UUID. The unique constraint prevents duplicate sends. Rows can be pruned after 30 days (add a cleanup job later, or a partial index with a trigger).

---

## 5. API Contract

### 5a. Push token registration

#### `POST /users/me/push-tokens`

Registers a device token for the authenticated user.

**Request:**
```json
{ "token": "ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]" }
```

**Response `201 Created`:**
```json
{ "id": "uuid" }
```

Upserts: if the token already exists for this user, returns 200 with the existing record. If the token exists for a different user (device re-used), reassigns it.

#### `DELETE /users/me/push-tokens/{token}`

Deregisters a specific device token. Call on logout or when `expo-notifications` fires a `DeviceNotRegistered` event on the client.

**Response:** `204 No Content`

Returns `204` even if the token did not exist (idempotent).

---

### 5b. Notification preferences

#### `GET /users/me/notification-preferences`

Returns all supported notification types with their current preferences. Missing DB rows are materialised with defaults.

**Response `200 OK`:**
```json
[
  { "type": "matchday",  "enabled": true },
  { "type": "pre_match", "enabled": false }
]
```

#### `PUT /users/me/notification-preferences/{type}`

Upserts preferences for a single notification type. `type` must be `matchday` or `pre_match`.

**Request:**
```json
{ "enabled": true }
```

**Response `200 OK`:** same shape as one element of the GET array.

---

### 5c. User timezone and silent window (extend existing `PATCH /users/me`)

Add three optional fields to the existing `PATCH /users/me` endpoint:

**Request additions:**
```json
{
  "timezone": "Europe/London",
  "silent_from": "22:00",
  "silent_until": "08:00"
}
```

Validation:
- `timezone` must be a valid IANA name (use `time.LoadLocation`).
- `silent_from` and `silent_until` must both be present or both absent.
- Both are `HH:MM` 24-hour strings.
- `silent_from == silent_until` is rejected (zero-length window).
- To clear the silent window, send `"silent_from": null, "silent_until": null`.

---

## 6. Notification Workers

The existing `internal/worker` package already establishes the pattern: a struct with a `Run(ctx context.Context) error` method that blocks on a ticker loop, launched as a goroutine from `cmd/server/main.go`. Notification workers follow the exact same pattern — they are **not** added to the existing `worker.Worker` because they have a different tick interval and different dependencies (Expo client instead of `MatchAPI`). Three new structs, each in `internal/notification/job/`, each started alongside the existing worker in `main.go`.

### 6a. Pre-match notifier

**Cadence:** every 5 minutes.

**Logic:**
1. Find fixtures where `kickoff_at BETWEEN now() + 2h25m AND now() + 2h35m` — a 10-minute window centred on 2.5h before kickoff.
2. For each fixture, run a single query that returns all eligible users:
   - Has at least one push token.
   - Has `pre_match` enabled (or no preference row).
   - Has NOT submitted a `score_prediction` for this fixture.
   - Has NOT been sent a `pre_match` notification for this fixture (no row in `notification_log`).
   - Is NOT currently in their silent window.
3. Group tokens by user, batch into groups of 100, call Expo.
4. Insert a `notification_log` row for each successfully queued user.
5. On Expo `DeviceNotRegistered` response per token: delete that token from `push_tokens`.

**Silent window evaluation (SQL):**

```sql
-- user is in silent window if:
-- both columns set AND (
--   wraps midnight: NOW_LOCAL >= silent_from OR NOW_LOCAL < silent_until
--   same day:       NOW_LOCAL >= silent_from AND NOW_LOCAL < silent_until
-- )
```

Implement as a SQL helper function or inline WHERE clause using `now() AT TIME ZONE u.timezone`.

### 6b. Tournament reminder notifier

**Cadence:** every 5 minutes.

**Logic:**
1. Find tournaments where the earliest fixture kickoff falls in the window `now() + 4h55m` to `now() + 5h05m`. Use the existing `tournaments.predictions_locked` flag or a direct MIN query on `fixtures.kickoff_at` per tournament.
2. For each such tournament, find users who:
   - Have at least one push token.
   - Have NOT explicitly disabled `tournament_reminder` (missing preference row = enabled).
   - Are NOT currently in their silent window.
   - Have NOT already been sent a `tournament_reminder` for this tournament (no row in `notification_log` with `reference_id = tournament_id`).
3. Send notification, insert `notification_log` rows.

`reference_id` in `notification_log` is the tournament UUID. The unique constraint ensures this fires at most once per user per tournament even if the job runs multiple times inside the 10-minute window.

### 6c. Matchday notifier

**Cadence:** every 5 minutes.

**Logic:**
1. Find users where it is currently between 12:00 and 12:05 in their local timezone:
   ```sql
   SELECT id, timezone FROM users
   WHERE EXTRACT(HOUR  FROM (now() AT TIME ZONE timezone)) = 12
     AND EXTRACT(MINUTE FROM (now() AT TIME ZONE timezone)) < 5;
   ```
2. For each such user, check:
   - Has at least one push token.
   - Has `matchday` enabled (or no preference row).
   - Has NOT been sent a `matchday` notification for today's date in their timezone (no row in `notification_log` with `reference_id = date_in_user_tz`).
   - Has at least one fixture in the next 24 hours with no `score_prediction` from this user.
   - Is NOT in their silent window (noon is the send time, so this is less likely to conflict, but still enforce it).
3. Send notification, insert log.

**Why 5-minute cadence + 5-minute window matches:** the job runs at t=0,5,10,… minutes. The window covers 12:00–12:05, so any job execution within that window fires once for a given user. The `notification_log` unique constraint ensures idempotency if the job somehow runs twice.

---

## 7. Expo Push Client

New package: `internal/notification/expo/`

```
expo/
  client.go      — HTTP client, Send(), batch splitting
  types.go       — request/response structs
  client_test.go — unit tests with mock HTTP server
```

`Send(ctx context.Context, messages []Message) ([]Receipt, error)` where `Message` is:

```go
type Message struct {
    To    string // ExponentPushToken[...]
    Title string
    Body  string
    Data  map[string]string // optional extra payload
}
```

Batch splitting: max 100 per request. On Expo error responses, surface `DeviceNotRegistered` errors per token so callers can clean up. Do not retry on `DeviceNotRegistered`; do retry (once) on transient HTTP errors.

---

## 8. Notification Copy

| Type                   | Title                          | Body                                                                                              |
|------------------------|--------------------------------|---------------------------------------------------------------------------------------------------|
| `matchday`             | "Predictions open today"       | "You have unpredicted matches today - get your picks in!"                                         |
| `pre_match`            | "Kick-off in 2.5 hours"        | "{{HomeTeam}} vs {{AwayTeam}} - predict before it's too late!"                                    |
| `tournament_reminder`  | "Tournament predictions close soon" | "{{TournamentName}} kicks off in 5 hours - fill in your tournament predictions before they lock!" |

Copy can be refined later. Store it as constants in a `messages.go` file inside `internal/notification/`.

---

## 9. Domain Types

`internal/domain/notification.go`:

```go
type NotificationType string

const (
    NotificationTypeMatchday  NotificationType = "matchday"
    NotificationTypePreMatch  NotificationType = "pre_match"
)

// AllNotificationTypes lists every supported type, used to materialise defaults.
var AllNotificationTypes = []NotificationType{NotificationTypeMatchday, NotificationTypePreMatch}

type PushToken struct {
    ID        string
    UserID    string
    Token     string
    CreatedAt time.Time
}

type NotificationPreference struct {
    UserID  string
    Type    NotificationType
    Enabled bool
}

// TimeOfDay represents an HH:MM wall-clock time with no date component.
// Used for users.silent_from / users.silent_until.
type TimeOfDay struct {
    Hour   int
    Minute int
}
```

---

## 10. Silent Window Logic

```
given: silent_from F, silent_until U, current local time T

if F == U:   no window (or reject at validation)
if F < U:    in window if T >= F && T < U          (same day, e.g. 13:00–14:00)
if F > U:    in window if T >= F || T < U          (wraps midnight, e.g. 22:00–08:00)
```

Implement as a pure Go function `IsInSilentWindow(t TimeOfDay, from, until TimeOfDay) bool` — easy to unit test.

---

## 11. Package Layout

```
internal/
  notification/
    expo/
      client.go
      client_test.go
      types.go
    job/
      prematch.go       — PreMatchNotifier struct, Run(ctx)
      prematch_test.go
      matchday.go       — MatchdayNotifier struct, Run(ctx)
      matchday_test.go
    messages.go         — notification copy constants
    silent_window.go    — IsInSilentWindow()
    silent_window_test.go
  repository/
    push_token.go
    notification_preference.go
    notification_log.go
    (+ _test.go files)
  service/
    notification.go     — NotificationService (preferences + token registration)
    notification_test.go
  domain/
    notification.go
  server/handler/
    notification.go     — HTTP handlers for tokens + preferences
    notification_test.go
```

---

## 12. Edge Cases

| Scenario | Handling |
|----------|----------|
| User has no timezone set | Defaults to `UTC` |
| Fixture kickoff changes after notification sent | No re-notification; acceptable |
| User predicts after receiving pre-match notification | Fine — no action needed |
| Expo `DeviceNotRegistered` | Delete token from `push_tokens` |
| User uninstalls app without logging out | Token cleaned up on first `DeviceNotRegistered` response |
| Multiple devices | Send to all active tokens for the user |
| No fixtures today | Matchday notification not sent |
| All fixtures already predicted | Matchday notification not sent |
| Job crashes mid-send | Partially logged; on restart, already-logged users skipped |
| Silent window `null` | Always notify (if type is enabled) |
| Silent window wraps midnight | Handled by the `F > U` branch in `IsInSilentWindow` |

---

## 13. Test Plan

- **`IsInSilentWindow`**: table-driven unit tests covering same-day, midnight-wrapping, and boundary cases.
- **Expo client**: unit tests with an `httptest.Server` mock covering batching, `DeviceNotRegistered` handling, and transient error retry.
- **Repository tests** (testcontainers): upsert token, reassign token, list by user, delete; upsert preference, get preferences with defaults; insert log, dedup constraint.
- **Service tests** (hand-written fakes): register token, update preferences with invalid timezone, get preferences materialises defaults.
- **Job tests**: inject fake repos and fake Expo client; assert correct users are targeted, log rows inserted, and `DeviceNotRegistered` tokens removed.
- **Handler tests**: happy path + error path for each of the 4 new endpoints.

---

## 14. Acceptance Criteria

1. Mobile app can register a push token; a second call with the same token is idempotent.
2. Mobile app can deregister a token on logout.
3. User can enable/disable each notification type via API.
4. User can set a silent window; invalid windows (one field missing, `from == until`) are rejected with 422.
5. `PATCH /users/me` accepts and validates a `timezone` field.
6. Pre-match: users without a prediction for a fixture receive a notification 2.5h before kickoff (±5 min); users who already predicted do not.
7. Matchday: users receive a notification at noon local time on days with unpredicted fixtures; users who predicted everything do not.
8. No user receives the same notification twice (enforced by `notification_log` unique constraint).
9. A `DeviceNotRegistered` Expo response causes the offending token to be deleted.
10. `make fmt && make lint && make test` all pass.
