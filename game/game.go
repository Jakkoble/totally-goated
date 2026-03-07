package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	level     *Level
	cameraY   float64
	Goat      Goat
	state     GameState
	menuPulse float64
	tick      int
	menuGoat  *ebiten.Image
	score     float64
	bestScore float64
}

type GameState int

const (
	GameMenu    GameState = iota
	GamePlaying GameState = iota
	GameOver
)

func NewGame() *Game {
	img, _, _ := ebitenutil.NewImageFromFile("assets/goat.png")
	s := loadSave()
	return &Game{
		state:     GameMenu,
		menuGoat:  img,
		bestScore: s.BestScore,
	}
}

func (g *Game) startGame() {
	g.level = NewLevel()
	g.Goat = *NewGoat(0, 0)
	g.cameraY = 0
	g.score = 0
	g.state = GamePlaying
}

func (g *Game) Update() error {
	g.tick++
	g.menuPulse += 0.03

	switch g.state {
	case GameMenu:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startGame()
		}

	case GamePlaying:
		g.Goat.Update(g.level, g.cameraY)
		g.level.Update()

		height := -g.Goat.Pos.Y
		if height > g.score {
			g.score = height
		}

		targetY := g.Goat.Pos.Y
		if targetY < g.cameraY {
			lerpSpeed := 0.08
			g.cameraY += (targetY - g.cameraY) * lerpSpeed
		}

		deathY := g.cameraY + DeathMargin
		if g.Goat.Pos.Y > deathY {
			if g.score > g.bestScore {
				g.bestScore = g.score
				saveBest(g.bestScore)
			}
			g.state = GameOver
		}

		g.level.GenerateUntil(g.cameraY - 1000)

	case GameOver:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
			inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = GameMenu
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case GameMenu:
		screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255})
		g.drawMenu(screen)
	case GamePlaying:
		g.drawGame(screen)
	case GameOver:
		g.drawGame(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	sw := float64(ScreenWidth)
	sh := float64(ScreenHeight)

	title := []string{
		` _____ ___ _____ _   _    _  __   __`,
		`|_   _/ _ \_   _/ \ | |  | | \ \ / /`,
		`  | || | | || |/ _ \| |  | |  \ V / `,
		`  | || |_| || / ___ \ |__| |__ | |  `,
		`  |_| \___/ |_/_/  \_\____|___)|_|  `,
		`                                     `,
		`   ____  ___    _  _____ _____ ____  `,
		`  / ___|/ _ \  / \|_   _| ____|  _ \ `,
		` | |  _| | | |/ _ \ | | |  _| | | | |`,
		` | |_| | |_| / ___ \| | | |___| |_| |`,
		`  \____|\___/_/   \_\_| |_____|____/ `,
	}

	titleY := int(sh * 0.05)
	for i, line := range title {
		textW := len(line) * 6
		x := int(sw)/2 - textW/2
		ebitenutil.DebugPrintAt(screen, line, x, titleY+i*16)
	}

	if g.menuGoat != nil {
		op := &ebiten.DrawImageOptions{}
		iw := float64(g.menuGoat.Bounds().Dx())
		ih := float64(g.menuGoat.Bounds().Dy())
		scale := 3.0
		bob := math.Sin(g.menuPulse*2) * 6
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(sw/2-iw*scale/2, sh*0.42-ih*scale/2+bob)
		screen.DrawImage(g.menuGoat, op)
	}

	centerText := func(s string, y int) {
		x := int(sw)/2 - len(s)*6/2
		ebitenutil.DebugPrintAt(screen, s, x, y)
	}

	cy := int(sh * 0.6)
	centerText("CONTROLS:", cy)

	type ctrl struct{ key, desc string }
	controls := []ctrl{
		{"CLICK + AIM", "Charge & aim your horn dash"},
		{"RELEASE", "Launch!"},
		{"A/D", "Air control while flying"},
	}

	maxKeyW := 0
	for _, c := range controls {
		if w := len(c.key) * 6; w > maxKeyW {
			maxKeyW = w
		}
	}

	maxDescW := 0
	for _, c := range controls {
		if w := len(c.desc) * 6; w > maxDescW {
			maxDescW = w
		}
	}
	arrowW := 18
	gap := 6
	totalW := maxKeyW + gap + arrowW + maxDescW
	blockX := int(sw)/2 - totalW/2

	for i, c := range controls {
		y := cy + 20 + i*16
		keyW := len(c.key) * 6
		ebitenutil.DebugPrintAt(screen, c.key, blockX+maxKeyW-keyW, y)
		ebitenutil.DebugPrintAt(screen, "->", blockX+maxKeyW+gap, y)
		ebitenutil.DebugPrintAt(screen, c.desc, blockX+maxKeyW+gap+arrowW, y)
	}
}

func (g *Game) drawGame(screen *ebiten.Image) {
	g.level.Draw(screen, g.cameraY)
	g.Goat.Draw(screen, g.cameraY)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d m", g.currentMeters()), 10, 10)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(ScreenWidth), float32(ScreenHeight),
		color.RGBA{0, 0, 0, 150}, false)
	ebitenutil.DebugPrintAt(screen, "GAME OVER", ScreenWidth/2-36, ScreenHeight/2-30)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Height: %d m", g.currentMeters()), ScreenWidth/2-50, ScreenHeight/2-10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Best: %d m", g.bestMeters()), ScreenWidth/2-40, ScreenHeight/2+6)

	if g.tick%60 < 40 {
		ebitenutil.DebugPrintAt(screen, "Click or Space to continue", ScreenWidth/2-100, ScreenHeight/2+30)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func toMeters(px float64) int {
	return int(px / PixelsPerMeter)
}

func (g *Game) currentMeters() int {
	return toMeters(g.score)
}

func (g *Game) bestMeters() int {
	return toMeters(g.bestScore)
}
