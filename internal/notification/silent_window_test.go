package notification

import (
	"testing"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

func ptr(h, m int) *domain.TimeOfDay {
	return &domain.TimeOfDay{Hour: h, Minute: m}
}

func TestIsInSilentWindow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		localTime *domain.TimeOfDay
		from      *domain.TimeOfDay
		until     *domain.TimeOfDay
		want      bool
	}{
		// Nil inputs.
		{
			name:      "nil from returns false",
			localTime: ptr(12, 0),
			from:      nil,
			until:     ptr(14, 0),
			want:      false,
		},
		{
			name:      "nil until returns false",
			localTime: ptr(12, 0),
			from:      ptr(10, 0),
			until:     nil,
			want:      false,
		},
		{
			name:      "both nil returns false",
			localTime: ptr(12, 0),
			from:      nil,
			until:     nil,
			want:      false,
		},

		// Same-day window (from < until): 13:00–14:00.
		{
			name:      "same-day: before window",
			localTime: ptr(12, 59),
			from:      ptr(13, 0),
			until:     ptr(14, 0),
			want:      false,
		},
		{
			name:      "same-day: at start boundary (inclusive)",
			localTime: ptr(13, 0),
			from:      ptr(13, 0),
			until:     ptr(14, 0),
			want:      true,
		},
		{
			name:      "same-day: inside window",
			localTime: ptr(13, 30),
			from:      ptr(13, 0),
			until:     ptr(14, 0),
			want:      true,
		},
		{
			name:      "same-day: at end boundary (exclusive)",
			localTime: ptr(14, 0),
			from:      ptr(13, 0),
			until:     ptr(14, 0),
			want:      false,
		},
		{
			name:      "same-day: after window",
			localTime: ptr(14, 1),
			from:      ptr(13, 0),
			until:     ptr(14, 0),
			want:      false,
		},

		// Midnight-wrapping window (from > until): 22:00–08:00.
		{
			name:      "midnight-wrap: before from",
			localTime: ptr(21, 59),
			from:      ptr(22, 0),
			until:     ptr(8, 0),
			want:      false,
		},
		{
			name:      "midnight-wrap: at from boundary (inclusive)",
			localTime: ptr(22, 0),
			from:      ptr(22, 0),
			until:     ptr(8, 0),
			want:      true,
		},
		{
			name:      "midnight-wrap: after midnight inside window",
			localTime: ptr(0, 0),
			from:      ptr(22, 0),
			until:     ptr(8, 0),
			want:      true,
		},
		{
			name:      "midnight-wrap: just before until boundary",
			localTime: ptr(7, 59),
			from:      ptr(22, 0),
			until:     ptr(8, 0),
			want:      true,
		},
		{
			name:      "midnight-wrap: at until boundary (exclusive)",
			localTime: ptr(8, 0),
			from:      ptr(22, 0),
			until:     ptr(8, 0),
			want:      false,
		},
		{
			name:      "midnight-wrap: midday outside window",
			localTime: ptr(12, 0),
			from:      ptr(22, 0),
			until:     ptr(8, 0),
			want:      false,
		},

		// Edge: single-minute same-day window 10:00–10:01.
		{
			name:      "single-minute window: inside",
			localTime: ptr(10, 0),
			from:      ptr(10, 0),
			until:     ptr(10, 1),
			want:      true,
		},
		{
			name:      "single-minute window: outside",
			localTime: ptr(10, 1),
			from:      ptr(10, 0),
			until:     ptr(10, 1),
			want:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := IsInSilentWindow(tc.localTime, tc.from, tc.until)
			if got != tc.want {
				t.Errorf("IsInSilentWindow(%v, %v, %v) = %v, want %v",
					tc.localTime, tc.from, tc.until, got, tc.want)
			}
		})
	}
}
