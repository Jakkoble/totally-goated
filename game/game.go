package game

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	_ "image/png"
)

type Game struct {
	level      *Level
	cameraY    float64
	Goat       Goat
	state      GameState
	menuPulse  float64
	tick       int
	menuGoat   *ebiten.Image
	score      float64
	bestScore  int
	bestMeters int
	bellCount  int
	bestBells  int
	background *ebiten.Image
	stars      *Starfield
	speedLines *SpeedLines
	shakeMag   float64
	shakeX     float64
	shakeY     float64
}

type GameState int

const (
	GameMenu GameState = iota
	GamePlaying
	GameOver
)

func NewGame() *Game {
	img := loadImageFromFS("assets/textures/goat.png")
	background := loadImageFromFS("assets/textures/background.png")

	s := loadSave()
	return &Game{
		state:      GameMenu,
		menuGoat:   img,
		bestScore:  s.BestScore,
		bestMeters: s.BestMeters,
		bestBells:  s.BestBells,
		background: background,
		stars:      NewStarfield(),
		speedLines: &SpeedLines{},
	}
}

func (g *Game) startGame() {
	g.level = NewLevel()
	g.Goat = *NewGoat(0, 0)
	g.cameraY = 0
	g.score = 0
	g.bellCount = 0
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
		g.Goat.Update(g.level, g.cameraY, g)
		g.level.Update()
		g.speedLines.Update(g.Goat.Vel)
		g.updateShake()

		height := -g.Goat.Pos.Y
		if height > g.score {
			oldMeters := int(g.score / PixelsPerMeter)
			g.score = height
			newMeters := int(g.score / PixelsPerMeter)
			if newMeters/100 > oldMeters/100 {
				sfxMilestone.Play()
			}
		}

		targetY := g.Goat.Pos.Y
		if targetY < g.cameraY {
			lerpSpeed := 0.08
			g.cameraY += (targetY - g.cameraY) * lerpSpeed
		}

		deathY := g.cameraY + DeathMargin
		if g.Goat.Pos.Y > deathY {
			if g.Goat.HasShield {
				g.Goat.HasShield = false
				sfxShieldBreak.Play()

				g.Goat.Pos.Y = g.cameraY
				g.Goat.Pos.X = 0
				g.Goat.Vel = Vec2{0, 0}
				g.Goat.State = StateCharging
				g.Goat.ChargeTime = 0
				g.Goat.SlowFallTimer = 5.0
			} else {
				sfxDeath.Play()
				total := g.totalScore()
				meters := g.currentMeters()
				saveBest(total, meters, g.bellCount)
				if total > g.bestScore {
					g.bestScore = total
				}
				if meters > g.bestMeters {
					g.bestMeters = meters
				}
				if g.bellCount > g.bestBells {
					g.bestBells = g.bellCount
				}
				g.state = GameOver
			}
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
	screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255})
	g.stars.Draw(screen, g.tick)
	g.drawBackground(screen)

	switch g.state {
	case GameMenu:
		g.drawMenu(screen)
	case GamePlaying:
		g.drawGame(screen)
	case GameOver:
		g.drawGame(screen)
		g.drawGameOver(screen)
	}
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	if g.background == nil {
		return
	}
	imgW := float64(g.background.Bounds().Dx())
	imgH := float64(g.background.Bounds().Dy())
	scale := float64(ScreenWidth) / imgW
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(0, float64(ScreenHeight)-imgH*scale)
	screen.DrawImage(g.background, op)
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

	credit := "Jakob Schwendinger & Jakob Wassertheurer"
	creditX := int(sw)/2 - len(credit)*6/2
	ebitenutil.DebugPrintAt(screen, credit, creditX, int(sh)-24)
}

func (g *Game) AddShake(amount float64) {
	g.shakeMag = math.Min(g.shakeMag+amount, 16.0)
}

func (g *Game) updateShake() {
	if g.shakeMag > 0.2 {
		g.shakeX = math.Sin(float64(g.tick)*1.1)*g.shakeMag + math.Cos(float64(g.tick)*2.3)*g.shakeMag*0.5
		g.shakeY = math.Cos(float64(g.tick)*1.7)*g.shakeMag + math.Sin(float64(g.tick)*3.1)*g.shakeMag*0.5
		g.shakeMag *= 0.88
	} else {
		g.shakeX = 0
		g.shakeY = 0
		g.shakeMag = 0
	}
}

func (g *Game) drawGame(screen *ebiten.Image) {
	g.level.Draw(screen, g.cameraY, g.shakeX, g.shakeY)
	g.Goat.Draw(screen, g.cameraY, g.shakeX, g.shakeY)
	g.speedLines.Draw(screen)
	drawHUD(screen, g)
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, float32(ScreenWidth), float32(ScreenHeight),
		color.RGBA{0, 0, 0, 180}, false)

	white := color.NRGBA{255, 255, 255, 220}
	gold := color.NRGBA{255, 210, 50, 220}
	dim := color.NRGBA{180, 180, 180, 180}

	cx := float64(ScreenWidth) / 2
	cy := float64(ScreenHeight)/2 - 20

	title := "GAME OVER"
	drawScaledText(screen, title, cx-float64(len(title))*6*1.5/2, cy-90, 1.5, color.NRGBA{255, 80, 80, 240})

	vector.FillRect(screen, float32(cx-100), float32(cy-55), 200, 1, color.NRGBA{255, 255, 255, 60}, true)

	heightStr := fmt.Sprintf("%d m", g.currentMeters())
	bellsStr := fmt.Sprintf("%d bells", g.bellCount)

	drawScaledText(screen, "Height:", cx-120, cy-40, 1.0, dim)
	drawScaledText(screen, heightStr, cx+20, cy-40, 1.0, white)

	drawScaledText(screen, "Bells:", cx-120, cy-20, 1.0, dim)
	drawScaledText(screen, bellsStr, cx+20, cy-20, 1.0, gold)

	vector.FillRect(screen, float32(cx-100), float32(cy+1), 200, 1, color.NRGBA{255, 255, 255, 60}, true)

	totalStr := fmt.Sprintf("%d", g.totalScore())
	drawScaledText(screen, "SCORE:", cx-120, cy+14, 1.3, white)
	drawScaledText(screen, totalStr, cx+20, cy+14, 1.3, white)

	vector.FillRect(screen, float32(cx-100), float32(cy+42), 200, 1, color.NRGBA{255, 255, 255, 40}, true)

	bestScoreStr := fmt.Sprintf("%d", g.bestScore)
	bestMetersStr := fmt.Sprintf("%d m", g.bestMeters)
	bestBellsStr := fmt.Sprintf("%d", g.bestBells)

	drawScaledText(screen, "BEST", cx-120, cy+52, 1.0, dim)
	drawScaledText(screen, "Score:", cx-120, cy+70, 1.0, dim)
	drawScaledText(screen, bestScoreStr, cx+20, cy+70, 1.0, dim)
	drawScaledText(screen, "Height:", cx-120, cy+86, 1.0, dim)
	drawScaledText(screen, bestMetersStr, cx+20, cy+86, 1.0, dim)
	drawScaledText(screen, "Bells:", cx-120, cy+102, 1.0, dim)
	drawScaledText(screen, bestBellsStr, cx+20, cy+102, 1.0, gold)

	if g.tick%60 < 40 {
		prompt := "Click or Space to continue"
		pw := float64(len(prompt)) * 6
		drawScaledText(screen, prompt, cx-pw/2, cy+130, 1.0, dim)
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

func (g *Game) totalScore() int {
	return g.currentMeters() + g.bellCount*BellScoreValue
}
