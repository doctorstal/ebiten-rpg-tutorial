package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"rpg-tutorial/animations"
	"rpg-tutorial/components"
	"rpg-tutorial/constants"
	"rpg-tutorial/entities"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

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

func NewGame() *Game {

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

	if err != nil {
		log.Fatal(err)
	}
	playerSpriteSheet := spritesheet.NewSpriteSheet(4, 7, 16)

	newEnemy := func(x, y float64, fp bool) *entities.Enemy {
		return &entities.Enemy{
			Sprite: &entities.Sprite{
				Img:         skeletonImg,
				X:           x,
				Y:           y,
				Width:       constants.TileSize,
				Height:      constants.TileSize,
				Spritesheet: playerSpriteSheet,
				Animations:  personAnimations(),
			},
			CombatComponent: components.NewEnemyCombat(3, 1, 30),
			FollowsPlayer:   fp,
		}
	}

	game := &Game{
		player: &entities.Player{
			Sprite: &entities.Sprite{
				Img:         playerImg,
				X:           160.0,
				Y:           120.0,
				Width:       constants.TileSize,
				Height:      constants.TileSize,
				Spritesheet: playerSpriteSheet,
				Animations:  personAnimations(),
			},
			CombatComponent: components.NewBasicCombat(5, 1),
			Health:          3,
		},
		enemies: []*entities.Enemy{
			newEnemy(50.0, 50.0, true),
			newEnemy(75.0, 75.0, true),
			newEnemy(150.0, 75.0, true),
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
			float64(tilemapJSON.Layers[0].Width*constants.TileSize),
			float64(tilemapJSON.Layers[0].Height*constants.TileSize),
		),
		colliders: []image.Rectangle{
			image.Rect(100, 100, 116, 116),
		},
	}
	return game
}

func personAnimations() map[entities.SpriteState]*animations.Animation {
	return map[entities.SpriteState]*animations.Animation{
		entities.Down:  animations.NewAnimation(0, 12, 4, 10.0),
		entities.Up:    animations.NewAnimation(1, 13, 4, 10.0),
		entities.Left:  animations.NewAnimation(2, 14, 4, 10.0),
		entities.Right: animations.NewAnimation(3, 15, 4, 10.0),
	}
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
		if enemy.FollowsPlayer && enemy.Dist(g.player.Sprite) < 5*constants.TileSize {
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

	}

	for _, potion := range g.potions {
		if !potion.Consumed && g.player.Rect().Overlaps(potion.Rect()) {
			g.player.Health += potion.AmtHeal
			potion.Consumed = true
		}
	}

	g.player.CombatComponent.Update()

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	cX, cY := ebiten.CursorPosition()
	cP := image.Point{cX - int(g.cam.X), cY - int(g.cam.Y)}

	deadEnemies := make(map[int]struct{})
	for idx, enemy := range g.enemies {
		enemy.CombatComponent.Update()
		if clicked && cP.In(enemy.Rect()) && enemy.Dist(g.player.Sprite) < 5*constants.TileSize {
			enemy.CombatComponent.Damage(g.player.CombatComponent.AttackPower())
			if enemy.CombatComponent.Health() <= 0 {
				deadEnemies[idx] = struct{}{}
			}
		}
		if enemy.Rect().Overlaps(g.player.Rect()) {
			if enemy.CombatComponent.Attack() {
				g.player.CombatComponent.Damage(enemy.CombatComponent.AttackPower())
			}
		}
	}

	if len(deadEnemies) > 0 {
		n := 0
		for idx, e := range g.enemies {
			if _, exists := deadEnemies[idx]; !exists {
				g.enemies[n] = e
				n++
			}
		}
		g.enemies = g.enemies[:n]
	}

	g.cam.FollowTarget(g.player.X+constants.TileSize/2, g.player.Y+constants.TileSize/2)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x66, 0x99, 255})

	for layerIdx, layer := range g.tilemapJSON.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := constants.TileSize * (index % layer.Width)
			y := constants.TileSize * (index / layer.Width)

			img := g.tilesets[layerIdx].Img(id)
			g.cam.Render(
				screen,
				img,
				float64(x),
				float64(y-img.Bounds().Dy()+constants.TileSize),
			)
		}
	}

	// draw player
	drawSprite(g.player.Sprite, screen, g.cam)

	for _, sprite := range g.enemies {
		drawSprite(sprite.Sprite, screen, g.cam)
	}

	for _, potion := range g.potions {
		if !potion.Consumed {
			drawSprite(potion.Sprite, screen, g.cam)
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
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d \n", g.player.CombatComponent.Health()))
}

func drawSprite(sprite *entities.Sprite, screen *ebiten.Image, cam *Camera) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(sprite.X, sprite.Y)

	activeAnim := sprite.ActiveAnimation()
	frame := 0
	if activeAnim != nil {
		frame = activeAnim.Frame()
	}
	var rect image.Rectangle
	if sprite.Spritesheet != nil {
		rect = sprite.Spritesheet.Rect(frame)
	} else {
		rect = image.Rect(0, 0, constants.TileSize, constants.TileSize)
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
