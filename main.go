package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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
	player  *Player
	enemies []*Enemy
	potions []*Potion
}

func (g *Game) Update() error {
	// react to key presses

	g.player.Frame += 1
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.X += 2
		g.player.Direction = 3
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.X -= 2
		g.player.Direction = 2
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Y -= 2
		g.player.Direction = 1
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Y += 2
		g.player.Direction = 0
	} else {
		g.player.Frame -= 1
	}
	g.player.Frame %= 256

	for _, enemy := range g.enemies {
		if enemy.FollowsPlayer {
			enemy.Frame += 1
			if enemy.X < g.player.X {
				enemy.X += 1
				enemy.Direction = 3
			} else if enemy.X > g.player.X {
				enemy.X -= 1
				enemy.Direction = 2
			} else if enemy.Y > g.player.Y {
				enemy.Y -= 1
				enemy.Direction = 1
			} else if enemy.Y < g.player.Y {
				enemy.Y += 1
				enemy.Direction = 0
			} else {
				enemy.Frame -= 1
			}
			enemy.Frame %= 256
		}
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

func DrawSprite(sprite *Sprite, screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sprite.X, sprite.Y)

	dx := 16 * sprite.Direction
	fy := 16 * (sprite.Frame / 4 % 4)

	screen.DrawImage(
		sprite.Img.SubImage(
			image.Rect(dx, fy, 16+dx, 16+fy),
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
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok, Bye!")
}
