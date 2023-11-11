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
	"math"
	"math/rand"
)

const (
	debug        = false
	baseX        = 100
	jumpingPower = 15
	gravity      = 1
	fontSize     = 20

	// game modes
	modeTitle    = 0
	modeGame     = 1
	modeGameover = 2

	speed    = 3
	interval = 100
)

// Game struct
type Game struct {
	mode  int
	count int

	score int

	pause bool
	speed int

	tatsus           [maxTatsuCount]*tatsu
	showTatsuIndexes []int
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
	g.count = 0
	g.score = 0
	g.speed = speed
	g.showTatsuIndexes = nil
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
		if g.isSPressed() {
			g.pause = !g.pause
		}

		if g.pause {
			return nil
		}

		g.speed = int(math.Max(float64(g.speed), float64(g.score/5)+3))

		g.count++

		for i, t := range g.tatsus {
			if t.visible {
				t.move(g.speed)
				if t.isOutOfScreen() {
					if !t.answered() {
						g.showTatsuIndexes = g.showTatsuIndexes[1:]
					}
					t.hide()
				}
			} else {
				if g.count-g.lastPaiX > minTatsuDist && g.count%interval == 0 && rand.Intn(interval/2) == 0 {
					g.lastPaiX = g.count
					t.show()
					g.showTatsuIndexes = append(g.showTatsuIndexes, i)
					break
				}
			}
		}

		g.ground.move(g.speed)

		if g.isSpacePressed() {
			if len(g.showTatsuIndexes) == 0 {
				g.mode = modeGameover
			} else {
				deleteIdx := g.showTatsuIndexes[0]
				pt := g.tatsus[deleteIdx]
				if pt.answer(true) {
					g.score += pt.point()
				} else {
					g.score -= pt.point()
				}

				g.showTatsuIndexes = g.showTatsuIndexes[1:]
			}
		}
		if g.isEnterPressed() {
			if len(g.showTatsuIndexes) == 0 {
				g.mode = modeGameover
			} else {
				deleteIdx := g.showTatsuIndexes[0]
				pt := g.tatsus[deleteIdx]
				if pt.answer(false) {
					g.score += pt.point()
				} else {
					g.score -= pt.point()
				}

				g.showTatsuIndexes = g.showTatsuIndexes[1:]
			}
		}
		if g.isKeyEscapePressed() {
			g.init()
			g.mode = modeTitle
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
	text.Draw(screen, fmt.Sprintf("Score: %d", g.score), arcadeFont, 20, 30, color.White)
	text.Draw(screen, fmt.Sprintf("Level: %d", g.level()), arcadeFont, 240, 30, color.White)
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
		titleX = 380
		titleY = 330
	)

	switch g.mode {
	case modeTitle:
		text.Draw(screen, "PRESS SPACE", arcadeFont, titleX, titleY, color.White)
	case modeGameover:
		text.Draw(screen, "GAME OVER", arcadeFont, titleX, titleY, color.White)
	}

	if g.pause {
		text.Draw(screen, "RESTART with s", arcadeFont, titleX, titleY, color.White)
	}
}

func (g *Game) drawTatsus(screen *ebiten.Image) {
	for _, t := range g.tatsus {
		if t.visible {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(t.x), float64(t.y))
			op.Filter = ebiten.FilterLinear
			screen.DrawImage(t.i, op)

			tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
			if err != nil {
				log.Fatal(err)
			}
			const dpi = 72
			arcadeFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    fontSize + 14,
				DPI:     dpi,
				Hinting: font.HintingFull,
			})

			text.Draw(screen, fmt.Sprintf("%d符", t.dummyFu), arcadeFont, t.x+(t.i.Bounds().Dx()/2)-30, t.y-30, color.White)
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

func (g *Game) isSPressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		return true
	}
	return false
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

func (g *Game) isKeyEscapePressed() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return true
	}
	return false
}

func (g *Game) level() int {
	return g.speed - 2
}
