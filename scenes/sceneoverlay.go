package scenes

import (
	"image/color"
	"rpg-tutorial/constants"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type SceneOvelay struct {
	Value float64
	gameScene Scene
	tmpImage  *ebiten.Image
}

func (s *SceneOvelay) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{33, 88, 99, 255})

	s.gameScene.Draw(s.tmpImage)
	cm := colorm.ColorM{}
	cm.ChangeHSV(0.0, 1.0, s.Value)
	colorm.DrawImage(screen, s.tmpImage, cm, nil)
}

func NewSceneOverlay(gameScene Scene) *SceneOvelay {
	return &SceneOvelay{
		gameScene: gameScene,
		tmpImage:  ebiten.NewImage(constants.ScreenWidth, constants.ScreenHeight),
		Value: 0.6,
	}
}
