package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	level   *Level
	cameraY float64
	Goat    Goat
}

func NewGame() *Game {
	return &Game{
		level: NewLevel(),
		Goat:  *NewGoat(0, 0),
	}
}

func (g *Game) Update() error {
	g.Goat.Update(g.level, g.cameraY)

	lerpSpeed := 0.08
	g.cameraY += (g.Goat.Pos.Y - g.cameraY) * lerpSpeed

	g.level.GenerateUntil(g.cameraY - 1000)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen)

	offsetX := float64(ScreenWidth) / 2
	offsetY := float64(ScreenHeight)/2 - g.cameraY

	for _, p := range g.level.Platforms {
		sx := float32(p.X + offsetX)
		sy := float32(p.Y + offsetY)
		vector.FillRect(screen, sx, sy, float32(p.W), float32(p.H), color.White, false)
	}

	g.Goat.Draw(screen, g.cameraY)
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
