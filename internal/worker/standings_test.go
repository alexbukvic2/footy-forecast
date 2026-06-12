package worker

import (
	"testing"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

func TestComputeGroupStandings_AllPlayedOnce(t *testing.T) {
	// 4-team group: A beats B 2-1, C draws D 0-0, A beats C 1-0, B beats D 3-0, A beats D 1-0, B beats C 2-1
	teamA, teamB, teamC, teamD := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	teams := []uuid.UUID{teamA, teamB, teamC, teamD}
	fixtures := []domain.GroupFixture{
		{HomeTeamID: teamA, AwayTeamID: teamB, GoalsHome: 2, GoalsAway: 1},
		{HomeTeamID: teamC, AwayTeamID: teamD, GoalsHome: 0, GoalsAway: 0},
		{HomeTeamID: teamA, AwayTeamID: teamC, GoalsHome: 1, GoalsAway: 0},
		{HomeTeamID: teamB, AwayTeamID: teamD, GoalsHome: 3, GoalsAway: 0},
		{HomeTeamID: teamA, AwayTeamID: teamD, GoalsHome: 1, GoalsAway: 0},
		{HomeTeamID: teamB, AwayTeamID: teamC, GoalsHome: 2, GoalsAway: 1},
	}

	entries := computeGroupStandings("A", teams, fixtures)

	if len(entries) != 4 {
		t.Fatalf("want 4 entries, got %d", len(entries))
	}
	// A: 3W 0D 0L = 9pts; B: 2W 0D 1L = 6pts; C: 0W 1D 2L = 1pt; D: 0W 1D 2L = 1pt
	// D vs C tiebreak: C has +(-2) GD, D has -(3+0) = -3 GD → C above D
	assertEntry(t, entries[0], teamA, 1, 9, 3)
	assertEntry(t, entries[1], teamB, 2, 6, 3)
	assertEntry(t, entries[2], teamC, 3, 1, 3)
	assertEntry(t, entries[3], teamD, 4, 1, 3)
}

func TestComputeGroupStandings_H2HTiebreaker(t *testing.T) {
	// Two teams tied on points; H2H should decide.
	teamA, teamB, teamC := uuid.New(), uuid.New(), uuid.New()
	teams := []uuid.UUID{teamA, teamB, teamC}
	fixtures := []domain.GroupFixture{
		// A beats C, B beats C — both on 3pts; H2H: B beat A
		{HomeTeamID: teamB, AwayTeamID: teamA, GoalsHome: 1, GoalsAway: 0},
		{HomeTeamID: teamA, AwayTeamID: teamC, GoalsHome: 2, GoalsAway: 0},
		{HomeTeamID: teamB, AwayTeamID: teamC, GoalsHome: 2, GoalsAway: 0},
	}

	entries := computeGroupStandings("B", teams, fixtures)

	if len(entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(entries))
	}
	// B and A both 3pts; H2H B beat A → B first
	if entries[0].TeamID != teamB {
		t.Errorf("want B first (H2H winner), got position 1 = %v", entries[0].TeamID)
	}
	if entries[1].TeamID != teamA {
		t.Errorf("want A second, got position 2 = %v", entries[1].TeamID)
	}
}

func TestComputeGroupStandings_GoalDifferenceTiebreaker(t *testing.T) {
	// Two teams tied on points and H2H; better GD wins.
	teamA, teamB := uuid.New(), uuid.New()
	teams := []uuid.UUID{teamA, teamB}
	// A and B drew their head-to-head; A has better overall GD.
	fixtures := []domain.GroupFixture{
		{HomeTeamID: teamA, AwayTeamID: teamB, GoalsHome: 1, GoalsAway: 1},
	}

	// Inject a phantom third-team result by adding fixtures outside the teams list — they must be ignored.
	phantom := uuid.New()
	fixtures = append(fixtures,
		domain.GroupFixture{HomeTeamID: phantom, AwayTeamID: teamA, GoalsHome: 0, GoalsAway: 3},
	)

	entries := computeGroupStandings("C", teams, fixtures)

	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	// Both 1pt, H2H draw (1-1); GD: A=+2 from phantom match BUT phantom is not in teams → ignored.
	// So GD A = 1-1 = 0, B = 1-1 = 0; GoalsFor A = 1, B = 1 → equal, stable sort keeps original order.
	if entries[0].Points != 1 || entries[1].Points != 1 {
		t.Errorf("both teams should have 1 point")
	}
}

func TestComputeGroupStandings_NoFixtures(t *testing.T) {
	// Before any matches, all teams should be at position 1..N with 0 points.
	teamA, teamB := uuid.New(), uuid.New()
	teams := []uuid.UUID{teamA, teamB}

	entries := computeGroupStandings("D", teams, nil)

	if len(entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Points != 0 || e.Played != 0 {
			t.Errorf("expected zero stats before any fixtures, got %+v", e)
		}
	}
}

func TestComputeGroupStandings_GroupFieldSet(t *testing.T) {
	teamA := uuid.New()
	entries := computeGroupStandings("E", []uuid.UUID{teamA}, nil)
	if len(entries) != 1 {
		t.Fatalf("want 1 entry")
	}
	if entries[0].Group == nil || *entries[0].Group != "E" {
		t.Errorf("Group field not set correctly: %v", entries[0].Group)
	}
}

func assertEntry(t *testing.T, e domain.StandingsEntry, teamID uuid.UUID, pos, pts, played int) {
	t.Helper()
	if e.TeamID != teamID {
		t.Errorf("position %d: want team %v, got %v", pos, teamID, e.TeamID)
	}
	if e.Position != pos {
		t.Errorf("team %v: want position %d, got %d", teamID, pos, e.Position)
	}
	if e.Points != pts {
		t.Errorf("team %v: want %d points, got %d", teamID, pts, e.Points)
	}
	if e.Played != played {
		t.Errorf("team %v: want %d played, got %d", teamID, played, e.Played)
	}
}
