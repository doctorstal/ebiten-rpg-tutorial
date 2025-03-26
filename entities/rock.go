package entities

import (
	"image"
	"math"
	"rpg-tutorial/animations"
	"rpg-tutorial/constants"
	"rpg-tutorial/resources"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Rock struct {
	*Sprite
	AmtDamage uint
	MoveSpeed float64
}

// HitRect implements AttackItem.
func (b *Rock) HitRect() *image.Rectangle {
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

// GetAnimator implements AttackItem.
func (b *Rock) GetAnimator() Animator {
	return b.Sprite
}
func (b *Rock) GetRenderer() Renderer {
	return b.Sprite.GetRenderer()
}

// ShouldRemove implements AttackItem.
func (b *Rock) ShouldRemove() bool {
	return b.MoveSpeed < 0.3 || b.state == Dead
}

// Update implements AttackItem.
func (b *Rock) Update() {
	b.UpdateAnimation()
	b.Forward(b.MoveSpeed)
	b.MoveSpeed *= 0.95
	b.Move()
}

func NewRock(loader *resource.Loader, x, y float64, dmg uint, dir int, speed float64) AttackItem {
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
			Img:         loader.LoadImage(resources.ImgRock).Data,
			Spritesheet: spritesheet.NewSpriteSheet(9, 1, 16),
			Animations: map[SpriteState]animations.Animation{
				Idle: animations.NewLoopAnimation(0, 3, 1, 1.0),
				Dead: animations.NewOneTimeAnimation(4, 8, 1, 1.0, true),
			},
			Direction: dir,
			state:     Idle,
			drawOpts:  opts,
		},
		AmtDamage: dmg,
		MoveSpeed: speed,
	}
}
