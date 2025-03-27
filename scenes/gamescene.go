package scenes

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"rpg-tutorial/camera"
	"rpg-tutorial/components"
	"rpg-tutorial/constants"
	"rpg-tutorial/entities"
	"rpg-tutorial/resources"
	"rpg-tutorial/state"
	"rpg-tutorial/tiled"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	resource "github.com/quasilyte/ebitengine-resource"
)

type GameScene struct {
	gameState       *state.GlobalGameState
	player          *entities.Player
	shadowImg       *ebiten.Image
	enemies         []*entities.Enemy
	staticAnimators []entities.Animator
	potions         []*entities.Potion
	tilemapJSON     *tiled.TilemapJSON
	tilesets        []tiled.Tileset
	tiledMap        *tiled.TiledMap
	cam             *camera.Camera
	colliders       []*image.Rectangle
	isLoaded        bool
	loader          *resource.Loader
}

// IsLoaded implements Scene.
func (g *GameScene) IsLoaded() bool {
	return g.isLoaded
}

func (g *GameScene) drawShadow(screen *ebiten.Image, spriteRect *image.Rectangle) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(
		float64(spriteRect.Dx())/float64(constants.TileSize),
		float64(spriteRect.Dy())/float64(constants.TileSize),
	)
	g.cam.Render(screen,
		g.shadowImg,
		float64(spriteRect.Min.X+spriteRect.Dx()/2-g.shadowImg.Bounds().Dx()/2),
		float64(spriteRect.Min.Y+spriteRect.Dy()-g.shadowImg.Bounds().Dy()/2),
		opts,
	)
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	viewRect := g.cam.ViewRect()
	screen.DrawImage(g.tiledMap.GroundImage(viewRect), nil)
	g.drawShadow(screen, g.player.Rect())
	for _, bomb := range g.player.AttackItems {
		g.drawShadow(screen, bomb.HitRect())
	}
	for _, enemy := range g.enemies {
		g.drawShadow(screen, enemy.Rect())
	}
	for _, potion := range g.potions {
		if !potion.Consumed {
			g.drawShadow(screen, potion.Rect())
		}
	}
	screen.DrawImage(g.tiledMap.ObjectsImage(viewRect), nil)

	renderers := make([]entities.Renderer, 0)
	addRenderer := func(r entities.Renderer) {
		if viewRect.Overlaps(*r.Rect()) {
			renderers = append(renderers, r)
		}
	}

	//
	// for layerIdx, layer := range g.tilemapJSON.Layers {
	// 	for index, id := range layer.Data {
	// 		if id == 0 {
	// 			continue
	// 		}
	// 		x := constants.TileSize * (index % layer.Width)
	// 		y := constants.TileSize * (index / layer.Width)
	//
	// 		img := g.tilesets[layerIdx].Img(id)
	// 		g.cam.Render(
	// 			screen,
	// 			img,
	// 			float64(x),
	// 			float64(y-img.Bounds().Dy()+constants.TileSize),
	// 		)
	// 	}
	// 	for _, obj := range layer.Objects {
	// 		img := g.tilesets[1].Img(obj.Gid)
	// 		g.cam.Render(
	// 			screen,
	// 			img,
	// 			obj.X,
	// 			obj.Y-float64(obj.Height),
	// 		)
	// 	}
	// }

	for _, animator := range g.staticAnimators {
		addRenderer(animator.GetRenderer())
	}

	for _, potion := range g.potions {
		if !potion.Consumed {
			addRenderer(potion.GetRenderer())

		}
	}

	for _, attackItem := range g.player.AttackItems {
		addRenderer(attackItem.GetRenderer())
	}

	for _, enemy := range g.enemies {
		addRenderer(enemy.GetRenderer())
	}

	// draw player

	addRenderer(g.player.GetRenderer())

	slices.SortFunc(renderers, func(r1, r2 entities.Renderer) int {
		if r1.Z() == r2.Z() {
			return r1.Rect().Max.Y - r2.Rect().Max.Y
		} else {
			return r1.Z() - r2.Z()
		}
	})
	for _, r := range renderers {
		g.cam.Render(
			screen,
			r.Image(),
			float64(r.Rect().Min.X),
			float64(r.Rect().Min.Y),
			r.DrawOpts(),
		)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d, Enemies Left: %d \n", g.player.CombatComponent.Health(), len(g.enemies)))
}

func (g *GameScene) FirstLoad() {

	shadowImg := g.loader.LoadImage(resources.ImgShadow).Data
	skeletonImg := g.loader.LoadImage(resources.ImgSkeleton).Data
	potionImg := g.loader.LoadImage(resources.ImgPotion).Data

	tiledMap, err := tiled.NewTiledMap("assets/maps/first.tmx")
	if err != nil {
		log.Fatal(err)
	}

	newEnemy := func(x, y float64, fp bool) *entities.Enemy {
		return &entities.Enemy{
			Character:     entities.NewCharacter(skeletonImg, x, y, components.NewEnemyCombat(3, 1, 30)),
			FollowsPlayer: fp,
		}
	}

	g.shadowImg = ebiten.NewImage(shadowImg.Bounds().Dx(), shadowImg.Bounds().Dy())
	var cm colorm.ColorM
	cm.Scale(1, 1, 1, 0.3)
	colorm.DrawImage(g.shadowImg, shadowImg, cm, nil)

	g.player = entities.NewPlayer(g.loader, g.gameState.SelectedHero, 360.0, 100.0)

	g.enemies = []*entities.Enemy{
		newEnemy(50.0, 50.0, true),
		newEnemy(75.0, 75.0, false),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, false),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, true),
		newEnemy(150.0, 75.0, true),
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
	
	g.tiledMap = tiledMap
	mapWidth := tiledMap.Width()
	mapHeight := tiledMap.Height()
	g.cam = camera.NewCamera(
		320,
		240,
		mapWidth,
		mapHeight,
		g.gameState,
	)

	colliders := make([]*image.Rectangle, 0)

	for _, objectRect := range g.tiledMap.ObjectRects() {
		colliders = append(colliders, objectRect)
	}
	g.colliders = colliders
	g.isLoaded = true

}

func (g *GameScene) Unload() {
	// TODO do proper unload
	g.staticAnimators = g.staticAnimators[:0]
	g.isLoaded = false
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
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.Direction = 0
		g.player.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.Direction = 1
		g.player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.Direction = 2
		g.player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.Direction = 3
		g.player.Dx = 2
	}

	g.player.UpdateState()
	g.player.UpdateAnimation()
	g.player.CheckCollision(g.colliders)
	g.player.Move()

	for _, enemy := range g.enemies {
		enemy.Dx = 0
		enemy.Dy = 0
		if enemy.FollowsPlayer && enemy.Dist(g.player.Sprite) < 5*constants.TileSize {
			if enemy.Y-g.player.Y < -1 {
				enemy.Direction = 0
				enemy.Dy = 1
			}
			if enemy.X-g.player.X > 1 {
				enemy.Direction = 1
				enemy.Dx = -1
			}
			if enemy.Y-g.player.Y > 1 {
				enemy.Direction = 2
				enemy.Dy = -1
			}
			if enemy.X-g.player.X < -1 {
				enemy.Direction = 3
				enemy.Dx = 1
			}

		} else {
			if rand.Float64() > 0.95 {
				enemy.Direction = rand.Intn(4)
				enemy.WonderingSpeed = 0.1 + 0.5*rand.Float64()
			}
			enemy.Forward(enemy.WonderingSpeed)
		}

		enemy.UpdateAnimation()
		enemy.UpdateState()
		enemy.CheckCollision(g.colliders)
		enemy.Move()
	}

	sn := 0
	for _, animator := range g.staticAnimators {
		finished := animator.UpdateAnimation()
		if !finished {
			g.staticAnimators[sn] = animator
			sn++
		}
	}
	g.staticAnimators = g.staticAnimators[:sn]

	for _, potion := range g.potions {
		if !potion.Consumed && g.player.Rect().Overlaps(*potion.Rect()) {
			g.player.CombatComponent.Heal(potion.AmtHeal)
			potion.Consumed = true
		}
	}

	g.player.CombatComponent.Update()

	for _, attackItem := range g.player.AttackItems {
		attackItem.Update()
	}

	playerAttacks := inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.CombatComponent.Attack()
	// playerAttacks := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && g.player.CombatComponent.Attack()
	if playerAttacks {
		g.player.UpdateState()
		g.player.AttackItems = append(g.player.AttackItems, g.player.NewBomb())
	}

	deadEnemies := make(map[int]*entities.Enemy)
	for idx, enemy := range g.enemies {
		enemy.CombatComponent.Update()
		for _, attackItem := range g.player.AttackItems {
			if enemy.Rect().Overlaps(*attackItem.HitRect()) {
				enemy.CombatComponent.Damage(attackItem.GetAmtDamage())
				attackItem.DoDamage()
				g.staticAnimators = append(g.staticAnimators, attackItem.GetAnimator())
				if enemy.CombatComponent.Health() <= 0 {
					enemy.Die()
					deadEnemies[idx] = enemy
					g.staticAnimators = append(g.staticAnimators, enemy.Sprite)
					g.colliders = append(g.colliders, enemy.Rect())
				}
			}
		}
		if !enemy.IsDead() && enemy.Rect().Overlaps(*g.player.Rect()) {
			if enemy.CombatComponent.Attack() {
				enemy.UpdateState()
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
	n := 0
	for _, b := range g.player.AttackItems {
		if remove := b.ShouldRemove(); !remove {
			g.player.AttackItems[n] = b
			n++
		}
	}
	g.player.AttackItems = g.player.AttackItems[:n]

	g.cam.FollowTarget(g.player.X+constants.TileSize/2, g.player.Y+constants.TileSize/2)

	if g.player.CombatComponent.Health() <= 0 {
		g.Unload()
		return LostSceneId
	}
	if len(g.enemies) == 0 {
		g.Unload()
		return WonSceneId
	}

	return GameSceneId
}

func NewGameScene(gameState *state.GlobalGameState, loader *resource.Loader) Scene {
	return &GameScene{
		gameState:       gameState,
		player:          nil,
		enemies:         make([]*entities.Enemy, 0),
		staticAnimators: make([]entities.Animator, 0),
		potions:         make([]*entities.Potion, 0),
		tilemapJSON:     nil,
		tilesets:        make([]tiled.Tileset, 0),
		cam:             nil,
		colliders:       make([]*image.Rectangle, 0),
		loader:          loader,
	}
}
