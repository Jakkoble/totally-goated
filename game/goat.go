package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GoatState int

const (
	StateStill GoatState = iota
	StateDashing
)

type Goat struct {
	Pos   Vec2
	Vel   Vec2
	State GoatState
	Image *ebiten.Image
}

func NewGoat(x, y float64) *Goat {
	goatImage, _, err := ebitenutil.NewImageFromFile("assets/goat.png")
	if err != nil {
		log.Fatal(err)
	}

	return &Goat{
		Pos:   Vec2{x, y},
		Vel:   Vec2{0, 0},
		Image: goatImage,
	}
}

func (g *Goat) Update() {
	dt := 1.0 / float64(ebiten.TPS())
	g.Vel.Y += Gravity * dt
	g.Pos = g.Pos.Add(g.Vel.Scale(dt))

	floorY := float64(ScreenHeight - g.Image.Bounds().Dy())
	if g.Pos.Y > floorY {
		g.Pos.Y = floorY
		g.Vel.Y = 0
		g.State = StateStill
	}

	if g.State == StateStill && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		target := Vec2{float64(mx), float64(my)}

		dir := target.Sub(g.Pos).Normalize()
		g.Vel = dir.Scale(500)
		g.State = StateDashing
	}
}

func (g *Goat) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.Pos.X, g.Pos.Y)
	screen.DrawImage(g.Image, op)
}
