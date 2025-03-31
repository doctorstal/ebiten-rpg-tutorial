package tiled

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path"
	"github.com/doctorstal/ebiten-rpg-tutorial/constants"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Img(id int) *ebiten.Image
}

type TilesetJSON struct {
	Path  string     `json:"image"`
	Tiles []TileJSON `json:"tiles"`
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid
	srcX := constants.TileSize * (id % 22)
	srcY := constants.TileSize * (id / 22)

	return u.img.SubImage(
		image.Rect(
			srcX,
			srcY,
			srcX+constants.TileSize,
			srcY+constants.TileSize,
		),
	).(*ebiten.Image)
}

type DynamicTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynamicTileset) Img(id int) *ebiten.Image {
	id -= d.gid
	return d.imgs[id]
}

func NewTileset(tilesetPath string, gid int) (Tileset, error) {

	contents, err := os.ReadFile(tilesetPath)
	if err != nil {
		return nil, err
	}
	var tilesetJSON TilesetJSON
	err = json.Unmarshal(contents, &tilesetJSON)
	if err != nil {
		return nil, err
	}

	if len(tilesetJSON.Tiles) > 0 {
		// return DynamicTileset

		dynTileset := DynamicTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		for _, tile := range tilesetJSON.Tiles {
			imgPath := path.Join("assets", "maps", "tilesets", tile.Path)
			img, _, err := ebitenutil.NewImageFromFile(imgPath)
			if err != nil {
				fmt.Println(imgPath)
				return nil, err
			}
			dynTileset.imgs = append(dynTileset.imgs, img)
		}
		return &dynTileset, nil
	} else {
		imgPath := path.Join("assets", "maps", "tilesets", tilesetJSON.Path)
		img, _, err := ebitenutil.NewImageFromFile(imgPath)
		if err != nil {

			return nil, err
		}
		uniTileset := UniformTileset{
			img: img,
			gid: gid,
		}
		return &uniTileset, nil
	}
}
