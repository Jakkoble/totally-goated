package game
import (
    "testing"
)
func TestInitAssets(t *testing.T) {
	// Assets are already initialized in TestMain for all tests
	// So we don't need to explicitly call InitAssets here.
	// But we can check if it loaded things properly:
	if assetsFS == nil {
		t.Errorf("expected assets to be loaded by TestMain")
	}
}
