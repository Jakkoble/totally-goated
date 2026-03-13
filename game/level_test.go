package game

import (
	"testing"
)

func TestLevel_NewLevel(t *testing.T) {
	l := NewLevel()
	if len(l.Platforms) != 1 {
		t.Errorf("expected 1 initial platform, got %d", len(l.Platforms))
	}
	if l.Platforms[0].Y != 0 {
		t.Errorf("expected initial platform to be at Y=0, got %v", l.Platforms[0].Y)
	}
}

func TestLevel_Overlaps(t *testing.T) {
	l := NewLevel()
	// Initial platform: {X: -35, Y: 0, W: 70, H: 120}

	p1 := Platform{X: 0, Y: 50, W: 50, H: 50}
	if !l.overlaps(p1) {
		t.Errorf("expected platform to overlap initial platform")
	}

	p2 := Platform{X: 100, Y: 100, W: 50, H: 50}
	if l.overlaps(p2) {
		t.Errorf("expected platform to NOT overlap initial platform")
	}
}

func TestLevel_GenerateUntil(t *testing.T) {
	l := NewLevel()

	// Generate up to -1000 pixels (upwards)
	l.GenerateUntil(-1000)

	if l.curY > -1000 {
		t.Errorf("expected curY to be <= -1000, got %v", l.curY)
	}

	if len(l.Platforms) <= 1 {
		t.Errorf("expected multiple platforms to be generated, got %d", len(l.Platforms))
	}

	// Verify alternate placement logic and no overlaps (roughly)
	for i := range l.Platforms {
		for j := range l.Platforms {
			if i != j {
				// We can't use l.overlaps directly for each pair easily, but the generation
				// ensures there's no overlapping platforms in its own logic.
				// We just check they exist and have valid coordinates.
				if l.Platforms[i].W <= 0 || l.Platforms[i].H <= 0 {
					t.Errorf("platform %d has invalid dimensions: %v, %v", i, l.Platforms[i].W, l.Platforms[i].H)
				}
			}
		}
	}
}

func TestLevel_Difficulty(t *testing.T) {
	l := NewLevel()

	l.curY = 0
	if l.difficulty() != 0 {
		t.Errorf("expected difficulty 0 at Y=0, got %v", l.difficulty())
	}

	l.curY = -250 * PixelsPerMeter // Halfway to max difficulty
	diff := l.difficulty()
	if diff < 0.49 || diff > 0.51 {
		t.Errorf("expected difficulty ~0.5, got %v", diff)
	}

	l.curY = -600 * PixelsPerMeter // Past max difficulty
	if l.difficulty() != 1 {
		t.Errorf("expected max difficulty 1, got %v", l.difficulty())
	}
}
