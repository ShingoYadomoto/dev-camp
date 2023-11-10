package main

import (
	twenty48 "github.com/hajimehoshi/ebiten/examples/2048/2048"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	initialScreen *ebiten.Image
)

type Game struct{}

func (g *Game) Update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// 画面クリア
	err := screen.Fill(color.White)
	if err != nil {
		log.Fatal(err)
	}

	// 初期画面の表示
	op := &ebiten.DrawImageOptions{}
	err = screen.DrawImage(initialScreen, op)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 描画はUpdateメソッドで行っているので、ここでは何もしない
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenWidth, screenHeight
}

func main() {
	game, err := twenty48.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	ebiten.SetWindowSize(twenty48.ScreenWidth, twenty48.ScreenHeight)
	ebiten.SetWindowTitle("2048 (Ebitengine Demo)")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
