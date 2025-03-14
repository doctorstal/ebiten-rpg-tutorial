package entities

import "rpg-tutorial/components"

type Enemy struct {
	*Sprite
	FollowsPlayer   bool
	CombatComponent components.Combat
}
