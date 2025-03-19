package main

import (
	"rpg-tutorial/scenes"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneMap      map[scenes.SceneId]scenes.Scene
	activeSceneId scenes.SceneId
}

func NewGame() *Game {
	startScene := scenes.NewStartScene()
	startScene.FirstLoad()
	startScene.OnEnter()
	gameScene := scenes.NewGameScene()
	game := &Game{
		sceneMap: map[scenes.SceneId]scenes.Scene{
			scenes.StartSceneId: startScene,
			scenes.GameSceneId:  gameScene,
			scenes.PauseSceneId: scenes.NewPauseScene(),
			scenes.WonSceneId:   scenes.NewEndScene(true),
			scenes.LostSceneId:  scenes.NewEndScene(false),
		},
		activeSceneId: scenes.StartSceneId,
	}

	return game
}

func (g *Game) Update() error {
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
	return 320, 240
}
