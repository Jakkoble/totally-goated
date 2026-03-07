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

var (
	rockTiles    [10]*ebiten.Image
	bouncyTiles  [10]*ebiten.Image
	iceTiles     [10]*ebiten.Image
	crumblyTiles [10]*ebiten.Image
	stickyTiles  [10]*ebiten.Image
)

func init() {
	loadTileSet(rockTiles[:], "assets/rock_tile_%d.png")
	loadTileSet(bouncyTiles[:], "assets/bouncy_tile_%d.png")
	loadTileSet(iceTiles[:], "assets/ice_tile_%d.png")
	loadTileSet(crumblyTiles[:], "assets/crumbly_tile_%d.png")
	loadTileSet(stickyTiles[:], "assets/sticky_tile_%d.png")
}

func loadTileSet(tiles []*ebiten.Image, pattern string) {
	for i := range tiles {
		path := fmt.Sprintf(pattern, i+1)
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			if rockTiles[i] != nil {
				tiles[i] = rockTiles[i]
			} else {
				log.Fatalf("failed to load %s: %v", path, err)
			}
			continue
		}
		tiles[i] = img
	}
}

type PlatformType int

const (
	PlatNormal PlatformType = iota
	PlatBouncy
	PlatIce
	PlatCrumbly
	PlatSticky
)

type Platform struct {
	X, Y, W, H float64
	TileIndex  int
	Type       PlatformType

	CrumbleTimer   float64
	CrumbleStarted bool
	Destroyed      bool
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

		platType := PlatNormal
		if l.index > 5 {
			roll := rand.Float64()
			switch {
			case roll < 0.12:
				platType = PlatBouncy
			case roll < 0.24:
				platType = PlatIce
			case roll < 0.36:
				platType = PlatCrumbly
			case roll < 0.44:
				platType = PlatSticky
			}
		}

		p := Platform{
			X:         px,
			Y:         l.curY,
			W:         w,
			H:         h,
			TileIndex: tileIdx,
			Type:      platType,
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
		if p.Destroyed {
			continue
		}
		sx := p.X + offsetX
		sy := p.Y + offsetY

		var tile *ebiten.Image
		switch p.Type {
		case PlatBouncy:
			tile = bouncyTiles[p.TileIndex]
		case PlatIce:
			tile = iceTiles[p.TileIndex]
		case PlatCrumbly:
			tile = crumblyTiles[p.TileIndex]
		case PlatSticky:
			tile = stickyTiles[p.TileIndex]
		default:
			tile = rockTiles[p.TileIndex]
		}

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

func (l *Level) Update() {
	dt := 1.0 / 60.0
	for i := range l.Platforms {
		p := &l.Platforms[i]
		if p.Type == PlatCrumbly && p.CrumbleStarted && !p.Destroyed {
			p.CrumbleTimer += dt
			if p.CrumbleTimer >= CrumbleTime {
				p.Destroyed = true
			}
		}
	}
}
