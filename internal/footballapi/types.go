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
	Events  []fixtureEvent `json:"events"`
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

type fixtureEvent struct {
	Player playerRef `json:"player"`
	Type   string    `json:"type"`   // "Goal", "Card", etc.
	Detail string    `json:"detail"` // "Normal Goal", "Own Goal", "Penalty", etc.
}

type teamResult struct {
	ID     int64 `json:"id"`
	Winner *bool `json:"winner"`
}

type playerRef struct {
	ID int64 `json:"id"`
}
