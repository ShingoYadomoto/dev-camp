package game

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	ScreenWidth  = 420
	ScreenHeight = 600
	boardSize    = 4
)

// Game represents a game state.
type Game struct {
	input      *Input
	board      *Board
	boardImage *ebiten.Image
}

// NewGame generates a new Game object.
func NewGame() (*Game, error) {
	g := &Game{
		input: NewInput(),
	}
	var err error
	g.board, err = NewBoard(boardSize)
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update(*ebiten.Image) error {
	g.input.Update()
	if err := g.board.Update(g.input); err != nil {
		return err
	}
	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		w, h := g.board.Size()
		g.boardImage, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)
}