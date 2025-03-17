package spritesheet

import "image"

type SpriteSheet struct {
	WidthInTiles  int
	HeightInTiles int
	TileSize      int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := index % s.WidthInTiles * s.TileSize
	y := index / s.WidthInTiles * s.TileSize
	return image.Rect(x, y, x+s.TileSize, y+s.TileSize)
}

// Creates new sprite sheet, takes width and height in tiles and tile size
func NewSpriteSheet(w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		w,
		h,
		t,
	}
}
