package game

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// fetchLeaderboard makes an async HTTP request to get the top scores
func (g *Game) fetchLeaderboard() {
	g.leaderboardMutex.Lock()
	g.leaderboardLoading = true
	g.leaderboard = nil
	g.leaderboardMutex.Unlock()

	go func() {
		resp, err := http.Get("/api/leaderboard")
		if err != nil {
			g.leaderboardMutex.Lock()
			g.leaderboardLoading = false
			g.leaderboardMutex.Unlock()
			return
		}
		defer resp.Body.Close()

		var entries []ScoreEntry
		if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
			g.leaderboardMutex.Lock()
			g.leaderboardLoading = false
			g.leaderboardMutex.Unlock()
			return
		}

		g.leaderboardMutex.Lock()
		g.leaderboard = entries
		g.leaderboardLoading = false
		g.leaderboardMutex.Unlock()
	}()
}

// submitScoreAsync posts the score asynchronously
func (g *Game) submitScoreAsync(name string, score int) {
	g.leaderboardMutex.Lock()
	g.leaderboardLoading = true
	g.leaderboardMutex.Unlock()

	entry := ScoreEntry{Name: name, Score: score}
	data, err := json.Marshal(entry)
	if err != nil {
		g.fetchLeaderboard() // fallback to fetch if serialize fails
		return
	}

	go func() {
		resp, err := http.Post("/api/score", "application/json", bytes.NewBuffer(data))
		if err == nil {
			resp.Body.Close()
		}

		// After submit finishes, trigger a fetch so the user sees the newly added score
		g.fetchLeaderboard()
	}()
}
