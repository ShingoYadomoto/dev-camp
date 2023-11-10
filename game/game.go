package game

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math/rand"
)

const (
	debug        = false
	baseX        = 100
	jumpingPower = 15
	gravity      = 1
	fontSize     = 10

	// game modes
	modeTitle    = 0
	modeGame     = 1
	modeGameover = 2

	speed    = 6
	interval = 15
)

// Game struct
type Game struct {
	mode  int
	count int

	score   int
	hiscore int

	tatsus           [maxTatsuCount]*tatsu
	targetTatsuIndex int
	lastPaiX         int
	ground           *ground
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
	g.lastPaiX = 0
	for i := 0; i < maxTatsuCount; i++ {
		image, correctFu, dummyFu := generateRandomTatsu()
		g.tatsus[i] = &tatsu{
			i:         image,
			correctFu: correctFu,
			dummyFu:   dummyFu,
		}
	}
	g.ground = &ground{y: groundY - 10}
}

// Update method
func (g *Game) Update() error {
	switch g.mode {
	case modeTitle:
		if g.isSpacePressed() {
			g.mode = modeGame
		}
	case modeGame:
		g.count++
		g.score = g.count / 5

		for _, t := range g.tatsus {
			if t.visible {
				t.move(speed)
				if t.isOutOfScreen() {
					g.mode = modeGameover
				}
			} else {
				if g.count-g.lastPaiX > minTatsuDist && g.count%interval == 0 && rand.Intn(10) == 0 {
					g.lastPaiX = g.count
					t.show()
					break
				}
			}
		}

		g.ground.move(speed)

		if g.isSpacePressed() && !g.tatsus[g.targetTatsuIndex].answer(true) {
			g.mode = modeGameover
			g.targetTatsuIndex++
		}
		if g.isEnterPressed() && !g.tatsus[g.targetTatsuIndex].answer(false) {
			g.mode = modeGameover
			g.targetTatsuIndex++
		}
	case modeGameover:
		if g.isSpacePressed() {
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
	var xs [maxTatsuCount]int
	var ys [maxTatsuCount]int

	if len(g.tatsus) > 0 {
		for i, t := range g.tatsus {
			xs[i] = t.x
			ys[i] = t.y
		}
	}

	if debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(
			"Tree1 x:%d, y:%d\nTree2 x:%d, y:%d\nTree3 x:%d, y:%d",
			xs[0],
			ys[0],
			xs[1],
			ys[1],
			xs[2],
			ys[2],
		))
	}

	g.drawGround(screen)
	g.drawTatsus(screen)

	var (
		titleX = 425
		titleY = 300
	)

	switch g.mode {
	case modeTitle:
		text.Draw(screen, "PRESS SPACE KEY", arcadeFont, titleX, titleY, color.Black)
	case modeGameover:
		text.Draw(screen, "GAME OVER", arcadeFont, titleX, titleY, color.Black)
	}
}

func (g *Game) drawTatsus(screen *ebiten.Image) {
	for _, t := range g.tatsus {
		if t.visible {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(t.x), float64(t.y))
			op.Filter = ebiten.FilterLinear
			screen.DrawImage(t.i, op)

			tt, err := opentype.Parse(fonts.PressStart2P_ttf)
			if err != nil {
				log.Fatal(err)
			}
			const dpi = 72
			arcadeFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    fontSize + 10,
				DPI:     dpi,
				Hinting: font.HintingFull,
			})

			text.Draw(screen, fmt.Sprint(t.dummyFu), arcadeFont, t.x+(t.i.Bounds().Dx()/2)-10, t.y-30, color.Black)
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

func (g *Game) isSpacePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	return false
}

func (g *Game) isEnterPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}
	return false
}

func (g *Game) hit() bool {
	for _, t := range g.tatsus {
		if t.visible {
			return false
		}
	}
	return false
}
