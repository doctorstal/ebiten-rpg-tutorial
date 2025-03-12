package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"rpg-tutorial/animations"
	"rpg-tutorial/entities"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const TileWidth = 16

type Game struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	animationFrame    int
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *TilemapJSON
	tilesets          []Tileset
	tilemapImg        *ebiten.Image
	cam               *Camera
	colliders         []image.Rectangle
}

func (g *Game) Update() error {
	g.animationFrame++
	// react to key presses

	g.player.Dx = 0
	g.player.Dy = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.Direction = 3
		g.player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.Direction = 2
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Direction = 1
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Direction = 0
		g.player.Dy = 2
	}

	g.player.CheckCollision(g.colliders)

	if g.player.Dx != 0 || g.player.Dy != 0 {
		g.player.Move()
	}
	activeAnim := g.player.ActiveAnimation()
	if activeAnim != nil {
		activeAnim.Update()
	}

	for _, enemy := range g.enemies {
		enemy.Dx = 0
		enemy.Dy = 0
		if enemy.FollowsPlayer && enemy.Dist(g.player.Sprite) < 5*TileWidth {
			if enemy.X < g.player.X {
				enemy.Direction = 3
				enemy.Dx = 1
			}
			if enemy.X > g.player.X {
				enemy.Direction = 2
				enemy.Dx = -1
			}
			if enemy.Y > g.player.Y {
				enemy.Direction = 1
				enemy.Dy = -1
			}
			if enemy.Y < g.player.Y {
				enemy.Direction = 0
				enemy.Dy = 1
			}

		} else {
			if rand.Float64() > 0.95 {
				enemy.Direction = rand.Intn(4)
			}
			enemy.Forward(1)
		}

		enemy.CheckCollision(g.colliders)
		enemy.Move()
		activeAnim := enemy.ActiveAnimation()
		if activeAnim != nil {
			activeAnim.Update()
		}

		if enemy.Rect().Overlaps(g.player.Rect()) {
			g.player.Health -= 1
		}
	}

	for _, potion := range g.potions {
		if !potion.Consumed && g.player.Rect().Overlaps(potion.Rect()) {
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

	for _, collider := range g.colliders {
		vector.StrokeRect(screen,
			float32(collider.Min.X)+float32(g.cam.X),
			float32(collider.Min.Y)+float32(g.cam.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d \n", g.player.Health))
}

func DrawSprite(sprite *entities.Sprite, screen *ebiten.Image, cam *Camera) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sprite.X, sprite.Y)

	// // dx := TileWidth * sprite.Direction
	// fy := TileWidth * (sprite.Frame / 4 % 4)

	activeAnim := sprite.ActiveAnimation()
	frame := 0
	if activeAnim != nil {
		frame = activeAnim.Frame()
	}
	var rect image.Rectangle
	if sprite.Spritesheet != nil {
		rect = sprite.Spritesheet.Rect(frame)
	} else {
		rect = image.Rect(0, 0, TileWidth, TileWidth)
	}
	cam.Render(
		screen,
		sprite.Img.SubImage(
			rect,
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
	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)
	personAnimations := map[entities.SpriteState]*animations.Animation{
		entities.Down:  animations.NewAnimation(0, 12, 4, 10.0, 0),
		entities.Up:    animations.NewAnimation(1, 13, 4, 10.0, 0),
		entities.Left:  animations.NewAnimation(2, 14, 4, 10.0, 0),
		entities.Right: animations.NewAnimation(3, 15, 4, 10.0, 0),
	}

	fmt.Println("Starting...")
	game := &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img:         playerImg,
				X:           160.0,
				Y:           120.0,
				Width:       TileWidth,
				Height:      TileWidth,
				Spritesheet: playerSpriteSheet,
				Animations:  personAnimations,
			},
			Health: 3,
		},
		enemies: []*entities.Enemy{
			{
				Sprite: &entities.Sprite{
					Img:         skeletonImg,
					X:           50.0,
					Y:           50.0,
					Width:       TileWidth,
					Height:      TileWidth,
					Spritesheet: playerSpriteSheet,
					Animations:  personAnimations,
				},
				FollowsPlayer: true,
			},
			{
				Sprite: &entities.Sprite{
					Img:         skeletonImg,
					X:           150.0,
					Y:           150.0,
					Width:       TileWidth,
					Height:      TileWidth,
					Spritesheet: playerSpriteSheet,
					Animations:  personAnimations,
				},
				FollowsPlayer: false,
			},
			{
				Sprite: &entities.Sprite{
					Img:         skeletonImg,
					X:           75.0,
					Y:           75.0,
					Width:       TileWidth,
					Height:      TileWidth,
					Spritesheet: playerSpriteSheet,
					Animations:  personAnimations,
				},
				FollowsPlayer: true,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{
					Img:    potionImg,
					X:      210.0,
					Y:      100.0,
					Width:  8.0,
					Height: 10.0,
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
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok, Bye!")
}
