package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TransitionScene struct {
	isLoaded     bool
	overlayScene *SceneOvelay
}

// Draw implements Scene.
func (s *TransitionScene) Draw(screen *ebiten.Image) {
	s.overlayScene.Draw(screen)
}

// FirstLoad implements Scene.
func (s *TransitionScene) FirstLoad() {
	s.isLoaded = true
}

// IsLoaded implements Scene.
func (s *TransitionScene) IsLoaded() bool {
	return s.isLoaded
}

// OnEnter implements Scene.
func (s *TransitionScene) OnEnter() {
	s.overlayScene.Value = 0
}

// OnExit implements Scene.
func (s *TransitionScene) OnExit() {
}

// Update implements Scene.
func (s *TransitionScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	s.overlayScene.Value +=0.05
	if s.overlayScene.Value>0.9 {
		return GameSceneId
	}
	return TransitionSceneId
}

func NewTransitionScene(gameScene Scene) Scene {
	return &TransitionScene{
		overlayScene: NewSceneOverlay(gameScene),
	}
}
