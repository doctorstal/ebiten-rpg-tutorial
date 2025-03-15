package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StartScene struct {
	isLoaded bool
}

// Draw implements Scene.
func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{33, 66, 99, 255})
	ebitenutil.DebugPrint(screen, "Press <Enter> to start.")
}

// FirstLoad implements Scene.
func (s *StartScene) FirstLoad() {
	s.isLoaded = true
}

// IsLoaded implements Scene.
func (s *StartScene) IsLoaded() bool {
	return s.isLoaded
}

// OnEnter implements Scene.
func (s *StartScene) OnEnter() {
}

// OnExit implements Scene.
func (s *StartScene) OnExit() {
}

// Update implements Scene.
func (s *StartScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return StartSceneId
}

func NewStartScene() Scene {
	return &StartScene{}
}
