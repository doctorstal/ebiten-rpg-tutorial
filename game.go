package main

import (
	"github.com/doctorstal/ebiten-rpg-tutorial/constants"
	"github.com/doctorstal/ebiten-rpg-tutorial/scenes"
	"github.com/doctorstal/ebiten-rpg-tutorial/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Game struct {
	sceneMap      map[scenes.SceneId]scenes.Scene
	activeSceneId scenes.SceneId
	gameState     *state.GlobalGameState
}

func NewGame(loader *resource.Loader) *Game {
	gameState := &state.GlobalGameState{
		SelectedHero: state.HeroSamurai,
	}

	startScene := scenes.NewStartScene(gameState, loader)
	startScene.FirstLoad()
	startScene.OnEnter()

	gameScene := scenes.NewGameScene(gameState, loader)
	game := &Game{
		sceneMap: map[scenes.SceneId]scenes.Scene{
			scenes.StartSceneId:      startScene,
			scenes.GameSceneId:       gameScene,
			scenes.PauseSceneId:      scenes.NewPauseScene(gameScene),
			scenes.TransitionSceneId: scenes.NewTransitionScene(gameScene),
			scenes.WonSceneId:        scenes.NewEndScene(gameScene, loader, true),
			scenes.LostSceneId:       scenes.NewEndScene(gameScene, loader, false),
		},
		activeSceneId: scenes.StartSceneId,
		gameState:     gameState,
	}

	return game
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.gameState.DebugMode = !g.gameState.DebugMode
	}
	nextSceneId := g.sceneMap[g.activeSceneId].Update()
	if nextSceneId == scenes.ExitSceneId {
		g.sceneMap[g.activeSceneId].OnExit()
		return ebiten.Termination
	}
	if nextSceneId != g.activeSceneId {
		nextScene := g.sceneMap[nextSceneId]
		// if not loaded load scene
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad()
		}
		g.sceneMap[g.activeSceneId].OnExit()
		nextScene.OnEnter()

		g.activeSceneId = nextSceneId
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
