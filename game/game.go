package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	level   *Level
	cameraY float64
	Goat    Goat
	state   GameState
}

type GameState int

const (
	GamePlaying GameState = iota
	GameOver
)

func NewGame() *Game {
	return &Game{
		level: NewLevel(),
		Goat:  *NewGoat(0, 0),
		state: GamePlaying,
	}
}

func (g *Game) Update() error {
	if g.state == GameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*g = *NewGame()
		}
		return nil
	}

	g.Goat.Update(g.level, g.cameraY)

	targetY := g.Goat.Pos.Y
	if targetY < g.cameraY {
		lerpSpeed := 0.08
		g.cameraY += (targetY - g.cameraY) * lerpSpeed
	}

	deathY := g.cameraY + DeathMargin
	if g.Goat.Pos.Y > deathY {
		g.state = GameOver
	}

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

	if g.state == GameOver {
		vector.FillRect(screen, 0, 0, float32(ScreenWidth), float32(ScreenHeight),
			color.RGBA{0, 0, 0, 150}, false)
		ebitenutil.DebugPrintAt(screen, "GAME OVER\nClick or Space to restart",
			ScreenWidth/2-80, ScreenHeight/2-10)
	}
}

func drawBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
