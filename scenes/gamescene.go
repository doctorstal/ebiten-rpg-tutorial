package scenes

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"rpg-tutorial/camera"
	"rpg-tutorial/components"
	"rpg-tutorial/constants"
	"rpg-tutorial/entities"
	"rpg-tutorial/spritesheet"
	"rpg-tutorial/tiled"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameScene struct {
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	deadEnemies       []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *tiled.TilemapJSON
	tilesets          []tiled.Tileset
	cam               *camera.Camera
	colliders         []image.Rectangle
	isLoaded          bool
}

// IsLoaded implements Scene.
func (g *GameScene) IsLoaded() bool {
	return g.isLoaded
}

func (g *GameScene) Draw(screen *ebiten.Image) {
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
		for _, obj := range layer.Objects {
			img := g.tilesets[1].Img(obj.Gid)
			g.cam.Render(
				screen,
				img,
				obj.X,
				obj.Y,
			)
		}
	}

	for _, bomb := range g.player.Bombs {
		drawSprite(bomb.Sprite, screen, g.cam)
	}

	for _, sprite := range g.enemies {
		drawSprite(sprite.Sprite, screen, g.cam)
	}
	for _, sprite := range g.deadEnemies {
		drawSprite(sprite.Sprite, screen, g.cam)
	}

	for _, potion := range g.potions {
		if !potion.Consumed {
			drawSprite(potion.Sprite, screen, g.cam)
		}
	}

	// draw player
	drawSprite(g.player.Sprite, screen, g.cam)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d \n", g.player.CombatComponent.Health()))
}

func drawSprite(sprite *entities.Sprite, screen *ebiten.Image, cam *camera.Camera) {
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

func (g *GameScene) FirstLoad() {

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/samurai.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	bombImg, _, err := ebitenutil.NewImageFromFile("assets/images/bomb.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/LifePot.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := tiled.NewTilemapJSON("assets/maps/first.tmj")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTilesets()

	if err != nil {
		log.Fatal(err)
	}

	newEnemy := func(x, y float64, fp bool) *entities.Enemy {
		return &entities.Enemy{
			Character:     entities.NewCharacter(skeletonImg, x, y, components.NewEnemyCombat(3, 1, 30)),
			FollowsPlayer: fp,
		}
	}

	g.player = entities.NewPlayer(playerImg, bombImg, 160.0, 100.0)
	g.enemies = []*entities.Enemy{
		newEnemy(50.0, 50.0, true),
		newEnemy(75.0, 75.0, false),
		newEnemy(150.0, 75.0, true),
	}
	g.potions = []*entities.Potion{
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
	}
	g.tilemapJSON = tilemapJSON
	g.tilesets = tilesets
	g.cam = camera.NewCamera(
		320,
		240,
		float64(tilemapJSON.Layers[0].Width*constants.TileSize),
		float64(tilemapJSON.Layers[0].Height*constants.TileSize),
	)

	colliders := make([]image.Rectangle, 0)
	for _, layer := range tilemapJSON.Layers {
		for _, obj := range layer.Objects {

			x := int(obj.X)
			y := int(obj.Y)
			colliders = append(colliders, image.Rect(
				x,
				y,
				x+obj.Width,
				y+obj.Height,
			))
		}

	}

	g.colliders = colliders
	g.isLoaded = true

}

func (g *GameScene) OnEnter() {
}

func (g *GameScene) OnExit() {
}

func (g *GameScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return PauseSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	// react to key presses

	g.player.Dx = 0
	g.player.Dy = 0
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.Direction = 3
		g.player.Dx = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.Direction = 2
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.Direction = 1
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.Direction = 0
		g.player.Dy = 2
	}

	g.player.UpdateAnimation()
	g.player.CheckCollision(g.colliders)
	g.player.Move()

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

		enemy.UpdateAnimation()
		enemy.CheckCollision(g.colliders)
		enemy.Move()
	}

	for _, potion := range g.potions {
		if !potion.Consumed && g.player.Rect().Overlaps(potion.Rect()) {
			g.player.CombatComponent.Heal(potion.AmtHeal)
			potion.Consumed = true
		}
	}

	g.player.CombatComponent.Update()

	playerAttacks := inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.CombatComponent.Attack()
	// playerAttacks := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && g.player.CombatComponent.Attack()
	if playerAttacks {
		g.player.UpdateAnimation()
		g.player.Bombs = append(g.player.Bombs, entities.NewBomb(g.player.BombImg, g.player.X, g.player.Y, g.player.CombatComponent.AttackPower()))
	}
	// cX, cY := ebiten.CursorPosition()
	// cP := image.Point{cX - int(g.cam.X), cY - int(g.cam.Y)}

	deadEnemies := make(map[int]*entities.Enemy)
	for idx, enemy := range g.enemies {
		enemy.CombatComponent.Update()
		firedBombs := make(map[int]struct{})
		// if playerAttacks && cP.In(enemy.Rect()) && enemy.Dist(g.player.Sprite) < 5*constants.TileSize {
		for bidx, bomb := range g.player.Bombs {
			if enemy.Rect().Overlaps(bomb.Rect()) {
				enemy.CombatComponent.Damage(bomb.AmtDamage)
				firedBombs[bidx] = struct{}{}
				if enemy.CombatComponent.Health() <= 0 {
					enemy.Die()
					deadEnemies[idx] = enemy
					g.deadEnemies = append(g.deadEnemies, enemy)
					g.colliders = append(g.colliders, enemy.Rect())
				}
			}
		}
		if len(firedBombs) > 0 {
			n := 0
			for idx, b := range g.player.Bombs {
				if _, fired := firedBombs[idx]; !fired {
					g.player.Bombs[n] = b
					n++
				}
			}
			g.player.Bombs = g.player.Bombs[:n]
		}
		if enemy.Rect().Overlaps(g.player.Rect()) {
			if enemy.CombatComponent.Attack() {
				enemy.UpdateAnimation()
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

	return GameSceneId
}

func NewGameScene() Scene {
	return &GameScene{
		player:            nil,
		playerSpriteSheet: nil,
		enemies:           make([]*entities.Enemy, 0),
		deadEnemies:       make([]*entities.Enemy, 0),
		potions:           make([]*entities.Potion, 0),
		tilemapJSON:       nil,
		tilesets:          make([]tiled.Tileset, 0),
		cam:               nil,
		colliders:         make([]image.Rectangle, 0),
	}
}
