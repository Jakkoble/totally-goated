package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PowerUpType int

const (
	PowerUpShield PowerUpType = iota
	PowerUpSuperDash
	PowerUpSlowFall
)

const (
	PowerUpRadius      = 14.0
	PowerUpSpawnChance = 0.12
	SlowFallDuration   = 5.0
	SuperDashMult      = 1.5
	SlowFallGravityMul = 0.45
)

type PowerUp struct {
	Pos       Vec2
	Type      PowerUpType
	Collected bool
	bobPhase  float64
}

func (pu *PowerUp) Update() {
	pu.bobPhase += 0.05
}

func (pu *PowerUp) Draw(screen *ebiten.Image, cameraY, shakeX, shakeY float64) {
	if pu.Collected {
		return
	}
	ox := float64(ScreenWidth)/2 + shakeX
	oy := float64(ScreenHeight)/2 - cameraY + shakeY

	sx := float32(pu.Pos.X + ox)
	sy := float32(pu.Pos.Y + oy + math.Sin(pu.bobPhase)*4)

	var fill color.NRGBA
	var symbol string
	switch pu.Type {
	case PowerUpShield:
		fill = color.NRGBA{80, 160, 255, 200}
		symbol = "S"
	case PowerUpSuperDash:
		fill = color.NRGBA{255, 100, 60, 200}
		symbol = "D"
	case PowerUpSlowFall:
		fill = color.NRGBA{100, 230, 120, 200}
		symbol = "F"
	}

	glow := fill
	glow.A = 60
	vector.FillCircle(screen, sx, sy, float32(PowerUpRadius+4), glow, true)

	vector.FillCircle(screen, sx, sy, float32(PowerUpRadius), fill, true)
	vector.StrokeCircle(screen, sx, sy, float32(PowerUpRadius), 2, color.NRGBA{255, 255, 255, 180}, true)

	ebitenutil.DebugPrintAt(screen, symbol, int(sx)-3, int(sy)-8)
}
