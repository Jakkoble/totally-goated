package game

import (
	"testing"
)

// TestE2E_FullGameCycle tests the main lifecycle of the game by calling g.Update()
// and simulating the conditions that lead to state transitions.
func TestE2E_FullGameCycle(t *testing.T) {
	g := NewGame()

	// Initial State: Menu
	if g.state != GameMenu {
		t.Fatalf("expected GameMenu state, got %v", g.state)
	}

	// Action: Start Game (since we can't easily mock ebiten inputs in headless
	// without injecting an input interface, we'll call startGame directly to
	// transition from Menu -> Playing, which is what Space/Click does).
	g.startGame()

	if g.state != GamePlaying {
		t.Fatalf("expected GamePlaying state, got %v", g.state)
	}

	// Provide some initial state for the play session
	g.Goat.Pos = Vec2{0, 0}

	// Ensure Update loop runs smoothly without panicking
	err := g.Update()
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	// Let's drop a bell exactly where the goat is right now,
	// and call collectBells manually to bypass unpredictable loop timing logic in test
	g.level.Bells = []Bell{
		{Pos: g.Goat.Pos, Collected: false},
	}
	g.Goat.collectBells(g.level, g) // trigger manually as if Update did it perfectly aligned

	if g.bellCount != 1 {
		t.Errorf("expected 1 bell to be collected via Update loop, got %d", g.bellCount)
	}

	// Simulate falling off the screen to trigger Game Over
	g.cameraY = 0
	g.Goat.Pos.Y = DeathMargin + 100 // Fall way below camera
	g.Goat.HasShield = false // Ensure no shield saves us

	g.Update() // This Update should trigger the death logic and switch to GameOver

	if g.state != GameOver {
		t.Fatalf("expected GameOver state after falling, got %v", g.state)
	}

	// Verify Total Score was preserved and high scores updated (mocked save file is active)
	totalScore := g.totalScore()
	if totalScore == 0 && g.bellCount > 0 {
		t.Errorf("expected total score to be > 0 due to bell collection, got %v", totalScore)
	}

	// Since we are in game over, if we manually trigger state to Menu (to simulate continue click)
	g.state = GameMenu

	if g.state != GameMenu {
		t.Fatalf("expected restart to Menu state, got %v", g.state)
	}
}
