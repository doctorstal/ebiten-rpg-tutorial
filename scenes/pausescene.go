package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	isLoaded     bool
	overlayScene *SceneOvelay
}

// Draw implements Scene.
func (s *PauseScene) Draw(screen *ebiten.Image) {
	s.overlayScene.Draw(screen)
	ebitenutil.DebugPrintAt(screen, "Press <Enter> to unpause.", 0, 20)
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}
	return PauseSceneId
}

func NewPauseScene(gameScene Scene) Scene {
	return &PauseScene{
		overlayScene: NewSceneOverlay(gameScene),
	}
}
