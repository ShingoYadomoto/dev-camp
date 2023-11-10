package main

import (
	_ "embed"
	"github.com/ShingoYadomoto/dev-camp/game"
	"github.com/hajimehoshi/ebiten/v2"
	_ "image/png"
	"log"
)

func main() {
	ebiten.SetWindowSize(game.ScreenX, game.ScreenY)
	ebiten.SetWindowTitle("Dinosaur Jump")
	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
