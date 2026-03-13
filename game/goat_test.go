package game

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// setupMockGame creates a basic mocked game environment for testing
func setupMockGame() *Game {
	// Create dummy images so we don't need real assets
	dummyImg := ebiten.NewImage(32, 32)
	dummyImg.Fill(color.White)

	goat := &Goat{
		Pos:             Vec2{0, 0},
		FacingDir:       1,
		WallPlatIdx:     -1,
		LastDashPlatIdx: -1,
		Image:           dummyImg,
		DashSpeedMod:    1.0,
		SquashX:         1.0,
		SquashY:         1.0,
	}

	game := &Game{
		level: &Level{
			Platforms: []Platform{},
			PowerUps:  []PowerUp{},
			Bells:     []Bell{},
		},
		Goat: *goat,
	}
	return game
}

func TestGoat_ChargingPercentage(t *testing.T) {
	dummyImg := ebiten.NewImage(32, 32)
	g := &Goat{
		Image:      dummyImg,
		ChargeTime: 0.5,
	}

	pct := g.ChargingPercentage()
	expected := 0.5 / FullChargeTime
	if pct != expected {
		t.Errorf("expected %v, got %v", expected, pct)
	}

	g.ChargeTime = FullChargeTime + 1.0
	if g.ChargingPercentage() != 1.0 {
		t.Errorf("expected max percentage 1.0, got %v", g.ChargingPercentage())
	}
}

func TestGoat_Gravity(t *testing.T) {
	game := setupMockGame()
	goat := &game.Goat
	goat.State = StateAir
	goat.Vel = Vec2{0, 0}

	goat.updateAir(game.level, game)

	if goat.Vel.Y != Gravity {
		t.Errorf("expected Y velocity to equal gravity %v, got %v", Gravity, goat.Vel.Y)
	}
}

func TestGoat_SlowFallGravity(t *testing.T) {
	game := setupMockGame()
	goat := &game.Goat
	goat.State = StateAir
	goat.Vel = Vec2{0, 0}
	goat.SlowFallTimer = 1.0 // Active slow fall

	// Needs some downward velocity for slow fall to kick in
	goat.Vel.Y = 1.0

	initialVelY := goat.Vel.Y
	goat.updateAir(game.level, game)

	expectedGrav := Gravity * SlowFallGravityMul
	if goat.Vel.Y != initialVelY+expectedGrav {
		t.Errorf("expected Y velocity %v, got %v", initialVelY+expectedGrav, goat.Vel.Y)
	}
}

func TestGoat_CollectBells(t *testing.T) {
	game := setupMockGame()
	goat := &game.Goat
	goat.Pos = Vec2{100, 100}

	game.level.Bells = []Bell{
		{Pos: Vec2{100, 100}, Collected: false}, // Same pos as goat
		{Pos: Vec2{500, 500}, Collected: false}, // Far away
	}

	goat.collectBells(game.level, game)

	if !game.level.Bells[0].Collected {
		t.Errorf("expected bell 0 to be collected")
	}
	if game.level.Bells[1].Collected {
		t.Errorf("expected bell 1 to NOT be collected")
	}
}

func TestGoat_ResolveCollisions(t *testing.T) {
	game := setupMockGame()
	goat := &game.Goat
	goat.Pos = Vec2{50, -5}
	goat.Vel = Vec2{0, 10} // Falling down

	_ = float64(goat.Image.Bounds().Dx())
	_ = float64(goat.Image.Bounds().Dy())

	// A platform right below the goat.
	// We want the bounding boxes to intersect so resolution happens.
	// Goat bounding box before resolution:
	// L: 50 - 16 = 34
	// R: 50 + 16 = 66
	// T: -5 - 16 = -21
	// B: -5 + 16 = 11
	// If platform is at Y=10, then B (11) > p.Y (10). Collision!
	game.level.Platforms = []Platform{
		{X: 0, Y: 10, W: 100, H: 20, Type: PlatNormal, Destroyed: false},
	}

	goat.resolveCollisions(game.level, game)

	// Because of attachment logic during top collisions (which are actually wall slides downwards in this game)
	// it pushes goat down + attaches to wall based on distance.
	// We check that it at least modified the Y or attached to a wall.
	if goat.State != StateWall {
		t.Errorf("expected goat to attach to wall after top impact, got state %v", goat.State)
	}

	// Because distance to left side is 50-0=50, and distance to right side is 100-50=50,
	// it will prefer right wall attachment since distance left < distance right check is `<` or `<=`.
	// Goat moves to X = 0 - gw/2 = -16, attaches to WallRight.
	// Let's just check that it moved out of the platform.
	if goat.Pos.X > game.level.Platforms[0].X && goat.Pos.X < game.level.Platforms[0].X+game.level.Platforms[0].W {
		t.Errorf("expected goat to be moved horizontally out of platform, got X=%v", goat.Pos.X)
	}
}

// Since init loading assets is tricky in isolated tests, we mock what we need.
func init() {
	// Provide a dummy initial audio context if not already done,
	// though Ebiten's audio usually handles lack of init gracefully in headless.
}
