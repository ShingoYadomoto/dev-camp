package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	initialScreen *ebiten.Image
)

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// 画面クリア
	screen.Fill(ebitenutil.ColorRGB8(255, 255, 255))

	// 初期画面の表示
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(initialScreen, op)

	return nil
}

func main() {
	// ウィンドウの作成
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Mahjong Game"); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// 初期画面の読み込み
	img, _, err := ebitenutil.NewImageFromFile("initial_screen.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	initialScreen = img
}
