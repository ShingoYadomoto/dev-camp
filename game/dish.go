package game

const (
	dishWidth  = 25
	dishHeghit = 50

	minDishDist  = 50
	maxDishCount = 3
)

type dish struct {
	x       int
	y       int
	visible bool
}

func (t *dish) move(speed int) {
	t.x -= speed
}

func (t *dish) show() {
	t.x = ScreenX
	t.y = groundY - dishHeghit
	t.visible = true
}

func (t *dish) hide() {
	t.visible = false
}

func (t *dish) isOutOfScreen() bool {
	return t.x < -50
}
