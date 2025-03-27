package state

type Hero int

const (
	HeroSamurai Hero = iota
	HeroRobot
	HeroSkeleton
)

type GlobalGameState struct {
	SelectedHero Hero
	DebugMode bool
}
