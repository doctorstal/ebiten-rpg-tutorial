package entities

import (
	"rpg-tutorial/components"
)

type Character struct {
	*Sprite
	CombatComponent components.Combat
}

func (c *Character) Move() {
	if !c.CombatComponent.Attacking() {
		c.Sprite.Move()
	}
}

func (c *Character) UpdateAnimation() {
	if !c.CombatComponent.Attacking() {
		c.Sprite.UpdateAnimation()
		return
	}
	c.state = AttackDown
	if c.Dy > 0 {
		c.state = AttackDown
	}
	if c.Dy < 0 {
		c.state = AttackUp
	}
	if c.Dx > 0 {
		c.state = AttackRight
	}
	if c.Dx < 0 {
		c.state = AttackLeft
	}
	animation := c.ActiveAnimation()
	if animation != nil {
		animation.Update()
	}
}

