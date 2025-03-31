package state

type Hero int

const (
	HeroSamurai Hero = iota
	HeroRobot
	HeroBoy
)

type GlobalGameState struct {
	SelectedHero Hero
	DebugMode    bool
}
