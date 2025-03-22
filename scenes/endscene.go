package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EndScene struct {
	victory      bool
	sceneOverlay *SceneOvelay
}

// Draw implements Scene.
func (e *EndScene) Draw(screen *ebiten.Image) {
	e.sceneOverlay.Draw(screen)
	var msg = "Press <Enter> to start over."
	if e.victory {
		msg = "You won!\n" + msg
	} else {
		msg = "You lost!\n" + msg
	}

	ebitenutil.DebugPrintAt(screen, msg, 0, 20)

}

func (e *EndScene) FirstLoad() {
}

func (e *EndScene) IsLoaded() bool {
	return true
}

func (e *EndScene) OnEnter() {
}

func (e *EndScene) OnExit() {
}

// Update implements Scene.
func (e *EndScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	if e.victory {
		return WonSceneId

	} else {
		return LostSceneId
	}
}

func NewEndScene(gameScene Scene, victory bool) Scene {
	return &EndScene{
		victory:      victory,
		sceneOverlay: NewSceneOverlay(gameScene),
	}
}
