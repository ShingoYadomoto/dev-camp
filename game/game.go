package game

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"
)

const (
	debug        = false
	baseX        = 100
	jumpingPower = 15
	gravity      = 1
	fontSize     = 10

	// tree kinds
	kindTreeSmall = 0
	kindTreeBig   = 1

	// game modes
	modeTitle    = 0
	modeGame     = 1
	modeGameover = 2

	speed    = 6
	interval = 10
)

var (
	dishImg    *ebiten.Image
	groundImg  *ebiten.Image
	arcadeFont font.Face
)

//go:embed images/dish.png
var byteDishImg []byte

//go:embed images/ground.png
var byteGroundImg []byte

func init() {
	rand.Seed(time.Now().UnixNano())

	img, _, err := image.Decode(bytes.NewReader(byteDishImg))
	if err != nil {
		log.Fatal(err)
	}
	dishImg = ebiten.NewImageFromImage(img)

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

// Game struct
type Game struct {
	mode      int
	count     int
	score     int
	hiscore   int
	dinosaurX int
	dinosaurY int
	gy        int
	dishes    [maxDishCount]*dish
	lastTreeX int
	ground    *ground
}

// NewGame method
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

// Init method
func (g *Game) init() {
	g.hiscore = g.score
	g.count = 0
	g.score = 0
	g.lastTreeX = 0
	g.gy = 0
	for i := 0; i < maxDishCount; i++ {
		g.dishes[i] = &dish{}
	}
	g.ground = &ground{y: groundY - 30}
}

// Update method
func (g *Game) Update() error {
	switch g.mode {
	case modeTitle:
		if g.isKeyJustPressed() {
			g.mode = modeGame
		}
	case modeGame:
		g.count++
		g.score = g.count / 5

		for _, t := range g.dishes {
			if t.visible {
				t.move(speed)
				if t.isOutOfScreen() {
					t.hide()
				}
			} else {
				if g.count-g.lastTreeX > minDishDist && g.count%interval == 0 && rand.Intn(10) == 0 {
					g.lastTreeX = g.count
					t.show()
					break
				}
			}
		}

		g.ground.move(speed)
	case modeGameover:
		if g.isKeyJustPressed() {
			g.init()
			g.mode = modeGame
		}
	}

	return nil
}

// Draw method
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	text.Draw(screen, fmt.Sprintf("Hisore: %d", g.hiscore), arcadeFont, 300, 20, color.Black)
	text.Draw(screen, fmt.Sprintf("Score: %d", g.score), arcadeFont, 500, 20, color.Black)
	var xs [3]int
	var ys [3]int

	if len(g.dishes) > 0 {
		for i, t := range g.dishes {
			xs[i] = t.x
			ys[i] = t.y
		}
	}

	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(
			"g.y: %d\nTree1 x:%d, y:%d\nTree2 x:%d, y:%d\nTree3 x:%d, y:%d",
			g.dinosaurY,
			xs[0],
			ys[0],
			xs[1],
			ys[1],
			xs[2],
			ys[2],
		))
	}

	g.drawGround(screen)
	g.drawDishes(screen)

	switch g.mode {
	case modeTitle:
		text.Draw(screen, "PRESS SPACE KEY", arcadeFont, 245, 240, color.Black)
	case modeGameover:
		text.Draw(screen, "GAME OVER", arcadeFont, 275, 240, color.Black)
	}
}

func (g *Game) drawDishes(screen *ebiten.Image) {
	for _, t := range g.dishes {
		if t.visible {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(t.x), float64(t.y))
			op.Filter = ebiten.FilterLinear
			screen.DrawImage(dishImg, op)
		}
	}
}

func (g *Game) drawGround(screen *ebiten.Image) {
	for i := 0; i < 14; i++ {
		x := float64(groundWidth * i)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, float64(g.ground.y))
		op.GeoM.Translate(float64(g.ground.x), 0.0)
		op.Filter = ebiten.FilterLinear
		screen.DrawImage(groundImg, op)
	}
}

// Layout method
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenX, ScreenY
}

func (g *Game) isKeyJustPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (g *Game) hit() bool {
	for _, t := range g.dishes {
		if t.visible {
			return false
		}
	}
	return false
}
