package tiled

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type TiledMap struct {
	groundImg   *image.NRGBA
	objectsImg  *image.NRGBA
	gameMap     *tiled.Map
	objectRects []*image.Rectangle
	doors       map[string]*image.Rectangle
}

func NewTiledMap(path string) (*TiledMap, error) {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFile(path)
	if err != nil {
		fmt.Printf("error parsing map: %s\n", err.Error())
		return nil, err
	}

	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		return nil, err
	}

	err = renderer.RenderVisibleLayers()
	if err != nil {
		fmt.Printf("layer unsupported for rendering: %s", err.Error())
		return nil, err
	}

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	gImg := renderer.Result
	renderer.Clear()

	err = renderer.RenderVisibleObjectGroups()
	if err != nil {
		fmt.Printf("object group unsupported for rendering: %s", err.Error())
		return nil, err
	}

	objImg := renderer.Result

	doors := make(map[string]*image.Rectangle)
	rects := make([]*image.Rectangle, 0)
	for _, og := range gameMap.ObjectGroups {
		for _, o := range og.Objects {
			if o.Type == "door" {
				rect := image.Rect(
					int(o.X),
					int(o.Y),
					int(o.X+o.Width),
					int(o.Y+o.Height),
				)
				doors[o.Properties.GetString("goto")] = &rect
			} else {
				rect := image.Rect(
					int(o.X),
					int(o.Y-o.Height),
					int(o.X+o.Width),
					int(o.Y),
				)
				rects = append(rects, &rect)

			}
		}
	}

	return &TiledMap{
		groundImg:   gImg,
		objectsImg:  objImg,
		gameMap:     gameMap,
		objectRects: rects,
		doors:       doors,
	}, err
}

func (t *TiledMap) GroundImage(rect image.Rectangle) *ebiten.Image {
	subImage := t.groundImg.SubImage(rect)
	return ebiten.NewImageFromImage(subImage)
}

func (t *TiledMap) ObjectsImage(rect image.Rectangle) *ebiten.Image {
	subImage := t.objectsImg.SubImage(rect)
	return ebiten.NewImageFromImage(subImage)
}

func (t *TiledMap) ObjectRects() []*image.Rectangle {
	return t.objectRects
}
func (t *TiledMap) Doors() map[string]*image.Rectangle {
	return t.doors
}

func (t *TiledMap) Width() float64 {
	return float64(t.gameMap.Width * t.gameMap.TileWidth)
}

func (t *TiledMap) Height() float64 {
	return float64(t.gameMap.Height * t.gameMap.TileHeight)
}
