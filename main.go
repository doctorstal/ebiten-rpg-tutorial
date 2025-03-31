package main

import (
	"embed"
	"fmt"
	"log"
	"github.com/doctorstal/ebiten-rpg-tutorial/resources"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//go:embed assets
var assets embed.FS

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Skeleton Bomber")
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	audioContext := audio.NewContext(44100)
	loader := resources.NewResourceLoader(assets, audioContext)
	game := NewGame(loader)

	fmt.Println("Starting...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok Bye!")
}
