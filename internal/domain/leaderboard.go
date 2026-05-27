package domain

import "github.com/google/uuid"

// LeaderboardEntry is a single row in a ranked leaderboard.
type LeaderboardEntry struct {
	Position          int
	UserID            uuid.UUID
	DisplayName       string
	ScorePts          int
	GroupTopScorerPts int
	TotalTopScorerPts int
	GroupWinnerPts    int
	PlayoffPts        int
	SemifinalistPts   int
	WinnerPts         int
	TotalPoints       int
}
