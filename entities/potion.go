package entities

import (
	"github.com/doctorstal/ebiten-rpg-tutorial/resources"

	resource "github.com/quasilyte/ebitengine-resource"
)

type Potion struct {
	*Sprite
	AmtHeal  uint
	Consumed bool
}

func NewPotion(x, y float64, loader *resource.Loader) *Potion {
	return &Potion{
		Sprite: &Sprite{
			Img:    loader.LoadImage(resources.ImgPotion).Data,
			X:      x,
			Y:      y,
			Width:  8.0,
			Height: 10.0,
		},
		AmtHeal: 1,
	}
}
