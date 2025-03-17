package entities

import (
	"rpg-tutorial/animations"
	"rpg-tutorial/constants"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bomb struct {
	*Sprite
	AmtDamage uint
}

func NewBomb(img *ebiten.Image, x, y float64, dmg uint) *Bomb {
	return &Bomb{
		Sprite: &Sprite{
			X:           x,
			Y:           y,
			Width:       constants.TileSize,
			Height:      constants.TileSize,
			Img:         img,
			Spritesheet: spritesheet.NewSpriteSheet(1, 3, 16),
			Animations: map[SpriteState]*animations.Animation{
				Idle: animations.NewAnimation(0, 2, 1, 10.0),
			},
			state: Idle,
		},
		AmtDamage: dmg,
	}
}
