package entities

import (
	"image"
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
			drawOpts: &ebiten.DrawImageOptions{},
		},
		CombatComponent: combat,
	}

}

func personAnimations() map[SpriteState]animations.Animation {
	return map[SpriteState]animations.Animation{
		Down:        animations.NewLoopAnimation(0, 12, 4, 10.0),
		Up:          animations.NewLoopAnimation(1, 13, 4, 10.0),
		Left:        animations.NewLoopAnimation(2, 14, 4, 10.0),
		Right:       animations.NewLoopAnimation(3, 15, 4, 10.0),
		AttackDown:  animations.NewLoopAnimation(0, 16, 16, 10.0),
		AttackUp:    animations.NewLoopAnimation(1, 17, 16, 10.0),
		AttackLeft:  animations.NewLoopAnimation(2, 18, 16, 10.0),
		AttackRight: animations.NewLoopAnimation(3, 19, 16, 10.0),
		Dead:        animations.NewOneTimeAnimation(0, 24, 24, 10.0, false),
	}
}

func (s *Character) Rect() image.Rectangle{
	return image.Rect(int(s.X), int(s.Y), int(s.X+s.Width), int(s.Y+s.Height))

}

func (c *Character) Move() {
	if !c.CombatComponent.Attacking() && !c.CombatComponent.Damaged() {
		c.Sprite.Move()
	}
}

func (c *Character) Die() {
	c.state = Dead
}

func (c *Character) UpdateState() {
	if !c.CombatComponent.Attacking() {
		if c.CombatComponent.Damaged() {
			c.drawOpts.ColorM.Scale(0, 0, 0, 1)
			c.drawOpts.ColorM.Translate(1, 1, 1, 0)
		} else {
			c.drawOpts.ColorM.Reset()
		}
		c.Sprite.UpdateState()
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
}
