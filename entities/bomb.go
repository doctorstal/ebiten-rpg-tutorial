package entities

import (
	"image"
	"rpg-tutorial/animations"
	"rpg-tutorial/constants"
	"rpg-tutorial/resources"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type DeadBombAnimator struct {
	animation   animations.Animation
	spritesheet *spritesheet.SpriteSheet
	img         *ebiten.Image
	drawOpts    *ebiten.DrawImageOptions
	rect        *image.Rectangle
}

// GetRenderer implements Animator.
func (b *DeadBombAnimator) GetRenderer() Renderer {
	frame := b.animation.Frame()
	frameRect := b.spritesheet.Rect(frame)
	z := 0
	if frame >= 6 {
		z = -1
	}

	return &BasicRenderer{
		img: b.img.SubImage(
			frameRect,
		).(*ebiten.Image),
		drawOpts: b.drawOpts,
		rect:     b.rect,
		z:        z,
	}
}

// UpdateAnimation implements Animator.
func (b *DeadBombAnimator) UpdateAnimation() bool {
	return b.animation.Update()
}

type Bomb struct {
	*Sprite
	AmtDamage uint
}

// HitRect implements AttackItem.
func (b *Bomb) HitRect() *image.Rectangle {
	return b.Rect()
}

// DoDamage implements AttackItem.
func (b *Bomb) DoDamage() {
	b.state = Dead
	b.drawOpts = &ebiten.DrawImageOptions{}
	b.drawOpts.ColorScale.SetA(0.8)
}

// GetAmtDamage implements AttackItem.
func (b *Bomb) GetAmtDamage() uint {
	return b.AmtDamage
}

// GetAnimator implements AttackItem.
func (b *Bomb) GetAnimator() Animator {
	if b.state == Dead {
		rect := image.Rect(int(b.X), int(b.Y), int(b.X+b.Width), int(b.Y+b.Height))
		return &DeadBombAnimator{
			animation:   b.Animations[Dead],
			spritesheet: b.Spritesheet,
			img:         b.Img,
			drawOpts:    b.drawOpts,
			rect:        &rect,
		}
	} else {
		return b.Sprite
	}
}

// ShouldRemove implements AttackItem.
func (b *Bomb) ShouldRemove() bool {
	return b.state == Dead
}

// Update implements AttackItem.
func (b *Bomb) Update() {
	b.UpdateAnimation()
}

func NewBomb(loader *resource.Loader, x, y float64, dmg uint) AttackItem {
	return &Bomb{
		Sprite: &Sprite{
			X:           x,
			Y:           y,
			Width:       constants.TileSize,
			Height:      constants.TileSize,
			Img:         loader.LoadImage(resources.ImgBomb).Data,
			Spritesheet: spritesheet.NewSpriteSheet(1, 8, 16),
			Animations: map[SpriteState]animations.Animation{
				Idle: animations.NewLoopAnimation(0, 2, 1, 10.0),
				Dead: animations.NewOneTimeAnimation(3, 7, 1, 5.0, false),
			},
			state: Idle,
		},
		AmtDamage: dmg,
	}
}
