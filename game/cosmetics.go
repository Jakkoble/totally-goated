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
	// Hats
	{1, "Top Hat", CosmeticHat, RarityCommon},
	{2, "Crown", CosmeticHat, RarityEpic},
	{3, "Beanie", CosmeticHat, RarityCommon},
	{4, "Halo", CosmeticHat, RarityLegendary},
	{8, "Propeller Hat", CosmeticHat, RarityEpic},
	{9, "Viking Helmet", CosmeticHat, RarityRare},
	{10, "Chef Hat", CosmeticHat, RarityCommon},
	{11, "Wizard Hat", CosmeticHat, RarityEpic},
	{12, "Cowboy Hat", CosmeticHat, RarityRare},
	{13, "Pirate Hat", CosmeticHat, RarityEpic},
	{14, "Party Hat", CosmeticHat, RarityCommon},
	{15, "King's Crown", CosmeticHat, RarityLegendary},
	{16, "Sombrero", CosmeticHat, RarityRare},
	{17, "Baseball Cap", CosmeticHat, RarityCommon},
	{18, "Devil Horns", CosmeticHat, RarityEpic},

	// Shoes
	{5, "Sneakers", CosmeticShoes, RarityCommon},
	{6, "Boots", CosmeticShoes, RarityRare},
	{7, "Golden Hooves", CosmeticShoes, RarityLegendary},
	{19, "Roller Skates", CosmeticShoes, RarityEpic},
	{20, "Slippers", CosmeticShoes, RarityCommon},
	{21, "Clown Shoes", CosmeticShoes, RarityRare},
	{22, "Ice Skates", CosmeticShoes, RarityRare},
	{23, "Flippers", CosmeticShoes, RarityEpic},
	{24, "Bunny Slippers", CosmeticShoes, RarityEpic},
	{25, "Rocket Boots", CosmeticShoes, RarityLegendary},
	{26, "Glass Slippers", CosmeticShoes, RarityLegendary},
	{27, "High Heels", CosmeticShoes, RarityRare},
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
	// All coordinates are top-left relative. So (0,0) is top-left of the 29x47 image space.
	// The center of the head is around x = 21.

	switch id {
	case 1: // Top Hat
		drawRectTransformed(screen, 14.5, 9.5, 15, 3, color.RGBA{30, 30, 30, 255}, transform) // Brim
		drawRectTransformed(screen, 16.5, -0.5, 11, 10, color.RGBA{30, 30, 30, 255}, transform) // Top
		drawRectTransformed(screen, 16.5, 7.5, 11, 2, color.RGBA{200, 50, 50, 255}, transform) // Red band
	case 2: // Crown
		// Base of crown
		drawRectTransformed(screen, 16.5, 8.5, 11, 4, color.RGBA{255, 215, 0, 255}, transform)
		// Points
		drawRectTransformed(screen, 16.5, 4.5, 3, 4, color.RGBA{255, 215, 0, 255}, transform)
		drawRectTransformed(screen, 20.5, 3.5, 3, 5, color.RGBA{255, 215, 0, 255}, transform)
		drawRectTransformed(screen, 24.5, 4.5, 3, 4, color.RGBA{255, 215, 0, 255}, transform)
	case 3: // Beanie
		drawRectTransformed(screen, 15.5, 4.5, 13, 8, color.RGBA{50, 100, 200, 255}, transform)
		// Pom pom
		drawRectTransformed(screen, 19.5, 0.5, 5, 5, color.RGBA{200, 200, 200, 255}, transform)
	case 4: // Halo
		drawRectTransformed(screen, 14.5, 3.5, 15, 2, color.RGBA{255, 255, 100, 200}, transform)
		drawRectTransformed(screen, 14.5, 4.5, 2, 2, color.RGBA{255, 255, 100, 200}, transform)
		drawRectTransformed(screen, 27.5, 4.5, 2, 2, color.RGBA{255, 255, 100, 200}, transform)
	case 8: // Propeller Hat
		drawRectTransformed(screen, 16.5, 7.5, 11, 4, color.RGBA{50, 200, 50, 255}, transform)
		drawRectTransformed(screen, 21.5, 4.5, 2, 3, color.RGBA{150, 150, 150, 255}, transform) // Pin
		drawRectTransformed(screen, 16.5, 3.5, 12, 1, color.RGBA{200, 50, 50, 255}, transform) // Blades
	case 9: // Viking Helmet
		drawRectTransformed(screen, 16.5, 6.5, 11, 5, color.RGBA{150, 150, 150, 255}, transform) // Helmet
		drawRectTransformed(screen, 14.5, 2.5, 2, 6, color.RGBA{220, 220, 200, 255}, transform) // Left Horn
		drawRectTransformed(screen, 27.5, 2.5, 2, 6, color.RGBA{220, 220, 200, 255}, transform) // Right Horn
	case 10: // Chef Hat
		drawRectTransformed(screen, 15.5, 5.5, 13, 6, color.RGBA{255, 255, 255, 255}, transform)
		drawRectTransformed(screen, 14.5, 1.5, 15, 6, color.RGBA{255, 255, 255, 255}, transform) // Poof
	case 11: // Wizard Hat
		drawRectTransformed(screen, 13.5, 9.5, 17, 2, color.RGBA{50, 50, 200, 255}, transform) // Brim
		drawRectTransformed(screen, 16.5, 4.5, 11, 5, color.RGBA{50, 50, 200, 255}, transform) // Base
		drawRectTransformed(screen, 18.5, 0.5, 7, 4, color.RGBA{50, 50, 200, 255}, transform) // Mid
		drawRectTransformed(screen, 20.5, -4.5, 3, 5, color.RGBA{50, 50, 200, 255}, transform) // Tip
		drawRectTransformed(screen, 19.5, 5.5, 2, 2, color.RGBA{255, 215, 0, 255}, transform) // Star
	case 12: // Cowboy Hat
		drawRectTransformed(screen, 12.5, 9.5, 19, 2, color.RGBA{139, 69, 19, 255}, transform) // Brim
		drawRectTransformed(screen, 15.5, 4.5, 13, 5, color.RGBA{139, 69, 19, 255}, transform) // Top
		drawRectTransformed(screen, 15.5, 7.5, 13, 1, color.RGBA{50, 50, 50, 255}, transform) // Band
	case 13: // Pirate Hat
		drawRectTransformed(screen, 13.5, 7.5, 17, 4, color.RGBA{40, 40, 40, 255}, transform)
		drawRectTransformed(screen, 15.5, 3.5, 13, 4, color.RGBA{40, 40, 40, 255}, transform)
		drawRectTransformed(screen, 20.5, 5.5, 3, 3, color.RGBA{255, 255, 255, 255}, transform) // Skull
	case 14: // Party Hat
		drawRectTransformed(screen, 16.5, 7.5, 11, 4, color.RGBA{200, 50, 200, 255}, transform)
		drawRectTransformed(screen, 18.5, 3.5, 7, 4, color.RGBA{50, 200, 200, 255}, transform)
		drawRectTransformed(screen, 20.5, 0.5, 3, 3, color.RGBA{200, 200, 50, 255}, transform)
		drawRectTransformed(screen, 20.5, -2.5, 3, 3, color.RGBA{255, 255, 255, 255}, transform) // Pom pom
	case 15: // King's Crown
		drawRectTransformed(screen, 15.5, 8.5, 13, 3, color.RGBA{255, 215, 0, 255}, transform) // Base
		drawRectTransformed(screen, 15.5, 4.5, 3, 4, color.RGBA{255, 215, 0, 255}, transform) // Point 1
		drawRectTransformed(screen, 20.5, 2.5, 3, 6, color.RGBA{255, 215, 0, 255}, transform) // Point 2
		drawRectTransformed(screen, 25.5, 4.5, 3, 4, color.RGBA{255, 215, 0, 255}, transform) // Point 3
		drawRectTransformed(screen, 16.0, 3.0, 2, 2, color.RGBA{255, 50, 50, 255}, transform) // Gem 1
		drawRectTransformed(screen, 21.0, 1.0, 2, 2, color.RGBA{50, 255, 50, 255}, transform) // Gem 2
		drawRectTransformed(screen, 26.0, 3.0, 2, 2, color.RGBA{50, 50, 255, 255}, transform) // Gem 3
		drawRectTransformed(screen, 16.5, 6.5, 11, 2, color.RGBA{200, 0, 0, 255}, transform) // Red cushion
	case 16: // Sombrero
		drawRectTransformed(screen, 10.5, 9.5, 23, 2, color.RGBA{210, 180, 140, 255}, transform) // Big brim
		drawRectTransformed(screen, 15.5, 6.5, 13, 3, color.RGBA{210, 180, 140, 255}, transform) // Top base
		drawRectTransformed(screen, 17.5, 3.5, 9, 3, color.RGBA{210, 180, 140, 255}, transform) // Top point
		drawRectTransformed(screen, 15.5, 7.5, 13, 1, color.RGBA{200, 50, 50, 255}, transform) // Red strip
	case 17: // Baseball Cap
		drawRectTransformed(screen, 16.5, 7.5, 11, 4, color.RGBA{50, 50, 200, 255}, transform) // Cap
		drawRectTransformed(screen, 24.5, 9.5, 6, 2, color.RGBA{50, 50, 200, 255}, transform) // Brim facing right
	case 18: // Devil Horns
		drawRectTransformed(screen, 16.5, 9.5, 11, 2, color.RGBA{200, 50, 50, 255}, transform) // Headband
		drawRectTransformed(screen, 16.5, 4.5, 2, 5, color.RGBA{255, 0, 0, 255}, transform) // Left horn
		drawRectTransformed(screen, 25.5, 4.5, 2, 5, color.RGBA{255, 0, 0, 255}, transform) // Right horn
	}
}

func drawShoes(screen *ebiten.Image, id int, transform ebiten.GeoM) {
	// The goat sprite is 29x47. The legs are roughly at x=4, 9, 14, 19 and go down to y=47.
	// All coordinates are top-left relative.

	switch id {
	case 5: // Sneakers
		drawRectTransformed(screen, 3.5, 43.5, 7, 3, color.RGBA{200, 50, 50, 255}, transform) // Back legs
		drawRectTransformed(screen, 3.5, 46.5, 7, 1, color.RGBA{255, 255, 255, 255}, transform)

		drawRectTransformed(screen, 13.5, 43.5, 7, 3, color.RGBA{200, 50, 50, 255}, transform) // Front legs
		drawRectTransformed(screen, 13.5, 46.5, 7, 1, color.RGBA{255, 255, 255, 255}, transform)
	case 6: // Boots
		drawRectTransformed(screen, 3.5, 39.5, 7, 8, color.RGBA{100, 50, 20, 255}, transform) // Back legs
		drawRectTransformed(screen, 13.5, 39.5, 7, 8, color.RGBA{100, 50, 20, 255}, transform) // Front legs
	case 7: // Golden Hooves
		drawRectTransformed(screen, 3.5, 44.5, 7, 3, color.RGBA{255, 215, 0, 255}, transform) // Back legs
		drawRectTransformed(screen, 13.5, 44.5, 7, 3, color.RGBA{255, 215, 0, 255}, transform) // Front legs
	case 19: // Roller Skates
		drawRectTransformed(screen, 2.5, 43.5, 9, 3, color.RGBA{200, 200, 200, 255}, transform) // Boot
		drawRectTransformed(screen, 12.5, 43.5, 9, 3, color.RGBA{200, 200, 200, 255}, transform)
		drawRectTransformed(screen, 3.5, 46.5, 3, 3, color.RGBA{255, 100, 100, 255}, transform) // Wheel
		drawRectTransformed(screen, 8.5, 46.5, 3, 3, color.RGBA{255, 100, 100, 255}, transform) // Wheel
		drawRectTransformed(screen, 13.5, 46.5, 3, 3, color.RGBA{255, 100, 100, 255}, transform) // Wheel
		drawRectTransformed(screen, 18.5, 46.5, 3, 3, color.RGBA{255, 100, 100, 255}, transform) // Wheel
	case 20: // Slippers
		drawRectTransformed(screen, 2.5, 45.5, 9, 2, color.RGBA{150, 100, 200, 255}, transform)
		drawRectTransformed(screen, 12.5, 45.5, 9, 2, color.RGBA{150, 100, 200, 255}, transform)
	case 21: // Clown Shoes
		drawRectTransformed(screen, 2.5, 43.5, 11, 4, color.RGBA{255, 50, 50, 255}, transform) // Back legs
		drawRectTransformed(screen, 12.5, 43.5, 11, 4, color.RGBA{255, 50, 50, 255}, transform) // Front legs
		drawRectTransformed(screen, 10.5, 42.5, 3, 3, color.RGBA{255, 255, 50, 255}, transform) // Bulb
		drawRectTransformed(screen, 20.5, 42.5, 3, 3, color.RGBA{255, 255, 50, 255}, transform) // Bulb
	case 22: // Ice Skates
		drawRectTransformed(screen, 3.5, 43.5, 7, 2, color.RGBA{50, 50, 50, 255}, transform) // Boot
		drawRectTransformed(screen, 13.5, 43.5, 7, 2, color.RGBA{50, 50, 50, 255}, transform)
		drawRectTransformed(screen, 2.5, 46.5, 9, 1, color.RGBA{200, 200, 255, 255}, transform) // Blade
		drawRectTransformed(screen, 12.5, 46.5, 9, 1, color.RGBA{200, 200, 255, 255}, transform)
		drawRectTransformed(screen, 6.5, 45.5, 1, 1, color.RGBA{200, 200, 200, 255}, transform) // Connect
		drawRectTransformed(screen, 16.5, 45.5, 1, 1, color.RGBA{200, 200, 200, 255}, transform) // Connect
	case 23: // Flippers
		drawRectTransformed(screen, 1.5, 45.5, 12, 2, color.RGBA{50, 200, 100, 255}, transform)
		drawRectTransformed(screen, 11.5, 45.5, 12, 2, color.RGBA{50, 200, 100, 255}, transform)
	case 24: // Bunny Slippers
		drawRectTransformed(screen, 2.5, 44.5, 9, 3, color.RGBA{255, 180, 200, 255}, transform)
		drawRectTransformed(screen, 12.5, 44.5, 9, 3, color.RGBA{255, 180, 200, 255}, transform)
		drawRectTransformed(screen, 8.5, 42.5, 2, 4, color.RGBA{255, 180, 200, 255}, transform) // Ears
		drawRectTransformed(screen, 18.5, 42.5, 2, 4, color.RGBA{255, 180, 200, 255}, transform)
	case 25: // Rocket Boots
		drawRectTransformed(screen, 3.5, 42.5, 7, 5, color.RGBA{150, 150, 150, 255}, transform) // Boot
		drawRectTransformed(screen, 13.5, 42.5, 7, 5, color.RGBA{150, 150, 150, 255}, transform)
		drawRectTransformed(screen, 4.5, 47.5, 5, 2, color.RGBA{255, 150, 0, 255}, transform) // Flame
		drawRectTransformed(screen, 14.5, 47.5, 5, 2, color.RGBA{255, 150, 0, 255}, transform)
		drawRectTransformed(screen, 5.5, 49.5, 3, 2, color.RGBA{255, 50, 0, 255}, transform) // Flame tip
		drawRectTransformed(screen, 15.5, 49.5, 3, 2, color.RGBA{255, 50, 0, 255}, transform)
	case 26: // Glass Slippers
		drawRectTransformed(screen, 3.5, 45.5, 7, 2, color.RGBA{200, 255, 255, 150}, transform)
		drawRectTransformed(screen, 13.5, 45.5, 7, 2, color.RGBA{200, 255, 255, 150}, transform)
		drawRectTransformed(screen, 4.5, 43.5, 3, 2, color.RGBA{200, 255, 255, 150}, transform) // Heel
		drawRectTransformed(screen, 14.5, 43.5, 3, 2, color.RGBA{200, 255, 255, 150}, transform)
	case 27: // High Heels
		drawRectTransformed(screen, 3.5, 46.5, 7, 1, color.RGBA{255, 50, 50, 255}, transform) // Sole
		drawRectTransformed(screen, 13.5, 46.5, 7, 1, color.RGBA{255, 50, 50, 255}, transform)
		drawRectTransformed(screen, 3.5, 43.5, 2, 3, color.RGBA{255, 50, 50, 255}, transform) // Spike
		drawRectTransformed(screen, 13.5, 43.5, 2, 3, color.RGBA{255, 50, 50, 255}, transform)
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

	// Weights: 40% Common, 30% Rare, 20% Epic, 10% Legendary
	if roll < 0.40 {
		rarity = RarityCommon
	} else if roll < 0.70 {
		rarity = RarityRare
	} else if roll < 0.90 {
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
