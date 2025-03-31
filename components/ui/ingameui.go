package ui

import (
	"github.com/doctorstal/ebiten-rpg-tutorial/resources"
	"github.com/doctorstal/ebiten-rpg-tutorial/spritesheet"
	"github.com/doctorstal/ebiten-rpg-tutorial/state"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type HeartBar struct {
	heartImg     *ebiten.Image
	spriteSheet  spritesheet.SpriteSheet
	playerHealth int
}

func (h *HeartBar) Update(playerHealth int) {
	h.playerHealth = playerHealth
}

func (h *HeartBar) Draw(screen *ebiten.Image) {
	img := h.heartImg.SubImage(
		h.spriteSheet.Rect(4),
	).(*ebiten.Image)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(10, 10)
	for range h.playerHealth {
		screen.DrawImage(img, opts)
		opts.GeoM.Translate(float64(img.Bounds().Dx()+2), 0)
	}
	img = h.heartImg.SubImage(
		h.spriteSheet.Rect(0),
	).(*ebiten.Image)
	for range 5-h.playerHealth {
		screen.DrawImage(img, opts)
		opts.GeoM.Translate(float64(img.Bounds().Dx()+2), 0)
	}
}

type IngameUi struct {
	heartBar HeartBar
	loader   *resource.Loader
	gameState *state.GlobalGameState
}

func (i *IngameUi) Update(playerHealth int) {
	i.heartBar.Update(playerHealth)
}

func (i *IngameUi) Draw(screen *ebiten.Image) {
	i.heartBar.Draw(screen)
}

func NewIngameUi(loader *resource.Loader, gameState *state.GlobalGameState) *IngameUi {
	return &IngameUi{
		heartBar: HeartBar{
			heartImg:    loader.LoadImage(resources.ImgHeart).Data,
			spriteSheet: *spritesheet.NewSpriteSheet(5, 1, 16),
		},
		loader: loader,
		gameState: gameState,
	}
}
