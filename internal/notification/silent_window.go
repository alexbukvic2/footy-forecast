// Package notification contains notification copy, silent window logic,
// and background job utilities.
package notification

import "github.com/alexbukvic2/footy-forecast/internal/domain"

// IsInSilentWindow reports whether localTime falls within the user's configured
// silent window [from, until). Returns false if from or until is nil.
//
// Same-day window (from < until): in window when from <= t < until.
// Midnight-wrapping (from > until): in window when t >= from OR t < until.
func IsInSilentWindow(localTime, from, until *domain.TimeOfDay) bool {
	if from == nil || until == nil {
		return false
	}

	t := localTime.Hour*60 + localTime.Minute
	f := from.Hour*60 + from.Minute
	u := until.Hour*60 + until.Minute

	if f < u {
		// Same-day window e.g. 13:00–14:00.
		return t >= f && t < u
	}
	// Wraps midnight e.g. 22:00–08:00.
	return t >= f || t < u
}
