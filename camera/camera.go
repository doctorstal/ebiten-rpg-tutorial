package camera

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	X, Y                                                  float64
	screenWith, screenHeight, tilemapWidth, tilemapHeight float64
}

func NewCamera(screenWidth, screenHeight, tilemapWith, tilemapHeight float64) *Camera {
	return &Camera{
		X:             0,
		Y:             0,
		screenWith:    screenWidth,
		screenHeight:  screenHeight,
		tilemapWidth:  tilemapWith,
		tilemapHeight: tilemapHeight,
	}
}

func (c *Camera) FollowTarget(targetX, targetY float64) {
	tx := -targetX + c.screenWith/2
	ty := -targetY + c.screenHeight/2

	if math.Abs(tx-c.X) > 3.0 {
		c.X += math.Floor((tx - c.X) / 15)
	} else {
		c.X += math.Copysign(1.0, tx)
	}

	if math.Abs(ty-c.Y) > 3.0 {
		c.Y += math.Floor((ty - c.Y) / 15)
	} else {
		c.Y += math.Copysign(1.0, ty)
	}

	c.X = math.Max(math.Min(c.X, 0), c.screenWith-c.tilemapWidth)
	c.Y = math.Max(math.Min(c.Y, 0), c.screenHeight-c.tilemapHeight)
}

func (c *Camera) Render(screen *ebiten.Image, subimage *ebiten.Image, x, y float64) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(
		float64(c.X+x),
		float64(c.Y+y),
	)

	screen.DrawImage(
		subimage,
		opts,
	)
}
