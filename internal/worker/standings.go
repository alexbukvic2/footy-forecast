package worker

import (
	"sort"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// teamPair is an ordered (team, opponent) key for head-to-head point lookups.
type teamPair struct{ team, opponent uuid.UUID }

// computeGroupStandings builds ordered standings from raw fixture results.
//
// Ranking criteria (in order):
//  1. Points (W=3, D=1, L=0)
//  2. Head-to-head points between the tied teams
//  3. Overall goal difference
//  4. Overall goals scored
func computeGroupStandings(
	groupLetter string,
	teams []uuid.UUID,
	fixtures []domain.GroupFixture,
) []domain.StandingsEntry {
	type stats struct {
		points, played, won, drawn, lost, goalsFor, goalsAgainst int
	}

	m := make(map[uuid.UUID]*stats, len(teams))
	for _, t := range teams {
		m[t] = &stats{}
	}

	// h2h[[a, b]] = points team a earned in matches against team b.
	h2h := make(map[teamPair]int)

	for _, f := range fixtures {
		h, hok := m[f.HomeTeamID]
		a, aok := m[f.AwayTeamID]
		if !hok || !aok {
			continue
		}
		h.played++
		h.goalsFor += f.GoalsHome
		h.goalsAgainst += f.GoalsAway
		a.played++
		a.goalsFor += f.GoalsAway
		a.goalsAgainst += f.GoalsHome
		switch {
		case f.GoalsHome > f.GoalsAway:
			h.points += 3
			h.won++
			a.lost++
			h2h[teamPair{f.HomeTeamID, f.AwayTeamID}] += 3
		case f.GoalsHome < f.GoalsAway:
			a.points += 3
			a.won++
			h.lost++
			h2h[teamPair{f.AwayTeamID, f.HomeTeamID}] += 3
		default:
			h.points++
			h.drawn++
			a.points++
			a.drawn++
			h2h[teamPair{f.HomeTeamID, f.AwayTeamID}]++
			h2h[teamPair{f.AwayTeamID, f.HomeTeamID}]++
		}
	}

	sorted := make([]uuid.UUID, len(teams))
	copy(sorted, teams)
	sort.SliceStable(sorted, func(i, j int) bool {
		si, sj := m[sorted[i]], m[sorted[j]]
		if si.points != sj.points {
			return si.points > sj.points
		}
		hi := h2h[teamPair{sorted[i], sorted[j]}]
		hj := h2h[teamPair{sorted[j], sorted[i]}]
		if hi != hj {
			return hi > hj
		}
		gdi := si.goalsFor - si.goalsAgainst
		gdj := sj.goalsFor - sj.goalsAgainst
		if gdi != gdj {
			return gdi > gdj
		}
		return si.goalsFor > sj.goalsFor
	})

	gl := groupLetter
	entries := make([]domain.StandingsEntry, len(sorted))
	for i, tid := range sorted {
		s := m[tid]
		entries[i] = domain.StandingsEntry{
			TeamID:       tid,
			Group:        &gl,
			Position:     i + 1,
			Points:       s.points,
			Played:       s.played,
			Won:          s.won,
			Drawn:        s.drawn,
			Lost:         s.lost,
			GoalsFor:     s.goalsFor,
			GoalsAgainst: s.goalsAgainst,
		}
	}
	return entries
}
