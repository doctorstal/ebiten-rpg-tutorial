package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EndScene struct {
	victory bool
}

// Draw implements Scene.
func (e *EndScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{33, 88, 99, 255})
	var msg = "Press <Enter> to start over."
	if e.victory {
		msg = "You won!\n" + msg
	} else {
		msg = "You lost!\n" + msg
	}

	ebitenutil.DebugPrint(screen, msg)

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

func NewEndScene(victory bool) Scene {
	return &EndScene{
		victory: victory,
	}
}
