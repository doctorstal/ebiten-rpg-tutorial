package tiled

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type TiledMap struct {
	groundImg  *image.NRGBA
	objectsImg *image.NRGBA
	gameMap    *tiled.Map
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

	return &TiledMap{
		groundImg:  gImg,
		objectsImg: objImg,
		gameMap:    gameMap,
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

func (t *TiledMap) ObjectRects() []image.Rectangle {
	rects := make([]image.Rectangle, 0)
	for _, og := range t.gameMap.ObjectGroups {
		for _, o := range og.Objects {
			rects = append(rects, image.Rect(
				int(o.X),
				int(o.Y-o.Height),
				int(o.X+o.Width),
				int(o.Y),
			))
		}
	}
	return rects
}
