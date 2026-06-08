package domain

import "time"

// NotificationType identifies a push notification category.
type NotificationType string

// Notification type constants.
const (
	NotificationTypeMatchday           NotificationType = "matchday"
	NotificationTypePreMatch           NotificationType = "pre_match"
	NotificationTypeTournamentReminder NotificationType = "tournament_reminder"
)

// AllNotificationTypes lists the types shown in GET /users/me/notification-preferences
// and for which defaults are materialised. tournament_reminder is intentionally
// excluded from this list so it doesn't appear as a user-facing preference,
// but it can still be toggled via PUT.
var AllNotificationTypes = []NotificationType{
	NotificationTypeMatchday,
	NotificationTypePreMatch,
}

// ValidNotificationTypes is the full set accepted by the preferences API.
var ValidNotificationTypes = []NotificationType{
	NotificationTypeMatchday,
	NotificationTypePreMatch,
	NotificationTypeTournamentReminder,
}

// PushToken represents a device push token registered to a user.
type PushToken struct {
	ID        string
	UserID    string
	Token     string
	CreatedAt time.Time
}

// NotificationPreference records whether a user has enabled a notification type.
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
