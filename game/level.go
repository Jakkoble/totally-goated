package game

import "math/rand/v2"

type Platform struct {
	X, Y, W, H float64
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

func (l *Level) GenerateUntil(targetY float64) {
	for l.curY > targetY {
		rise := 50.0 + float64(l.index)*0.5
		l.curY -= rise

		w := 55.0
		h := 65.0

		spread := float64(l.index) * 0.1
		var px float64
		base := 150.0 + spread
		jitterX := (rand.Float64() - 0.5) * 30.0
		if l.goRight {
			px = base - w/2 + jitterX
		} else {
			px = -base - w/2 + jitterX
		}

		jitterY := (rand.Float64() - 0.5) * 20

		l.Platforms = append(l.Platforms, Platform{
			X: px,
			Y: l.curY + jitterY,
			W: w,
			H: h,
		})

		l.goRight = !l.goRight
		l.index++
	}
}
