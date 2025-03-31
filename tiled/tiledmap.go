package tiled

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type Door struct {
	Rect      *image.Rectangle
	Direction string
}

type Enemy struct {
	Rect          *image.Rectangle
	Kind          string
	FollorsPlayer bool
}

type TiledMap struct {
	groundImg   *image.NRGBA
	objectsImg  *image.NRGBA
	gameMap     *tiled.Map
	objectRects []*image.Rectangle
	doors       map[string]*Door
	enemies     []*Enemy
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

	doors := make(map[string]*Door)
	enemies := make([]*Enemy, 0)
	rects := make([]*image.Rectangle, 0)
	for _, og := range gameMap.ObjectGroups {
		for _, o := range og.Objects {
			switch o.Type {
			case "Door":
				rect := image.Rect(
					int(o.X),
					int(o.Y),
					int(o.X+o.Width),
					int(o.Y+o.Height),
				)
				doors[o.Properties.GetString("goto")] = &Door{
					Rect:      &rect,
					Direction: o.Properties.GetString("direction"),
				}
			case "Enemy":
				rect := image.Rect(
					int(o.X),
					int(o.Y-o.Height),
					int(o.X+o.Width),
					int(o.Y),
				)
				enemies = append(enemies, &Enemy{
					Rect:          &rect,
					Kind:          o.Properties.GetString("kind"),
					FollorsPlayer: o.Properties.GetBool("follows_player"),
				})
				o.Visible = false
			default:
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

	err = renderer.RenderVisibleObjectGroups()
	if err != nil {
		fmt.Printf("object group unsupported for rendering: %s", err.Error())
		return nil, err
	}
	objImg := renderer.Result

	return &TiledMap{
		groundImg:   gImg,
		objectsImg:  objImg,
		gameMap:     gameMap,
		objectRects: rects,
		doors:       doors,
		enemies:     enemies,
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
func (t *TiledMap) Doors() map[string]*Door {
	return t.doors
}

func (t *TiledMap) Enemies() []*Enemy {
	return t.enemies
}

func (t *TiledMap) Width() float64 {
	return float64(t.gameMap.Width * t.gameMap.TileWidth)
}

func (t *TiledMap) Height() float64 {
	return float64(t.gameMap.Height * t.gameMap.TileHeight)
}
