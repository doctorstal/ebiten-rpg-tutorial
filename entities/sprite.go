package entities

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Img       *ebiten.Image
	X, Y      float64
	Direction int
	Frame     int
}

func (s *Sprite) Dist(other *Sprite) float64 {
	return math.Sqrt(math.Pow(s.X-other.X, 2) + math.Pow(s.Y-other.Y, 2))
}

func (s *Sprite) Move(d float64) {
	switch s.Direction {
	case 0:
		s.Y += d
	case 1:
		s.Y -= d
	case 2:
		s.X -= d
	case 3:
		s.X += d
	}
	s.Frame = (s.Frame + 1) % 256
}

