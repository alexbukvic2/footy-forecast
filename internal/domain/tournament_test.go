package domain_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

func TestTournamentStatus_Valid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		status domain.TournamentStatus
		want   bool
	}{
		{"upcoming", domain.TournamentStatusUpcoming, true},
		{"in_progress", domain.TournamentStatusInProgress, true},
		{"concluded", domain.TournamentStatusConcluded, true},
		{"unknown lowercase", domain.TournamentStatus("future"), false},
		{"empty string", domain.TournamentStatus(""), false},
		{"casing matters", domain.TournamentStatus("Upcoming"), false},
	}

	for _, tc := range cases {
		t.Run(
			tc.name, func(t *testing.T) {
				t.Parallel()
				require.Equal(t, tc.want, tc.status.Valid())
			},
		)
	}
}
