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
	Down SpriteState = iota
	Up
	Left
	Right
)

type Sprite struct {
	Img                 *ebiten.Image
	X, Y, Width, Height float64
	Direction           int
	Animations          map[SpriteState]*animations.Animation
	Spritesheet         *spritesheet.SpriteSheet
	Dx, Dy              float64
}

func (s *Sprite) Dist(other *Sprite) float64 {
	return math.Sqrt(math.Pow(s.X+s.Width/2-other.X-other.Width/2, 2) + math.Pow(s.Y+s.Height/2-other.Y-other.Height/2, 2))
}

func (s *Sprite) Move() {
	s.X += s.Dx
	s.Y += s.Dy
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
	return image.Rect(int(s.X), int(s.Y), int(s.X+s.Width), int(s.Y+s.Height))
}

func (s *Sprite) CheckCollision(colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(s.Rect().Add(image.Point{int(s.Dx), 0})) {
			s.Dx = 0
		}
		if collider.Overlaps(s.Rect().Add(image.Point{0, int(s.Dy)})) {
			s.Dy = 0
		}
	}
}

func (p *Sprite) ActiveAnimation() *animations.Animation {
	if p.Dy > 0 {
		return p.Animations[Down]
	}
	if p.Dy < 0 {
		return p.Animations[Up]
	}
	if p.Dx > 0 {
		return p.Animations[Right]
	}
	if p.Dx < 0 {
		return p.Animations[Left]
	}
	return nil
}
