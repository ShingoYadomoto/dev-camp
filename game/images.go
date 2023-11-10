package game

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/draw"
	"log"
	"math/rand"
	"time"
)

var (
	tatsuImg   *ebiten.Image
	groundImg  *ebiten.Image
	arcadeFont font.Face
)

//go:embed images/tenbou.png
var byteGroundImg []byte

//go:embed images/pai/*
var files embed.FS

func joinImages(srcImages ...image.Image) image.Image {
	width, height := 0, 0
	for _, img := range srcImages {
		rct := img.Bounds()
		width += rct.Dx()
		if rct.Dy() > height {
			height = rct.Dy()
		}
	}

	dstImage := image.NewRGBA(image.Rect(0, 0, width, height))

	offset := 0
	for _, img := range srcImages {
		srcRect := img.Bounds()
		rect := image.Rect(offset, height-srcRect.Dy(), offset+srcRect.Dx(), height)

		draw.Draw(
			dstImage,
			rect,
			img,
			image.Point{0, 0},
			draw.Over,
		)
		offset += srcRect.Dx()
	}

	return dstImage
}

func init() {
	rand.Seed(time.Now().UnixNano())

	img, _, err := image.Decode(bytes.NewReader(byteGroundImg))
	if err != nil {
		log.Fatal(err)
	}
	groundImg = ebiten.NewImageFromImage(img)

	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	arcadeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func generateRandomTatsu() *ebiten.Image {
	manzu11, _ := files.ReadFile("images/pai/manzu/vertical/1.png")

	img, _, err := image.Decode(bytes.NewReader(manzu11))
	if err != nil {
		log.Fatal(err)
	}

	manzu11H, _ := files.ReadFile("images/pai/manzu/horizontal/1.png")

	img2, _, err := image.Decode(bytes.NewReader(manzu11H))
	if err != nil {
		log.Fatal(err)
	}

	tatsuImg = ebiten.NewImageFromImage(joinImages([]image.Image{img, img2, img}...))

	return tatsuImg
}
