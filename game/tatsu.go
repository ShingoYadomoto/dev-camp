package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	tatsuWidth  = 66
	tatsuHeghit = 99

	minTatsuDist  = 50
	maxTatsuCount = 50
)

type tatsu struct {
	i *ebiten.Image

	x       int
	y       int
	visible bool

	dummyFu   uint8
	correctFu uint8
}

func (t *tatsu) move(speed int) {
	t.x -= speed
}

func (t *tatsu) show() {
	t.x = ScreenX
	t.y = groundY - tatsuHeghit
	t.visible = true
}

func (t *tatsu) hide() {
	t.visible = false
}

func (t *tatsu) isOutOfScreen() bool {
	return t.x < -50
}
