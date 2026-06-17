package worker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// runAIAnalysis fires as a detached goroutine after a fixture finishes.
// It fetches all leagues for the fixture, then concurrently generates and
// persists AI analysis for each league.
func (w *Worker) runAIAnalysis(fixtureID uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	leagueIDs, err := w.repo.ListLeaguesForFixture(ctx, fixtureID)
	cancel()
	if err != nil {
		w.logger.Error("ai_analysis: list leagues", "fixture_id", fixtureID, "err", err)
		return
	}
	if len(leagueIDs) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, lid := range leagueIDs {
		wg.Add(1)
		go func(leagueID uuid.UUID) {
			defer wg.Done()
			w.runAIAnalysisForLeague(fixtureID, leagueID)
		}(lid)
	}
	wg.Wait()
}

// runAIAnalysisForLeague generates and persists AI analysis for a single league.
func (w *Worker) runAIAnalysisForLeague(fixtureID, leagueID uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	input, err := w.repo.GetFixtureAnalysisInput(ctx, fixtureID, leagueID)
	if err != nil {
		w.logger.Error("ai_analysis: get input", "fixture_id", fixtureID, "league_id", leagueID, "err", err)
		return
	}
	if len(input.Predictions) == 0 {
		return
	}

	prompt := buildAnalysisPrompt(input)

	var analysis string
	backoff := []time.Duration{2 * time.Second, 4 * time.Second}
	for attempt := 0; attempt <= 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				w.logger.Warn(
					"ai_analysis: context done before retry",
					"fixture_id",
					fixtureID,
					"league_id",
					leagueID,
					"attempt",
					attempt,
				)
				return
			case <-time.After(backoff[attempt-1]):
			}
		}
		analysis, err = w.analyser.Analyse(ctx, prompt)
		if err == nil {
			break
		}
		w.logger.Warn(
			"ai_analysis: bedrock call failed",
			"fixture_id",
			fixtureID,
			"league_id",
			leagueID,
			"attempt",
			attempt+1,
			"err",
			err,
		)
	}
	if err != nil {
		w.logger.Error("ai_analysis: all attempts failed", "fixture_id", fixtureID, "league_id", leagueID, "err", err)
		return
	}

	if err := w.repo.UpsertFixtureAnalysis(ctx, fixtureID, leagueID, analysis); err != nil {
		w.logger.Error("ai_analysis: upsert failed", "fixture_id", fixtureID, "league_id", leagueID, "err", err)
	}
}

// buildAnalysisPrompt formats the fixture and prediction data into the LLM prompt.
func buildAnalysisPrompt(input domain.FixtureAnalysisInput) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Match: %s vs %s", input.HomeTeamName, input.AwayTeamName)
	if input.GroupLetter != nil {
		fmt.Fprintf(&sb, " (Group %s)", *input.GroupLetter)
	}
	fmt.Fprintf(&sb, ", Round: %s\n", input.Round)

	if input.GoalsHome != nil && input.GoalsAway != nil {
		fmt.Fprintf(
			&sb,
			"Final score: %s %d – %d %s\n",
			input.HomeTeamName,
			*input.GoalsHome,
			*input.GoalsAway,
			input.AwayTeamName,
		)
	} else {
		sb.WriteString("Final score: not yet available\n")
	}

	if len(input.Predictions) == 0 {
		sb.WriteString("No predictions were submitted for this match.\n")
	} else {
		fmt.Fprintf(&sb, "\nPlayer predictions (%d total):\n", len(input.Predictions))
		for _, p := range input.Predictions {
			if p.GoalsHome != nil && p.GoalsAway != nil {
				pts := "?"
				if p.Points != nil {
					pts = fmt.Sprintf("%d", *p.Points)
				}
				fmt.Fprintf(&sb, "- %s predicted %d–%d (points: %s)\n", p.DisplayName, *p.GoalsHome, *p.GoalsAway, pts)
			} else {
				fmt.Fprintf(&sb, "- %s: no prediction\n", p.DisplayName)
			}
		}
	}

	sb.WriteString(
		`Generate 4-5 bullet points about this match. Rules: use percentages or counts (e.g. "3/16", "only 18%"), name players only when there are 1 or 2 of them (e.g. "only Andy got it right", "only Andy and Demijan got it right") — if 3 or more, use the count only, keep each bullet to one short punchy sentence, no narrative prose. Examples of the style: "only 18% got the exact score", "majority predicted 2–0, missing the away goal", "nobody predicted a home loss". Do not use em dashes.`,
	)

	return sb.String()
}
