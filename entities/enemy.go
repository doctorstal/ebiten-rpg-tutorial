package entities

import (
	"rpg-tutorial/components"
	"rpg-tutorial/resources"

	resource "github.com/quasilyte/ebitengine-resource"
)

type Enemy struct {
	*Character
	FollowsPlayer  bool
	WonderingSpeed float64
}

func (e *Enemy) IsDead() bool {
	return e.state == Dead
}

func NewEnemy(x, y float64, fp bool, loader *resource.Loader) *Enemy {
	return &Enemy{
		Character: NewCharacter(
			loader.LoadImage(resources.ImgSkeleton).Data,
			x,
			y,
			components.NewEnemyCombat(3, 1, 30),
		),
		FollowsPlayer: fp,
	}
}
