package entities

import (
	"image"
	"math"
	"rpg-tutorial/animations"
	"rpg-tutorial/constants"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnergyBall struct {
	*Sprite
	AmtDamage uint
	MoveSpeed float64
}

// HitRect implements AttackItem.
func (e *EnergyBall) HitRect() image.Rectangle {
	return e.Rect()
}

// DoDamage implements AttackItem.
func (e *EnergyBall) DoDamage() {
	e.state = Dead
}

// GetAmtDamage implements AttackItem.
func (e *EnergyBall) GetAmtDamage() uint {
	return e.AmtDamage
}

// GetSprite implements AttackItem.
func (e *EnergyBall) GetSprite() *Sprite {
	return e.Sprite
}

// ShouldRemove implements AttackItem.
func (e *EnergyBall) ShouldRemove() bool {
	return e.MoveSpeed < 5 || e.state == Dead
}

// Update implements AttackItem.
func (e *EnergyBall) Update() {
	e.UpdateAnimation()
	e.Forward(e.MoveSpeed)
	e.MoveSpeed *= 0.96
	e.Move()
}

func NewEnergyBall(img *ebiten.Image, x, y float64, dmg uint, dir int, speed float64) AttackItem {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-constants.TileSize*0.5, -constants.TileSize*0.5)
	opts.GeoM.Rotate(float64(2+dir) * 0.5 * math.Pi)
	opts.GeoM.Translate(constants.TileSize*0.5, constants.TileSize*0.5)
	return &EnergyBall{
		Sprite: &Sprite{
			X:           x,
			Y:           y,
			Width:       constants.TileSize,
			Height:      constants.TileSize,
			Img:         img,
			Spritesheet: spritesheet.NewSpriteSheet(4, 1, 16),
			Animations: map[SpriteState]animations.Animation{
				Idle: animations.NewLoopAnimation(0, 3, 1, 1.0),
				Dead: animations.NewOneTimeAnimation(0, 3, 1, 2.0, true),
			},
			Direction: dir,
			state:     Idle,
			drawOpts:  opts,
		},
		AmtDamage: dmg,
		MoveSpeed: speed,
	}
}
