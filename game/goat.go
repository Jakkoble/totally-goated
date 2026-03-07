package game

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GoatState int

const (
	StateAir GoatState = iota
	StateWall
	StateCharging
)

type WallSide int

const (
	WallNone WallSide = iota
	WallLeft
	WallRight
)

type Goat struct {
	Pos          Vec2
	Vel          Vec2
	State        GoatState
	Wall         WallSide
	WallPlatIdx  int
	ClingTimer   int
	FacingDir    float64
	Image        *ebiten.Image
	DashSpeedMod float64

	ChargeTime float64
}

func NewGoat(x, y float64) *Goat {
	goatImage, _, err := ebitenutil.NewImageFromFile("assets/goat.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Goat{
		Pos:          Vec2{x, y},
		FacingDir:    1,
		WallPlatIdx:  -1,
		Image:        goatImage,
		DashSpeedMod: 1.0,
	}
}

func (g *Goat) Update(level *Level, cameraY float64) {
	switch g.State {
	case StateAir:
		g.updateAir(level)
	case StateWall:
		g.updateWall(level)
	case StateCharging:
		g.updateCharging(level, cameraY)
	}
}

func (g *Goat) updateAir(level *Level) {
	g.Vel.Y += Gravity
	if g.Vel.Y > MaxFallSpeed {
		g.Vel.Y = MaxFallSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.Vel.X -= AirControlAccel
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.Vel.X += AirControlAccel
	}

	g.Vel.X *= DashDrag

	g.Pos = g.Pos.Add(g.Vel)
	g.resolveCollisions(level)
}

func (g *Goat) updateWall(level *Level) {
	if g.WallPlatIdx >= 0 && g.WallPlatIdx < len(level.Platforms) {
		if level.Platforms[g.WallPlatIdx].Destroyed {
			g.detachFromWall()
			return
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.State = StateCharging
		g.ChargeTime = 0
		return
	}

	p := &level.Platforms[g.WallPlatIdx]

	slideSpeed := WallSlideSpeed
	if p.Type == PlatIce {
		slideSpeed = IceSlideSpeed
	}

	g.Vel.Y = slideSpeed
	g.Vel.X = 0
	g.Pos.Y += g.Vel.Y
	g.ClingTimer++

	g.checkStillOnWall(level)
}

func (g *Goat) updateCharging(level *Level, cameraY float64) {
	if g.WallPlatIdx >= 0 && g.WallPlatIdx < len(level.Platforms) {
		if level.Platforms[g.WallPlatIdx].Destroyed {
			g.detachFromWall()
			return
		}
	}
	dt := 1.0 / float64(ebiten.TPS())
	g.ChargeTime += dt
	if g.ChargeTime > FullChargeTime {
		g.ChargeTime = FullChargeTime
	}

	slideSpeed := WallSlideSpeed
	if g.WallPlatIdx >= 0 && g.WallPlatIdx < len(level.Platforms) {
		if level.Platforms[g.WallPlatIdx].Type == PlatIce {
			slideSpeed = IceSlideSpeed
		}
	}
	g.Pos.Y += slideSpeed
	g.ClingTimer++

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		worldX := float64(mx) - float64(ScreenWidth)/2
		worldY := float64(my) - float64(ScreenHeight)/2 + cameraY
		dir := (Vec2{worldX, worldY}).Sub(g.Pos).Normalize()

		speed := Lerp(DashMinSpeed, DashMaxSpeed, g.ChargingPercentage())
		speed *= g.DashSpeedMod
		pct := g.ChargingPercentage()
		switch {
		case pct > 0.8:
			sfxDashPower15.Play()
		case pct > 0.4:
			sfxDashPower10.Play()
		default:
			sfxDashPower5.Play()
		}
		sfxChargeRelease.Play()

		g.Vel = dir.Scale(speed)
		g.State = StateAir
		g.Wall = WallNone
		g.WallPlatIdx = -1
		return
	}

	g.checkStillOnWall(level)
}

func (g *Goat) checkStillOnWall(level *Level) {
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

func (g *Goat) resolveCollisions(level *Level) {
	gw := float64(g.Image.Bounds().Dx())
	gh := float64(g.Image.Bounds().Dy())

	for idx := range level.Platforms {
		p := &level.Platforms[idx]

		if p.Destroyed {
			continue
		}

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
				g.attachToWall(WallRight, idx, level)
			} else {
				g.Pos.X = p.X + p.W + gw/2
				g.Pos.Y = p.Y + gh/2
				g.attachToWall(WallLeft, idx, level)
			}
			return

		case minPen == penL:
			g.Pos.X = p.X - gw/2
			g.attachToWall(WallRight, idx, level)
			return

		case minPen == penR:
			g.Pos.X = p.X + p.W + gw/2
			g.attachToWall(WallLeft, idx, level)
			return

		case minPen == penB:
			g.Pos.Y = p.Y + p.H + gh/2
			g.Vel.Y = math.Abs(g.Vel.Y) * 0.15
		}
	}
}

func (g *Goat) attachToWall(side WallSide, platIdx int, level *Level) {
	p := &level.Platforms[platIdx]

	if p.Type == PlatBouncy {
		sfxBounce.Play()
		if side == WallLeft || side == WallRight {
			g.Vel.X = -g.Vel.X * BouncyReflect
		}
		if g.Vel.Y > 0 {
			g.Vel.Y = -g.Vel.Y * BouncyReflect
		}
		g.State = StateAir
		g.Wall = WallNone
		g.WallPlatIdx = -1
		return
	}

	switch p.Type {
	case PlatIce:
		sfxIce.Play()
	case PlatCrumbly:
		if !p.CrumbleStarted {
			p.CrumbleStarted = true
			sfxCrumble.Play()
		}
		sfxWallHit.Play()
	default:
		sfxWallHit.Play()
	}

	g.DashSpeedMod = 1.0
	if p.Type == PlatSticky {
		g.DashSpeedMod = StickyDashMult
	}

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
	offsetX := float64(ScreenWidth) / 2
	offsetY := float64(ScreenHeight)/2 - cameraY

	if g.State == StateCharging {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Charge %d%%", int(g.ChargingPercentage()*100)), 0, 0)
		g.drawDashAimLine(screen, offsetX, offsetY, cameraY)
		g.drawDashTrajectory(screen, offsetX, offsetY, cameraY)
		g.drawChargeCircle(screen, offsetX, offsetY)
		g.drawChargeBar(screen, offsetX, offsetY)
	}

	op := &ebiten.DrawImageOptions{}
	iw := float64(g.Image.Bounds().Dx())
	ih := float64(g.Image.Bounds().Dy())
	if g.FacingDir < 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(iw, 0)
	}
	op.GeoM.Translate(g.Pos.X-iw/2+offsetX, g.Pos.Y-ih/2+offsetY)
	screen.DrawImage(g.Image, op)
}

func (g *Goat) ChargingPercentage() float64 {
	return Clamp(g.ChargeTime/FullChargeTime, 0, 1)
}

func (g *Goat) drawChargeBar(screen *ebiten.Image, offsetX, offsetY float64) {
	pct := g.ChargingPercentage()

	barW := float32(60)
	barH := float32(5)
	barX := float32(g.Pos.X+offsetX) - barW/2
	barY := float32(g.Pos.Y+offsetY) + float32(g.Image.Bounds().Dy())/2 + 10

	vector.FillRect(screen, barX, barY, barW, barH, color.RGBA{50, 50, 50, 200}, false)

	fillColor := chargingColor(pct)
	vector.FillRect(screen, barX, barY, barW*float32(pct), barH, fillColor, false)
	vector.StrokeRect(screen, barX, barY, barW, barH, 1, color.RGBA{255, 255, 255, 255}, false)
}

func (g *Goat) drawChargeCircle(screen *ebiten.Image, offsetX, offsetY float64) {
	pct := g.ChargingPercentage()

	cx := float32(g.Pos.X + offsetX)
	cy := float32(g.Pos.Y + offsetY)

	fillColor := chargingColor(pct)
	fillColor.A = 100

	bounds := g.Image.Bounds()
	maxDimension := max(bounds.Dx(), bounds.Dy())
	radius := Lerp(float64(maxDimension)*.8, float64(maxDimension)*1.5, pct)

	vector.FillCircle(screen, cx, cy, float32(radius), fillColor, true)
}

func chargingColor(pct float64) color.NRGBA {
	r := uint8(pct * 255)
	g := uint8((1 - pct) * 255)
	return color.NRGBA{r, g, 0, 255}
}

func (g *Goat) drawDashAimLine(screen *ebiten.Image, offsetX, offsetY, cameraY float64) {
	mx, my := ebiten.CursorPosition()
	worldX := float64(mx) - float64(ScreenWidth)/2
	worldY := float64(my) - float64(ScreenHeight)/2 + cameraY

	dir := (Vec2{worldX, worldY}).Sub(g.Pos).Normalize()
	if dir.Len() < 0.0001 {
		if g.Wall == WallLeft {
			dir = Vec2{1, -0.35}.Normalize()
		} else {
			dir = Vec2{-1, -0.35}.Normalize()
		}
	}

	power := g.ChargingPercentage()
	lineLen := 40.0 + power*70.0

	sx := g.Pos.X + offsetX
	sy := g.Pos.Y + offsetY
	endX := sx + dir.X*lineLen
	endY := sy + dir.Y*lineLen

	for i := 0; i < 8; i++ {
		t1 := float64(i) / 8.0
		t2 := (float64(i) + 0.5) / 8.0
		x1 := sx + (endX-sx)*t1
		y1 := sy + (endY-sy)*t1
		x2 := sx + (endX-sx)*t2
		y2 := sy + (endY-sy)*t2

		r := uint8(255)
		gc := uint8(Lerp(255, 80, power))
		b := uint8(Lerp(200, 30, power))

		vector.StrokeLine(
			screen,
			float32(x1), float32(y1),
			float32(x2), float32(y2),
			float32(2+power*2),
			color.RGBA{r, gc, b, 200},
			true,
		)
	}

	arrowSz := 6.0 + power*5.0
	perp := Vec2{-dir.Y, dir.X}
	ax1x := endX + (-dir.X+perp.X*0.5)*arrowSz
	ax1y := endY + (-dir.Y+perp.Y*0.5)*arrowSz
	ax2x := endX + (-dir.X-perp.X*0.5)*arrowSz
	ax2y := endY + (-dir.Y-perp.Y*0.5)*arrowSz

	vector.StrokeLine(screen, float32(endX), float32(endY), float32(ax1x), float32(ax1y),
		float32(2+power*2), color.RGBA{255, 200, 50, 255}, true)
	vector.StrokeLine(screen, float32(endX), float32(endY), float32(ax2x), float32(ax2y),
		float32(2+power*2), color.RGBA{255, 200, 50, 255}, true)
}

func (g *Goat) drawDashTrajectory(screen *ebiten.Image, offsetX, offsetY, cameraY float64) {
	mx, my := ebiten.CursorPosition()
	worldX := float64(mx) - float64(ScreenWidth)/2
	worldY := float64(my) - float64(ScreenHeight)/2 + cameraY

	dir := (Vec2{worldX, worldY}).Sub(g.Pos).Normalize()
	if dir.Len() < 0.0001 {
		if g.Wall == WallLeft {
			dir = Vec2{1, -0.35}.Normalize()
		} else {
			dir = Vec2{-1, -0.35}.Normalize()
		}
	}

	speed := Lerp(DashMinSpeed, DashMaxSpeed, g.ChargingPercentage())
	speed *= g.DashSpeedMod
	vel := dir.Scale(speed)
	pos := g.Pos

	for i := 0; i < 40; i++ {
		pos = pos.Add(vel)
		vel = vel.Scale(DashDrag)

		if vel.Len() < 3.0 {
			vel.Y += Gravity
		}
		if i%3 == 0 {
			sx := pos.X + offsetX
			sy := pos.Y + offsetY
			alpha := uint8(Clamp(float64(120-i*3), 20, 120))
			sz := float32(Clamp(3.0-float64(i)*0.05, 1, 3))
			vector.FillCircle(screen, float32(sx), float32(sy), sz, color.RGBA{255, 255, 255, alpha}, true)
		}
	}
}
