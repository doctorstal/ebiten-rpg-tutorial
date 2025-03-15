package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	isLoaded bool
}

// Draw implements Scene.
func (s *PauseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{33, 88, 99, 255})
	ebitenutil.DebugPrint(screen, "Press <Enter> to unpause.")
}

// FirstLoad implements Scene.
func (s *PauseScene) FirstLoad() {
	s.isLoaded = true
}

// IsLoaded implements Scene.
func (s *PauseScene) IsLoaded() bool {
	return s.isLoaded
}

// OnEnter implements Scene.
func (s *PauseScene) OnEnter() {
}

// OnExit implements Scene.
func (s *PauseScene) OnExit() {
}

// Update implements Scene.
func (s *PauseScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return PauseSceneId
}

func NewPauseScene() Scene {
	return &PauseScene{}
}
