package entities

type Enemy struct {
	*Character
	FollowsPlayer   bool
}
