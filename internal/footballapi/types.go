// Package footballapi implements the worker.MatchAPI interface against api-sports.io v3.
package footballapi

// fixtureResponse is the top-level JSON envelope for GET /fixtures.
type fixtureResponse struct {
	Response []fixtureItem `json:"response"`
}

type fixtureItem struct {
	Fixture fixtureDetails `json:"fixture"`
	League  fixtureLeague  `json:"league"`
	Goals   goalsDetails   `json:"goals"`
	Teams   teamsDetails   `json:"teams"`
}

type fixtureLeague struct {
	Round string `json:"round"`
}

type fixtureDetails struct {
	ID     int64         `json:"id"`
	Date   string        `json:"date"` // RFC3339 kickoff timestamp
	Status fixtureStatus `json:"status"`
}

type fixtureStatus struct {
	Short string `json:"short"`
}

type goalsDetails struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type teamsDetails struct {
	Home teamResult `json:"home"`
	Away teamResult `json:"away"`
}

type teamResult struct {
	ID     int64 `json:"id"`
	Winner *bool `json:"winner"`
}

// topScorersResponse is the top-level JSON envelope for GET /players/topscorers.
type topScorersResponse struct {
	Response []topScorerItem `json:"response"`
}

type topScorerItem struct {
	Player playerRef    `json:"player"`
	Stats  []playerStat `json:"statistics"`
}

type playerRef struct {
	ID int64 `json:"id"`
}

type playerStat struct {
	Goals goalsStat `json:"goals"`
}

type goalsStat struct {
	Total *int `json:"total"`
}
