package game

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var rockTiles [10]*ebiten.Image

func init() {
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("assets/rock_tile_%d.png", i+1)
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			log.Fatalf("failed to load %s: %v", path, err)
		}
		rockTiles[i] = img
	}
}

type Platform struct {
	X, Y, W, H float64
	TileIndex   int
}

type Level struct {
	Platforms []Platform
	curY      float64
	goRight   bool
	index     int
}

func NewLevel() *Level {
	l := &Level{
		Platforms: make([]Platform, 0, 300),
	}
	l.Platforms = append(l.Platforms, Platform{
		X: -35, Y: 0, W: 70, H: 120,
	})
	return l
}

func (l *Level) overlaps(p Platform) bool {
	for _, e := range l.Platforms {
		if p.X < e.X+e.W && p.X+p.W > e.X && p.Y < e.Y+e.H && p.Y+p.H > e.Y {
			return true
		}
	}
	return false
}

func (l *Level) GenerateUntil(targetY float64) {
	for l.curY > targetY {
		w := 40.0 + rand.Float64()*80.0
		h := 60.0 + rand.Float64()*70.0
		tileIdx := rand.IntN(10)

		// Vertical gap: -30 to +30 px between consecutive platforms
		gap := -30.0 + rand.Float64()*60.0
		l.curY -= h + gap

		spread := float64(l.index) * 0.1
		base := 130.0 + spread
		jitterX := (rand.Float64() - 0.5) * 30.0

		var px float64
		if l.goRight {
			px = base - w/2 + jitterX
		} else {
			px = -base - w/2 + jitterX
		}

		p := Platform{
			X:         px,
			Y:         l.curY,
			W:         w,
			H:         h,
			TileIndex: tileIdx,
		}

		if !l.overlaps(p) {
			l.Platforms = append(l.Platforms, p)
		}

		l.goRight = !l.goRight
		l.index++
	}
}

func (l *Level) Draw(screen *ebiten.Image, cameraY float64) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 50, A: 255})

	offsetX := float64(ScreenWidth) / 2
	offsetY := float64(ScreenHeight)/2 - cameraY

	for _, p := range l.Platforms {
		sx := p.X + offsetX
		sy := p.Y + offsetY

		tile := rockTiles[p.TileIndex]
		tw, th := tile.Bounds().Dx(), tile.Bounds().Dy()

		cropW := int(p.W)
		cropH := int(p.H)
		if cropW > tw {
			cropW = tw
		}
		if cropH > th {
			cropH = th
		}
		cropX := (tw - cropW) / 2
		sub := tile.SubImage(image.Rect(cropX, 0, cropX+cropW, cropH)).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(sx, sy)
		screen.DrawImage(sub, op)
	}
}
