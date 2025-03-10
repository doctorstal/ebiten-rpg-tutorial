package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"rpg-tutorial/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const TileWidth = 16

type Game struct {
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
	tilemapImg  *ebiten.Image
	cam         *Camera
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

	g.cam.FollowTarget(g.player.X+TileWidth/2, g.player.Y+TileWidth/2)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x66, 0x99, 255})

	for layerIdx, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := TileWidth * (index % layer.Width)
			y := TileWidth * (index / layer.Width)

			img := g.tilesets[layerIdx].Img(id)
			g.cam.Render(
				screen,
				img,
				float64(x),
				float64(y-img.Bounds().Dy()+TileWidth),
			)
		}
	}

	// draw player
	DrawSprite(g.player.Sprite, screen, g.cam)

	for _, sprite := range g.enemies {
		DrawSprite(sprite.Sprite, screen, g.cam)
	}

	for _, potion := range g.potions {
		if !potion.Consumed {
			DrawSprite(potion.Sprite, screen, g.cam)
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d \n", g.player.Health))
}

func DrawSprite(sprite *entities.Sprite, screen *ebiten.Image, cam *Camera) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sprite.X, sprite.Y)

	dx := TileWidth * sprite.Direction
	fy := TileWidth * (sprite.Frame / 4 % 4)

	cam.Render(
		screen,
		sprite.Img.SubImage(
			image.Rect(dx, fy, TileWidth+dx, TileWidth+fy),
		).(*ebiten.Image),
		sprite.X,
		sprite.Y,
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

	tilesets, err := tilemapJSON.GenTilesets()

	fmt.Println(tilesets)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting...")
	game := &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img: playerImg,
				X:   160.0,
				Y:   120.0,
			},
			Health: 3,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   50.0,
					Y:   50.0,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   150.0,
					Y:   150.0,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img: skeletonImg,
					X:   75.0,
					Y:   75.0,
				},
				FollowsPlayer: true,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img: potionImg,
					X:   210.0,
					Y:   100.0,
				},
				AmtHeal: 1,
			},
		},
		tilemapJSON: tilemapJSON,
		tilesets:    tilesets,
		tilemapImg:  tilemapImg,
		cam: NewCamera(
			320,
			240,
			float64(tilemapJSON.Layers[0].Width*TileWidth),
			float64(tilemapJSON.Layers[0].Height*TileWidth),
		),
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok, Bye!")
}
