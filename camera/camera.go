package camera

import (
	"image"
	"image/color"
	"math"
	"github.com/doctorstal/ebiten-rpg-tutorial/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Camera struct {
	X, Y                                             float64
	screenWidth, screenHeight, roomWidth, roomHeight float64
	gameState                                        *state.GlobalGameState
}

func NewCamera(screenWidth, screenHeight, roomWidth, roomHeight float64, gameState *state.GlobalGameState) *Camera {
	return &Camera{
		X:            0,
		Y:            0,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		roomWidth:    roomWidth,
		roomHeight:   roomHeight,
		gameState:    gameState,
	}
}

func (c *Camera) UpdateRoomSize(width, height float64) {
	c.roomWidth = width
	c.roomHeight = height
}

func (c *Camera) GoToTarget(targetX, targetY float64) {
	c.X = -targetX + c.screenWidth/2
	c.Y = -targetY + c.screenHeight/2
	c.X = math.Max(math.Min(c.X, 0), c.screenWidth-c.roomWidth)
	c.Y = math.Max(math.Min(c.Y, 0), c.screenHeight-c.roomHeight)
}

func (c *Camera) FollowTarget(targetX, targetY float64) {
	tx := -targetX + c.screenWidth/2
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

	c.X = math.Max(math.Min(c.X, 0), c.screenWidth-c.roomWidth)
	c.Y = math.Max(math.Min(c.Y, 0), c.screenHeight-c.roomHeight)
}

func (c *Camera) Render(screen *ebiten.Image, subimage *ebiten.Image, x, y float64, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(
		c.X+x,
		c.Y+y,
	)

	screen.DrawImage(
		subimage,
		opts,
	)
	opts.GeoM.Translate(
		-c.X-x,
		-c.Y-y,
	)
	if c.gameState.DebugMode {
		vector.StrokeRect(screen, float32(c.X+x), float32(c.Y+y), float32(subimage.Bounds().Dx()), float32(subimage.Bounds().Dy()), 1.0, color.RGBA{255, 0, 0, 255}, false)
	}

}

func (c *Camera) ViewRect() image.Rectangle {
	return image.Rect(int(-c.X), int(-c.Y), int(-c.X+c.screenWidth), int(-c.Y+c.screenHeight))
}
