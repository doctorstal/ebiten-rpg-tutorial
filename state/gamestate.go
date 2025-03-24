package state

type Hero int

const (
	Samurai Hero = iota
	Robot
	Skeleton
)

type GlobalGameState struct {
	SelectedHero Hero
}
