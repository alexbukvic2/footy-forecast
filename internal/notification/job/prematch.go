// Package job contains background notification jobs.
package job

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/notification"
	"github.com/alexbukvic2/footy-forecast/internal/notification/expo"
)

// PreMatchNotifier sends notifications 2.5h before kickoff to users who have not
// yet predicted the fixture.
type PreMatchNotifier struct {
	pool   *db.Pool
	expo   *expo.Client
	logger *slog.Logger
}

// NewPreMatchNotifier constructs a PreMatchNotifier.
func NewPreMatchNotifier(
	pool *db.Pool,
	expoClient *expo.Client,
	logger *slog.Logger,
) *PreMatchNotifier {
	return &PreMatchNotifier{pool: pool, expo: expoClient, logger: logger}
}

// Run blocks on a 5-minute ticker until ctx is cancelled. It always returns nil.
func (j *PreMatchNotifier) Run(ctx context.Context) error {
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

type preMatchTarget struct {
	userID    uuid.UUID
	token     string
	homeTeam  string
	awayTeam  string
	fixtureID uuid.UUID
}

func (j *PreMatchNotifier) run(ctx context.Context) {
	// Find fixtures kicking off in 2h25m–2h35m that have eligible users (one row
	// per user+fixture+token combination).
	const q = `
SELECT
    u.id          AS user_id,
    pt.token      AS token,
    ht.name       AS home_team,
    at.name       AS away_team,
    f.id          AS fixture_id
FROM fixtures f
JOIN teams ht ON ht.id = f.home_team_id
JOIN teams at ON at.id = f.away_team_id
JOIN users u ON TRUE
JOIN push_tokens pt ON pt.user_id = u.id
WHERE f.kickoff_at BETWEEN now() + INTERVAL '2 hours 25 minutes'
                        AND now() + INTERVAL '2 hours 35 minutes'
  AND f.prediction_locked = FALSE
  -- user has not predicted this fixture
  AND NOT EXISTS (
      SELECT 1 FROM score_predictions sp
      WHERE sp.fixture_id = f.id AND sp.user_id = u.id
  )
  -- pre_match enabled (missing row = enabled)
  AND NOT EXISTS (
      SELECT 1 FROM notification_preferences np
      WHERE np.user_id = u.id AND np.type = 'pre_match' AND np.enabled = FALSE
  )
  -- not already sent for this fixture
  AND NOT EXISTS (
      SELECT 1 FROM notification_log nl
      WHERE nl.user_id = u.id AND nl.type = 'pre_match' AND nl.reference_id = f.id::text
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
  )
ORDER BY f.id, u.id`

	rows, err := j.pool.Query(ctx, q)
	if err != nil {
		j.logger.ErrorContext(ctx, "pre_match: query targets", "err", err)
		return
	}
	defer rows.Close()

	var targets []preMatchTarget
	for rows.Next() {
		var t preMatchTarget
		if err := rows.Scan(&t.userID, &t.token, &t.homeTeam, &t.awayTeam, &t.fixtureID); err != nil {
			j.logger.ErrorContext(ctx, "pre_match: scan target", "err", err)
			return
		}
		targets = append(targets, t)
	}
	if err := rows.Err(); err != nil {
		j.logger.ErrorContext(ctx, "pre_match: rows error", "err", err)
		return
	}

	if len(targets) == 0 {
		return
	}

	msgs := make([]expo.Message, 0, len(targets))
	for _, t := range targets {
		msgs = append(
			msgs, expo.Message{
				To:    t.token,
				Title: notification.PreMatchTitle,
				Body:  fmt.Sprintf(notification.PreMatchBodyFmt, t.homeTeam, t.awayTeam),
			},
		)
	}

	receipts, err := j.expo.Send(ctx, msgs)
	if err != nil {
		j.logger.ErrorContext(ctx, "pre_match: expo send", "err", err)
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
			j.logger.WarnContext(ctx, "pre_match: expo error receipt", "status", r.Status, "msg", r.Message)
			continue
		}
		if err := j.insertLog(ctx, t.userID, "pre_match", t.fixtureID.String()); err != nil {
			j.logger.ErrorContext(ctx, "pre_match: insert log", "user_id", t.userID, "err", err)
		}
	}
}

func (j *PreMatchNotifier) deleteToken(
	ctx context.Context,
	token string,
) {
	const q = `DELETE FROM push_tokens WHERE token = $1`
	if _, err := j.pool.Exec(ctx, q, token); err != nil {
		j.logger.ErrorContext(ctx, "pre_match: delete token", "token", token, "err", err)
	}
}

func (j *PreMatchNotifier) insertLog(
	ctx context.Context,
	userID uuid.UUID,
	typ, referenceID string,
) error {
	const q = `
INSERT INTO notification_log (user_id, type, reference_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, type, reference_id) DO NOTHING`
	_, err := j.pool.Exec(ctx, q, userID, typ, referenceID)
	return err
}
