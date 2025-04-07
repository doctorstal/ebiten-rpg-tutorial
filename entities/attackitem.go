package entities

import "image"

type AttackItem interface {
	GetRenderers() []Renderer
	GetAnimator() Animator
	Update()
	HitRect() *image.Rectangle
	GetAmtDamage() uint
	DoDamage()
	ShouldRemove() bool
}
