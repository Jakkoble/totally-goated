package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type CosmeticType string

const (
	CosmeticHat   CosmeticType = "Hat"
	CosmeticShoes CosmeticType = "Shoes"
)

type Rarity int

const (
	RarityCommon Rarity = iota
	RarityRare
	RarityEpic
	RarityLegendary
)

func (r Rarity) String() string {
	switch r {
	case RarityCommon:
		return "Common"
	case RarityRare:
		return "Rare"
	case RarityEpic:
		return "Epic"
	case RarityLegendary:
		return "Legendary"
	}
	return "Unknown"
}

func (r Rarity) Color() color.NRGBA {
	switch r {
	case RarityCommon:
		return color.NRGBA{200, 200, 200, 255} // Gray
	case RarityRare:
		return color.NRGBA{50, 150, 255, 255} // Blue
	case RarityEpic:
		return color.NRGBA{200, 50, 255, 255} // Purple
	case RarityLegendary:
		return color.NRGBA{255, 215, 0, 255} // Gold
	}
	return color.NRGBA{255, 255, 255, 255}
}

type Cosmetic struct {
	ID     int
	Name   string
	Type   CosmeticType
	Rarity Rarity
}

var AllCosmetics = []Cosmetic{
	{1, "Top Hat", CosmeticHat, RarityCommon},
	{2, "Crown", CosmeticHat, RarityEpic},
	{3, "Beanie", CosmeticHat, RarityCommon},
	{4, "Halo", CosmeticHat, RarityLegendary},
	{5, "Sneakers", CosmeticShoes, RarityCommon},
	{6, "Boots", CosmeticShoes, RarityRare},
	{7, "Golden Hooves", CosmeticShoes, RarityLegendary},
}

func getCosmeticByID(id int) *Cosmetic {
	for i := range AllCosmetics {
		if AllCosmetics[i].ID == id {
			return &AllCosmetics[i]
		}
	}
	return nil
}

// DrawCosmetics draws the equipped cosmetics onto the given screen using the exact transform applied to the goat.
// This is called from Goat.Draw.
func DrawCosmetics(screen *ebiten.Image, equipped map[string]int, transform ebiten.GeoM) {
	if hatID, ok := equipped[string(CosmeticHat)]; ok {
		drawHat(screen, hatID, transform)
	}
	if shoesID, ok := equipped[string(CosmeticShoes)]; ok {
		drawShoes(screen, shoesID, transform)
	}
}

func drawHat(screen *ebiten.Image, id int, transform ebiten.GeoM) {
	// The goat sprite is 29x47. The head rests between x=14 and x=28, and the top of the head is at y=12.
	// We translate by -14.5, -23.5. So the top of the head is at y = 12 - 23.5 = -11.5.
	// The center of the head is around x = 21 - 14.5 = 6.5.

	switch id {
	case 1: // Top Hat
		drawRectTransformed(screen, 0, -14, 15, 3, color.RGBA{30, 30, 30, 255}, transform) // Brim
		drawRectTransformed(screen, 2, -24, 11, 10, color.RGBA{30, 30, 30, 255}, transform) // Top
		drawRectTransformed(screen, 2, -16, 11, 2, color.RGBA{200, 50, 50, 255}, transform) // Red band
	case 2: // Crown
		// Base of crown
		drawRectTransformed(screen, 2, -15, 11, 4, color.RGBA{255, 215, 0, 255}, transform)
		// Points
		drawRectTransformed(screen, 2, -19, 3, 4, color.RGBA{255, 215, 0, 255}, transform)
		drawRectTransformed(screen, 6, -20, 3, 5, color.RGBA{255, 215, 0, 255}, transform)
		drawRectTransformed(screen, 10, -19, 3, 4, color.RGBA{255, 215, 0, 255}, transform)
	case 3: // Beanie
		drawRectTransformed(screen, 1, -19, 13, 8, color.RGBA{50, 100, 200, 255}, transform)
		// Pom pom
		drawRectTransformed(screen, 5, -23, 5, 5, color.RGBA{200, 200, 200, 255}, transform)
	case 4: // Halo
		drawRectTransformed(screen, 0, -20, 15, 2, color.RGBA{255, 255, 100, 200}, transform)
		drawRectTransformed(screen, 0, -19, 2, 2, color.RGBA{255, 255, 100, 200}, transform)
		drawRectTransformed(screen, 13, -19, 2, 2, color.RGBA{255, 255, 100, 200}, transform)
	}
}

func drawShoes(screen *ebiten.Image, id int, transform ebiten.GeoM) {
	// The goat sprite is 29x47. The legs are roughly at x=4, 9, 14, 19 and go down to y=47.
	// We translate by -14.5, -23.5. So the bottom of the feet are at y = 47 - 23.5 = 23.5.
	// The legs are at x = 4-14.5=-10.5, 9-14.5=-5.5, 14-14.5=-0.5, 19-14.5=4.5.

	switch id {
	case 5: // Sneakers
		drawRectTransformed(screen, -11, 20, 7, 3, color.RGBA{200, 50, 50, 255}, transform) // Back legs (covers 2 back legs)
		drawRectTransformed(screen, -11, 23, 7, 1, color.RGBA{255, 255, 255, 255}, transform)

		drawRectTransformed(screen, -1, 20, 7, 3, color.RGBA{200, 50, 50, 255}, transform) // Front legs (covers 2 front legs)
		drawRectTransformed(screen, -1, 23, 7, 1, color.RGBA{255, 255, 255, 255}, transform)
	case 6: // Boots
		drawRectTransformed(screen, -11, 16, 7, 8, color.RGBA{100, 50, 20, 255}, transform) // Back legs
		drawRectTransformed(screen, -1, 16, 7, 8, color.RGBA{100, 50, 20, 255}, transform) // Front legs
	case 7: // Golden Hooves
		drawRectTransformed(screen, -11, 21, 7, 3, color.RGBA{255, 215, 0, 255}, transform) // Back legs
		drawRectTransformed(screen, -1, 21, 7, 3, color.RGBA{255, 215, 0, 255}, transform) // Front legs
	}
}

var whiteImage *ebiten.Image

func init() {
	whiteImage = ebiten.NewImage(1, 1)
	whiteImage.Fill(color.White)
}

// Helper to draw a rectangle after applying a GeoM transformation
func drawRectTransformed(screen *ebiten.Image, x, y, w, h float32, clr color.Color, transform ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w), float64(h))
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Concat(transform)
	op.ColorScale.ScaleWithColor(clr)

	screen.DrawImage(whiteImage, op)
}

// RollGamblingBox rolls a random rarity and selects an item from it.
func RollGamblingBox() Cosmetic {
	roll := rand.Float64()
	var rarity Rarity

	// Weights: 60% Common, 25% Rare, 10% Epic, 5% Legendary
	if roll < 0.60 {
		rarity = RarityCommon
	} else if roll < 0.85 {
		rarity = RarityRare
	} else if roll < 0.95 {
		rarity = RarityEpic
	} else {
		rarity = RarityLegendary
	}

	var candidates []Cosmetic
	for _, c := range AllCosmetics {
		if c.Rarity == rarity {
			candidates = append(candidates, c)
		}
	}

	if len(candidates) == 0 {
		// Fallback to any item if no items of that rarity exist
		return AllCosmetics[rand.Intn(len(AllCosmetics))]
	}

	return candidates[rand.Intn(len(candidates))]
}
