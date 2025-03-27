package entities

import (
	"image"
	"math"
	"math/rand"
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
	deadFrame   int
}

// GetRenderer implements Animator.
func (b *DeadBombAnimator) GetRenderer() Renderer {
	frame := b.animation.Frame()
	z := 0
	if frame >= 15 {
		z = -1
		frame = b.deadFrame
	}
	frameRect := b.spritesheet.Rect(frame)

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
	b.animation.Update()
	return false
}

type Bomb struct {
	*Sprite
	AmtDamage uint
	loader    *resource.Loader
}

// HitRect implements AttackItem.
func (b *Bomb) HitRect() *image.Rectangle {
	return b.Rect()
}

// DoDamage implements AttackItem.
func (b *Bomb) DoDamage() {
	b.state = Dead
	b.drawOpts = &ebiten.DrawImageOptions{}
	b.drawOpts.GeoM.Translate(-16, -16)
	b.drawOpts.GeoM.Rotate( 2 * math.Pi/float64(1+rand.Intn(4)))
	b.drawOpts.GeoM.Translate(16, 16)
	sound := b.loader.LoadAudio(resources.SoundExplosion).Player
	sound.Rewind()
	sound.Play()
}

// GetAmtDamage implements AttackItem.
func (b *Bomb) GetAmtDamage() uint {
	return b.AmtDamage
}

// GetAnimator implements AttackItem.
func (b *Bomb) GetAnimator() Animator {
	if b.state == Dead {
		rect := image.Rect(int(b.X-8), int(b.Y-8), int(b.X+b.Width+8), int(b.Y+b.Height+8))
		img := b.loader.LoadImage(resources.ImgBombExplosion).Data
		return &DeadBombAnimator{
			animation:   animations.NewOneTimeAnimation(0, 15, 1, 1.0, false),
			spritesheet: spritesheet.NewSpriteSheet(4, 5, 32),
			img:         img,
			drawOpts:    b.drawOpts,
			rect:        &rect,
			deadFrame:   16 + rand.Intn(4),
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
		loader:    loader,
	}
}
