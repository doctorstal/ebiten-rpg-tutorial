package entities

import "rpg-tutorial/components"

type Player struct {
	*Sprite
	Health          uint
	CombatComponent components.Combat
}
