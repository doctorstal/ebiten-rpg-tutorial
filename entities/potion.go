package entities

type Potion struct {
	*Sprite
	AmtHeal  uint
	Consumed bool
}
