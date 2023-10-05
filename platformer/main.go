package main

import (
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	Frames       []image.Image // Note: simple but ruins any efficency you might get from tilesheets.
	CurrentFrame int
	Speed        int // Note: not all sprites loop, this doesn't account for that.
}

type Game struct {
	CurrentFrame int
	Sprites      map[string]*Sprite // Note: this is efficient but means all sprites on screen have same timing.
	Player       *Player
}

type PlayerState int

const (
	PlayerStateIdle PlayerState = iota
	PlayerStateRunning
	PlayerStateShooting
)

type FacingDirection int

const (
	FacingDirectionRight FacingDirection = iota
	FacingDirectionLeft
)

type Player struct {
	State  PlayerState
	Facing FacingDirection
	X      float64
	Y      float64 // Note: in a real platformer, you would want the bottom of the screen to be y == 0, not the top.
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.Player.X--
		g.Player.Facing = FacingDirectionLeft
		g.Player.State = PlayerStateRunning
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.Player.X++
		g.Player.Facing = FacingDirectionRight
		g.Player.State = PlayerStateRunning
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player.State = PlayerStateShooting
	} else {
		g.Player.State = PlayerStateIdle
	}
	// Move game and sprites to next frame.
	g.CurrentFrame++
	for _, sprite := range g.Sprites {
		if g.CurrentFrame%sprite.Speed == 0 {
			sprite.CurrentFrame++
			if sprite.CurrentFrame >= len(sprite.Frames) {
				sprite.CurrentFrame = 0
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Flip image.
	if g.Player.Facing == FacingDirectionLeft {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(16), 0)
	}
	// Change position - note that the anchor here is top left.
	op.GeoM.Translate(g.Player.X, g.Player.Y)
	spriteName := "player_idle"
	switch g.Player.State {
	case PlayerStateRunning:
		spriteName = "player_running"
	case PlayerStateShooting:
		spriteName = "player_shooting"
	}
	sprite := g.Sprites[spriteName]
	screen.DrawImage(sprite.Frames[sprite.CurrentFrame].(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 256, 240
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("tilemap_packed.png")
	if err != nil {
		panic(err)
	}
	// Split the image into a 2D array of tiles, to make access easier later on.
	tiles := [10][6]image.Image{}
	for x := 0; x < len(tiles); x++ {
		for y := 0; y < len(tiles[x]); y++ {
			tiles[x][y] = img.SubImage(image.Rectangle{
				Min: image.Point{x * 16, y * 16},
				Max: image.Point{(x * 16) + 16, (y * 16) + 16},
			})
		}
	}
	// Create named sprites, often re-using existing tiles.
	sprites := map[string]*Sprite{
		"player_running": {
			Frames: []image.Image{tiles[0][4], tiles[1][4]},
			Speed:  10,
		},
		"player_idle": {
			Frames: []image.Image{tiles[0][4]},
			Speed:  1,
		},
		"player_shooting": {
			Frames: []image.Image{tiles[0][4], tiles[2][4]},
			Speed:  5,
		},
	}
	// Double window size for sanity's sake.
	ebiten.SetWindowSize(768, 720)
	game := &Game{
		Sprites: sprites,
		Player: &Player{
			X: 50,
			Y: 190,
		},
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
