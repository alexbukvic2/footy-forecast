package domain

// FixtureAnalysisInput is the data the AI job reads before calling Bedrock.
type FixtureAnalysisInput struct {
	HomeTeamName string
	AwayTeamName string
	Round        string
	GroupLetter  *string
	GoalsHome    *int
	GoalsAway    *int
	Predictions  []AnalysisPrediction
}

// AnalysisPrediction is one player's score prediction, used in the AI prompt.
type AnalysisPrediction struct {
	DisplayName string
	GoalsHome   *int
	GoalsAway   *int
	Points      *int
}
