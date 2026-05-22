package domain

import "testing"

func TestUserStatus_Valid(t *testing.T) {
	tests := []struct {
		status UserStatus
		want   bool
	}{
		{UserStatusActive, true},
		{UserStatusSuspended, true},
		{"", false},
		{"unknown", false},
	}
	for _, tc := range tests {
		if got := tc.status.Valid(); got != tc.want {
			t.Errorf("UserStatus(%q).Valid() = %v, want %v", tc.status, got, tc.want)
		}
	}
}
