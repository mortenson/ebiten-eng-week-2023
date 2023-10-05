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

type Game struct{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		pos[0]--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		pos[0]++
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		pos[1]--
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		pos[1]++
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
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
