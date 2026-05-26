package domain

import "github.com/google/uuid"

// LeaderboardEntry is a single row in a ranked leaderboard.
type LeaderboardEntry struct {
	Position     int
	UserID       uuid.UUID
	DisplayName  string
	ScorePoints  int
	PlayerPoints int
	TeamPoints   int
	TotalPoints  int
}
