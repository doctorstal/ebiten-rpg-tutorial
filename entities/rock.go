package entities

import (
	"image"
	"math"
	"rpg-tutorial/animations"
	"rpg-tutorial/constants"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type Rock struct {
	*Sprite
	AmtDamage uint
	MoveSpeed float64
}

// HitRect implements AttackItem.
func (b *Rock) HitRect() image.Rectangle {
	return b.Rect()
}

// DoDamage implements AttackItem.
func (b *Rock) DoDamage() {
	b.state = Dead
}

// GetAmtDamage implements AttackItem.
func (b *Rock) GetAmtDamage() uint {
	return b.AmtDamage
}

// GetSprite implements AttackItem.
func (b *Rock) GetSprite() *Sprite {
	return b.Sprite
}

// ShouldRemove implements AttackItem.
func (b *Rock) ShouldRemove() bool {
	return b.MoveSpeed < 0.1 || b.state == Dead
}

// Update implements AttackItem.
func (b *Rock) Update() {
	b.UpdateAnimation()
	b.Forward(b.MoveSpeed)
	b.MoveSpeed *= 0.95
	b.Move()
}

func NewRock(img *ebiten.Image, x, y float64, dmg uint, dir int, speed float64) AttackItem {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-constants.TileSize*0.5, -constants.TileSize*0.5)
	opts.GeoM.Rotate(float64(2+dir) * 0.5 * math.Pi)
	opts.GeoM.Translate(constants.TileSize*0.5, constants.TileSize*0.5)
	return &Rock{
		Sprite: &Sprite{
			X:           x,
			Y:           y,
			Width:       constants.TileSize,
			Height:      constants.TileSize,
			Img:         img,
			Spritesheet: spritesheet.NewSpriteSheet(4, 1, 16),
			Animations: map[SpriteState]animations.Animation{
				Idle: animations.NewLoopAnimation(0, 3, 1, 1.0),
				Dead: animations.NewSingleFrameAnimation(0),
			},
			Direction: dir,
			state:     Idle,
			drawOpts:  opts,
		},
		AmtDamage: dmg,
		MoveSpeed: speed,
	}
}
