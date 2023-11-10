package game

const (
	tatsuWidth  = 66
	tatsuHeghit = 99

	minTatsuDist  = 50
	maxTatsuCount = 3
)

type tatsu struct {
	x       int
	y       int
	visible bool
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
