package entities

import "image"

type AttackItem interface {
	GetRenderer() Renderer
	GetAnimator() Animator
	Update()
	HitRect() *image.Rectangle
	GetAmtDamage() uint
	DoDamage()
	ShouldRemove() bool
}
