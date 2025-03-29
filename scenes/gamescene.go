package scenes

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"rpg-tutorial/camera"
	"rpg-tutorial/components/ui"
	"rpg-tutorial/constants"
	"rpg-tutorial/entities"
	"rpg-tutorial/resources"
	"rpg-tutorial/state"
	"rpg-tutorial/world"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	resource "github.com/quasilyte/ebitengine-resource"
)

type GameScene struct {
	gameState *state.GlobalGameState

	world     *world.WorldState
	roomState *world.RoomState

	shadowImg *ebiten.Image
	cam       *camera.Camera
	isLoaded  bool
	loader    *resource.Loader
	ingameUi  *ui.IngameUi
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
	screen.DrawImage(g.roomState.TiledMap.GroundImage(viewRect), nil)
	g.drawShadow(screen, g.roomState.Player.Rect())
	for _, attackItem := range g.roomState.Player.AttackItems {
		g.drawShadow(screen, attackItem.HitRect())
	}
	for _, enemy := range g.roomState.Enemies {
		g.drawShadow(screen, enemy.Rect())
	}
	for _, potion := range g.roomState.Potions {
		if !potion.Consumed {
			g.drawShadow(screen, potion.Rect())
		}
	}
	screen.DrawImage(g.roomState.TiledMap.ObjectsImage(viewRect), nil)

	renderers := make([]entities.Renderer, 0)
	addRenderer := func(r entities.Renderer) {
		if viewRect.Overlaps(*r.Rect()) {
			renderers = append(renderers, r)
		}
	}

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

	for _, animator := range g.roomState.StaticAnimators {
		addRenderer(animator.GetRenderer())
	}

	for _, potion := range g.roomState.Potions {
		if !potion.Consumed {
			addRenderer(potion.GetRenderer())

		}
	}

	for _, attackItem := range g.roomState.Player.AttackItems {
		addRenderer(attackItem.GetRenderer())
	}

	for _, enemy := range g.roomState.Enemies {
		addRenderer(enemy.GetRenderer())
	}

	// draw player

	addRenderer(g.roomState.Player.GetRenderer())

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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Player Health: %d, Enemies Left: %d \n", g.roomState.Player.CombatComponent.Health(), len(g.roomState.Enemies)))
	g.ingameUi.Draw(screen)
	if g.gameState.DebugMode {
		for _, c := range g.roomState.TiledMap.Doors() {
			x := float64(c.Min.X)
			y := float64(c.Min.Y)

			vector.StrokeRect(screen, float32(g.cam.X+x), float32(g.cam.Y+y), float32(c.Dx()), float32(c.Dy()), 1.0, color.RGBA{25, 200, 0, 255}, false)
		}
		for _, c := range g.roomState.Colliders {
			x := float64(c.Min.X)
			y := float64(c.Min.Y)

			vector.StrokeRect(screen, float32(g.cam.X+x), float32(g.cam.Y+y), float32(c.Dx()), float32(c.Dy()), 1.0, color.RGBA{255, 0, 0, 255}, false)
		}
	}
}

func (g *GameScene) FirstLoad() {

	shadowImg := g.loader.LoadImage(resources.ImgShadow).Data
	g.shadowImg = ebiten.NewImage(shadowImg.Bounds().Dx(), shadowImg.Bounds().Dy())
	var cm colorm.ColorM
	cm.Scale(1, 1, 1, 0.3)
	colorm.DrawImage(g.shadowImg, shadowImg, cm, nil)

	player := entities.NewPlayer(g.loader, g.gameState.SelectedHero, 360.0, 100.0)

	g.world = world.NewWorldState(g.loader, "second.tmx")
	g.roomState = g.world.LoadRoom("first.tmx", player)

	mapWidth := g.roomState.TiledMap.Width()
	mapHeight := g.roomState.TiledMap.Height()
	g.cam = camera.NewCamera(
		320,
		240,
		mapWidth,
		mapHeight,
		g.gameState,
	)

	g.isLoaded = true

}

func (g *GameScene) Unload() {
	// TODO do proper unload
	g.roomState.StaticAnimators = g.roomState.StaticAnimators[:0]
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

	g.roomState.Player.Dx = 0
	g.roomState.Player.Dy = 0
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.roomState.Player.Direction = 0
		g.roomState.Player.Dy = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.roomState.Player.Direction = 1
		g.roomState.Player.Dx = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.roomState.Player.Direction = 2
		g.roomState.Player.Dy = -2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.roomState.Player.Direction = 3
		g.roomState.Player.Dx = 2
	}

	// Go through doors
	for name, d := range g.roomState.TiledMap.Doors() {
		if d.Overlaps(g.roomState.Player.Rect().Add(image.Point{int(g.roomState.Player.Dx), int(g.roomState.Player.Dy)})) {
			g.roomState = g.world.LoadRoom(name, g.roomState.Player)
			g.cam.UpdateRoomSize(g.roomState.TiledMap.Width(), g.roomState.TiledMap.Height())
			g.cam.GoToTarget(g.roomState.Player.X, g.roomState.Player.Y)
			return PauseSceneId
			// return TransitionSceneId
		}
	}

	g.roomState.Player.UpdateState()
	g.roomState.Player.UpdateAnimation()
	g.roomState.Player.CheckCollision(g.roomState.Colliders)
	g.roomState.Player.Move()

	for _, enemy := range g.roomState.Enemies {
		enemy.Dx = 0
		enemy.Dy = 0
		if enemy.FollowsPlayer && enemy.Dist(g.roomState.Player.Sprite) < 5*constants.TileSize {
			if enemy.Y-g.roomState.Player.Y < -1 {
				enemy.Direction = 0
				enemy.Dy = 1
			}
			if enemy.X-g.roomState.Player.X > 1 {
				enemy.Direction = 1
				enemy.Dx = -1
			}
			if enemy.Y-g.roomState.Player.Y > 1 {
				enemy.Direction = 2
				enemy.Dy = -1
			}
			if enemy.X-g.roomState.Player.X < -1 {
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
		enemy.CheckCollision(g.roomState.Colliders)
		enemy.Move()
	}

	sn := 0
	for _, animator := range g.roomState.StaticAnimators {
		finished := animator.UpdateAnimation()
		if !finished {
			g.roomState.StaticAnimators[sn] = animator
			sn++
		}
	}
	g.roomState.StaticAnimators = g.roomState.StaticAnimators[:sn]

	for _, potion := range g.roomState.Potions {
		if !potion.Consumed && g.roomState.Player.Rect().Overlaps(*potion.Rect()) {
			g.roomState.Player.CombatComponent.Heal(potion.AmtHeal)
			potion.Consumed = true
		}
	}

	g.roomState.Player.CombatComponent.Update()

	for _, attackItem := range g.roomState.Player.AttackItems {
		attackItem.Update()
	}

	playerAttacks := inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.roomState.Player.CombatComponent.Attack()
	if playerAttacks {
		g.roomState.Player.UpdateState()
		g.roomState.Player.AttackItems = append(g.roomState.Player.AttackItems, g.roomState.Player.NewBomb())
	}

	deadEnemies := make(map[int]*entities.Enemy)
	for idx, enemy := range g.roomState.Enemies {
		enemy.CombatComponent.Update()
		for _, attackItem := range g.roomState.Player.AttackItems {
			if enemy.Rect().Overlaps(*attackItem.HitRect()) {
				enemy.CombatComponent.Damage(attackItem.GetAmtDamage())
				attackItem.DoDamage()
				g.roomState.StaticAnimators = append(g.roomState.StaticAnimators, attackItem.GetAnimator())
				if enemy.CombatComponent.Health() <= 0 {
					enemy.Die()
					deadEnemies[idx] = enemy
					g.roomState.StaticAnimators = append(g.roomState.StaticAnimators, enemy.Sprite)
					g.roomState.Colliders = append(g.roomState.Colliders, enemy.Rect())
				}
			}
		}
		if !enemy.IsDead() && enemy.Rect().Overlaps(*g.roomState.Player.Rect()) {
			if enemy.CombatComponent.Attack() {
				enemy.UpdateState()
				g.roomState.Player.CombatComponent.Damage(enemy.CombatComponent.AttackPower())
			}
		}
	}

	if len(deadEnemies) > 0 {
		n := 0
		for idx, e := range g.roomState.Enemies {
			if _, exists := deadEnemies[idx]; !exists {
				g.roomState.Enemies[n] = e
				n++
			}
		}
		g.roomState.Enemies = g.roomState.Enemies[:n]
	}
	n := 0
	for _, b := range g.roomState.Player.AttackItems {
		if remove := b.ShouldRemove(); !remove {
			g.roomState.Player.AttackItems[n] = b
			n++
		}
	}
	g.roomState.Player.AttackItems = g.roomState.Player.AttackItems[:n]

	g.ingameUi.Update(g.roomState.Player.CombatComponent.Health())

	g.cam.FollowTarget(g.roomState.Player.X+constants.TileSize/2, g.roomState.Player.Y+constants.TileSize/2)

	if g.roomState.Player.CombatComponent.Health() <= 0 {
		g.Unload()
		return LostSceneId
	}
	if len(g.roomState.Enemies) == 0 {
		g.Unload()
		return WonSceneId
	}

	return GameSceneId
}

func NewGameScene(gameState *state.GlobalGameState, loader *resource.Loader) Scene {
	return &GameScene{
		gameState: gameState,
		cam:       nil,
		loader:    loader,
		ingameUi:  ui.NewIngameUi(loader, gameState),
	}
}
