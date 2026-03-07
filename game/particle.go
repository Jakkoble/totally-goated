package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Particle struct {
	Pos  Vec2
	Vel  Vec2
	Life float64
}

type ParticleSystem struct {
	Particles []Particle
}

func (ps *ParticleSystem) Emit(pos, vel Vec2) {
	if vel.Len() < 3 || rand.Float64() > 0.4 {
		return
	}
	ps.Particles = append(ps.Particles, Particle{
		Pos: Vec2{
			pos.X + (rand.Float64()-0.5)*6,
			pos.Y + (rand.Float64()-0.5)*6,
		},
		Vel:  Vec2{(rand.Float64() - 0.5) * 0.5, -rand.Float64() * 0.3},
		Life: 0.4 + rand.Float64()*0.3,
	})
}

func (ps *ParticleSystem) Update() {
	dt := 1.0 / float64(ebiten.TPS())
	n := 0
	for i := range ps.Particles {
		p := &ps.Particles[i]
		p.Life -= dt
		if p.Life <= 0 {
			continue
		}
		p.Pos = p.Pos.Add(p.Vel)
		ps.Particles[n] = *p
		n++
	}
	ps.Particles = ps.Particles[:n]
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, cameraY float64) {
	ox := float64(ScreenWidth) / 2
	oy := float64(ScreenHeight)/2 - cameraY

	for i := range ps.Particles {
		p := &ps.Particles[i]
		a := uint8(p.Life * 300)
		vector.FillCircle(screen,
			float32(p.Pos.X+ox), float32(p.Pos.Y+oy),
			float32(1+p.Life*3),
			color.NRGBA{220, 200, 170, a}, true)
	}
}
