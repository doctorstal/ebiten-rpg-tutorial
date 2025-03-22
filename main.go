package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Skeleton Bomber")
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	fmt.Println("Starting...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ok Bye!")
}
