package worker

import (
	"strings"
	"testing"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

func intPtr(v int) *int { return &v }

func TestBuildAnalysisPrompt(t *testing.T) {
	group := "A"
	tests := []struct {
		name     string
		input    domain.FixtureAnalysisInput
		contains []string
	}{
		{
			name: "finished group match with predictions",
			input: domain.FixtureAnalysisInput{
				HomeTeamName: "Brazil",
				AwayTeamName: "Germany",
				Round:        "Group Stage - 1",
				GroupLetter:  &group,
				GoalsHome:    intPtr(2),
				GoalsAway:    intPtr(1),
				Predictions: []domain.AnalysisPrediction{
					{DisplayName: "alice", GoalsHome: intPtr(2), GoalsAway: intPtr(1), Points: intPtr(3)},
					{DisplayName: "bob", GoalsHome: intPtr(1), GoalsAway: intPtr(0), Points: intPtr(0)},
				},
			},
			contains: []string{
				"Brazil", "Germany", "Group A", "Group Stage - 1",
				"2 – 1",
				"alice", "bob",
				"bullet points",
			},
		},
		{
			name: "no predictions",
			input: domain.FixtureAnalysisInput{
				HomeTeamName: "France",
				AwayTeamName: "Spain",
				Round:        "Final",
				GoalsHome:    intPtr(1),
				GoalsAway:    intPtr(0),
				Predictions:  nil,
			},
			contains: []string{
				"No predictions were submitted",
				"France", "Spain",
			},
		},
		{
			name: "result not yet available",
			input: domain.FixtureAnalysisInput{
				HomeTeamName: "Argentina",
				AwayTeamName: "England",
				Round:        "Semi-Final",
			},
			contains: []string{
				"not yet available",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			prompt := buildAnalysisPrompt(tc.input)
			for _, want := range tc.contains {
				if !strings.Contains(prompt, want) {
					t.Errorf("prompt missing %q\nprompt:\n%s", want, prompt)
				}
			}
		})
	}
}
