package entities

type Enemy struct {
	*Character
	FollowsPlayer bool
}

func (e *Enemy) IsDead() bool {
	return e.state == Dead
}
