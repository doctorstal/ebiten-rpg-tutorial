package entities

import "image"

type AttackItem interface {
	GetSprite() *Sprite
	Update()
	HitRect() image.Rectangle
	GetAmtDamage() uint
	DoDamage()
	ShouldRemove() bool
}
