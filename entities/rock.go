package entities

import (
	"image"
	"math"
	"math/rand"

	"github.com/doctorstal/ebiten-rpg-tutorial/animations"
	"github.com/doctorstal/ebiten-rpg-tutorial/constants"
	"github.com/doctorstal/ebiten-rpg-tutorial/resources"
	"github.com/doctorstal/ebiten-rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type DeadRockAnimator struct {
	animation   animations.Animation
	spritesheet *spritesheet.SpriteSheet
	img         *ebiten.Image
	drawOpts    *ebiten.DrawImageOptions
	rect        *image.Rectangle
	deadFrame   int
}

// GetRenderer implements Animator.
func (b *DeadRockAnimator) GetRenderers() Renderer {
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
func (b *DeadRockAnimator) UpdateAnimation() bool {
	b.animation.Update()
	return false
}

type Rock struct {
	*Sprite
	AmtDamage uint
	MoveSpeed float64
	loader    *resource.Loader
}

// HitRect implements AttackItem.
func (b *Rock) HitRect() *image.Rectangle {
	return b.Rect()
}

// DoDamage implements AttackItem.
func (b *Rock) DoDamage() {
	b.state = Dead
	expSound := b.loader.LoadAudio(resources.SoundRockSmash).Player
	expSound.Rewind()
	expSound.Play()
}

// GetAmtDamage implements AttackItem.
func (b *Rock) GetAmtDamage() uint {
	return b.AmtDamage
}

// GetAnimator implements AttackItem.
func (b *Rock) GetAnimator() Animator {
	if b.state == Dead {
		rect := image.Rect(int(b.X-7), int(b.Y-7), int(b.X+b.Width+7), int(b.Y+b.Height+7))
		img := b.loader.LoadImage(resources.ImgRockExplosion).Data
		return &DeadBombAnimator{
			animation:   animations.NewOneTimeAnimation(9, 14, 1, 3.0, false),
			spritesheet: spritesheet.NewSpriteSheet(14, 1, 30),
			img:         img,
			drawOpts:    GetRotationOpts(32.0, (b.Direction+2)%4),
			rect:        &rect,
			deadFrame:   16 + rand.Intn(4),
		}
	} else {
		return b.Sprite
	}
}

func (b *Rock) GetRenderers() []Renderer {
	return b.Sprite.GetRenderers()
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
	opts := GetRotationOpts(16.0, dir)
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
		loader:    loader,
	}
}

func GetRotationOpts(frameWidth float64, dir int) *ebiten.DrawImageOptions {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-frameWidth*0.5, -frameWidth*0.5)
	opts.GeoM.Rotate(float64(2+dir) * 0.5 * math.Pi)
	opts.GeoM.Translate(frameWidth*0.5, frameWidth*0.5)
	return opts
}
