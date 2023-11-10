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
	dishImg    *ebiten.Image
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
		height = rct.Dy()
		srcImages = append(srcImages, img)
	}

	dstImage := image.NewRGBA(image.Rect(0, 0, width, height))

	offset := 0
	for _, img := range srcImages {
		srcRect := img.Bounds()
		draw.Draw(
			dstImage,
			image.Rect(offset, 0, offset+srcRect.Dx(), srcRect.Dy()),
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

	//manzu11, _ := files.ReadFile("images/pai/manzu/horizontal/1.png")
	manzu11, _ := files.ReadFile("images/pai/manzu/vertical/1.png")

	img, _, err := image.Decode(bytes.NewReader(manzu11))
	if err != nil {
		log.Fatal(err)
	}

	dishImg = ebiten.NewImageFromImage(joinImages([]image.Image{img, img, img}...))

	img, _, err = image.Decode(bytes.NewReader(byteGroundImg))
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
