package game

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GoatState int

const (
	StateAir GoatState = iota
	StateWall
	StateDashing
)

type WallSide int

const (
	WallNone WallSide = iota
	WallLeft
	WallRight
)

type Goat struct {
	Pos         Vec2
	Vel         Vec2
	State       GoatState
	Wall        WallSide
	WallPlatIdx int
	ClingTimer  int
	FacingDir   float64
	Image       *ebiten.Image
}

func NewGoat(x, y float64) *Goat {
	goatImage, _, err := ebitenutil.NewImageFromFile("assets/goat.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Goat{
		Pos:         Vec2{x, y},
		FacingDir:   1,
		WallPlatIdx: -1,
		Image:       goatImage,
	}
}

func (g *Goat) Update(level *Level, cameraY float64) {
	switch g.State {
	case StateAir:
		g.updateAir(level)
	case StateWall:
		g.updateWall(level, cameraY)
	case StateDashing:
		g.updateDashing(level)
	}
}

func (g *Goat) updateAir(level *Level) {
	g.Vel.Y += Gravity
	if g.Vel.Y > MaxFallSpeed {
		g.Vel.Y = MaxFallSpeed
	}

	g.Pos = g.Pos.Add(g.Vel)
	g.resolveCollisions(level)
}

func (g *Goat) updateWall(level *Level, cameraY float64) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		worldX := float64(mx) - float64(ScreenWidth)/2
		worldY := float64(my) - float64(ScreenHeight)/2 + cameraY
		dir := (Vec2{worldX, worldY}).Sub(g.Pos).Normalize()

		g.Vel = dir.Scale(DashSpeed)
		g.State = StateDashing
		g.Wall = WallNone
		g.WallPlatIdx = -1
		return
	}

	g.Vel.Y = WallSlideSpeed
	g.Vel.X = 0
	g.Pos.Y += g.Vel.Y
	g.ClingTimer++

	gh := float64(g.Image.Bounds().Dy())
	if g.WallPlatIdx >= 0 && g.WallPlatIdx < len(level.Platforms) {
		p := &level.Platforms[g.WallPlatIdx]
		stillOn := g.Pos.Y+gh/2 >= p.Y && g.Pos.Y-gh/2 <= p.Y+p.H
		if !stillOn {
			g.detachFromWall()
		}
	} else {
		g.detachFromWall()
	}
}

func (g *Goat) updateDashing(level *Level) {
	g.Pos = g.Pos.Add(g.Vel)
	g.Vel = g.Vel.Scale(DashDrag)

	if g.Vel.Len() < 3.0 {
		g.State = StateAir
	}

	g.resolveCollisions(level)
}

func (g *Goat) resolveCollisions(level *Level) {
	gw := float64(g.Image.Bounds().Dx())
	gh := float64(g.Image.Bounds().Dy())

	for idx := range level.Platforms {
		p := &level.Platforms[idx]

		goatL := g.Pos.X - gw/2
		goatR := g.Pos.X + gw/2
		goatT := g.Pos.Y - gh/2
		goatB := g.Pos.Y + gh/2

		if goatR <= p.X || goatL >= p.X+p.W || goatB <= p.Y || goatT >= p.Y+p.H {
			continue
		}

		penL := goatR - p.X
		penR := (p.X + p.W) - goatL
		penT := goatB - p.Y
		penB := (p.Y + p.H) - goatT

		minPen := math.Min(math.Min(penL, penR), math.Min(penT, penB))

		switch {
		case minPen == penT && g.Vel.Y >= 0:
			g.Pos.Y = p.Y - gh/2

			distL := g.Pos.X - p.X
			distR := (p.X + p.W) - g.Pos.X
			if distL < distR {
				g.Pos.X = p.X - gw/2
				g.Pos.Y = p.Y + gh/2
				g.attachToWall(WallRight, idx)
			} else {
				g.Pos.X = p.X + p.W + gw/2
				g.Pos.Y = p.Y + gh/2
				g.attachToWall(WallLeft, idx)
			}
			return

		case minPen == penL:
			g.Pos.X = p.X - gw/2
			g.attachToWall(WallRight, idx)
			return

		case minPen == penR:
			g.Pos.X = p.X + p.W + gw/2
			g.attachToWall(WallLeft, idx)
			return

		case minPen == penB:
			g.Pos.Y = p.Y + p.H + gh/2
			g.Vel.Y = math.Abs(g.Vel.Y) * 0.15
		}
	}
}

func (g *Goat) attachToWall(side WallSide, platIdx int) {
	g.State = StateWall
	g.Wall = side
	g.WallPlatIdx = platIdx
	g.Vel = Vec2{}
	g.ClingTimer = 0

	if side == WallLeft {
		g.FacingDir = -1
	} else {
		g.FacingDir = 1
	}
}

func (g *Goat) detachFromWall() {
	g.State = StateAir
	g.Wall = WallNone
	g.WallPlatIdx = -1
}

func (g *Goat) Draw(screen *ebiten.Image, cameraY float64) {
	op := &ebiten.DrawImageOptions{}
	offsetX := float64(ScreenWidth) / 2
	offsetY := float64(ScreenHeight)/2 - cameraY
	iw := float64(g.Image.Bounds().Dx())
	ih := float64(g.Image.Bounds().Dy())
	op.GeoM.Translate(g.Pos.X-iw/2+offsetX, g.Pos.Y-ih/2+offsetY)
	screen.DrawImage(g.Image, op)
}
