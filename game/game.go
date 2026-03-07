package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	level   *Level
	cameraY float64
}

func NewGame() *Game {
	return &Game{level: NewLevel()}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.cameraY -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.cameraY += 5
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	offsetX := float64(ScreenWidth) / 2
	offsetY := float64(ScreenHeight)/2 - g.cameraY

	for _, p := range g.level.Platforms {
		sx := float32(p.X + offsetX)
		sy := float32(p.Y + offsetY)
		vector.FillRect(screen, sx, sy, float32(p.W), float32(p.H), color.White, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
