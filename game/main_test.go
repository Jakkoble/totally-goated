package game

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Initialize real assets for all tests
	InitAssets(os.DirFS(".."))

	// Create a temp save file so tests loading stats don't use real saves
	saveFile = "test_save.json"
	defer os.Remove(saveFile)

	os.Exit(m.Run())
}
