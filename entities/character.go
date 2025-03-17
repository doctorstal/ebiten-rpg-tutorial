package entities

import (
	"rpg-tutorial/animations"
	"rpg-tutorial/components"
	"rpg-tutorial/constants"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
	*Sprite
	CombatComponent components.Combat
}

func NewCharacter(img *ebiten.Image, x, y float64, combat components.Combat) *Character {
	return &Character{
		Sprite: &Sprite{
			Img:         img,
			X:           x,
			Y:           y,
			Width:       constants.TileSize,
			Height:      constants.TileSize,
			Spritesheet: spritesheet.NewSpriteSheet(4, 7, 16),
			Animations:  personAnimations(),
		},
		CombatComponent: combat,
	}

}

func personAnimations() map[SpriteState]*animations.Animation {
	return map[SpriteState]*animations.Animation{
		Down:        animations.NewAnimation(0, 12, 4, 10.0),
		Up:          animations.NewAnimation(1, 13, 4, 10.0),
		Left:        animations.NewAnimation(2, 14, 4, 10.0),
		Right:       animations.NewAnimation(3, 15, 4, 10.0),
		AttackDown:  animations.NewAnimation(0, 16, 16, 10.0),
		AttackUp:    animations.NewAnimation(1, 17, 16, 10.0),
		AttackLeft:  animations.NewAnimation(2, 18, 16, 10.0),
		AttackRight: animations.NewAnimation(3, 19, 16, 10.0),
		Dead:        animations.NewAnimation(24, 27, 1, 10.0),
	}
}

func (c *Character) Move() {
	if !c.CombatComponent.Attacking() {
		c.Sprite.Move()
	}
}

func (c *Character) Die() {
	c.state = Dead
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
