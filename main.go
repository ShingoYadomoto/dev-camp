package main

import (
	"github.com/ShingoYadomoto/dev-camp/game"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("麻雀")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
