package game

import (
	"testing"
)

func TestSave(t *testing.T) {
	// Use a temporary file for tests so we don't destroy actual save data
	tempDir := t.TempDir()
	originalSaveFile := saveFile
	saveFile = tempDir + "/test_save.json"
	defer func() {
		saveFile = originalSaveFile // Restore after test
	}()

	// Initial load should return empty struct
	s1 := loadSave()
	if s1.BestScore != 0 || s1.BestMeters != 0 || s1.BestBells != 0 {
		t.Errorf("expected empty save data initially, got %v", s1)
	}

	// Save new best
	saveBest(100, 50, 10)

	s2 := loadSave()
	if s2.BestScore != 100 || s2.BestMeters != 50 || s2.BestBells != 10 {
		t.Errorf("expected {100, 50, 10}, got %v", s2)
	}

	// Save lower values, shouldn't overwrite bests
	saveBest(50, 20, 5)
	s3 := loadSave()
	if s3.BestScore != 100 || s3.BestMeters != 50 || s3.BestBells != 10 {
		t.Errorf("expected lower values to not overwrite bests, got %v", s3)
	}

	// Save higher values partially
	saveBest(150, 40, 15)
	s4 := loadSave()
	if s4.BestScore != 150 || s4.BestMeters != 50 || s4.BestBells != 15 {
		t.Errorf("expected partial overwrite, got %v", s4)
	}
}
