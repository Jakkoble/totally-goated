package game

import "math/rand/v2"

type Platform struct {
	X, Y, W, H float64
}

type Level struct {
	Platforms []Platform
	Height    float64
}

func NewLevel() *Level {
	l := &Level{Platforms: make([]Platform, 0, 300)}
	l.generate()
	return l
}

func (l *Level) generate() {
	l.Platforms = append(l.Platforms, Platform{
		X: -35, Y: 0, W: 70, H: 120,
	})

	curY := 0.0
	goRight := true

	for i := 0; i < 200; i++ {
		rise := 50.0 + float64(i)*0.5
		curY -= rise

		w := 55.0
		h := 65.0

		spread := float64(i) * 0.1
		var px float64
		base := 150.0 + spread
		jitterX := (rand.Float64() - 0.5) * 30.0
		if goRight {
			px = base - w/2 + jitterX
		} else {
			px = -base - w/2 + jitterX
		}

		jitterY := (rand.Float64() - 0.5) * 20

		l.Platforms = append(l.Platforms, Platform{
			X: px,
			Y: curY + jitterY,
			W: w,
			H: h,
		})

		goRight = !goRight
	}

	l.Height = curY - 200
}
