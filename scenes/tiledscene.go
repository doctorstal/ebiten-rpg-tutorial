package scenes

import (
	"fmt"
	"image"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

type TiledScene struct {
	loaded bool
	img    *image.NRGBA
}

// Draw implements Scene.
func (t *TiledScene) Draw(screen *ebiten.Image) {
	img := t.img.SubImage(image.Rect(0, 0, 320, 240))
	screen.DrawImage(ebiten.NewImageFromImage(img), nil)
}

// FirstLoad implements Scene.
func (t *TiledScene) FirstLoad() {
	// Parse .tmx file.
	gameMap, err := tiled.LoadFile("assets/maps/first.tmx")
	if err != nil {
		fmt.Printf("error parsing map: %s\n", err.Error())
		os.Exit(2)
	}

	// gameMap.Layers

	fmt.Println(gameMap)

	// You can also render the map to an in-memory image for direct
	// use with the default Renderer, or by making your own.
	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}

	// Render just layer 0 to the Renderer.
	err = renderer.RenderVisibleLayers()
	if err != nil {
		fmt.Printf("layer unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}

	// Get a reference to the Renderer's output, an image.NRGBA struct.
	t.img = renderer.Result

	// Clear the render result after copying the output if separation of
	// layers is desired.
	renderer.Clear()

	// And so on. You can also export the image to a file by using the
	// Renderer's Save functions.
}

// IsLoaded implements Scene.
func (t *TiledScene) IsLoaded() bool {
	return t.loaded
}

// OnEnter implements Scene.
func (t *TiledScene) OnEnter() {

}

// OnExit implements Scene.
func (t *TiledScene) OnExit() {

}

// Update implements Scene.
func (t *TiledScene) Update() SceneId {
	return StartSceneId
}

func NewTiledScene() Scene {
	return &TiledScene{}
}
