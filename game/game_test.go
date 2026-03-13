package game

import (
	"testing"
)

func TestGame_StartGame(t *testing.T) {
	g := NewGame()

	if g.state != GameMenu {
		t.Errorf("expected initial state to be GameMenu, got %v", g.state)
	}

	g.startGame()

	if g.state != GamePlaying {
		t.Errorf("expected state after start to be GamePlaying, got %v", g.state)
	}

	if g.level == nil {
		t.Errorf("expected level to be initialized")
	}

	if g.score != 0 {
		t.Errorf("expected score to be reset, got %v", g.score)
	}
}

func TestGame_CurrentMetersAndScore(t *testing.T) {
	g := NewGame()
	g.startGame()

	// Setting score (which tracks height natively in game)
	g.score = 300 // 10 meters at 30 pixels per meter

	if g.currentMeters() != 10 {
		t.Errorf("expected 10 meters, got %d", g.currentMeters())
	}

	g.bellScore = 25 // 5 bells * 5 score
	total := g.totalScore()

	if total != 35 { // 10 + 25
		t.Errorf("expected total score to be 35, got %d", total)
	}
}

func TestGame_OnBellCollected(t *testing.T) {
	g := NewGame()
	g.startGame()

	pos := Vec2{0, 0}

	// Collect a bell
	g.OnBellCollected(pos)

	if g.bellCount != 1 {
		t.Errorf("expected 1 bell collected, got %d", g.bellCount)
	}

	if g.bellScore != BellScoreValue {
		t.Errorf("expected bell score %d, got %d", BellScoreValue, g.bellScore)
	}

	if g.comboCount != 1 {
		t.Errorf("expected combo count 1, got %d", g.comboCount)
	}

	// Collect another bell immediately for a combo (x2)
	g.OnBellCollected(pos)

	if g.comboCount != 2 {
		t.Errorf("expected combo count 2, got %d", g.comboCount)
	}

	if g.bellScore != BellScoreValue + BellScoreValue*2 {
		t.Errorf("expected total bell score %d, got %d", BellScoreValue + BellScoreValue*2, g.bellScore)
	}
}

func TestGame_UpdateCombo(t *testing.T) {
	g := NewGame()
	g.startGame()

	g.comboCount = 2
	g.comboTimer = 0.01 // Minimal timer

	// Normally Update decreases the timer based on TPS. We can just call updateCombo with the mock tick manually.
	g.updateCombo() // Decrements by 1/TPS (~0.016)

	if g.comboCount != 0 {
		t.Errorf("expected combo count to reset to 0 after timer expiry, got %d", g.comboCount)
	}
}

func TestGame_SkyColor(t *testing.T) {
	g := NewGame()

	// Low altitude
	g.cameraY = 0
	lowColor := g.skyColor()
	if lowColor.R != 30 || lowColor.G != 30 || lowColor.B != 50 {
		t.Errorf("expected low altitude color {30, 30, 50}, got %v", lowColor)
	}

	// High altitude
	g.cameraY = -500 * PixelsPerMeter
	highColor := g.skyColor()
	if highColor.R != 50 || highColor.G != 20 || highColor.B != 70 {
		t.Errorf("expected high altitude color {50, 20, 70}, got %v", highColor)
	}

	// Mid altitude
	g.cameraY = -250 * PixelsPerMeter
	midColor := g.skyColor()
	if midColor.R != 40 || midColor.G != 25 || midColor.B != 60 {
		t.Errorf("expected mid altitude color {40, 25, 60}, got %v", midColor)
	}
}
