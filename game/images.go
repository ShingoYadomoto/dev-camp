package game

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/draw"
	"log"
	"math"
	"math/rand"
	"path"
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

func generateRandomTatsu() (image *ebiten.Image, correctFu, dummyFu uint8) {
	img, cfu, dfu := generateRandomTatsuImage()

	return ebiten.NewImageFromImage(img), cfu, dfu
}

// generate random tatsu

type (
	TatsuType uint8
)

const (
	TatsuTypeShuntsu TatsuType = iota + 1
	TatsuTypeAnko
	TatsuTypeMinko
	TatsuTypeAnkan
	TatsuTypeMinkan
	TatsuTypeHead
	TatsuTypeMachi
)

func RandomTatsuType() TatsuType {
	rand.Seed(time.Now().UnixNano())
	types := []TatsuType{TatsuTypeShuntsu, TatsuTypeAnko, TatsuTypeMinko, TatsuTypeAnkan, TatsuTypeMinkan, TatsuTypeHead, TatsuTypeMachi}
	return types[rand.Intn(len(types))]
}

type (
	PaiType  uint8
	PaiIndex uint8
	PaiDir   uint8

	Pai struct {
		Type  PaiType
		Index PaiIndex
		Dir   PaiDir
	}
)

func (p Pai) Image() image.Image {
	pathList := []string{"images", "pai"}

	switch p.Type {
	case PaiTypeZi:
		pathList = append(pathList, "zi")
	case PaiTypeManzu:
		pathList = append(pathList, "manzu")
	case PaiTypeSozu:
		pathList = append(pathList, "sozu")
	case PaiTypePinzu:
		pathList = append(pathList, "pinzu")
	}

	switch p.Dir {
	case PaiDirHorizontal:
		pathList = append(pathList, "horizontal")
	case PaiDirVertical:
		pathList = append(pathList, "vertical")
	}

	pathList = append(pathList, fmt.Sprintf("%d.png", p.Index))

	imgByte, err := files.ReadFile(path.Join(pathList...))
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatal(err)
	}

	return img
}

const (
	PaiTypeZi PaiType = iota + 1
	PaiTypeManzu
	PaiTypeSozu
	PaiTypePinzu
)

func RandomPaiType() PaiType {
	rand.Seed(time.Now().UnixNano())
	types := []PaiType{PaiTypeZi, PaiTypeManzu, PaiTypeSozu, PaiTypePinzu}
	return types[rand.Intn(len(types))]
}

const (
	PaiDirHorizontal PaiDir = iota + 1
	PaiDirVertical
)

const (
	PaiIndex1 PaiIndex = iota + 1
	PaiIndex2
	PaiIndex3
	PaiIndex4
	PaiIndex5
	PaiIndex6
	PaiIndex7
	PaiIndex8
	PaiIndex9
)

func RandomPaiIndex() PaiIndex {
	rand.Seed(time.Now().UnixNano())
	indexes := []PaiIndex{
		PaiIndex1, PaiIndex2, PaiIndex3, PaiIndex4, PaiIndex5,
		PaiIndex6, PaiIndex7, PaiIndex8, PaiIndex9,
	}
	return indexes[rand.Intn(len(indexes))]
}

func generateRandomTatsuImage() (img image.Image, correctFu, dummyFu uint8) {
	var (
		pl = make([]Pai, 0, 4)
		fu uint8

		setShuntsu = func() {
			var (
				paiType    PaiType
				startIndex PaiIndex
			)

			for {
				paiType = RandomPaiType()
				if paiType != PaiTypeZi {
					break
				}
			}

			for {
				startIndex = RandomPaiIndex()
				if PaiIndex7 >= startIndex {
					break
				}
			}

			for i := 0; i < 3; i++ {
				pl = append(pl, Pai{
					Type:  paiType,
					Index: startIndex + PaiIndex(i),
					Dir:   PaiDirVertical,
				})
			}
		}

		setKotsuOrKantsu = func(tatsuType TatsuType) (isYaochu bool) {
			var (
				paiType = RandomPaiType()
				index   PaiIndex
			)

			for {
				index = RandomPaiIndex()
				if paiType != PaiTypeZi {
					break
				}

				if index <= PaiIndex7 {
					break
				}
			}

			if paiType == PaiTypeZi || (index == PaiIndex1 || index == PaiIndex9) {
				isYaochu = true
			}

			paiNum := 0
			switch tatsuType {
			case TatsuTypeAnko, TatsuTypeMinko:
				paiNum = 3
			case TatsuTypeAnkan, TatsuTypeMinkan:
				paiNum = 4
			default:
				panic("invalid tatsu type")
			}

			for i := 0; i < paiNum; i++ {
				dir := PaiDirVertical
				isFulou := tatsuType == TatsuTypeMinko || tatsuType == TatsuTypeMinkan
				if isFulou && i == 1 {
					dir = PaiDirHorizontal
				}

				pl = append(pl, Pai{
					Type:  paiType,
					Index: index,
					Dir:   dir,
				})
			}
			return
		}

		setHead = func() (isYakuhai bool) {
			var (
				paiType = RandomPaiType()
				index   PaiIndex
			)

			for {
				index = RandomPaiIndex()
				if paiType != PaiTypeZi {
					break
				}

				if index <= PaiIndex7 {
					break
				}
			}

			// ToDo: 自風場風によって変える
			if paiType == PaiTypeZi {
				isYakuhai = true
			}

			for i := 0; i < 2; i++ {
				pl = append(pl, Pai{
					Type:  paiType,
					Index: index,
					Dir:   PaiDirVertical,
				})
			}

			return
		}

		setMachi = func() (isRyanmen bool) {
			setShuntsu()

			rand.Seed(time.Now().UnixNano())
			indexes := make([]PaiIndex, len(pl))
			for i, p := range pl {
				indexes[i] = p.Index
			}

			removeIndex := rand.Intn(len(indexes))

			switch removeIndex {
			case 0:
				if pl[2].Index != PaiIndex9 {
					isRyanmen = true
				}
			case 2:
				if pl[0].Index != PaiIndex1 {
					isRyanmen = true
				}
			}

			pl = append(pl[:removeIndex], pl[removeIndex+1:]...)

			return
		}
	)

	tatsuType := RandomTatsuType()
	switch tatsuType {
	case TatsuTypeShuntsu:
		setShuntsu()
	case TatsuTypeMinko:
		isYaochu := setKotsuOrKantsu(tatsuType)
		if isYaochu {
			fu = 4
		} else {
			fu = 2
		}
	case TatsuTypeAnko:
		isYaochu := setKotsuOrKantsu(tatsuType)
		if isYaochu {
			fu = 8
		} else {
			fu = 4
		}
	case TatsuTypeMinkan:
		isYaochu := setKotsuOrKantsu(tatsuType)
		if isYaochu {
			fu = 16
		} else {
			fu = 8
		}
	case TatsuTypeAnkan:
		isYaochu := setKotsuOrKantsu(tatsuType)
		if isYaochu {
			fu = 32
		} else {
			fu = 16
		}
	case TatsuTypeHead:
		isYakuhai := setHead()
		if isYakuhai {
			fu = 2
		}
	case TatsuTypeMachi:
		isRyanmen := setMachi()
		if !isRyanmen {
			fu = 2
		}
	}

	images := make([]image.Image, len(pl))
	for i, p := range pl {
		images[i] = p.Image()
	}

	dummyFu = func() uint8 {
		rand.Seed(time.Now().UnixNano())
		d := math.Pow(2, float64(rand.Intn(6)))
		if d == 1 {
			return 0
		}
		return uint8(d)
	}()

	return joinImages(images...), fu, dummyFu
}
