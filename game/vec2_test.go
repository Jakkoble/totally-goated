package game

import (
	"math"
	"testing"
)

func TestVec2_Add(t *testing.T) {
	v1 := Vec2{1, 2}
	v2 := Vec2{3, 4}
	res := v1.Add(v2)
	if res.X != 4 || res.Y != 6 {
		t.Errorf("expected {4, 6}, got %v", res)
	}
}

func TestVec2_Sub(t *testing.T) {
	v1 := Vec2{5, 6}
	v2 := Vec2{3, 4}
	res := v1.Sub(v2)
	if res.X != 2 || res.Y != 2 {
		t.Errorf("expected {2, 2}, got %v", res)
	}
}

func TestVec2_Scale(t *testing.T) {
	v := Vec2{2, 3}
	res := v.Scale(2)
	if res.X != 4 || res.Y != 6 {
		t.Errorf("expected {4, 6}, got %v", res)
	}
}

func TestVec2_Len(t *testing.T) {
	v := Vec2{3, 4}
	res := v.Len()
	if res != 5 {
		t.Errorf("expected 5, got %v", res)
	}
}

func TestVec2_Normalize(t *testing.T) {
	v := Vec2{3, 4}
	res := v.Normalize()
	if res.X != 0.6 || res.Y != 0.8 {
		t.Errorf("expected {0.6, 0.8}, got %v", res)
	}

	zero := Vec2{0, 0}
	resZero := zero.Normalize()
	if resZero.X != 0 || resZero.Y != 0 {
		t.Errorf("expected {0, 0}, got %v", resZero)
	}
}

func TestVec2_Lerp(t *testing.T) {
	v1 := Vec2{0, 0}
	v2 := Vec2{10, 10}
	res := v1.Lerp(v2, 0.5)
	if res.X != 5 || res.Y != 5 {
		t.Errorf("expected {5, 5}, got %v", res)
	}
}

func TestClamp(t *testing.T) {
	if res := Clamp(5, 0, 10); res != 5 {
		t.Errorf("expected 5, got %v", res)
	}
	if res := Clamp(-5, 0, 10); res != 0 {
		t.Errorf("expected 0, got %v", res)
	}
	if res := Clamp(15, 0, 10); res != 10 {
		t.Errorf("expected 10, got %v", res)
	}
}

func TestLerp(t *testing.T) {
	if res := Lerp(0, 10, 0.5); res != 5 {
		t.Errorf("expected 5, got %v", res)
	}
	if res := Lerp(0, 10, 0); res != 0 {
		t.Errorf("expected 0, got %v", res)
	}
	if res := Lerp(0, 10, 1); res != 10 {
		t.Errorf("expected 10, got %v", res)
	}
}

func TestAngleFromDir(t *testing.T) {
	dir1 := Vec2{1, 0}
	if math.Abs(AngleFromDir(dir1)-0) > 0.0001 {
		t.Errorf("expected 0, got %v", AngleFromDir(dir1))
	}

	dir2 := Vec2{0, 1}
	if math.Abs(AngleFromDir(dir2)-math.Pi/2) > 0.0001 {
		t.Errorf("expected Pi/2, got %v", AngleFromDir(dir2))
	}
}
