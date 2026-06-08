package job

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/notification"
	"github.com/alexbukvic2/footy-forecast/internal/notification/expo"
)

// MatchdayNotifier sends a daily notification at noon (user local time) on days
// that have at least one unpredicted fixture.
type MatchdayNotifier struct {
	pool   *db.Pool
	expo   *expo.Client
	logger *slog.Logger
}

// NewMatchdayNotifier constructs a MatchdayNotifier.
func NewMatchdayNotifier(pool *db.Pool, expoClient *expo.Client, logger *slog.Logger) *MatchdayNotifier {
	return &MatchdayNotifier{pool: pool, expo: expoClient, logger: logger}
}

// Run blocks on a 5-minute ticker until ctx is cancelled. It always returns nil.
func (j *MatchdayNotifier) Run(ctx context.Context) error {
	tick := time.NewTicker(5 * time.Minute)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			j.run(ctx)
		}
	}
}

type matchdayTarget struct {
	userID  uuid.UUID
	token   string
	dateKey string
}

func (j *MatchdayNotifier) run(ctx context.Context) {
	const q = `
SELECT
    u.id                                                      AS user_id,
    pt.token                                                  AS token,
    to_char(now() AT TIME ZONE u.timezone, 'YYYY-MM-DD')    AS date_key
FROM users u
JOIN push_tokens pt ON pt.user_id = u.id
WHERE
    -- it is currently noon (12:00–12:05) in the user's local timezone
    EXTRACT(HOUR   FROM (now() AT TIME ZONE u.timezone)) = 12
    AND EXTRACT(MINUTE FROM (now() AT TIME ZONE u.timezone)) < 5
    -- matchday notifications are enabled (missing row = enabled)
    AND NOT EXISTS (
        SELECT 1 FROM notification_preferences np
        WHERE np.user_id = u.id AND np.type = 'matchday' AND np.enabled = FALSE
    )
    -- not already sent today (keyed on local date string)
    AND NOT EXISTS (
        SELECT 1 FROM notification_log nl
        WHERE nl.user_id = u.id
          AND nl.type = 'matchday'
          AND nl.reference_id = to_char(now() AT TIME ZONE u.timezone, 'YYYY-MM-DD')
    )
    -- user has at least one unpredicted fixture in the next 24 hours
    AND EXISTS (
        SELECT 1 FROM fixtures f
        WHERE f.kickoff_at BETWEEN now() AND now() + INTERVAL '24 hours'
          AND f.prediction_locked = FALSE
          AND NOT EXISTS (
              SELECT 1 FROM score_predictions sp
              WHERE sp.fixture_id = f.id AND sp.user_id = u.id
          )
    )
    -- not in silent window
    AND NOT (
        u.silent_from IS NOT NULL
        AND u.silent_until IS NOT NULL
        AND CASE
            WHEN (EXTRACT(HOUR FROM u.silent_from)*60 + EXTRACT(MINUTE FROM u.silent_from))
                 < (EXTRACT(HOUR FROM u.silent_until)*60 + EXTRACT(MINUTE FROM u.silent_until))
            THEN
                (EXTRACT(HOUR FROM (now() AT TIME ZONE u.timezone))*60 + EXTRACT(MINUTE FROM (now() AT TIME ZONE u.timezone)))
                    >= (EXTRACT(HOUR FROM u.silent_from)*60 + EXTRACT(MINUTE FROM u.silent_from))
                AND
                (EXTRACT(HOUR FROM (now() AT TIME ZONE u.timezone))*60 + EXTRACT(MINUTE FROM (now() AT TIME ZONE u.timezone)))
                    < (EXTRACT(HOUR FROM u.silent_until)*60 + EXTRACT(MINUTE FROM u.silent_until))
            ELSE
                (EXTRACT(HOUR FROM (now() AT TIME ZONE u.timezone))*60 + EXTRACT(MINUTE FROM (now() AT TIME ZONE u.timezone)))
                    >= (EXTRACT(HOUR FROM u.silent_from)*60 + EXTRACT(MINUTE FROM u.silent_from))
                OR
                (EXTRACT(HOUR FROM (now() AT TIME ZONE u.timezone))*60 + EXTRACT(MINUTE FROM (now() AT TIME ZONE u.timezone)))
                    < (EXTRACT(HOUR FROM u.silent_until)*60 + EXTRACT(MINUTE FROM u.silent_until))
        END
    )`

	rows, err := j.pool.Query(ctx, q)
	if err != nil {
		j.logger.ErrorContext(ctx, "matchday: query targets", "err", err)
		return
	}
	defer rows.Close()

	var targets []matchdayTarget
	for rows.Next() {
		var t matchdayTarget
		if err := rows.Scan(&t.userID, &t.token, &t.dateKey); err != nil {
			j.logger.ErrorContext(ctx, "matchday: scan target", "err", err)
			return
		}
		targets = append(targets, t)
	}
	if err := rows.Err(); err != nil {
		j.logger.ErrorContext(ctx, "matchday: rows error", "err", err)
		return
	}

	if len(targets) == 0 {
		return
	}

	msgs := make([]expo.Message, 0, len(targets))
	for _, t := range targets {
		msgs = append(msgs, expo.Message{
			To:    t.token,
			Title: notification.MatchdayTitle,
			Body:  notification.MatchdayBody,
		})
	}

	receipts, err := j.expo.Send(ctx, msgs)
	if err != nil {
		j.logger.ErrorContext(ctx, "matchday: expo send", "err", err)
		return
	}

	for i, r := range receipts {
		if i >= len(targets) {
			break
		}
		t := targets[i]
		if r.Status == "error" && r.Details != nil && r.Details.Error == expo.ErrDeviceNotRegistered {
			j.deleteToken(ctx, t.token)
			continue
		}
		if r.Status != "ok" {
			j.logger.WarnContext(ctx, "matchday: expo error receipt", "status", r.Status, "msg", r.Message)
			continue
		}
		if err := j.insertLog(ctx, t.userID, "matchday", t.dateKey); err != nil {
			j.logger.ErrorContext(ctx, "matchday: insert log", "user_id", t.userID, "err", err)
		}
	}
}

func (j *MatchdayNotifier) deleteToken(ctx context.Context, token string) {
	const q = `DELETE FROM push_tokens WHERE token = $1`
	if _, err := j.pool.Exec(ctx, q, token); err != nil {
		j.logger.ErrorContext(ctx, "matchday: delete token", "token", token, "err", err)
	}
}

func (j *MatchdayNotifier) insertLog(ctx context.Context, userID uuid.UUID, typ, referenceID string) error {
	const q = `
INSERT INTO notification_log (user_id, type, reference_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, type, reference_id) DO NOTHING`
	_, err := j.pool.Exec(ctx, q, userID, typ, referenceID)
	return err
}
