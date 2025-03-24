package entities

import (
	"image"
	"rpg-tutorial/animations"
	"rpg-tutorial/constants"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bomb struct {
	*Sprite
	AmtDamage uint
}

// HitRect implements AttackItem.
func (b *Bomb) HitRect() image.Rectangle {
	return b.Rect()
}

// DoDamage implements AttackItem.
func (b *Bomb) DoDamage() {
	b.state = Dead
}

// GetAmtDamage implements AttackItem.
func (b *Bomb) GetAmtDamage() uint {
	return b.AmtDamage
}

// GetSprite implements AttackItem.
func (b *Bomb) GetSprite() *Sprite {
	return b.Sprite
}

// ShouldRemove implements AttackItem.
func (b *Bomb) ShouldRemove() bool {
	return b.state == Dead
}

// Update implements AttackItem.
func (b *Bomb) Update() {
	b.UpdateAnimation()
}

func NewBomb(img *ebiten.Image, x, y float64, dmg uint) AttackItem {
	return &Bomb{
		Sprite: &Sprite{
			X:           x,
			Y:           y,
			Width:       constants.TileSize,
			Height:      constants.TileSize,
			Img:         img,
			Spritesheet: spritesheet.NewSpriteSheet(1, 7, 16),
			Animations: map[SpriteState]animations.Animation{
				Idle: animations.NewLoopAnimation(0, 2, 1, 10.0),
				Dead: animations.NewOneTimeAnimation(3, 6, 1, 10.0, false),
			},
			state: Idle,
		},
		AmtDamage: dmg,
	}
}
