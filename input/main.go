package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var img *ebiten.Image
var pos [2]float64 // x, y

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("tunasprite.png")
	if err != nil {
		panic(err)
	}
}

var heldKeys map[ebiten.Key]int

type Game struct{}

func (g *Game) Update() error {
	// init in the map.
	_, ok := heldKeys[ebiten.KeyArrowLeft]
	if !ok {
		heldKeys[ebiten.KeyArrowLeft] = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		// set in held keys state
		heldKeys[ebiten.KeyArrowLeft] += 1 // unsafe
		// pos[0] -= 2
	} else {
		// reset state
		heldKeys[ebiten.KeyArrowLeft] = 0 // delete?
	}
	// long press?
	if heldKeys[ebiten.KeyArrowLeft] > 120 {
		pos[0] -= 20
		heldKeys[ebiten.KeyArrowLeft] = 0
	} else if heldKeys[ebiten.KeyArrowLeft] > 30 {
		pos[0] -= 2
		heldKeys[ebiten.KeyArrowLeft] = 0
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		pos[0] += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		pos[1] -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		pos[1] += 2
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pos[0], pos[1])
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	heldKeys = map[ebiten.Key]int{}
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
