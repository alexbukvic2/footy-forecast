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

// TournamentReminderNotifier sends a one-time notification per tournament, 5h
// before the earliest fixture, to all users who have not disabled it.
type TournamentReminderNotifier struct {
	pool   *db.Pool
	expo   *expo.Client
	logger *slog.Logger
}

// NewTournamentReminderNotifier constructs a TournamentReminderNotifier.
func NewTournamentReminderNotifier(pool *db.Pool, expoClient *expo.Client, logger *slog.Logger) *TournamentReminderNotifier {
	return &TournamentReminderNotifier{pool: pool, expo: expoClient, logger: logger}
}

// Run blocks on a 5-minute ticker until ctx is cancelled. It always returns nil.
func (j *TournamentReminderNotifier) Run(ctx context.Context) error {
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

type tournamentReminderTarget struct {
	userID         uuid.UUID
	token          string
	tournamentID   uuid.UUID
	tournamentName string
}

func (j *TournamentReminderNotifier) run(ctx context.Context) {
	const q = `
SELECT
    u.id              AS user_id,
    pt.token          AS token,
    t.id              AS tournament_id,
    t.name            AS tournament_name
FROM tournaments t
JOIN (
    SELECT tournament_id, MIN(kickoff_at) AS earliest_kickoff
    FROM fixtures
    WHERE is_demo = FALSE
    GROUP BY tournament_id
) earliest ON earliest.tournament_id = t.id
JOIN users u ON TRUE
JOIN push_tokens pt ON pt.user_id = u.id
WHERE
    earliest.earliest_kickoff BETWEEN now() + INTERVAL '4 hours 55 minutes'
                                  AND now() + INTERVAL '5 hours 5 minutes'
    -- tournament_reminder enabled (missing row = enabled)
    AND NOT EXISTS (
        SELECT 1 FROM notification_preferences np
        WHERE np.user_id = u.id AND np.type = 'tournament_reminder' AND np.enabled = FALSE
    )
    -- not already sent for this tournament
    AND NOT EXISTS (
        SELECT 1 FROM notification_log nl
        WHERE nl.user_id = u.id
          AND nl.type = 'tournament_reminder'
          AND nl.reference_id = t.id::text
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
ORDER BY t.id, u.id`

	rows, err := j.pool.Query(ctx, q)
	if err != nil {
		j.logger.ErrorContext(ctx, "tournament_reminder: query targets", "err", err)
		return
	}
	defer rows.Close()

	var targets []tournamentReminderTarget
	for rows.Next() {
		var t tournamentReminderTarget
		if err := rows.Scan(&t.userID, &t.token, &t.tournamentID, &t.tournamentName); err != nil {
			j.logger.ErrorContext(ctx, "tournament_reminder: scan target", "err", err)
			return
		}
		targets = append(targets, t)
	}
	if err := rows.Err(); err != nil {
		j.logger.ErrorContext(ctx, "tournament_reminder: rows error", "err", err)
		return
	}

	if len(targets) == 0 {
		return
	}

	msgs := make([]expo.Message, 0, len(targets))
	for _, t := range targets {
		msgs = append(msgs, expo.Message{
			To:    t.token,
			Title: notification.TournamentReminderTitle,
			Body:  fmt.Sprintf(notification.TournamentReminderBodyFmt, t.tournamentName),
		})
	}

	receipts, err := j.expo.Send(ctx, msgs)
	if err != nil {
		j.logger.ErrorContext(ctx, "tournament_reminder: expo send", "err", err)
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
			j.logger.WarnContext(ctx, "tournament_reminder: expo error receipt", "status", r.Status, "msg", r.Message)
			continue
		}
		if err := j.insertLog(ctx, t.userID, "tournament_reminder", t.tournamentID.String()); err != nil {
			j.logger.ErrorContext(ctx, "tournament_reminder: insert log", "user_id", t.userID, "err", err)
		}
	}
}

func (j *TournamentReminderNotifier) deleteToken(ctx context.Context, token string) {
	const q = `DELETE FROM push_tokens WHERE token = $1`
	if _, err := j.pool.Exec(ctx, q, token); err != nil {
		j.logger.ErrorContext(ctx, "tournament_reminder: delete token", "token", token, "err", err)
	}
}

func (j *TournamentReminderNotifier) insertLog(ctx context.Context, userID uuid.UUID, typ, referenceID string) error {
	const q = `
INSERT INTO notification_log (user_id, type, reference_id)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, type, reference_id) DO NOTHING`
	_, err := j.pool.Exec(ctx, q, userID, typ, referenceID)
	return err
}
