package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const TileWidth = 16

type Sprite struct {
	Img       *ebiten.Image
	X, Y      float64
	Direction int
	Frame     int
}

type Player struct {
	*Sprite
	Health uint
}

type Enemy struct {
	*Sprite
	FollowsPlayer bool
}

type Potion struct {
	*Sprite
	AmtHeal  uint
	Consumed bool
}

type Game struct {
	player      *Player
	enemies     []*Enemy
	potions     []*Potion
	tilemapJSON *TilemapJSON
	tilemapImg  *ebiten.Image
}

func (g *Game) Update() error {
	// react to key presses

	playerMoves := true
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.Direction = 3
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.Direction = 2
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Direction = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Direction = 0
	} else {
		playerMoves = false
	}
	if playerMoves {
		g.player.Move(2)
	}

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer && enemy.Dist(g.player.Sprite) < 5*TileWidth {
			if enemy.X < g.player.X {
				enemy.Direction = 3
			} else if enemy.X > g.player.X {
				enemy.Direction = 2
			} else if enemy.Y > g.player.Y {
				enemy.Direction = 1
			} else if enemy.Y < g.player.Y {
				enemy.Direction = 0
			}

		} else {
			if rand.Float64() > 0.95 {
				enemy.Direction = rand.Intn(4)
			}
		}
		enemy.Move(1)
	}

	for _, potion := range g.potions {
		if !potion.Consumed && g.player.X > potion.X {
			g.player.Health += potion.AmtHeal
			potion.Consumed = true
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x66, 0x99, 255})

	opts := &ebiten.DrawImageOptions{}
	for _, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			x := TileWidth * (index % layer.Width)
			y := TileWidth * (index / layer.Width)

			srcX := TileWidth * ((id - 1) % 22)
			srcY := TileWidth * ((id - 1) / 22)

			opts.GeoM.Translate(float64(x), float64(y))

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+TileWidth, srcY+TileWidth)).(*ebiten.Image),
				opts,
			)
			opts.GeoM.Reset()
		}
	}

	// draw player
	DrawSprite(g.player.Sprite, screen)

	for _, sprite := range g.enemies {
		DrawSprite(sprite.Sprite, screen)
	}

	for _, potion := range g.potions {
		if !potion.Consumed {
			DrawSprite(potion.Sprite, screen)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d \n", g.player.Health))
}

func (s *Sprite) Dist(other *Sprite) float64 {
	return math.Sqrt(math.Pow(s.X-other.X, 2) + math.Pow(s.Y-other.Y, 2))
}

func (s *Sprite) Move(d float64) {
	switch s.Direction {
	case 0:
		s.Y += d
	case 1:
		s.Y -= d
	case 2:
		s.X -= d
	case 3:
		s.X += d
	}
	s.Frame = (s.Frame + 1) % 256
}

func DrawSprite(sprite *Sprite, screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sprite.X, sprite.Y)

	dx := TileWidth * sprite.Direction
	fy := TileWidth * (sprite.Frame / 4 % 4)

	screen.DrawImage(
		sprite.Img.SubImage(
			image.Rect(dx, fy, TileWidth+dx, TileWidth+fy),
		).(*ebiten.Image),
		opts,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/samurai.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/LifePot.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/TilesetFloor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTilemapJSON("assets/maps/first.tmj")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting...")
	game := &Game{
		player: &Player{
			Sprite: &Sprite{
				Img: playerImg,
				X:   100.0,
				Y:   100.0,
			},
			Health: 3,
		},
		enemies: []*Enemy{
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   50.0,
					Y:   50.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   150.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &Sprite{
					Img: skeletonImg,
					X:   75.0,
					Y:   75.0,
				},
				FollowsPlayer: true,
			},
		},
		potions: []*Potion{
			{
				Sprite: &Sprite{
					Img: potionImg,
					X:   210.0,
					Y:   100.0,
				},
				AmtHeal: 1,
			},
		},
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok, Bye!")
}
