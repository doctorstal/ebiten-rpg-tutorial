package entities

type Enemy struct {
	*Character
	FollowsPlayer bool
	WonderingSpeed float64
}

func (e *Enemy) IsDead() bool {
	return e.state == Dead
}
