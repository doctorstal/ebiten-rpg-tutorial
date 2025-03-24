package entities

import (
	"image"
	"math"
	"rpg-tutorial/animations"
	"rpg-tutorial/spritesheet"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteState uint8

const (
	Idle SpriteState = iota
	Down
	Up
	Left
	Right
	AttackDown
	AttackUp
	AttackLeft
	AttackRight
	Dead
)

type Sprite struct {
	Img                 *ebiten.Image
	X, Y, Width, Height float64
	Direction           int
	Animations          map[SpriteState]animations.Animation
	Spritesheet         *spritesheet.SpriteSheet
	Dx, Dy              float64
	state               SpriteState
}

func (s *Sprite) Dist(other *Sprite) float64 {
	return math.Sqrt(math.Pow(s.X+s.Width/2-other.X-other.Width/2, 2) + math.Pow(s.Y+s.Height/2-other.Y-other.Height/2, 2))
}

func (s *Sprite) NormalizeSpeed() {
	if s.Dx != 0 && s.Dy != 0 {
		s.Dx /= math.Sqrt2
		s.Dy /= math.Sqrt2
	}
}
func (s *Sprite) Move() {
	s.NormalizeSpeed()
	s.X += s.Dx
	s.Y += s.Dy
	if s.X < 0 {
		s.X = 0
	}
	if s.Y < 0 {
		s.Y = 0
	}
	// TODO Add upper bounds
}
func (s *Sprite) Forward(d float64) {
	switch s.Direction {
	case 0:
		s.Dy = d
	case 1:
		s.Dy = -d
	case 2:
		s.Dx = -d
	case 3:
		s.Dx = d
	}
}

func (s *Sprite) Rect() image.Rectangle {
	return image.Rect(int(s.X), int(s.Y+0.5*s.Height), int(s.X+s.Width), int(s.Y+s.Height))
}

func (s *Sprite) CheckCollision(colliders []image.Rectangle) {
	// moveX := s.X + s.Width/2 + s.Dx + math.Copysign(s.Width/2, s.Dx)
	// moveY := s.Y + s.Height/2 + s.Dy + math.Copysign(s.Height/2, s.Dy)

	// xRect := image.Rect(int(s.X+s.Width/2), int(s.Y), int(s.X+s.Width/2+1), int(s.Y+s.Height)).Add(image.Point{int(s.Dx + math.Copysign(s.Width/2, s.Dx)), 0})
	// yRect := image.Rect(int(s.X), int(s.Y+s.Height/2), int(s.X+s.Width), int(s.Y+s.Height/2+1)).Add(image.Point{0, int(s.Dy + math.Copysign(s.Height/2, s.Dy))})

	xRect := image.Rect(int(s.X+s.Width/2), int(s.Y+0.5*s.Height), int(s.X+s.Width/2+1), int(s.Y+s.Height)).Add(image.Point{int(s.Dx + math.Copysign(s.Width/2, s.Dx)), 0})
	yRect := image.Rect(int(s.X), int(s.Y+0.75*s.Height), int(s.X+s.Width), int(s.Y+0.75*s.Height+1)).Add(image.Point{0, int(s.Dy + math.Copysign(0.25*s.Height, s.Dy))})

	for _, collider := range colliders {
		// if (image.Point{int(moveX), int(moveY)}).In(collider) {
		// 	s.Dx = 0
		// 	s.Dy = 0
		// }
		if collider.Overlaps(xRect) {
			s.Dx = 0
		}
		if collider.Overlaps(yRect) {
			s.Dy = 0
		}
	}
}

func (s *Sprite) UpdateState() {
	s.state = Idle
	if s.Dy > 0 {
		s.state = Down
	}
	if s.Dy < 0 {
		s.state = Up
	}
	if s.Dx > 0 {
		s.state = Right
	}
	if s.Dx < 0 {
		s.state = Left
	}
}

func (s *Sprite) UpdateAnimation() {
	animation := s.ActiveAnimation()
	if animation != nil {
		animation.Update()
	}
}

func (s *Sprite) ActiveAnimation() animations.Animation {
	if anim, ok := s.Animations[s.state]; ok {
		return anim
	}
	return nil
}
