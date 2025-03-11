package entities

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img                 *ebiten.Image
	X, Y, Width, Height float64
	Direction           int
	Frame               int
	Dx, Dy              float64
}

func (s *Sprite) Dist(other *Sprite) float64 {
	return math.Sqrt(math.Pow(s.X-other.X, 2) + math.Pow(s.Y-other.Y, 2))
}

func (s *Sprite) Move() {
	s.X += s.Dx
	s.Y += s.Dy
	s.Frame = (s.Frame + 1) % 256
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
